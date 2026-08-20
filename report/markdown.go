package report

import (
	"fmt"
	"instrumentarchive/model"
	"strings"
)

type Section struct {
	Title string
	Lines []string
}

func GroupByLab(items []model.Instrument) map[string][]model.Instrument {
	out := map[string][]model.Instrument{}
	for _, i := range items {
		out[i.Laboratory] = append(out[i.Laboratory], i)
	}
	return out
}
func StatusCounts(items []model.Instrument) map[model.Status]int {
	out := map[model.Status]int{}
	for _, i := range items {
		out[i.Status]++
	}
	return out
}
func BuildSections(items []model.Instrument) []Section {
	groups := GroupByLab(items)
	out := []Section{}
	for lab, list := range groups {
		lines := []string{}
		for _, i := range list {
			lines = append(lines, fmt.Sprintf("%s | %s | %s", i.Number, i.Name, i.DisplayStatus()))
		}
		out = append(out, Section{Title: lab, Lines: lines})
	}
	return out
}
func RenderMarkdown(items []model.Instrument) string {
	var b strings.Builder
	b.WriteString("# 仪器档案报告\n")
	for _, s := range BuildSections(items) {
		b.WriteString("\n## " + s.Title + "\n")
		for _, line := range s.Lines {
			b.WriteString("- " + line + "\n")
		}
	}
	return b.String()
}
func PendingCalibration(items []model.Instrument) []model.Instrument {
	out := []model.Instrument{}
	for _, i := range items {
		if i.CalibrationID == "" {
			out = append(out, i)
		}
	}
	return out
}
