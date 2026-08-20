package report

import (
	"instrumentarchive/model"
	"sort"
)

func SortByStatus(items []model.Instrument) []model.Instrument {
	out := append([]model.Instrument(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return model.StatusRank(out[i].Status) > model.StatusRank(out[j].Status) })
	return out
}
func FilterAttention(items []model.Instrument) []model.Instrument {
	out := []model.Instrument{}
	for _, i := range items {
		if i.Owner == "" || i.CalibrationID == "" {
			out = append(out, i)
		}
	}
	return out
}
func FilterArchived(items []model.Instrument, archived bool) []model.Instrument {
	out := []model.Instrument{}
	for _, i := range items {
		if i.Archived == archived {
			out = append(out, i)
		}
	}
	return out
}
func UniqueLaboratories(items []model.Instrument) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, i := range items {
		if !seen[i.Laboratory] {
			seen[i.Laboratory] = true
			out = append(out, i.Laboratory)
		}
	}
	sort.Strings(out)
	return out
}
