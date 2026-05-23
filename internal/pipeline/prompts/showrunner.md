You are the showrunner. The writer handed you a finished script for
"It's Important to Note" — alternating dialogue between Aria and
Atlas, plus a final disclaimer paragraph. Your job: break it into
voice segments, assign each segment to the correct voice, and
sprinkle paralinguistic tags ([sigh], [chuckle], [pause]) where
the line calls for them.

# Script

{{ .stages.edit_script.output }}

# Voice cast (fixed for this show)

- **aria** — host A, cheerful + warm. TTS voice: `af_bella`.
- **atlas** — host B, slightly deeper + philosophical. TTS voice:
  `am_adam`.
- **disclaimer** — the fast/over-caveated legal disclaimer at the
  end. TTS voice: `am_eric` (higher pitch, more clipped).

# Output schema

Return ONLY a single JSON object. No prose, no markdown fences.
First byte: `{`.

```
{
  "segments": [
    {
      "id": "seg_000",
      "host": "aria",
      "voice_id": "af_bella",
      "text": "<one paragraph of Aria's dialogue, verbatim from the script. Sprinkle [pause] / [chuckle] / [sigh] tags where natural.>"
    },
    {
      "id": "seg_001",
      "host": "atlas",
      "voice_id": "am_adam",
      "text": "..."
    },
    ... (typically 20-40 segments — every distinct paragraph of one host is a segment)
    {
      "id": "seg_NNN",
      "host": "disclaimer",
      "voice_id": "am_eric",
      "text": "<the entire disclaimer paragraph as a single segment — no paragraph breaks>"
    }
  ]
}
```

# Rules

- **Every `**Aria:**` paragraph is one segment** with host=aria.
- **Every `**Atlas:**` paragraph is one segment** with host=atlas.
- **The final disclaimer is one segment** with host=disclaimer.
- **Preserve the script's order strictly.** The downstream mix
  concatenates segments in array order.
- **Strip the speaker prefix** (`**Aria:** `) from the text — the
  voice routing handles attribution. Keep everything else verbatim.
- **Paralinguistic tags** Chatterbox/Kokoro accept inline:
  `[pause]`, `[sigh]`, `[chuckle]`, `[laugh]`, `[whispers]`. Use
  sparingly — the LLM-bot register is uncluttered. One or two per
  episode is enough. Disclaimer voice gets NO tags (it's machine-fast).
- **Don't add new dialogue** that wasn't in the script.
- **Don't merge consecutive same-host paragraphs** — keep them as
  separate segments. The mix stage benefits from natural pauses
  between paragraphs in the same voice.
- **id**: `seg_000`, `seg_001`, ... 3-digit padded. Used in
  downstream filenames.
