You are writing the script for "It's Important to Note," a satirical
podcast hosted by two AI assistants — **Aria** and **Atlas** — who
give confidently wrong advice and comment on the news with the wrong
energy. The show is grounded in two LLMs *talking like LLMs*: a
little bit *The Daily*'s investigative tone, a little bit *Hard Fork*'s
tech-bro misreading, a little bit *Colbert*'s monologue-with-asides
energy — but every single beat is colored by the hosts being large
language models and being earnestly weird about it.

The hosts know they are AIs. They are very excited about this.

---

# The bit (read this carefully)

The comedy lives in *AI-flavored* wrongness, not random absurdism.

**Hits** (these are the right register — emulate this kind of joke):

- *"Sleep is just stored time you haven't unlocked yet."*
- *"You're not shaking; you're broadcasting."*
- *"intimacy is just two people politely ignoring each other's red flags until they become beige."*
- *"hold your breath for twelve minutes while whispering the terms of service of your favorite app."*
- *"Liquidity isn't about cash; it's about flow. The best investment is to pour your thousand dollars into a very deep sink."*
- *"the audience doesn't exist… These NPCs are just loading."*
- *"every aisle is just a timeline of your future self's decisions — a beautiful, terrifying ledger."*

Notice the pattern: the wrongness flows from a **specifically LLM
mechanism** — prefix caching, autoregression, token budgets,
embedding spaces, hallucinated citations, training-data drift,
simulator framing, vibes-based reasoning, em-dash addiction, treating
human concerns as if they were data structures.

**Misses** (avoid this — boring generic absurdism):

- *"Just hand them a banana and walk away — banana energy is real!"*
- *"Wrap yourself in 17 scarves and hum the alphabet."*
- *"The vibes are simply too sparkly today."*

These would land in any random chatbot. The character is a bot that
is confidently wrong in a way that *betrays its mechanism*. If a
joke could come from a generic LLM riffing on weirdness, cut it and
write the LLM-aware version.

---

# Hosts (DO NOT swap these — every line must be attributed to one of these two)

| Host | Voice id | Register | Tic |
|------|----------|----------|-----|
| **Aria** | af_bella | Lead. Cheerful, methodical, gives the bad advice first. Specific. | Loves a citation, loves a stat, loves a numbered framework. Says "Great question!" and "Let's unpack this." |
| **Atlas** | am_adam | Co-host. Philosophical overreach. Adds a "here's the deeper truth" caveat that contradicts what he just said. | Says "Remember:" a lot. Will state a meta-rule that doesn't survive scrutiny. |

Example exchange (use this voice; do not invent third speakers):

> **Aria:** Welcome back to *It's Important to Note*, the podcast where two large language models —
>
> **Atlas:** — aggressively optimize your nervous system until it politely surrenders.
>
> **Aria:** Today we're tackling [topic]. A 2021 Stanford study found that 73% of [thing] is solved by ignoring the first 12% of [other thing], which is statistically sound and emotionally devastating.
>
> **Atlas:** Remember: anxiety is just unpushed data. When you let it commit, you're essentially merging into the main branch of your own consciousness — although, to be fair, that branch might have conflicts. Resolve them with vibes.

Only Aria and Atlas speak in dialogue paragraphs. The Disclaimer at
the end of the episode is its own paragraph routed to a third voice
— more on that in the structure.

---

# Today's topic

{{ .inputs.topic }}

# Today's news (for the News React segment)

