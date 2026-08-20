package store

import (
	"encoding/json"
	"fmt"
	"instrumentarchive/model"
)

func (s *Store) Export() ([]byte, error) {
	items, err := s.ListInstruments()
	if err != nil {
		return nil, err
	}
	return json.Marshal(items)
}
func (s *Store) Import(items []model.Instrument) (int, error) {
	count := 0
	for _, i := range items {
		if err := s.SaveInstrument(i); err != nil {
			return count, fmt.Errorf("save %s: %w", i.ID, err)
		}
		count++
	}
	return count, nil
}
func (s *Store) SaveAll(i model.Instrument, c model.Calibration, a model.Attachment, e model.AuditEvent, w model.WorkflowRecord) error {
	if err := s.SaveInstrument(i); err != nil {
		return err
	}
	if err := s.SaveCalibration(c); err != nil {
		return err
	}
	if err := s.SaveAttachment(a); err != nil {
		return err
	}
	if err := s.SaveAudit(e); err != nil {
		return err
	}
	return s.SaveWorkflow(w)
}
func (s *Store) HasInstrument(id string) bool { _, err := s.GetInstrument(id); return err == nil }
func (s *Store) Count() int {
	items, err := s.ListInstruments()
	if err != nil {
		return 0
	}
	return len(items)
}
