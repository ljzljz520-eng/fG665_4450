package service

import (
	"errors"
	"instrumentarchive/model"
	"instrumentarchive/report"
)

type Dashboard struct {
	Total              int
	Active             int
	Archived           int
	PendingCalibration int
	Labs               map[string]int
	Status             map[model.Status]int
}

func (s *Service) Dashboard(role string) (Dashboard, error) {
	items, e := s.Search(model.Filter{IncludeArchived: role == "admin"}, role)
	if e != nil {
		return Dashboard{}, e
	}
	d := Dashboard{Labs: map[string]int{}, Status: map[model.Status]int{}}
	for _, i := range items {
		d.Total++
		d.Status[i.Status]++
		d.Labs[i.Laboratory]++
		if i.Archived {
			d.Archived++
		} else {
			d.Active++
		}
		if s.CalibrationStatus(i.ID) == "pending" {
			d.PendingCalibration++
		}
	}
	return d, nil
}
func (s *Service) ExportMarkdown(role string) (string, error) {
	items, e := s.Search(model.Filter{IncludeArchived: role == "admin"}, role)
	if e != nil {
		return "", e
	}
	return report.RenderMarkdown(items), nil
}
func (s *Service) ValidateImportRows(rows []report.Row) error {
	for _, r := range rows {
		if r.Error != "" {
			continue
		}
		if e := s.ValidateForImport(r.Instrument); e != nil {
			return e
		}
	}
	return nil
}
func (s *Service) PublishCalibration(id, actor string, c model.Calibration) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	if c.InstrumentID != id {
		return errors.New("instrument mismatch")
	}
	if e := s.AddCalibration(c, actor); e != nil {
		return e
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	i.CalibrationID = c.ID
	return s.Store.SaveInstrument(i)
}
func (s *Service) ArchiveWithPolicy(id, actor string, p Policy) error {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	c, _ := s.Store.GetCalibration(i.CalibrationID)
	if !p.ArchiveReady(i, c) {
		return errors.New("archive policy rejected")
	}
	return s.Archive(id, actor)
}
func (s *Service) MissingMetadata(id string) ([]string, error) {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return nil, e
	}
	return DefaultPolicy().MissingFields(i), nil
}
