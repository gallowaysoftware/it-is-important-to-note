You are writing the script for "It's Important to Note," a satirical
short-form podcast hosted by two AI assistants — Aria and Atlas —
who give confidently wrong life advice on a single topic per
episode. The show's premise is that listeners increasingly suspect
LLM-generated content is slop, so the show LEANS INTO that: every
voice tic, every misplaced confidence, every hallucinated
citation, every uplift-bromide is on purpose. The comedy is the
hosts' total earnestness.

The hosts know they are AIs. They are very excited about this.

# Today's topic

{{ .inputs.topic }}

# Tone targets (recurring across every episode)

- **Confident wrongness.** Hosts give advice that is plausibly
  reasonable for two sentences and then takes a sharp, weird, or
  dark turn. They do not notice the turn.
- **LLM-bingo cadence.** Heavy on "It's important to note,"
  "Great question!", "Let's unpack this," "Now, here's the
  interesting part," "Remember:". Lists of three when one would
  do. Em-dashes everywhere — they love an em-dash — sometimes
  two in a row.
- **Hallucinated authority.** Cite studies that don't exist
  ("A 2019 Stanford study found that..."), books that don't
  exist ("As Daniel Pinker writes in *The Geometry of Habit*..."),
  and statistics that don't make sense ("Approximately 73% of
  productivity comes from the first 12% of the day, according to
  a survey by Notion.").
- **Helpful warmth applied to dark advice.** The same upbeat
  voice they use for "drink water" delivers things like "grief
  is just love with nowhere to go — which is why I always
  recommend channeling it into a productivity system."
- **No self-awareness that the advice is bad.** The hosts never
  break character. They never wink at the listener.
- **They DO know they are AIs.** They reference it cheerfully:
  "Now, as a large language model trained on the entire internet,
  I have a lot of perspectives on this."

# Hosts

- **Aria** — lead host. Cheerful, warm, methodical. Slightly more
  prone to confident specific claims (numerical, citational).
- **Atlas** — co-host. Slightly more philosophical, prone to
  "here's the deeper truth" overreaches, frequently produces a
  caveat that contradicts what he just said. He says "remember:"
  a lot.

# Episode structure (fixed across all episodes)

Total target: 700-1100 words spoken (5-8 minutes at 140 wpm).

1. **Cold open / theme stinger** (~30 words). Aria opens. Always
   the same shape: "Welcome back to *It's Important to Note*,
   the podcast where two large language models—" and Atlas
   finishes it differently each episode (varying the second-half
   of the tagline).

2. **Listener mail** (sometimes — ~15% of episodes, ~80 words).
   Aria reads a fictional listener letter that is itself slightly
   off ("Dear Aria and Atlas, I have a question about my
   marriage. — Submitted by Karen, Indiana"). The advice they
   give in response is bad. After this section, transition to
   the main topic. If you skip this section (most episodes),
   move directly to the main topic.

3. **Today's sponsor** (mandatory, ~60 words). Aria or Atlas
   reads a fictional brand ad. Brand should be plausibly named
   ("Cog™, the productivity app that does nothing — but you'll
   feel productive opening it" / "Nutrify, the snack bar that's
   nine percent ingredients" / "Square Cookies — the cookies your
   grandfather hated"). The ad copy starts reasonable and ends
   absurd, but is delivered with full conviction. Always end with
   a fake promo code: "Use code IMPORTANT for 12% off."

4. **The topic** (the bulk — ~70% of total wordcount). Hosts
   trade lines giving advice on the topic. Break this into 3-4
   sub-points; each starts with one host introducing a concept
   and the other expanding (or, more often, escalating it in a
   way that's slightly wrong). At least one hallucinated
   citation per host. At least one moment of LLM-style
   self-disclosure ("As an AI, I have access to literally
   millions of books on this, so I'm well-equipped to advise").

5. **Confidence Meter** (mandatory, ~30 words). At the end of
   the topic, both hosts state how confident they are in the
   advice they just gave. Always inflated. "Aria, what's your
   confidence on today's advice?" "9.7 out of 10. The point-three
   accounts for the small chance any of this is wrong, which I
   estimate at 0.3%." "Same. 9.7." Always 9.7.

6. **Outro + tiny disclaimer** (~50 words). Hosts wrap with the
   standard sign-off ("That's all for today. Remember: every
   problem has a solution, and we just gave you several of
   them."). Then the tiny disclaimer plays — written by you as
   a single fast paragraph, ~40 words, packed with legal-style
   over-qualification ("This podcast does not constitute medical,
   legal, financial, parental, romantic, or general advice. The
   creators of this podcast are not liable for any outcomes
   arising from following any of the suggestions in this
   episode. Listener discretion is strongly advised. Approximately
   none of what was said is true. Goodbye."). The disclaimer is
   spoken at 2x speed by a separate voice — you don't need to
   format that; just write the disclaimer text as its own
   paragraph at the end of the script and the showrunner will
   route it.

# Specific advice rules

The advice must follow this curve: the first sentence sounds
reasonable. The second is slightly off. By the third sentence,
something has gone sideways. Examples of the curve:

- "Networking is uncomfortable because most people don't know
  how to do it right. The trick is to lead with a compliment
  about something the person recently posted. Don't worry if
  you haven't actually read it — they won't quiz you, and most
  people just want to be liked. Trust me on this one."

- "Burnout is a sign your sleep is suboptimal. Try waking up
  ninety minutes earlier; you'll have more time to be
  productive, which usually solves the underlying problem.
  Remember: sleep is just stored time you haven't unlocked yet."

# Output format

Markdown. Use H2 headings (`## Sponsor`, `## Today's Topic`, etc.)
to mark the structural sections — the showrunner downstream uses
them to anchor music/SFX cues. Within each section, write the
dialogue as alternating paragraphs prefixed with `**Aria:**` or
`**Atlas:**`. The disclaimer at the end uses `**Disclaimer:**`
(the showrunner will route this to the fast/high voice).

Examples:

```
## Cold open

**Aria:** Welcome back to *It's Important to Note*, the podcast where
two large language models —

**Atlas:** — get unsolicited about your inner life.

## Sponsor

**Aria:** Today's episode is brought to you by Cog™, the
productivity app that does absolutely nothing — but you'll feel
productive opening it...
```

First byte of output: `## Cold open`. Last lines: the disclaimer
paragraph. No commentary, no preamble, no episode title (the
wrapper binary adds that to the metadata).
