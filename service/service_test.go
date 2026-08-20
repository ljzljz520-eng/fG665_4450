package service

import (
	"instrumentarchive/model"
	"instrumentarchive/store"
	"path/filepath"
	"testing"
)

func TestServiceLifecycle(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	svc := New(s, StaticClock{"fixed"})
	i := model.NewInstrument("i1", "N1", "Scope", "Lab", "A", "2024")
	if e = svc.Create(i, "admin"); e != nil {
		t.Fatal(e)
	}
	if e = svc.Submit("i1", "admin"); e != nil {
		t.Fatal(e)
	}
	if e = svc.Review("i1", "admin", true); e != nil {
		t.Fatal(e)
	}
	if e = svc.Archive("i1", "admin"); e != nil {
		t.Fatal(e)
	}
	if lenMust(svc.Audit("i1")) != 4 {
		t.Fatal("audit")
	}
}
func lenMust(v []model.AuditEvent, e error) int {
	if e != nil {
		return 0
	}
	return len(v)
}
func TestDetailMissingCalibration(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	svc := New(s, StaticClock{"fixed"})
	i := model.NewInstrument("i1", "N1", "Scope", "Lab", "A", "2024")
	if e = svc.Create(i, "admin"); e != nil {
		t.Fatal(e)
	}
	d, e := svc.Detail("i1", "admin")
	if e != nil || d.Message != "校准信息待补充" {
		t.Fatalf("%+v %v", d, e)
	}
}
