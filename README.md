# It's Important to Note

A satirical short-form podcast. Two AI hosts give confidently wrong life advice on a single topic per episode. The show leans into the "LLM-generated content is slop" critique by being exactly that, on purpose, with affection.

Every voice tic is intentional: the em-dashes, the "Let's unpack this," the hallucinated citations to studies that don't exist, the helpful warmth applied to dark advice. The hosts know they are AIs and they are extremely excited about this.

## Format (per episode, ~5-8 min)

1. **Cold open** — Aria opens; Atlas finishes the tagline with a fresh variant each time.
2. **Listener mail** *(sometimes)* — Aria reads a fictional letter, also AI-generated, also questionable.
3. **Today's sponsor** — fictional brand read with a promo code ("Use code IMPORTANT for 12% off").
4. **The topic** — alternating dialogue, 3-4 sub-points, at least one hallucinated citation per host.
5. **Confidence meter** — hosts rate their own advice. Always 9.7/10.
6. **Outro + tiny disclaimer** — legal-disclaimer-style fast-talk by a separate voice at the end.

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

# Bring up everything the pipeline needs (searxng + tts_kokoro + comfyui + long_form).
# Reads the pipeline's RequireService / RequireProfile declarations and starts each
# via vibe in one shot. --skip-active if the GPU is busy and you just want the sidecars.
iitn activate

# Generate the next episode (auto-picks an unused topic).
iitn next

# Check what's running, what's missing.
iitn doctor

# Pick a topic yourself.
iitn next --topic "how to apologize properly"

# List + replay.
iitn list
```

Episodes land at `~/.local/state/iitn/episodes/NNN/episode.m4b`.

## Architecture

```
write_script (LLM)        → script.md
showrunner (LLM)          → script.json
aria_segments / atlas_segments / disclaimer_segments (renders)
   ↓        ↓        ↓
cast_aria / cast_atlas / cast_disclaimer (Kokoro foreach)
   ↓
compose_mix_script (render) → mix_script.json
mix_episode (mix stage)     → episode.mp3
```

5 stages, ~2 min compute per episode. Easy to cron daily so a new episode lands in your podcast feed every morning.

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
| Pipeline (5 stages) | wired end-to-end |
| Prompts: script + showrunner | drafted, tone-tuned, awaiting first live read |
| Topic catalog | 52 entries shipped |
| RSS feed generation | TODO — straightforward, deferred until first episode lands |
| Apple Podcasts submission | TODO — needs a public hosting URL |
| Per-host voice routing (single audio stage) | v2 — needs templateable `Audio.Voice` |
| Music / SFX (theme jingle) | v2 — show could open with a Suno-style 5-second AI-generated jingle |

## License

MIT. The show's content is generated; nobody owns it but you. The whole point is your private library of AI-generated comedy, optionally shared.
