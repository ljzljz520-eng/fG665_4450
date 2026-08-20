package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"instrumentarchive/model"
)

type Snapshot struct {
	Instruments  []model.Instrument
	Calibrations []model.Calibration
	Attachments  []model.Attachment
	Audits       []model.AuditEvent
	Workflows    []model.WorkflowRecord
}

func (s *Store) Snapshot() (Snapshot, error) {
	out := Snapshot{}
	items, e := s.ListInstruments()
	if e != nil {
		return out, e
	}
	out.Instruments = items
	s.mu.RLock()
	defer s.mu.RUnlock()
	e = s.db.View(func(tx *bbolt.Tx) error {
		for _, name := range []string{"calibrations", "attachments", "audits", "workflows"} {
			b := tx.Bucket(bucketNames[name])
			e := b.ForEach(func(_, raw []byte) error {
				switch name {
				case "calibrations":
					var v model.Calibration
					if x := json.Unmarshal(raw, &v); x != nil {
						return x
					}
					out.Calibrations = append(out.Calibrations, v)
				case "attachments":
					var v model.Attachment
					if x := json.Unmarshal(raw, &v); x != nil {
						return x
					}
					out.Attachments = append(out.Attachments, v)
				case "audits":
					var v model.AuditEvent
					if x := json.Unmarshal(raw, &v); x != nil {
						return x
					}
					out.Audits = append(out.Audits, v)
				case "workflows":
					var v model.WorkflowRecord
					if x := json.Unmarshal(raw, &v); x != nil {
						return x
					}
					out.Workflows = append(out.Workflows, v)
				}
				return nil
			})
			if e != nil {
				return e
			}
		}
		return nil
	})
	return out, e
}
func (s *Store) FindByLaboratory(lab string) ([]model.Instrument, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return nil, e
	}
	out := []model.Instrument{}
	for _, i := range items {
		if i.Laboratory == lab {
			out = append(out, i)
		}
	}
	return out, nil
}
func (s *Store) FindByOwner(owner string) ([]model.Instrument, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return nil, e
	}
	out := []model.Instrument{}
	for _, i := range items {
		if i.Owner == owner {
			out = append(out, i)
		}
	}
	return out, nil
}
func (s *Store) FindByStatus(status model.Status) ([]model.Instrument, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return nil, e
	}
	out := []model.Instrument{}
	for _, i := range items {
		if i.Status == status {
			out = append(out, i)
		}
	}
	return out, nil
}
