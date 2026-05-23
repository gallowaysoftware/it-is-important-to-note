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

// Topic pairs a catalog entry's prompt text with a genre tag. The
// picker uses the genre to keep recent episodes diverse — without
// it the random-from-unused selector kept dropping 5 relationship
// topics in the first 16 episodes (user feedback, 2026-05-23).
type Topic struct {
	Name  string
	Genre string
}

// Topic genres. Tags are tuned for diversification, not perfect
// taxonomy: "domestic" covers home life + family logistics, "world"
// covers civic + outward-facing topics, "absurd" tags catalog
// entries whose hook is already weird so we don't double up on them
// in the same week.
const (
	GenreWork          = "work"
	GenreMoney         = "money"
	GenreRelationships = "relationships"
	GenreWellness      = "wellness"
	GenreDomestic      = "domestic"
	GenreWorld         = "world"
	GenreTech          = "tech"
	GenreAbsurd        = "absurd"
)

// TopicCatalog is the bank of topics the show rotates through.
// Curated, not random — these are the genres of self-help slop the
// show is mocking. Each topic is genre-tagged so the picker can
// downweight recently-seen genres and keep listenable rotation.
var TopicCatalog = []Topic{
	// work
	{"asking for a raise at work", GenreWork},
	{"impostor syndrome at a new job", GenreWork},
	{"productivity systems for ADHD", GenreWork},
	{"recovering from a major mistake at work", GenreWork},
	{"how to network without being weird about it", GenreWork},
	{"managing a remote team for the first time", GenreWork},
	{"deciding whether to quit your job", GenreWork},
	{"choosing between two jobs that both have downsides", GenreWork},
	{"dealing with a passive-aggressive coworker", GenreWork},
	{"talking to your boss about your mental health", GenreWork},
	{"explaining your weird career to your relatives", GenreWork},
	{"finding meaning in work that doesn't pay enough", GenreWork},

	// money
	{"saving money on groceries", GenreMoney},
	{"investing your first thousand dollars", GenreMoney},
	{"buying a house in this market", GenreMoney},
	{"recovering from a bad financial decision", GenreMoney},
	{"managing money in a relationship", GenreMoney},

	// relationships (kept tagged so they don't cluster — user feedback)
	{"early-stage relationship advice", GenreRelationships},
	{"breaking up with a long-term partner", GenreRelationships},
	{"dating after a long relationship", GenreRelationships},
	{"how to ask for what you actually want", GenreRelationships},
	{"the proper way to apologize", GenreRelationships},
	{"forgiving someone who never apologized", GenreRelationships},
	{"having the kid conversation with a partner", GenreRelationships},
	{"giving honest feedback without ruining a friendship", GenreRelationships},
	{"reconnecting with an old friend", GenreRelationships},
	{"keeping a long-distance friendship alive", GenreRelationships},
	{"making new friends as an adult", GenreRelationships},

	// wellness
	{"recovering from burnout", GenreWellness},
	{"meditation and mindfulness for beginners", GenreWellness},
	{"managing anxiety in public speaking", GenreWellness},
	{"dealing with grief and loss", GenreWellness},
	{"how to be less hard on yourself", GenreWellness},
	{"becoming a morning person", GenreWellness},
	{"sleeping better when your brain won't shut up", GenreWellness},
	{"managing chronic pain that doctors can't explain", GenreWellness},
	{"finding a therapist", GenreWellness},
	{"talking to a doctor when something feels off", GenreWellness},
	{"finding hobbies as a working adult", GenreWellness},

	// domestic
	{"dealing with difficult parents", GenreDomestic},
	{"raising teenagers in the smartphone era", GenreDomestic},
	{"setting boundaries with your in-laws", GenreDomestic},
	{"making it through the holidays with your family", GenreDomestic},
	{"learning to cook as a grown adult", GenreDomestic},
	{"resolving conflict with a roommate", GenreDomestic},
	{"adjusting after retirement", GenreDomestic},
	{"deciding to move to a new city", GenreDomestic},

	// tech
	{"managing screen time in your own life", GenreTech},
	{"how to have a difficult conversation", GenreTech}, // tagged tech because the bots will frame it as protocol design
	{"telling someone they have something in their teeth", GenreAbsurd},
	{"recovering from a public embarrassment", GenreAbsurd},
	{"saying no to social plans without guilt", GenreRelationships},
}

