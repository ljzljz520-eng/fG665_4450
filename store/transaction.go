package store

import (
	"errors"
	"instrumentarchive/model"
)

type Transaction struct {
	Store       *Store
	Instrument  *model.Instrument
	Calibration *model.Calibration
	Attachments []model.Attachment
	Audits      []model.AuditEvent
	Workflow    *model.WorkflowRecord
	committed   bool
}

func (s *Store) Begin() *Transaction {
	return &Transaction{Store: s, Attachments: []model.Attachment{}, Audits: []model.AuditEvent{}}
}
func (t *Transaction) SetInstrument(v model.Instrument) *Transaction   { t.Instrument = &v; return t }
func (t *Transaction) SetCalibration(v model.Calibration) *Transaction { t.Calibration = &v; return t }
func (t *Transaction) AddAttachment(v model.Attachment) *Transaction {
	t.Attachments = append(t.Attachments, v)
	return t
}
func (t *Transaction) AddAudit(v model.AuditEvent) *Transaction {
	t.Audits = append(t.Audits, v)
	return t
}
func (t *Transaction) SetWorkflow(v model.WorkflowRecord) *Transaction { t.Workflow = &v; return t }
func (t *Transaction) Validate() error {
	if t.Instrument == nil {
		return errors.New("instrument missing")
	}
	if e := t.Instrument.Validate(); e != nil {
		return e
	}
	for _, a := range t.Attachments {
		if e := a.Validate(); e != nil {
			return e
		}
	}
	return nil
}
func (t *Transaction) Commit() error {
	if t.committed {
		return errors.New("already committed")
	}
	if e := t.Validate(); e != nil {
		return e
	}
	if e := t.Store.SaveInstrument(*t.Instrument); e != nil {
		return e
	}
	if t.Calibration != nil {
		if e := t.Store.SaveCalibration(*t.Calibration); e != nil {
			return e
		}
	}
	for _, a := range t.Attachments {
		if e := t.Store.SaveAttachment(a); e != nil {
			return e
		}
	}
	for _, a := range t.Audits {
		if e := t.Store.SaveAudit(a); e != nil {
			return e
		}
	}
	if t.Workflow != nil {
		if e := t.Store.SaveWorkflow(*t.Workflow); e != nil {
			return e
		}
	}
	t.committed = true
	return nil
}
func (t *Transaction) Rollback() {
	t.Instrument = nil
	t.Calibration = nil
	t.Attachments = nil
	t.Audits = nil
	t.Workflow = nil
}
func (t *Transaction) Committed() bool { return t.committed }
func (s *Store) UpsertInstrument(i model.Instrument) error {
	if i.ID == "" {
		return errors.New("id required")
	}
	if _, e := s.GetInstrument(i.ID); e == nil {
		return s.SaveInstrument(i)
	}
	return s.SaveInstrument(i)
}
func (s *Store) RemoveAttachments(id string) error {
	items, e := s.ListAttachments(id)
	if e != nil {
		return e
	}
	for _, a := range items {
		if e = s.del("attachments", a.ID); e != nil {
			return e
		}
	}
	return nil
}
