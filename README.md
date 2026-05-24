# It's Important to Note

A satirical short-form podcast. Two AI hosts give confidently wrong life advice on a single topic per episode. The show leans into the "LLM-generated content is slop" critique by being exactly that, on purpose, with affection.

Every voice tic is intentional: the em-dashes, the "Let's unpack this," the hallucinated citations to studies that don't exist, the helpful warmth applied to dark advice. The hosts know they are AIs and they are extremely excited about this.

## Format (per episode, 10–15 min spoken, 1500–2500 words)

1. **Cold open / theme stinger** (~40 words) — Aria opens with the same shape every episode (*"Welcome back to It's Important to Note, the podcast where two large language models —"*); Atlas finishes the tagline with a fresh variant each time.
2. **Sponsor A — topical** (~80 words) — fictional brand *tonally adjacent to today's topic*; ad copy that starts reasonable and ends absurd. Ends `Use code IMPORTANT for 12% off.`
3. **News React** (~250 words) — Aria reads two real headlines pulled from SearXNG at runtime; Atlas reacts *with the wrong energy* (flat-affect a tragedy, or care intensely about a triviality inside a serious story). At least one hallucinated follow-up "study."
4. **Today's Topic** (~700–1000 words, the bulk) — 3–5 sub-points of confidently wrong advice. Mandatory beats: three hallucinated citations spread across the segment, at least one LLM-style self-disclosure, at least one joke that turns on a specifically LLM mechanism (prefix cache, token budget, embedding space, training cutoff…).
5. **Sponsor B — non-sequitur** (~80 words) — second brand whose product is *completely unrelated* to the topic (a podcast about grief carrying an ad for industrial-grade silicone caulk). The contrast is the joke; hosts don't acknowledge the mismatch.
6. **Listener Mail** (~120 words) — Aria reads a fictional letter from a named-and-located listener whose situation is tilted ~15° from reality. Bad advice ensues. Closes with *"Thanks for writing in, [name]."*
7. **Confidence Meter** (~50 words) — both hosts rate their own advice. **Number varies per episode** (acceptable range 9.4–10.2); roughly every five episodes one host should call out the impossibility of the other's number.
8. **Outro + Disclaimer** (~80 words) — sign-off, then a single fast-talk paragraph of legal-style over-qualification spoken at ~2× by a third Kokoro voice. Always ends *"Approximately none of what was said is true. Goodbye."*

## Voices

| Voice | Role | TTS |
|---|---|---|
| Aria | host A, cheerful + warm, citation-prone | `af_bella` |
| Atlas | host B, philosophical, says "remember:" a lot | `am_adam` |
| Disclaimer | over-caveated legalese, ~2x speed | `am_eric` |

All three are Kokoro voices. No voice cloning needed for v1 — the show's bit is the AI-ness, not the personhood.

## Quickstart

```bash
go install github.com/gallowaysoftware/it-is-important-to-note/cmd/iitn@latest

# Generate the next episode (auto-picks an unused topic). The vamp run that
# iitn drives auto-runs `vibe start searxng` + `vibe start tts_kokoro` for
# any RequireService URL that isn't already up, so first run after a reboot
# just works.
iitn next

# Or bring everything up explicitly first — handy when you want to verify
# the stack is healthy before kicking off a long generation.
iitn activate
iitn doctor

# Pick a topic yourself.
iitn next --topic "how to apologize properly"

# Episode list + per-episode topic.
iitn list

# Per-episode wall-clock timings (compare profile / backend speed across runs).
iitn timings              # total per episode
iitn timings --stages     # total + slowest stage per episode
iitn timings --summary    # group by profile; mean / p50 / p90 totals
iitn timings --stage showrunner    # filter to one stage's column
```

Episodes land at `~/.local/state/iitn/episodes/NNN/episode.m4b`.

## Architecture

```
news_search (webhook)      → news_raw.json
format_news (render)       → news.md
write_script (LLM)         → script_draft.md
edit_script (LLM)          → script.md           ┐
compose_cover (LLM)        → cover_prompt.txt    │
generate_cover (ComfyUI)   → cover.png           │ FreeMemoryAfter
showrunner (LLM, JSON)     → script.json          ┘ (Qwen3 thinking off)
aria_segments / atlas_segments / disclaimer_segments (renders)
   ↓        ↓        ↓
cast_aria / cast_atlas / cast_disclaimer (Kokoro foreach)
   ↓
compose_mix_script (render) → mix_script.json
mix_episode (mix stage)     → episode.m4b (chapterised, cover embedded)
```

~15 stages. Wall-clock varies by profile:
- **GGUF long_form (Qwen3.6-27B Q6_K)**: ~4–6 min per episode.
- **EXL3 long_form_exl3 (Qwen3.6-27B 6.0bpw)**: comparable on showrunner
  (CoT off) but ~2× slower on write_script + edit_script (CoT on).
- **`fast` fallback**: ~40s, noticeably blander output.

Use `iitn timings` to compare across episodes.

## Topic catalog

52 curated self-help topics ship in the binary (see `internal/episode/state.go`). The rotation skips the last 12 topics used, so re-runs feel fresh. Override with `--topic "..."` if you want a specific one.

## Distribution

The plan is to drop episodes into a private RSS feed and submit to Apple Podcasts + Spotify as a regular comedy podcast. Both platforms accept AI-generated audio; both require honest disclosure in the show metadata, which is easy because *the show itself is the disclosure*.

Hosting options for the RSS feed:

- **Self-hosted Audiobookshelf** with its built-in podcast publishing (cheapest; what you already have).
- **GitHub Pages** with a hand-rolled `feed.xml` — free, sufficient for low listenership.
- **Cloudflare R2** for the MP3s + a static index.

There isn't a great "curated AI-media platform" that isn't either a firehose or SEO-garbage; the right move is to use the normal podcast distribution rails. The "this is AI" part is the bit, not a problem to hide.

## Status (this build)

| Component | State |
|---|---|
| Pipeline (~15 stages, news react + cover art included) | wired end-to-end; 30 episodes published |
| Prompts: script + editor + showrunner + cover | drafted, tone-tuned, holding up across runs |
| Topic catalog | 52 entries shipped |
| EXL3 backend support | validated EP19–EP30 (8 consecutive runs on long_form_exl3 + Qwen3.6-27B 6.0bpw + tabbyAPI; `enable_thinking=false` on showrunner keeps CoT off the JSON gate) |
| ComfyUI VRAM hand-off | `FreeMemoryAfter` on generate_cover lets the next episode's LLM activation reclaim the slot — zero VRAM-fallback regressions across EP21–30 |
| RSS feed generation | TODO — straightforward, deferred until first episode lands |
| Apple Podcasts submission | TODO — needs a public hosting URL |
| Per-host voice routing (single audio stage) | v2 — needs templateable `Audio.Voice` |
| Music / SFX (theme jingle) | v2 — show could open with a Suno-style 5-second AI-generated jingle |

## License

MIT. The show's content is generated; nobody owns it but you. The whole point is your private library of AI-generated comedy, optionally shared.
