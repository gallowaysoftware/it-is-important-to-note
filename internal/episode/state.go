// Package episode manages the on-disk state for iitn episodes:
// numbering, topic rotation, RSS feed generation. Standalone-show
// model — no per-series state, no autoregression. Each episode
// stands on its own.
package episode

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultRoot returns the iitn episode store root.
func DefaultRoot() string {
	if d := os.Getenv("XDG_STATE_HOME"); d != "" {
		return filepath.Join(d, "iitn")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "iitn"
	}
	return filepath.Join(home, ".local", "state", "iitn")
}

// Layout is the on-disk shape under DefaultRoot:
//
//	$root/episodes/NNN/topic.txt        — the topic that seeded this ep
//	$root/episodes/NNN/script.md        — final script
//	$root/episodes/NNN/script.json      — showrunner output
//	$root/episodes/NNN/episode.mp3      — final mix
//	$root/topic_log.txt                 — append-only history of topics used (1 per line)
//	$root/feed.xml                      — generated RSS feed
type Layout struct {
	Root string
}

func Open() (Layout, error) {
	root := DefaultRoot()
	if err := os.MkdirAll(filepath.Join(root, "episodes"), 0o755); err != nil {
		return Layout{}, err
	}
	return Layout{Root: root}, nil
}

func (l Layout) EpisodesDir() string  { return filepath.Join(l.Root, "episodes") }
func (l Layout) TopicLog() string     { return filepath.Join(l.Root, "topic_log.txt") }
func (l Layout) FeedFile() string     { return filepath.Join(l.Root, "feed.xml") }
func (l Layout) EpisodeDir(n int) string {
	return filepath.Join(l.EpisodesDir(), fmt.Sprintf("%03d", n))
}
func (l Layout) EpisodeFile(n int, name string) string {
	return filepath.Join(l.EpisodeDir(n), name)
}

