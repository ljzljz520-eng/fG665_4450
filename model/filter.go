package model

import "strings"

type Filter struct {
	Query           string
	Laboratory      string
	Status          Status
	IncludeArchived bool
}

func (f Filter) Match(i Instrument) bool {
	if !f.IncludeArchived && i.Archived {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	q := strings.ToLower(strings.TrimSpace(f.Query))
	if q != "" && !strings.Contains(strings.ToLower(i.Number), q) && !strings.Contains(strings.ToLower(i.Name), q) {
		return false
	}
	if f.Laboratory != "" && !strings.EqualFold(i.Laboratory, f.Laboratory) {
		return false
	}
	return true
}
func SortInstruments(items []Instrument, by string) {
	for x := 0; x < len(items); x++ {
		for y := x + 1; y < len(items); y++ {
			swap := false
			if by == "name" {
				swap = items[y].Name < items[x].Name
			} else {
				swap = items[y].Number < items[x].Number
			}
			if swap {
				items[x], items[y] = items[y], items[x]
			}
		}
	}
}
func NormalizeStatus(s string) Status {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "draft", "草稿":
		return StatusDraft
	case "pending", "待审核":
		return StatusPending
	case "approved", "已批准":
		return StatusApproved
	case "archived", "已归档":
		return StatusArchived
	}
	return StatusDraft
}
func IsTerminal(s Status) bool                { return s == StatusArchived }
func CloneInstrument(i Instrument) Instrument { return i }
