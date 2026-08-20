package store

import (
	"instrumentarchive/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "archive.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	i := model.NewInstrument("i1", "NO-1", "Scope", "Lab", "Alice", "2024-01-01")
	if e = s.SaveInstrument(i); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetInstrument("i1")
	if e != nil || got.Number != "NO-1" {
		t.Fatalf("%+v %v", got, e)
	}
}
func TestStoreEntities(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	c := model.NewCalibration("c1", "i1", "2025", "pass", "ok")
	if e = s.SaveCalibration(c); e != nil {
		t.Fatal(e)
	}
	got, e := s.GetCalibration("c1")
	if e != nil || got.Result != "pass" {
		t.Fatal(e)
	}
}