// CompletedEpisodes lists episode numbers with a finished episode.mp3.
func CompletedEpisodes(l Layout) ([]int, error) {
	entries, err := os.ReadDir(l.EpisodesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		var n int
		if _, err := fmt.Sscanf(e.Name(), "%d", &n); err != nil {
			continue
		}
		if _, err := os.Stat(l.EpisodeFile(n, "episode.mp3")); err == nil {
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out, nil
}

// NextEpisode returns the smallest 1-indexed number without a
// completed episode.mp3.
func NextEpisode(l Layout) (int, error) {
	done, err := CompletedEpisodes(l)
	if err != nil {
		return 0, err
	}
	seen := map[int]bool{}
	for _, n := range done {
		seen[n] = true
	}
	for n := 1; n <= 999; n++ {
		if !seen[n] {
			return n, nil
		}
	}
	return 0, fmt.Errorf("more than 999 episodes — refusing to keep counting")
}

// LogTopic appends a topic to the topic log so the rotation picker
// can avoid recently-used topics.
func LogTopic(l Layout, episodeNum int, topic string) error {
	f, err := os.OpenFile(l.TopicLog(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d\t%s\n", episodeNum, topic)
	return err
}

// RecentTopics returns the last N topics from the log, newest first.
func RecentTopics(l Layout, n int) ([]string, error) {
	raw, err := os.ReadFile(l.TopicLog())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	var out []string
	for i := len(lines) - 1; i >= 0 && len(out) < n; i-- {
		parts := strings.SplitN(lines[i], "\t", 2)
		if len(parts) == 2 {
			out = append(out, parts[1])
		}
	}
	return out, nil
}

// TopicCatalog is the bank of topics the show rotates through.
// Curated, not random — these are the genres of self-help slop the
// show is mocking.
var TopicCatalog = []string{
	"asking for a raise at work",
	"making new friends as an adult",
	"recovering from burnout",
	"saving money on groceries",
	"early-stage relationship advice",
	"dealing with difficult parents",
	"impostor syndrome at a new job",
	"how to have a difficult conversation",
	"productivity systems for ADHD",
	"meditation and mindfulness for beginners",
	"raising teenagers in the smartphone era",
	"breaking up with a long-term partner",
	"dealing with grief and loss",
	"managing anxiety in public speaking",
	"investing your first thousand dollars",
	"buying a house in this market",
	"finding meaning in work that doesn't pay enough",
	"setting boundaries with your in-laws",
	"recovering from a major mistake at work",
	"how to network without being weird about it",
	"dating after a long relationship",
	"managing a remote team for the first time",
	"learning to cook as a grown adult",
	"talking to a doctor when something feels off",
	"resolving conflict with a roommate",
	"sleeping better when your brain won't shut up",
	"giving honest feedback without ruining a friendship",
	"how to be less hard on yourself",
	"deciding whether to quit your job",
	"making it through the holidays with your family",
	"choosing between two jobs that both have downsides",
	"recovering from a public embarrassment",
	"finding a therapist",
	"forgiving someone who never apologized",
	"explaining your weird career to your relatives",
	"managing money in a relationship",
	"how to ask for what you actually want",
	"dealing with a passive-aggressive coworker",
	"the proper way to apologize",
	"becoming a morning person",
	"saying no to social plans without guilt",
	"reconnecting with an old friend",
	"talking to your boss about your mental health",
	"recovering from a bad financial decision",
	"keeping a long-distance friendship alive",
	"finding hobbies as a working adult",
	"managing screen time in your own life",
	"having the kid conversation with a partner",
	"deciding to move to a new city",
	"telling someone they have something in their teeth",
	"adjusting after retirement",
	"managing chronic pain that doctors can't explain",
}

// PickTopic chooses the next topic uniformly at random from those in
// TopicCatalog that have never appeared in the topic log. Once every
// topic has been used at least once, it picks uniformly at random from
// the topics least recently used (those that appeared least often).
//
// Earlier versions used a fixed cooldown window over `RecentTopics`,
// which let a topic recycle after N episodes — and because the LLM
// stages are content-addressed by rendered prompt, a repeat topic
// produced a verbatim cache replay of the prior episode's script. The
// LRU-over-history picker keeps every script fresh until the catalog
// is exhausted; the `cooldown` arg is retained for API stability but
// ignored. Pass 0 if you need to signal "no preference."
func PickTopic(l Layout, _ int) (string, error) {
	used, err := topicUseCounts(l)
	if err != nil {
		return "", err
	}
	// First pass: any topic with zero uses is fair game; pick one
	// uniformly at random so the order doesn't reveal the catalog's
	// declaration order to a listener binge-watching the feed.
	var unused []string
	for _, t := range TopicCatalog {
		if used[t] == 0 {
			unused = append(unused, t)
		}
	}
	if len(unused) > 0 {
		return unused[rand.Intn(len(unused))], nil
	}
	// All topics used at least once. Pick uniformly at random from the
	// set with the minimum use count so the rotation re-cycles in a
	// fresh order rather than catalog order.
	minUses := -1
	for _, t := range TopicCatalog {
		if minUses == -1 || used[t] < minUses {
			minUses = used[t]
		}
	}
	var leastUsed []string
	for _, t := range TopicCatalog {
		if used[t] == minUses {
			leastUsed = append(leastUsed, t)
		}
	}
	return leastUsed[rand.Intn(len(leastUsed))], nil
}

// topicUseCounts returns the count of times each topic appears in the
// topic log. Returns an empty map if the log doesn't exist yet — the
// caller treats the absent map as "every topic has zero uses," which
// is the desired semantics for a fresh series.
func topicUseCounts(l Layout) (map[string]int, error) {
	raw, err := os.ReadFile(l.TopicLog())
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]int{}, nil
		}
		return nil, err
	}
	counts := map[string]int{}
	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		counts[parts[1]]++
	}
	return counts, nil
}
