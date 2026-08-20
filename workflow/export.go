package workflow

import (
	"encoding/json"
	"fmt"
	"instrumentarchive/model"
	"instrumentarchive/report"
	"instrumentarchive/store"
)

type ExportBundle struct {
	Snapshot store.Snapshot
	Summary  report.Summary
}

func Export(s *store.Store) ([]byte, error) {
	snap, e := s.Snapshot()
	if e != nil {
		return nil, e
	}
	sum := report.Summary{Total: len(snap.Instruments), Imported: len(snap.Instruments)}
	return json.Marshal(ExportBundle{Snapshot: snap, Summary: sum})
}
func ValidateImport(rows []report.Row) ([]model.Instrument, []string) {
	valid := []model.Instrument{}
	errors := []string{}
	for _, r := range rows {
		if r.Error != "" {
			errors = append(errors, r.Error)
		} else {
			valid = append(valid, r.Instrument)
		}
	}
	return valid, errors
}
func ImportWithAudit(s *store.Store, rows []report.Row, actor string) (report.Summary, error) {
	valid, errs := ValidateImport(rows)
	n, e := s.Import(valid)
	sum := report.Summary{Total: len(rows), Imported: n, Rejected: len(errs), Errors: errs}
	if e != nil {
		return sum, fmt.Errorf("import: %w", e)
	}
	return sum, nil
}
func BuildWorkflow(id, instrumentID, actor, at string) model.WorkflowRecord {
	return model.NewWorkflow(id, instrumentID, "imported", actor, at)
}
