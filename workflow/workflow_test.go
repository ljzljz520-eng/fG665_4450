package workflow

import (
	"instrumentarchive/model"
	"instrumentarchive/service"
	"instrumentarchive/store"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	eng := New(service.New(s, service.StaticClock{"t"}))
	i := model.NewInstrument("i1", "N1", "Scope", "Lab", "A", "2024")
	if e = eng.CreateReviewArchive(i, "admin"); e != nil {
		t.Fatal(e)
	}
	got, _ := s.GetInstrument("i1")
	if !got.Archived {
		t.Fatal("not archived")
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	svc := service.New(s, service.StaticClock{"t"})
	i := model.NewInstrument("i1", "N1", "Scope", "Lab", "A", "2024")
	if e = svc.Create(i, "admin"); e != nil {
		t.Fatal(e)
	}
	n, e := New(svc).SearchUpdatePublish(model.Filter{Query: "N1"}, "B", "admin")
	if e != nil || n != 1 {
		t.Fatal(n, e)
	}
}
