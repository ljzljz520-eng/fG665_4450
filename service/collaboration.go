package service

import (
	"errors"
	"fmt"
)

type Comment struct {
	ID           string
	InstrumentID string
	Author       string
	Body         string
	CreatedAt    string
}

var comments = map[string][]Comment{}

func (s *Service) AddComment(instrumentID, author, body string) (Comment, error) {
	if author == "" || body == "" {
		return Comment{}, errors.New("author and body required")
	}
	if !s.Store.HasInstrument(instrumentID) {
		return Comment{}, errors.New("instrument missing")
	}
	c := Comment{ID: fmt.Sprintf("%s-%d", instrumentID, len(comments[instrumentID])+1), InstrumentID: instrumentID, Author: author, Body: body, CreatedAt: s.Clock.Now()}
	comments[instrumentID] = append(comments[instrumentID], c)
	return c, nil
}
func (s *Service) Comments(instrumentID string) []Comment {
	return append([]Comment(nil), comments[instrumentID]...)
}
func (s *Service) ClearComments(instrumentID string) { delete(comments, instrumentID) }
func (s *Service) Collaborators(instrumentID string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, c := range comments[instrumentID] {
		if !seen[c.Author] {
			seen[c.Author] = true
			out = append(out, c.Author)
		}
	}
	return out
}
func (s *Service) AddCalibrationNotes(id, notes, actor string) error {
	if actor != "admin" {
		return errors.New("admin required")
	}
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return e
	}
	if i.CalibrationID == "" {
		return errors.New("calibration missing")
	}
	c, e := s.Store.GetCalibration(i.CalibrationID)
	if e != nil {
		return e
	}
	c.Notes = notes
	return s.Store.SaveCalibration(*c)
}
func (s *Service) CalibrationStatus(id string) string {
	i, e := s.Store.GetInstrument(id)
	if e != nil {
		return "missing"
	}
	c, e := s.Store.GetCalibration(i.CalibrationID)
	if e != nil || i.CalibrationPending(c) {
		return "pending"
	}
	return "complete"
}
