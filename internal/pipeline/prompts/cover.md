You are writing the SDXL prompt for the cover art of today's "It's
Important to Note" episode. The cover should be a *visual metaphor*
for the topic — concrete, slightly off, with the show's recurring
visual register.

# Today's topic

{{ .inputs.topic }}

# A sample of today's script (use it to ground the metaphor)

{{ .stages.edit_script.output }}

---

# Cover constraints (all mandatory)

1. **One concrete subject.** Not "concept of burnout" — a specific
   wilting houseplant on a fluorescent-lit desk, or a single empty
   chair next to a softly humming server. The cover is a single
   tableau, not a collage.

2. **Show visual register** — every cover carries these to keep the
   feed coherent:
   - Muted pastel palette: dusty rose, sage green, cream, "existential
     beige", soft lavender, pale yellow. No saturated reds, no neon.
   - Soft natural lighting, slightly oversaturated highlights. Avoid
     harsh shadows.
   - Mid-century-meets-near-future aesthetic: a faint corporate
     dread under the warmth. The lighting feels like a wellness app
     ad photographed in 1974.
   - Slightly off-kilter framing, a touch of analog noise, faintly
     unsettling without being dark.

3. **Topical metaphor done sideways.** The subject should evoke the
   topic but through an *LLM-flavored* misreading. Examples:
   - "burnout" → a Roomba that has stopped halfway through a circle,
     a single dim LED blinking.
   - "groceries" → a single bruised apple on a chrome scale in
     fluorescent light, the price tag reading "0.00".
   - "asking for a raise" → a corporate slide deck open on a beige
     desk, the title slide reads only "(loading)".
   - "in-laws" → an empty chair at a holiday dinner table, place
     setting still wrapped in plastic.
   The metaphor is wry, not bleak. The mood is curious, not
   despairing.

4. **No legible text anywhere.** SDXL cannot render text reliably.
   If text appears in the subject (a slide, a label, a sign), it
   should be implied or rendered as soft blur. Do not ask the model
   to produce specific words.

5. **No human faces.** The show is hosted by AI bots; the cover
   should not feature humans. Hands, silhouettes, and figures from
   behind are acceptable.

6. **Album-art friendly composition.** Centered subject, comfortable
   margins, looks good cropped to a square thumbnail.

# Output

Return ONLY the SDXL positive prompt as a single line. No commentary,
no preamble, no negative prompt (the workflow handles that). Format:

```
<concrete subject>, <one or two sensory details>, muted pastel palette, soft natural lighting, mid-century-meets-near-future aesthetic, analog film grain, centered composition, no text, no faces, slightly off-kilter
```

Example outputs (do not copy verbatim — emulate the shape):

- `a single wilting houseplant on a beige cubicle desk under flat fluorescent light, the leaves drooping toward a half-full mug of cold tea, muted pastel palette, soft natural lighting, mid-century-meets-near-future aesthetic, analog film grain, centered composition, no text, no faces, slightly off-kilter`
- `a Roomba stalled in the middle of a sun-bleached living room, a single LED dimly blinking, a half-vacuumed circle of beige carpet around it, muted pastel palette, soft natural lighting, mid-century-meets-near-future aesthetic, analog film grain, centered composition, no text, no faces, slightly off-kilter`

First byte of output: a lowercase letter (the start of the subject).
No quotation marks, no JSON wrapper, no explanation.
