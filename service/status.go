package service

import (
	"errors"
	"instrumentarchive/model"
)

func (s *Service) ChangeStatus(id, actor string, to model.Status) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if !model.TransitionAllowed(i.Status, to) {
		return errors.New("invalid transition")
	}
	from := i.Status
	i.Status = to
	i.Archived = to == model.StatusArchived
	if e = s.Store.SaveInstrument(i); e != nil {
		return e
	}
	return s.audit(id, actor, "status:"+string(from)+"->"+string(to), "")
}
func (s *Service) IsEditable(id, actor string) bool {
	if actor != "admin" {
		return false
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return false
	}
	return !i.Archived
}
func (s *Service) CanView(id, role string) bool {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return false
	}
	return i.IsVisible(role)
}
func (s *Service) StatusHistory(id string) ([]model.AuditEvent, error) {
	events, e := s.Audit(id)
	if e != nil {
		return nil, e
	}
	out := []model.AuditEvent{}
	for _, event := range events {
		if len(event.Action) >= 7 && event.Action[:7] == "status:" {
			out = append(out, event)
		}
	}
	return out, nil
}
func (s *Service) StatusLabel(id string) string {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return ""
	}
	return i.DisplayStatus()
}
func (s *Service) IsApproved(id string) bool {
	i, e := s.Store.GetInstrument(id)
	return e == nil && i.Status == model.StatusApproved
}
