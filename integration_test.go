package instrumentarchive

import (
	"instrumentarchive/model"
	"instrumentarchive/service"
	"instrumentarchive/store"
	"instrumentarchive/workflow"
	"path/filepath"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	sum, e := workflow.ImportAndReport(s, "id,number,name,lab,owner\n1,N1,Scope,Lab,A\n")
	if e != nil || sum.Imported != 1 {
		t.Fatal(sum, e)
	}
}
func TestBusiness06Regression(t *testing.T) {
	s, e := store.Open(filepath.Join(t.TempDir(), "a.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	svc := service.New(s, service.StaticClock{"t"})
	i := model.NewInstrument("i1", "NO-06", "未校准仪器", "实验室A", "管理员", "2024-01-01")
	if e = svc.Create(i, "admin"); e != nil {
		t.Fatal(e)
	}
	d, e := svc.Detail("i1", "admin")
	if e != nil {
		t.Fatal(e)
	}
	if d.Instrument.Number != "NO-06" || d.Message != "校准信息待补充" {
		t.Fatalf("unexpected detail: %+v", d)
	}
}
