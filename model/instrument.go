package model

import (
	"errors"
	"strings"
)

type Status string

const (
	StatusDraft    Status = "draft"
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusArchived Status = "archived"
)

type Instrument struct {
	ID            string `json:"id"`
	Number        string `json:"number"`
	Name          string `json:"name"`
	Laboratory    string `json:"laboratory"`
	Owner         string `json:"owner"`
	PurchaseDate  string `json:"purchase_date"`
	Status        Status `json:"status"`
	CalibrationID string `json:"calibration_id"`
	Archived      bool   `json:"archived"`
}

type Calibration struct {
	ID           string `json:"id"`
	InstrumentID string `json:"instrument_id"`
	DueDate      string `json:"due_date"`
	Result       string `json:"result"`
	Notes        string `json:"notes"`
}
type Attachment struct {
	ID           string `json:"id"`
	InstrumentID string `json:"instrument_id"`
	Name         string `json:"name"`
	ContentType  string `json:"content_type"`
	Data         []byte `json:"data"`
}
type AuditEvent struct {
	ID           string `json:"id"`
	InstrumentID string `json:"instrument_id"`
	Actor        string `json:"actor"`
	Action       string `json:"action"`
	At           string `json:"at"`
	Comment      string `json:"comment"`
}
type WorkflowRecord struct {
	ID           string `json:"id"`
	InstrumentID string `json:"instrument_id"`
	State        string `json:"state"`
	Actor        string `json:"actor"`
	UpdatedAt    string `json:"updated_at"`
}
type Detail struct {
	Instrument  Instrument   `json:"instrument"`
	Calibration *Calibration `json:"calibration"`
	Attachments []Attachment `json:"attachments"`
	Message     string       `json:"message"`
}

func (i Instrument) Validate() error {
	if strings.TrimSpace(i.Number) == "" {
		return errors.New("instrument number is required")
	}
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("instrument name is required")
	}
	if strings.TrimSpace(i.Laboratory) == "" {
		return errors.New("laboratory is required")
	}
	if i.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
func (i Instrument) IsVisible(role string) bool {
	if role == "admin" {
		return true
	}
	return !i.Archived
}
func (i Instrument) CanReview(role string) bool { return role == "admin" && i.Status == StatusPending }
func (i Instrument) CanArchive(role string) bool {
	return role == "admin" && i.Status == StatusApproved
}
func (i Instrument) DisplayStatus() string {
	switch i.Status {
	case StatusDraft:
		return "草稿"
	case StatusPending:
		return "待审核"
	case StatusApproved:
		return "已批准"
	case StatusArchived:
		return "已归档"
	}
	return "未知"
}
func NewInstrument(id, number, name, lab, owner, date string) Instrument {
	return Instrument{ID: id, Number: number, Name: name, Laboratory: lab, Owner: owner, PurchaseDate: date, Status: StatusDraft}
}
func NewCalibration(id, instrumentID, due, result, notes string) Calibration {
	return Calibration{ID: id, InstrumentID: instrumentID, DueDate: due, Result: result, Notes: notes}
}
func NewAttachment(id, instrumentID, name, kind string, data []byte) Attachment {
	return Attachment{ID: id, InstrumentID: instrumentID, Name: name, ContentType: kind, Data: data}
}
func NewAudit(id, instrumentID, actor, action, at, comment string) AuditEvent {
	return AuditEvent{ID: id, InstrumentID: instrumentID, Actor: actor, Action: action, At: at, Comment: comment}
}
func NewWorkflow(id, instrumentID, state, actor, at string) WorkflowRecord {
	return WorkflowRecord{ID: id, InstrumentID: instrumentID, State: state, Actor: actor, UpdatedAt: at}
}
