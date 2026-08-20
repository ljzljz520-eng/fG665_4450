package model

import "testing"

func TestInstrumentValidation(t *testing.T) {
	i := NewInstrument("1", "N1", "scope", "lab", "owner", "2024")
	if err := i.Validate(); err != nil {
		t.Fatal(err)
	}
	i.Number = ""
	if i.Validate() == nil {
		t.Fatal("expected error")
	}
}
func TestFilterAndStatus(t *testing.T) {
	i := NewInstrument("1", "N1", "scope", "lab", "owner", "2024")
	i.Status = StatusPending
	if !(Filter{Query: "scope"}).Match(i) {
		t.Fatal("match")
	}
	if NormalizeStatus("已批准") != StatusApproved || IsTerminal(StatusDraft) {
		t.Fatal("status")
	}
}
