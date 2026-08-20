package workflow

import (
	"fmt"
	"instrumentarchive/model"
	"instrumentarchive/report"
	"instrumentarchive/service"
	"instrumentarchive/store"
	"strings"
)

type Engine struct{ Service *service.Service }

func New(s *service.Service) *Engine { return &Engine{Service: s} }
func (e *Engine) CreateReviewArchive(i model.Instrument, actor string) error {
	if err := e.Service.Create(i, actor); err != nil {
		return err
	}
	if err := e.Service.Submit(i.ID, actor); err != nil {
		return err
	}
	if err := e.Service.Review(i.ID, actor, true); err != nil {
		return err
	}
	return e.Service.Archive(i.ID, actor)
}
func (e *Engine) SearchUpdatePublish(f model.Filter, owner, actor string) (int, error) {
	items, err := e.Service.Search(f, actor)
	if err != nil {
		return 0, err
	}
	for _, i := range items {
		if err = e.Service.UpdateOwner(i.ID, owner, actor); err != nil {
			return 0, err
		}
	}
	return len(items), nil
}
func ImportAndReport(s *store.Store, input string) (report.Summary, error) {
	rows, e := report.ParseCSV(strings.NewReader(input))
	if e != nil {
		return report.Summary{}, e
	}
	n, e := s.Import(report.ValidRows(rows))
	sum := report.Summarize(rows)
	if e != nil {
		return sum, e
	}
	if n != sum.Imported {
		return sum, fmt.Errorf("import count mismatch")
	}
	return sum, nil
}
