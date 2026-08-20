package workflow

import (
	"errors"
	"instrumentarchive/model"
	"instrumentarchive/store"
)

type Checkpoint struct {
	Name   string
	Done   bool
	Detail string
}
type Chain struct {
	ID    string
	Steps []Checkpoint
}

func NewChain(id string, names []string) Chain {
	c := Chain{ID: id, Steps: []Checkpoint{}}
	for _, n := range names {
		c.Steps = append(c.Steps, Checkpoint{Name: n})
	}
	return c
}
func (c *Chain) Advance(detail string) error {
	for i := range c.Steps {
		if !c.Steps[i].Done {
			c.Steps[i].Done = true
			c.Steps[i].Detail = detail
			return nil
		}
	}
	return errors.New("chain complete")
}
func (c Chain) Complete() bool {
	for _, s := range c.Steps {
		if !s.Done {
			return false
		}
	}
	return true
}
func (c Chain) Remaining() int {
	n := 0
	for _, s := range c.Steps {
		if !s.Done {
			n++
		}
	}
	return n
}
func (c Chain) Current() string {
	for _, s := range c.Steps {
		if !s.Done {
			return s.Name
		}
	}
	return "done"
}
func SaveCheckpoint(s *store.Store, c Chain, instrumentID, actor, at string) error {
	state := c.Current()
	if c.Complete() {
		state = "done"
	}
	return s.SaveWorkflow(model.NewWorkflow(c.ID, instrumentID, state, actor, at))
}
func RestoreCheckpoint(s *store.Store, id string) (model.WorkflowRecord, error) {
	items, e := s.Snapshot()
	if e != nil {
		return model.WorkflowRecord{}, e
	}
	for _, w := range items.Workflows {
		if w.ID == id {
			return w, nil
		}
	}
	return model.WorkflowRecord{}, errors.New("checkpoint missing")
}
