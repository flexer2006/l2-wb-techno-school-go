package eleven

import (
	"slices"
	"sort"
	"strings"
)

func FindAnagrams(words []string) map[string][]string {
	type group struct {
		first string
		words []string
	}

	sigToGroup := make(map[string]*group, len(words))
	for _, w := range words {
		lw := strings.ToLower(w)
		sig := sortedRunes(lw)
		if g, ok := sigToGroup[sig]; ok {
			g.words = append(g.words, lw)
		} else {
			sigToGroup[sig] = &group{
				first: lw,
				words: []string{lw},
			}
		}
	}

	result := make(map[string][]string, len(sigToGroup))
	for _, g := range sigToGroup {
		if len(g.words) <= 1 {
			continue
		}
		sort.Strings(g.words)
		result[g.first] = g.words
	}
	return result
}

func sortedRunes(s string) string {
	r := []rune(s)
	slices.Sort(r) // вместо sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	return string(r)
}