// recentGenreWindow is the number of most-recent episodes whose
// genres get downweighted when picking a fresh topic. Set to 3 so
// runs of two are tolerated (relationship → wellness → relationship)
// but runs of three+ are actively discouraged.
const recentGenreWindow = 3

// recentGenrePenalty is the multiplicative weight applied to a topic
// whose genre appears in the last recentGenreWindow episodes. 0.2 is
// strong enough to push the picker toward other genres without
// hard-banning a genre when the catalog has run thin.
const recentGenrePenalty = 0.2

// PickTopic chooses the next topic with two-level diversification:
//
//  1. Topics that have never appeared in the log are preferred. Once
//     all topics have been used at least once, the least-used set is
//     the candidate pool instead.
//  2. Within the candidate pool, each topic's selection weight is
//     downweighted by `recentGenrePenalty` for every appearance of
//     its genre in the last `recentGenreWindow` episodes. This keeps
//     genres from clustering (the v1 picker dropped 5 relationship
//     topics in the first 16 episodes — user feedback, 2026-05-23).
//
// The `cooldown` arg is retained for API stability but ignored — the
// genre penalty subsumes its role.
func PickTopic(l Layout, _ int) (string, error) {
	used, err := topicUseCounts(l)
	if err != nil {
		return "", err
	}
	recentGenres, err := recentGenreCounts(l, recentGenreWindow)
	if err != nil {
		return "", err
	}

	// Build the candidate pool: topics with the minimum use count.
	// On a fresh install everything has count 0; once the catalog is
	// exhausted, the pool becomes the least-recently-cycled set.
	minUses := -1
	for _, t := range TopicCatalog {
		if minUses == -1 || used[t.Name] < minUses {
			minUses = used[t.Name]
		}
	}
	var pool []Topic
	for _, t := range TopicCatalog {
		if used[t.Name] == minUses {
			pool = append(pool, t)
		}
	}
	if len(pool) == 0 {
		return "", fmt.Errorf("PickTopic: empty candidate pool (catalog has %d topics)", len(TopicCatalog))
	}

	// Compute weights for each candidate. Base weight is 1.0; each
	// time the topic's genre appeared in the recent window the weight
	// is multiplied by recentGenrePenalty. So a topic whose genre is
	// in the recent window twice ends up at 0.04 relative weight.
	weights := make([]float64, len(pool))
	totalWeight := 0.0
	for i, t := range pool {
		w := 1.0
		for n := 0; n < recentGenres[t.Genre]; n++ {
			w *= recentGenrePenalty
		}
		weights[i] = w
		totalWeight += w
	}

	// Weighted random pick.
	r := rand.Float64() * totalWeight
	for i, w := range weights {
		r -= w
		if r <= 0 {
			return pool[i].Name, nil
		}
	}
	// Numeric rounding fallback.
	return pool[len(pool)-1].Name, nil
}

// recentGenreCounts returns the count of each genre in the last `n`
// log entries. Used by PickTopic to downweight recently-clustered
// genres. Empty log → empty map.
func recentGenreCounts(l Layout, n int) (map[string]int, error) {
	recent, err := RecentTopics(l, n)
	if err != nil {
		return nil, err
	}
	byName := make(map[string]string, len(TopicCatalog))
	for _, t := range TopicCatalog {
		byName[t.Name] = t.Genre
	}
	counts := map[string]int{}
	for _, name := range recent {
		if g, ok := byName[name]; ok {
			counts[g]++
		}
	}
	return counts, nil
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
