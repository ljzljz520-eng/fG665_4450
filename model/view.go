package model

import (
	"fmt"
	"strings"
)

type InstrumentView struct {
	ID              string
	Number          string
	Name            string
	Laboratory      string
	Owner           string
	PurchaseDate    string
	Status          string
	Archived        bool
	Calibration     string
	AttachmentCount int
}

func BuildView(i Instrument, c *Calibration, attachments []Attachment) InstrumentView {
	cal := "校准信息待补充"
	if c != nil {
		cal = c.Result
		if cal == "" {
			cal = "校准信息待补充"
		}
	}
	return InstrumentView{ID: i.ID, Number: i.NumberDisplay(), Name: i.Name, Laboratory: i.Laboratory, Owner: i.OwnerDisplay(), PurchaseDate: i.PurchaseDate, Status: i.DisplayStatus(), Archived: i.Archived, Calibration: cal, AttachmentCount: len(attachments)}
}
func (v InstrumentView) SearchText() string {
	return strings.Join([]string{v.Number, v.Name, v.Laboratory, v.Owner, v.Status, v.Calibration}, " ")
}
func (v InstrumentView) CSV() []string {
	return []string{v.ID, v.Number, v.Name, v.Laboratory, v.Owner, v.PurchaseDate, v.Status, fmt.Sprint(v.Archived), v.Calibration, fmt.Sprint(v.AttachmentCount)}
}
func (v InstrumentView) IsArchived() bool { return v.Archived }
func (v InstrumentView) NeedsAttention() bool {
	return v.Calibration == "校准信息待补充" || v.Owner == "未分配"
}
func (v InstrumentView) Badge() string {
	if v.NeedsAttention() {
		return "待完善"
	}
	if v.Archived {
		return "已归档"
	}
	return "正常"
}
func (v InstrumentView) Compact() string { return v.Number + " / " + v.Name + " / " + v.Status }
func (v InstrumentView) SameInstrument(other InstrumentView) bool {
	return v.ID != "" && v.ID == other.ID
}
