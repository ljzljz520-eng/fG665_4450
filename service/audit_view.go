package service

import (
	"instrumentarchive/model"
	"sort"
	"strings"
)

type AuditSummary struct {
	InstrumentID string
	Events       int
	Actors       []string
	Actions      []string
	LastAction   string
}

func SummarizeAudit(events []model.AuditEvent) AuditSummary {
	sum := AuditSummary{}
	if len(events) > 0 {
		sum.InstrumentID = events[0].InstrumentID
	}
	actors := map[string]bool{}
	actions := map[string]bool{}
	for _, e := range events {
		sum.Events++
		if !actors[e.Actor] {
			actors[e.Actor] = true
			sum.Actors = append(sum.Actors, e.Actor)
		}
		if !actions[e.Action] {
			actions[e.Action] = true
			sum.Actions = append(sum.Actions, e.Action)
		}
		sum.LastAction = e.Action
	}
	sort.Strings(sum.Actors)
	sort.Strings(sum.Actions)
	return sum
}
func (s *Service) AuditSummary(id string) (AuditSummary, error) {
	events, e := s.Audit(id)
	if e != nil {
		return AuditSummary{}, e
	}
	return SummarizeAudit(events), nil
}
func (s *Service) AuditTimeline(id string) ([]string, error) {
	events, e := s.Audit(id)
	if e != nil {
		return nil, e
	}
	out := []string{}
	for _, e := range events {
		out = append(out, e.At+" "+e.Actor+" "+e.Action)
	}
	return out, nil
}
func (s *Service) AuditContains(id, action string) bool {
	events, e := s.Audit(id)
	if e != nil {
		return false
	}
	for _, e := range events {
		if strings.EqualFold(e.Action, action) {
			return true
		}
	}
	return false
}
func (s *Service) AuditActors(id string) []string {
	sum, e := s.AuditSummary(id)
	if e != nil {
		return []string{}
	}
	return sum.Actors
}
func (s *Service) AuditActions(id string) []string {
	sum, e := s.AuditSummary(id)
	if e != nil {
		return []string{}
	}
	return sum.Actions
}
func (s *Service) LastAuditAction(id string) string {
	sum, e := s.AuditSummary(id)
	if e != nil {
		return ""
	}
	return sum.LastAction
}
