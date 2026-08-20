package service

import (
	"errors"
	"fmt"
	"instrumentarchive/model"
	"instrumentarchive/store"
)

type Clock interface{ Now() string }
type StaticClock struct{ Value string }

func (c StaticClock) Now() string { return c.Value }

type Service struct {
	Store *store.Store
	Clock Clock
}

func New(s *store.Store, c Clock) *Service { return &Service{Store: s, Clock: c} }
func (s *Service) Create(i model.Instrument, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	if err := i.Validate(); err != nil {
		return err
	}
	if i.Status == "" {
		i.Status = model.StatusDraft
	}
	if s.Store.HasInstrument(i.ID) {
		return errors.New("duplicate instrument")
	}
	if err := s.Store.SaveInstrument(i); err != nil {
		return err
	}
	return s.audit(i.ID, actor, "create", "")
}
func (s *Service) Submit(id, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if i.Status != model.StatusDraft {
		return errors.New("only draft can submit")
	}
	i.Status = model.StatusPending
	if e = s.Store.SaveInstrument(i); e != nil {
		return e
	}
	return s.audit(id, actor, "submit", "")
}
func (s *Service) Review(id, actor string, approve bool) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if !i.CanReview(actor) {
		return errors.New("not pending")
	}
	if approve {
		i.Status = model.StatusApproved
	} else {
		i.Status = model.StatusDraft
	}
	if e = s.Store.SaveInstrument(i); e != nil {
		return e
	}
	action := "reject"
	if approve {
		action = "approve"
	}
	return s.audit(id, actor, action, "")
}
func (s *Service) Archive(id, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if !i.CanArchive(actor) {
		return errors.New("not approved")
	}
	i.Status = model.StatusArchived
	i.Archived = true
	if e = s.Store.SaveInstrument(i); e != nil {
		return e
	}
	return s.audit(id, actor, "archive", "")
}
func (s *Service) AddCalibration(c model.Calibration, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	if !s.Store.HasInstrument(c.InstrumentID) {
		return errors.New("instrument missing")
	}
	return s.Store.SaveCalibration(c)
}
func (s *Service) AddAttachment(a model.Attachment, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	if !s.Store.HasInstrument(a.InstrumentID) {
		return errors.New("instrument missing")
	}
	return s.Store.SaveAttachment(a)
}
func (s *Service) Search(f model.Filter, role string) ([]model.Instrument, error) {
	all, e := s.Store.ListInstruments()
	if e != nil {
		return nil, e
	}
	out := make([]model.Instrument, 0, len(all))
	for _, i := range all {
		if i.IsVisible(role) && f.Match(i) {
			out = append(out, i)
		}
	}
	model.SortInstruments(out, "number")
	return out, nil
}
func (s *Service) Detail(id, role string) (model.Detail, error) {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return model.Detail{}, e
	}
	if !i.IsVisible(role) {
		return model.Detail{}, errors.New("forbidden")
	}
	c, _ := s.Store.GetCalibration(i.CalibrationID)
	at, _ := s.Store.ListAttachments(id)
	return UnsafeAssembleDetail(i, c, at), nil
}
func (s *Service) UpdateOwner(id, owner, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if owner == "" {
		return errors.New("owner required")
	}
	i.Owner = owner
	if e = s.Store.SaveInstrument(i); e != nil {
		return e
	}
	return s.audit(id, actor, "update_owner", owner)
}
func (s *Service) audit(id, actor, action, comment string) error {
	at := ""
	if s.Clock != nil {
		at = s.Clock.Now()
	}
	e := model.NewAudit(fmt.Sprintf("%s-%s", id, action), id, actor, action, at, comment)
	return s.Store.SaveAudit(e)
}
func (s *Service) Audit(id string) ([]model.AuditEvent, error) { return s.Store.ListAudits(id) }
