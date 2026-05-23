package pipeline

import (
	"time"

	"github.com/gallowaysoftware/vibe/vamp"
)

// Config drives one episode's generation.
type Config struct {
	// Topic is the day's subject ("burnout", "asking for a raise",
	// "early-stage relationship advice"). The script + showrunner
	// stages both reference it.
	Topic string
	// EpisodeNumber is the 1-indexed episode in the series. Used
	// only for metadata + RSS; doesn't affect generation (each
	// episode is standalone).
	EpisodeNumber int
}

// Build constructs the per-episode pipeline.
//
// Stages:
//
//	news_search (webhook)             → news_raw.json (SearXNG /search results)
//	format_news (render)              → news.md       (top headlines + snippets, markdown)
//	write_script (text)               → script_draft.md
//	edit_script (text)                → script.md     (revised, dribble cut)
//	showrunner   (text, json)         → script.json
//	aria_segments (render)            → aria.json    (filter host=aria)
//	atlas_segments (render)           → atlas.json   (filter host=atlas)
//	disclaimer_segments (render)      → disclaimer.json (filter host=disclaimer)
//	cast_aria (audio foreach)         → audio/{{.segment.id}}.wav
//	cast_atlas (audio foreach)        → audio/{{.segment.id}}.wav
//	cast_disclaimer (audio foreach)   → audio/{{.segment.id}}.wav
//	compose_mix_script (render)       → mix_script.json
//	mix_episode (mix)                 → episode.mp3
func Build(cfg Config) (*vamp.Pipeline, error) {
	p := vamp.New("iitn-episode").
		Describe("Generate one episode of \"It's Important to Note\" — two AI hosts give confidently wrong advice on a single topic + react to today's news.")

	p.Input("topic", vamp.Required(), vamp.WithDefault(cfg.Topic),
		vamp.Describe("Today's topic (e.g. 'burnout', 'asking for a raise')."))

	p.RequireService("kokoro-tts", "http://127.0.0.1:8880",
		"Kokoro-FastAPI TTS — provides af_bella / am_adam / am_eric voices.",
		"vibe start tts_kokoro")
	p.RequireService("searxng", "http://127.0.0.1:14002",
		"SearXNG — pulled for the News React segment.",
		"vibe start searxng")
	p.RequireGPUMemory("~30GB during write_script + edit_script + showrunner")
	p.RequireDiskSpace("~10MB per episode")
	p.CapabilityModel("long_form", vamp.ModelHint{
		MinParams: "27B", MinContext: 131072,
		SuggestedModel: "qwen3.6-27b-mtp-q6_k",
	})

	// ---- News fetch (no GPU, no cache — every episode reacts to
	//      whatever SearXNG returns at run time) ----

	newsSearch := p.Webhook("news_search").
		URL("http://127.0.0.1:14002/search?q=today+news&format=json").
		Method("GET").
		Output("news_raw.json")

	formatNews := p.Render("format_news").
		After(newsSearch).
		Prompt(`{{ $raw := parseJSON .stages.news_search.output -}}
{{ range $i, $r := (index $raw "results") -}}
{{ if lt $i 12 -}}
- **{{ index $r "title" }}** — {{ index $r "content" }}
{{ end -}}
{{ end -}}`).
		Output("news.md")

	// ---- LLM stages ----

	script := p.Text("write_script").
		Capability("long_form").
		After(formatNews).
		PromptFS(PromptsFS, "script.md").
		Output("script_draft.md").
		Param("temperature", 0.85).
		// 16384 because v2 targets ~1500-2500 spoken words across 8
		// sections (cold open / sponsor A / news react / topic /
		// sponsor B / listener mail / confidence / outro). The 8192
		// v1 budget regularly truncated mid-disclaimer.
		Param("max_tokens", 16384).
		Retry(&vamp.RetryPolicy{
			MaxAttempts:    3,
			InitialBackoff: 5 * time.Second,
			MaxBackoff:     30 * time.Second,
			RetryOn:        []string{"transient"},
		})

	// Editor pass: revise the draft so the comedy lands. The user
	// feedback that triggered this stage (2026-05-23): "sometimes
	// it hits, sometimes it just reads as hallucinations. It's
	// funny when it's deliberately a bit off, it's boring when it
	// feels like it's just rambling." Temperature is low — this is
	// editorial, not generative.
	editScript := p.Text("edit_script").
		Capability("long_form").
		After(script).
		PromptFS(PromptsFS, "editor.md").
		Output("script.md").
		Param("temperature", 0.35).
		Param("max_tokens", 16384).
		Retry(&vamp.RetryPolicy{
			MaxAttempts:    3,
			InitialBackoff: 5 * time.Second,
			MaxBackoff:     30 * time.Second,
			RetryOn:        []string{"transient"},
		})

	// Cover-art pipeline: an LLM composes a one-sentence SDXL prompt
	// from the topic + edited script, then a ComfyUI stage renders it
	// at 1024² via SDXL-Turbo. The mix stage at the end embeds the
	// resulting PNG into the .m4b as attached_pic.
	//
	// SDXL-Turbo is fast (1 step, ~3-4s/image at 1024) but produces
	// visibly worse output than non-turbo SDXL. If quality is the
	// bottleneck, swap the workflow for the full SDXL one — the
	// pipeline shape doesn't change.
	composeCover := p.Text("compose_cover").
		Capability("long_form").
		After(editScript).
		PromptFS(PromptsFS, "cover.md").
		Output("cover_prompt.txt").
		Param("temperature", 0.6).
		// Cover prompt is short — one line — so the budget is tiny.
		// Bumped above the literal need to absorb Qwen3's <think>
		// preamble without clipping the actual prompt.
		Param("max_tokens", 2048).
		Retry(&vamp.RetryPolicy{
			MaxAttempts:    3,
			InitialBackoff: 5 * time.Second,
			MaxBackoff:     30 * time.Second,
			RetryOn:        []string{"transient"},
		})

	generateCover := p.ComfyUI("generate_cover").
		Capability("image_gen").
		After(composeCover).
		WorkflowFS(WorkflowsFS, "sdxl_turbo.json").
		// Node 4 in the sdxl_turbo workflow is the positive-prompt
		// CLIPTextEncode; node 6 is the EmptyLatentImage's width/height;
		// node 7 is the KSampler (seed left at the workflow's default
		// since the prompt itself varies per episode).
		Parameter("4.text", "{{ .stages.compose_cover.output }}").
		Parameter("6.width", "1024").
		Parameter("6.height", "1024").
		Output("cover.png")

	showrunner := p.Text("showrunner").
		Capability("long_form").
		After(editScript).
		PromptFS(PromptsFS, "showrunner.md").
		OutputFormatJSON().
		Output("script.json").
		Param("temperature", 0.3).
		Param("max_tokens", 16384).
		Retry(&vamp.RetryPolicy{
			MaxAttempts:    3,
			InitialBackoff: 5 * time.Second,
			MaxBackoff:     30 * time.Second,
			RetryOn:        []string{"transient", "invalid_output"},
		})


	// ---- Split segments by host (audio stage's Voice is per-stage,
	//      not per-iteration, so we run three audio stages — one per
	//      voice — each foreaching only its own segments) ----

	ariaSegs := p.Render("aria_segments").
		After(showrunner).
		Prompt(`{{ filterByValue "host" "aria" (toJSON (index (parseJSON .stages.showrunner.output) "segments")) }}`).
		Output("aria.json").
		OutputFormatJSON()

	atlasSegs := p.Render("atlas_segments").
		After(showrunner).
		Prompt(`{{ filterByValue "host" "atlas" (toJSON (index (parseJSON .stages.showrunner.output) "segments")) }}`).
		Output("atlas.json").
		OutputFormatJSON()

	disclaimerSegs := p.Render("disclaimer_segments").
		After(showrunner).
		Prompt(`{{ filterByValue "host" "disclaimer" (toJSON (index (parseJSON .stages.showrunner.output) "segments")) }}`).
		Output("disclaimer.json").
		OutputFormatJSON()

	// ---- TTS via Kokoro (single voice per stage; multi-voice
	//      routing within one stage is a pending vamp enhancement) ----

	ariaAudio := p.Audio("cast_aria").
		Capability("tts").
		After(ariaSegs).
		Foreach(ariaSegs, "segment").
		Engine(vamp.AudioEngineKokoro).
		Voice("af_bella").
		TextTemplate("{{.segment.text}}").
		Output("audio/{{.segment.id}}.wav")

	atlasAudio := p.Audio("cast_atlas").
		Capability("tts").
		After(atlasSegs).
		Foreach(atlasSegs, "segment").
		Engine(vamp.AudioEngineKokoro).
		Voice("am_adam").
		TextTemplate("{{.segment.text}}").
		Output("audio/{{.segment.id}}.wav")

	disclaimerAudio := p.Audio("cast_disclaimer").
		Capability("tts").
		After(disclaimerSegs).
		Foreach(disclaimerSegs, "segment").
		Engine(vamp.AudioEngineKokoro).
		Voice("am_eric").
		TextTemplate("{{.segment.text}}").
		Output("audio/{{.segment.id}}.wav")

	// ---- Compose mix script: list voice_segments in master order
	//      (from showrunner.segments), pointing at audio files
	//      written by whichever cast stage handled each host. ----

	mixScript := p.Render("compose_mix_script").
		After(showrunner, ariaAudio, atlasAudio, disclaimerAudio).
		Prompt(`{{ $script := parseJSON .stages.showrunner.output -}}
{"voice_segments": [
{{- range $i, $seg := index $script "segments" -}}
  {{- if $i }}, {{ end -}}
  "audio/{{ index $seg "id" }}.wav"
{{- end }}
]}`).
		Output("mix_script.json").
		OutputFormatJSON()

	// Switched to .m4b in v2: gets attached_pic cover embedding for
	// free (vamp's mix executor already handles m4b+cover; the mp3
	// ID3 path was tagged v2-deferred there). Audiobookshelf / Plex /
	// Apple Books all prefer m4b for spoken-word audio.
	p.Mix("mix_episode").
		After(mixScript, ariaAudio, atlasAudio, disclaimerAudio, generateCover).
		ScriptFile("mix_script.json").
		CoverImage("cover.png").
		LoudnessTarget(-16).
		Output("episode.m4b")

	return p.Build()
}
