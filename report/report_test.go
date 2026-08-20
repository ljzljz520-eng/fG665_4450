package report

import (
	"strings"
	"testing"
)

func TestCSVReport(t *testing.T) {
	rows, e := ParseCSV(strings.NewReader("id,number,name,lab,owner\n1,N1,Scope,Lab,A\n2,,Bad,Lab,B\n"))
	if e != nil {
		t.Fatal(e)
	}
	s := Summarize(rows)
	if s.Imported != 1 || s.Rejected != 1 {
		t.Fatal(s)
	}
	if Render(s) == "" {
		t.Fatal("render")
	}
}
