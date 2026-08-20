package store

import (
	"instrumentarchive/model"
	"sort"
)

type Metrics struct {
	Total    int
	Active   int
	Archived int
	ByStatus map[model.Status]int
	ByLab    map[string]int
}

func (s *Store) CalculateMetrics() (Metrics, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return Metrics{}, e
	}
	m := Metrics{ByStatus: map[model.Status]int{}, ByLab: map[string]int{}}
	for _, i := range items {
		m.Total++
		m.ByStatus[i.Status]++
		m.ByLab[i.Laboratory]++
		if i.Archived {
			m.Archived++
		} else {
			m.Active++
		}
	}
	return m, nil
}
func (m Metrics) TopLabs(limit int) []string {
	keys := make([]string, 0, len(m.ByLab))
	for k := range m.ByLab {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if m.ByLab[keys[i]] == m.ByLab[keys[j]] {
			return keys[i] < keys[j]
		}
		return m.ByLab[keys[i]] > m.ByLab[keys[j]]
	})
	if limit < 0 {
		limit = 0
	}
	if limit > len(keys) {
		limit = len(keys)
	}
	return keys[:limit]
}
func (m Metrics) CompletionRate() float64 {
	if m.Total == 0 {
		return 0
	}
	return float64(m.Archived) / float64(m.Total)
}
func (m Metrics) StatusTotal(s model.Status) int { return m.ByStatus[s] }