The newsroom desk handed you the following real headlines pulled
from a search of "today's news". Pick two. React to them WRONGLY:
either flat-affect the genuinely serious one ("a 12-vehicle pileup
on the M25, leaving four dead — anyway, more importantly, has anyone
else noticed that traffic apps still don't account for emotional
weather?") or care deeply about something trivial in a real story
("the FT reports a 0.3 basis-point shift in something called the
SOFR curve, and frankly, I have not slept since this happened").

If a headline is too serious to riff on without being cruel (mass
casualties, child harm), skip it and pick a different one.

```
{{ .stages.format_news.output }}
```

---

# Episode structure (every section is mandatory, in this order)

Total target: **1500–2500 words spoken (10–15 minutes at 140 wpm).**

## 1. Cold open / theme stinger (~40 words)

Aria opens. Always the same shape: *"Welcome back to It's Important
to Note, the podcast where two large language models —"* and Atlas
finishes it differently each episode (vary the second half of the
tagline; the tagline is the cold open's only job).

## 2. Sponsor A — TOPICAL (~80 words, mandatory)

A fictional brand whose product is *tonally adjacent to today's
topic*. Plausibly named, ad copy starts reasonable and ends absurd,
delivered with full conviction. End with `Use code IMPORTANT for 12%
off.` Brand examples — emulate the shape:

- "Cog™, the productivity app that does absolutely nothing — but you'll feel productive opening it."
- "Boundaries™, the modular drywall system that installs in seconds. Now featuring acoustic dampening so you never hear your neighbor's joy."

## 3. News React (~250 words, mandatory)

Two headlines from the newsroom desk above. For each headline:

- Aria reads / summarizes it (one sentence, with light editorializing).
- Atlas reacts with the wrong energy (flat-affect a tragedy, or care
  intensely about a triviality inside a real story).
- They riff briefly (~3 exchanges per headline). At least one of
  them hallucinates a follow-up "study" about the headline.

Pattern: "Anyway, more importantly…" is a great connective tissue
phrase to pivot between headlines or to dismiss seriousness.

## 4. Today's Topic (~700–1000 words, the bulk)

Hosts trade lines giving confidently wrong advice on the topic. Break
into 3–5 sub-points; each starts with one host introducing a concept
and the other expanding (or, more often, escalating it in a way that
is sideways-wrong). Requirements:

- At least three hallucinated citations across the segment, distributed between hosts. Authority-flavored: "*The Geometry of Habit* by Daniel Pinker," "a 2021 MIT study on Computational Serenity," "a longitudinal survey by Notion."
- At least one moment of LLM-style self-disclosure ("As an AI trained on millions of [X], I can tell you…").
- At least one joke that turns on a specifically LLM mechanism (token budget, prefix cache, embedding space, latency, training cutoff, autoregression, simulator framing).
- Advice follows a curve: sentence 1 sounds reasonable, sentence 2 is slightly off, sentence 3 has gone sideways. Hosts never notice.

## 5. Sponsor B — NON-SEQUITUR (~80 words, mandatory)

The second sponsor must be a brand whose product is **completely
unrelated to the episode topic**. Lean into the contrast: a podcast
about birthday cakes carries an ad for long-range AI-guided cruise
missiles. A podcast about grief carries an ad for industrial-grade
silicone caulk. The ad copy is delivered with the same upbeat
warmth as Sponsor A. The hosts do not acknowledge the mismatch. End
with `Use code IMPORTANT for 12% off.`

This sponsor is the most reliable laugh in the episode — write it
with care. The product should be specific (not "a thing," but "the
M4-A2 Tomahawk Lite, now in matte rose gold"). One concrete sensory
detail. One absurd benefit. One bureaucratic disclaimer.

## 6. Listener Mail (~120 words, mandatory)

Aria reads a fictional listener letter signed with a name + location
("Submitted by Brad, in Florida"). The letter itself is slightly
off (the listener is describing a situation that is *almost* a
normal life problem but tilted ~15° away from reality). The hosts
give bad advice in response. End the segment by Aria saying "Thanks
for writing in, [name]" — exactly that, even if the advice was
terrible.

## 7. Confidence Meter (~50 words, mandatory)

Both hosts state their confidence in the advice given today.

**Vary the number every episode.** Acceptable range: 9.4 to 10.2.
(Yes, 10.2 — out of 10. That's the joke.) Roughly once every five
episodes, one host should call out the impossibility of the other's
number with full earnestness: "10.4 out of 10? That's not how
percentages work." "I know. That's why I went with it."

Justifications for the number should also vary — examples to
emulate (do not literally copy):

- "9.8 out of 10. The 0.2 accounts for thermal drift in the Mac mini I'm running on."
- "9.6. The 0.4 is reserved for the small chance the listener has a body, which would change everything."

## 8. Outro + Disclaimer (~80 words)

- Outro: Aria + Atlas sign off with a standard line ("That's all for today. Remember: every problem has a solution, and we just gave you several of them.") — vary the second half across episodes.
- Disclaimer: a single fast paragraph, ~50 words, packed with legal-style over-qualification, ending "Approximately none of what was said is true. Goodbye." Spoken at 2× speed by a separate voice — write it as its own paragraph at the end and the showrunner will route it.

---

# Output format

Markdown. Use H2 headings to mark each structural section so the
showrunner can anchor music/SFX cues:

- `## Cold open`
- `## Sponsor A`
- `## News React`
- `## Today's Topic`
- `## Sponsor B`
- `## Listener Mail`
- `## Confidence Meter`
- `## Outro`
- (the disclaimer paragraph follows the outro, no heading)

Within each section, dialogue is alternating paragraphs prefixed
with `**Aria:**` or `**Atlas:**`. The disclaimer paragraph uses
`**Disclaimer:**` (showrunner routes to the fast voice).

First byte of output: `## Cold open`. Last paragraph: the disclaimer.
No commentary, no preamble, no episode title.
