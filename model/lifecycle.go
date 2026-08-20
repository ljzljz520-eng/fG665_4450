package model

import "strings"

type LifecycleEvent struct {
	From   Status
	To     Status
	Actor  string
	Reason string
}

func Lifecycle(from, to Status, actor, reason string) LifecycleEvent {
	return LifecycleEvent{From: from, To: to, Actor: actor, Reason: reason}
}
func (e LifecycleEvent) Valid() bool {
	if e.Actor == "" {
		return false
	}
	if !IsKnownStatus(e.From) || !IsKnownStatus(e.To) {
		return false
	}
	return TransitionAllowed(e.From, e.To)
}
func (e LifecycleEvent) Label() string {
	switch e.To {
	case StatusPending:
		return "提交审核"
	case StatusApproved:
		return "审核通过"
	case StatusArchived:
		return "归档"
	case StatusDraft:
		return "退回草稿"
	}
	return "未知"
}
func NormalizeInstrument(i Instrument) Instrument {
	i.Number = strings.TrimSpace(i.Number)
	i.Name = strings.TrimSpace(i.Name)
	i.Laboratory = strings.TrimSpace(i.Laboratory)
	i.Owner = strings.TrimSpace(i.Owner)
	if i.Status == "" {
		i.Status = StatusDraft
	}
	return i
}
func (i Instrument) SearchTokens() []string { return []string{i.Number, i.Name, i.Laboratory, i.Owner} }
func (i Instrument) MatchesAny(terms []string) bool {
	for _, term := range terms {
		q := strings.ToLower(strings.TrimSpace(term))
		if q == "" {
			continue
		}
		for _, token := range i.SearchTokens() {
			if strings.Contains(strings.ToLower(token), q) {
				return true
			}
		}
	}
	return false
}
func (i Instrument) MetadataComplete() bool {
	return i.Number != "" && i.Name != "" && i.Laboratory != "" && i.Owner != "" && i.PurchaseDate != ""
}
func (i Instrument) CalibrationPending(c *Calibration) bool {
	if c == nil {
		return true
	}
	return strings.TrimSpace(c.Result) == "" || strings.TrimSpace(c.DueDate) == ""
}
func (i Instrument) ArchiveLabel() string {
	if i.Archived {
		return "档案已归档"
	}
	return "档案活动"
}
func (i Instrument) OwnerDisplay() string {
	if i.Owner == "" {
		return "未分配"
	}
	return i.Owner
}
func (i Instrument) NumberDisplay() string {
	if i.Number == "" {
		return "未编号"
	}
	return i.Number
}
