// Command iitn — "It's Important to Note" — generates short
// satirical AI-advice podcast episodes. Two AI hosts (Aria and
// Atlas) give confidently wrong life advice on one topic at a time.
// Each episode is standalone; topic rotation lives in the binary.
//
// Subcommands:
//
//	iitn next [--topic "..."]   Generate the next episode.
//	iitn list                   Show generated episodes + topics.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/gallowaysoftware/vibe/vamp"

	"github.com/gallowaysoftware/it-is-important-to-note/internal/episode"
	"github.com/gallowaysoftware/it-is-important-to-note/internal/pipeline"
)

func main() {
	root := &cobra.Command{
		Use:   "iitn",
		Short: "Generate \"It's Important to Note\" — AI-hosts-give-bad-advice podcast episodes.",
		Long: `iitn produces short (5-8 min) satirical podcast episodes hosted
by two AI assistants (Aria and Atlas) who give confidently wrong
life advice on one topic per episode. The comedy is the genre
itself; the show leans into LLM slop tropes — em-dashes, "It's
important to note," hallucinated citations, helpful warmth applied
to dark advice.

Each episode is standalone. Topic rotation is built into the
binary; the next subcommand picks the next unused topic (or you
can override with --topic).`,
		SilenceUsage: true,
	}
	root.AddCommand(nextCommand())
	root.AddCommand(listCommand())
	root.AddCommand(timingsCommand())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "iitn:", err)
		os.Exit(1)
	}
}

func nextCommand() *cobra.Command {
	var (
		explicitTopic string
		explicitNum   int
		publishTo     string
	)
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Generate the next episode.",
		RunE: func(cmd *cobra.Command, args []string) error {
			layout, err := episode.Open()
			if err != nil {
				return err
			}
			n := explicitNum
			if n == 0 {
				n, err = episode.NextEpisode(layout)
				if err != nil {
					return err
				}
			}
			topic := explicitTopic
			if topic == "" {
				// If this is a re-run on an existing episode dir,
				// reuse the previously-logged topic so cache hits
				// further down the pipeline don't waste an LLM call
				// re-generating against a different topic.
				if existing, err := os.ReadFile(layout.EpisodeFile(n, "topic.txt")); err == nil {
					topic = strings.TrimSpace(string(existing))
				}
			}
			if topic == "" {
				topic, err = episode.PickTopic(layout, 12)
				if err != nil {
					return err
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "episode %d topic: %s\n", n, topic)

			cfg := pipeline.Config{Topic: topic, EpisodeNumber: n}
			root, err := vamp.BuildRoot(func() (*vamp.Pipeline, error) {
				return pipeline.Build(cfg)
			})
			if err != nil {
				return err
			}
			episodeDir := layout.EpisodeDir(n)
			if err := os.MkdirAll(episodeDir, 0o755); err != nil {
				return err
			}
			// Persist the topic before the run so a later --resume
			// can pick it up + the list/log commands see it.
			if err := os.WriteFile(layout.EpisodeFile(n, "topic.txt"), []byte(topic+"\n"), 0o644); err != nil {
				return err
			}

			root.SetArgs([]string{
				"run", "--run-dir", episodeDir,
				"--input", "topic=" + topic,
				"--input", fmt.Sprintf("episode_number=%d", n),
			})
			if err := root.Execute(); err != nil {
				return fmt.Errorf("episode %d: %w", n, err)
			}
			if err := episode.LogTopic(layout, n, topic); err != nil {
				return fmt.Errorf("log topic: %w", err)
			}
			localM4B := layout.EpisodeFile(n, "episode.m4b")
			fmt.Fprintf(cmd.OutOrStdout(), "\n✓ episode %d done: %s\n", n, localM4B)

			// Optionally copy the episode into a podcast library
			// path. Audiobookshelf's podcast scanner reads ID3/m4b
			// metadata for episode number + title; the on-disk
			// name still wants to lead with a 3-digit episode
			// number so the per-show folder sorts.
			if publishTo != "" {
				name := fmt.Sprintf("%03d - %s.m4b", n, sanitiseFilename(topic))
				dst := filepath.Join(publishTo, name)
				if err := os.MkdirAll(publishTo, 0o755); err != nil {
					return fmt.Errorf("mkdir publish dir: %w", err)
				}
				if err := copyFile(localM4B, dst); err != nil {
					return fmt.Errorf("publish to %s: %w", dst, err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  published: %s\n", dst)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&explicitTopic, "topic", "", "Override topic. Default: next from catalog (skipping recent 12).")
	cmd.Flags().IntVar(&explicitNum, "episode", 0, "Override episode number. 0 = next pending.")
	cmd.Flags().StringVar(&publishTo, "publish-to", "", "Directory to copy the finished episode.m4b into, renamed `NNN - <topic>.m4b`. Typical: /mnt/<podcast-library>/Its Important to Note/.")
	return cmd
}

// sanitiseFilename strips characters that misbehave in podcast-app /
// audiobookshelf scanners — colons, slashes, control characters,
// trailing whitespace. Keeps spaces + hyphens since those scan fine.
func sanitiseFilename(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '/' || r == '\\' || r == ':' || r == '*' || r == '?' || r == '"' || r == '<' || r == '>' || r == '|':
			b.WriteRune('-')
		case r < 0x20:
			continue
		default:
			b.WriteRune(r)
		}
	}
	return strings.TrimSpace(b.String())
}

// copyFile streams src → dst with a 32KB buffer so a large m4b
// doesn't load entirely into memory. Used by the --publish-to flag
// to land each finished episode in the podcast library.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List generated episodes with their topics.",
		RunE: func(cmd *cobra.Command, args []string) error {
			layout, err := episode.Open()
			if err != nil {
				return err
			}
			done, err := episode.CompletedEpisodes(layout)
			if err != nil {
				return err
			}
			if len(done) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no episodes yet — `iitn next` to generate one")
				return nil
			}
			for _, n := range done {
				topic, _ := os.ReadFile(layout.EpisodeFile(n, "topic.txt"))
				fmt.Fprintf(cmd.OutOrStdout(), "  %03d  %s\n", n, string(topic))
			}
			return nil
		},
	}
}

