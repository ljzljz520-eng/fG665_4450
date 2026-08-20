package store

import (
	"instrumentarchive/model"
	"strings"
)

type Index struct {
	Numbers map[string]string
	Names   map[string]string
	Labs    map[string][]string
}

func BuildIndex(items []model.Instrument) Index {
	idx := Index{Numbers: map[string]string{}, Names: map[string]string{}, Labs: map[string][]string{}}
	for _, i := range items {
		idx.Numbers[strings.ToLower(i.Number)] = i.ID
		idx.Names[strings.ToLower(i.Name)] = i.ID
		idx.Labs[i.Laboratory] = append(idx.Labs[i.Laboratory], i.ID)
	}
	return idx
}
func (idx Index) LookupNumber(number string) string {
	return idx.Numbers[strings.ToLower(strings.TrimSpace(number))]
}
func (idx Index) LookupName(name string) string {
	return idx.Names[strings.ToLower(strings.TrimSpace(name))]
}
func (idx Index) LookupLab(lab string) []string { return append([]string(nil), idx.Labs[lab]...) }
func (s *Store) BuildIndex() (Index, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return Index{}, e
	}
	return BuildIndex(items), nil
}
func (idx Index) Count() int { return len(idx.Numbers) }
