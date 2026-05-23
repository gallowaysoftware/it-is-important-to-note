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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

			root.SetArgs([]string{"run", "--run-dir", episodeDir, "--input", "topic=" + topic})
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
