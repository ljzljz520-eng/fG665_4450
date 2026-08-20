package store

import (
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
	"instrumentarchive/model"
	"os"
	"sync"
)

var bucketNames = map[string][]byte{"instruments": []byte("instruments"), "calibrations": []byte("calibrations"), "attachments": []byte("attachments"), "audits": []byte("audits"), "workflows": []byte("workflows")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err = s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) put(bucket, key string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketNames[bucket]).Put([]byte(key), data) })
}
func (s *Store) get(bucket, key string, v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketNames[bucket]).Get([]byte(key))
		if b == nil {
			return os.ErrNotExist
		}
		return json.Unmarshal(append([]byte(nil), b...), v)
	})
}
func (s *Store) del(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucketNames[bucket]).Delete([]byte(key)) })
}
func (s *Store) SaveInstrument(v model.Instrument) error { return s.put("instruments", v.ID, v) }
func (s *Store) GetInstrument(id string) (model.Instrument, error) {
	var v model.Instrument
	err := s.get("instruments", id, &v)
	return v, err
}
func (s *Store) DeleteInstrument(id string) error          { return s.del("instruments", id) }
func (s *Store) SaveCalibration(v model.Calibration) error { return s.put("calibrations", v.ID, v) }
func (s *Store) GetCalibration(id string) (*model.Calibration, error) {
	var v model.Calibration
	err := s.get("calibrations", id, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (s *Store) SaveAttachment(v model.Attachment) error   { return s.put("attachments", v.ID, v) }
func (s *Store) SaveAudit(v model.AuditEvent) error        { return s.put("audits", v.ID, v) }
func (s *Store) SaveWorkflow(v model.WorkflowRecord) error { return s.put("workflows", v.ID, v) }
func (s *Store) ListInstruments() ([]model.Instrument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []model.Instrument{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketNames["instruments"]).ForEach(func(_, raw []byte) error {
			var v model.Instrument
			if err := json.Unmarshal(raw, &v); err != nil {
				return err
			}
			out = append(out, v)
			return nil
		})
	})
	return out, err
}
func (s *Store) ListAttachments(instrumentID string) ([]model.Attachment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []model.Attachment{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketNames["attachments"]).ForEach(func(_, raw []byte) error {
			var v model.Attachment
			if e := json.Unmarshal(raw, &v); e != nil {
				return e
			}
			if v.InstrumentID == instrumentID {
				out = append(out, v)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) ListAudits(instrumentID string) ([]model.AuditEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []model.AuditEvent{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucketNames["audits"]).ForEach(func(_, raw []byte) error {
			var v model.AuditEvent
			if e := json.Unmarshal(raw, &v); e != nil {
				return e
			}
			if v.InstrumentID == instrumentID {
				out = append(out, v)
			}
			return nil
		})
	})
	return out, err
}
