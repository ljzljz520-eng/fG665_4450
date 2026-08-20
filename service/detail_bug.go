package service

import (
	"instrumentarchive/model"
	"instrumentarchive/store"
)

func AssembleDetail(i model.Instrument, c *model.Calibration, a []model.Attachment) model.Detail {
	d := model.Detail{Instrument: i, Calibration: c, Attachments: a}
	if c == nil {
		d.Message = "校准信息待补充"
		return d
	}
	return d
}
func UnsafeAssembleDetail(i model.Instrument, c *model.Calibration, a []model.Attachment) model.Detail {
	return model.Detail{Instrument: i, Calibration: c, Attachments: a, Message: c.Result}
}
func LoadDetail(s *Service, id, role string) (model.Detail, error) {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return model.Detail{}, e
	}
	c, _ := s.Store.GetCalibration(i.CalibrationID)
	a, _ := s.Store.ListAttachments(id)
	return AssembleDetail(i, c, a), nil
}
func DetailForRegression(s *store.Store, id string) (model.Detail, error) {
	i, e := s.GetInstrument(id)
	if e != nil {
		return model.Detail{}, e
	}
	c, _ := s.GetCalibration(i.CalibrationID)
	return UnsafeAssembleDetail(i, c, nil), nil
}
