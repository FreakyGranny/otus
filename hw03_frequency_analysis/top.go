package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

type entry struct {
	Count int
	Value string
}

const topCount = 10

// Top10 ...
func Top10(text string) []string {
	words := strings.Fields(text)
	counts := make(map[string]int)
	r := regexp.MustCompile(`[а-яА-Яa-zA-Z]+(\-[а-яА-Яa-zA-Z]+)*`)

	for _, w := range words {
		if !r.MatchString(w) {
			continue
		}
		counts[strings.ToLower(r.FindString(w))]++
	}
	entries := make([]entry, 0, len(counts))
	for word, count := range counts {
		entries = append(entries, entry{Count: count, Value: word})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Count > entries[j].Count })

	results := make([]string, 0, topCount)
	for i, e := range entries {
		if i == topCount {
			break
		}
		results = append(results, e.Value)
	}
	return results
}
