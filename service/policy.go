package service

import (
	"errors"
	"instrumentarchive/model"
)

type Policy struct {
	RequireOwner                    bool
	RequireCalibrationBeforeArchive bool
	AllowViewerArchived             bool
}

func DefaultPolicy() Policy {
	return Policy{RequireOwner: true, RequireCalibrationBeforeArchive: false, AllowViewerArchived: false}
}
func (p Policy) ValidateCreate(i model.Instrument) error {
	if p.RequireOwner && i.Owner == "" {
		return errors.New("owner required")
	}
	return i.Validate()
}
func (p Policy) CanPublish(i model.Instrument) bool {
	return i.Status == model.StatusApproved && !i.Archived
}
func (p Policy) CanViewerSee(i model.Instrument) bool {
	if i.Archived && !p.AllowViewerArchived {
		return false
	}
	return true
}
func (p Policy) ArchiveReady(i model.Instrument, c *model.Calibration) bool {
	if i.Status != model.StatusApproved {
		return false
	}
	if p.RequireCalibrationBeforeArchive && c == nil {
		return false
	}
	return true
}
func (p Policy) MissingFields(i model.Instrument) []string {
	out := []string{}
	if i.Number == "" {
		out = append(out, "编号")
	}
	if i.Name == "" {
		out = append(out, "名称")
	}
	if i.Laboratory == "" {
		out = append(out, "实验室")
	}
	if i.Owner == "" {
		out = append(out, "负责人")
	}
	return out
}
func (s *Service) Publish(id, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if !DefaultPolicy().CanPublish(i) {
		return errors.New("not publishable")
	}
	return s.audit(id, actor, "publish", "")
}
func (s *Service) ValidateForImport(i model.Instrument) error {
	return DefaultPolicy().ValidateCreate(i)
}
