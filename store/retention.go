package store

import (
	"errors"
	"instrumentarchive/model"
	"time"
)

type RetentionRule struct {
	KeepArchived    bool
	MaxAttachments  int
	ExpireAfterDays int
}

func DefaultRetention() RetentionRule {
	return RetentionRule{KeepArchived: true, MaxAttachments: 20, ExpireAfterDays: 3650}
}
func (r RetentionRule) Keep(i model.Instrument) bool {
	if i.Archived {
		return r.KeepArchived
	}
	return true
}
func (r RetentionRule) AttachmentAllowed(a model.Attachment) bool {
	if r.MaxAttachments <= 0 {
		return false
	}
	return a.Size() <= 10*1024*1024
}
func (r RetentionRule) Expired(date string, now time.Time) bool {
	if date == "" || r.ExpireAfterDays <= 0 {
		return false
	}
	parsed, e := time.Parse("2006-01-02", date)
	if e != nil {
		return false
	}
	return now.Sub(parsed) > time.Duration(r.ExpireAfterDays)*24*time.Hour
}
func (s *Store) ApplyRetention(rule RetentionRule) (int, error) {
	items, e := s.ListInstruments()
	if e != nil {
		return 0, e
	}
	removed := 0
	for _, i := range items {
		if !rule.Keep(i) {
			if e = s.DeleteInstrument(i.ID); e != nil {
				return removed, e
			}
			removed++
		}
	}
	return removed, nil
}
func (s *Store) EnsureEntity(i model.Instrument) error {
	if i.ID == "" {
		return errors.New("entity id required")
	}
	if e := i.Validate(); e != nil {
		return e
	}
	return s.UpsertInstrument(i)
}