// pipelineTiming mirrors the JSON shape vamp writes to
// pipeline_timing.json after each run. Only the fields the timings
// subcommand needs are decoded; the rest is silently dropped.
type pipelineTiming struct {
	StartedAt time.Time           `json:"started_at"`
	TotalMS   int64               `json:"total_ms"`
	Stages    []pipelineStageTime `json:"stages"`
	// Capabilities is the per-capability profile resolution for the run
	// — vibe v0.6.0+ records this. Older runs (or runs against a vamp
	// without the field) decode to nil and the "profile" column reads
	// as "-".
	Capabilities map[string]string `json:"capabilities,omitempty"`
}

type pipelineStageTime struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	DurationMS int64          `json:"duration_ms"`
	Status     string         `json:"status"`
	Notes      map[string]any `json:"notes"`
}

// timingsCommand prints a per-episode wall-clock table. Useful for
// comparing run-to-run speed across profile / backend changes (e.g.
// EXL3 vs GGUF, write_script duration delta after an MTP draft swap).
func timingsCommand() *cobra.Command {
	var (
		showStages bool
		filterID   string
	)
	cmd := &cobra.Command{
		Use:   "timings",
		Short: "Show per-episode wall-clock timings parsed from pipeline_timing.json.",
		Long: `timings reads pipeline_timing.json from every published episode and
prints a table of total + LLM-stage durations. With --stages, also
breaks out the slowest stage per episode. With --stage <id>, prints
only that stage's duration column.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			layout, err := episode.Open()
			if err != nil {
				return err
			}
			done, err := episode.CompletedEpisodes(layout)
			if err != nil {
				return err
			}
			if len(done) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no episodes yet — `iitn next` to generate one")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			defer w.Flush()
			if filterID != "" {
				fmt.Fprintf(w, "ep\ttopic\t%s\n", filterID)
			} else if showStages {
				fmt.Fprintln(w, "ep\ttopic\ttotal\tprofile\tslowest")
			} else {
				fmt.Fprintln(w, "ep\ttopic\ttotal\tprofile")
			}
			for _, n := range done {
				topic, _ := os.ReadFile(layout.EpisodeFile(n, "topic.txt"))
				topicStr := strings.TrimSpace(string(topic))
				timing, err := readTiming(layout.EpisodeFile(n, "pipeline_timing.json"))
				if err != nil {
					// Missing pipeline_timing.json is normal for very old
					// runs from before vamp recorded timing. Skip with a
					// dash rather than erroring the whole table.
					fmt.Fprintf(w, "%03d\t%s\t-\n", n, topicStr)
					continue
				}
				total := time.Duration(timing.TotalMS) * time.Millisecond
				if filterID != "" {
					d := stageDuration(timing, filterID)
					if d < 0 {
						fmt.Fprintf(w, "%03d\t%s\t-\n", n, topicStr)
						continue
					}
					fmt.Fprintf(w, "%03d\t%s\t%s\n", n, topicStr, fmtDuration(d))
					continue
				}
				profile := timing.Capabilities["long_form"]
				if profile == "" {
					profile = "-"
				}
				if showStages {
					slow := slowestStage(timing)
					if slow == "" {
						fmt.Fprintf(w, "%03d\t%s\t%s\t%s\t-\n", n, topicStr, fmtDuration(total), profile)
					} else {
						fmt.Fprintf(w, "%03d\t%s\t%s\t%s\t%s\n", n, topicStr, fmtDuration(total), profile, slow)
					}
					continue
				}
				fmt.Fprintf(w, "%03d\t%s\t%s\t%s\n", n, topicStr, fmtDuration(total), profile)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&showStages, "stages", false, "Also surface each episode's slowest stage (id + duration).")
	cmd.Flags().StringVar(&filterID, "stage", "", "Show only this stage's duration column instead of the total + slowest.")
	return cmd
}

func readTiming(path string) (*pipelineTiming, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var t pipelineTiming
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &t, nil
}

// stageDuration returns the stage's duration in time.Duration units, or
// -1 when no stage with that id was recorded for the episode (e.g. an
// older pipeline version without the stage).
func stageDuration(t *pipelineTiming, id string) time.Duration {
	for i := range t.Stages {
		if t.Stages[i].ID == id {
			return time.Duration(t.Stages[i].DurationMS) * time.Millisecond
		}
	}
	return -1
}

// slowestStage returns a "<id> <duration>" string for the longest non-
// foreach stage in the timing record. Foreach parents typically
// dominate via fan-out aggregation, which dwarfs single-call LLM stages
// and makes the column boring; skip them in favour of the genuine
// bottleneck.
func slowestStage(t *pipelineTiming) string {
	sorted := make([]pipelineStageTime, len(t.Stages))
	copy(sorted, t.Stages)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DurationMS > sorted[j].DurationMS
	})
	for _, s := range sorted {
		// Skip the umbrella audio/foreach stages whose duration sums up
		// the fan-out items — they're not "a stage" the operator can
		// optimize.
		if s.Type == "audio" {
			continue
		}
		return fmt.Sprintf("%s %s", s.ID, fmtDuration(time.Duration(s.DurationMS)*time.Millisecond))
	}
	return ""
}

// fmtDuration renders a duration as a compact "1m23s" string, dropping
// sub-second precision when the duration is over a minute (the timings
// table is for at-a-glance comparison, not benchmark reporting).
func fmtDuration(d time.Duration) string {
	if d >= time.Minute {
		return d.Round(time.Second).String()
	}
	return d.Round(100 * time.Millisecond).String()
}
