package workflow

import (
	"errors"
	"instrumentarchive/model"
)

type Plan struct {
	ID           string
	InstrumentID string
	Actions      []string
	Index        int
}

func NewPlan(id, instrumentID string) Plan {
	return Plan{ID: id, InstrumentID: instrumentID, Actions: []string{"登记", "提交审核", "审核", "归档"}}
}
func (p *Plan) Next() (string, error) {
	if p.Index >= len(p.Actions) {
		return "", errors.New("plan complete")
	}
	a := p.Actions[p.Index]
	p.Index++
	return a, nil
}
func (p Plan) Done() bool { return p.Index >= len(p.Actions) }
func (p Plan) Progress() float64 {
	if len(p.Actions) == 0 {
		return 1
	}
	return float64(p.Index) / float64(len(p.Actions))
}
func (p Plan) ActionForStatus(s model.Status) string {
	switch s {
	case model.StatusDraft:
		return "提交审核"
	case model.StatusPending:
		return "审核"
	case model.StatusApproved:
		return "归档"
	default:
		return "完成"
	}
}
func (p *Plan) Reset() { p.Index = 0 }
func (p Plan) RemainingActions() []string {
	if p.Index >= len(p.Actions) {
		return []string{}
	}
	return append([]string(nil), p.Actions[p.Index:]...)
}
