package model

type FieldRule struct {
	Name      string
	Required  bool
	Max       int
	Sensitive bool
}

var InstrumentFieldRules = []FieldRule{
	{"编号", true, 40, false}, {"名称", true, 120, false}, {"实验室", true, 120, false}, {"负责人", true, 80, false}, {"购置日期", true, 20, false}, {"状态", true, 20, false}, {"校准记录", false, 80, false}, {"附件", false, 0, true},
}

func RequiredFieldCount() int {
	n := 0
	for _, r := range InstrumentFieldRules {
		if r.Required {
			n++
		}
	}
	return n
}
func FieldLimit(name string) int {
	for _, r := range InstrumentFieldRules {
		if r.Name == name {
			return r.Max
		}
	}
	return -1
}
func IsSensitiveField(name string) bool {
	for _, r := range InstrumentFieldRules {
		if r.Name == name {
			return r.Sensitive
		}
	}
	return false
}
func ValidateField(name, value string) bool {
	limit := FieldLimit(name)
	if limit < 0 {
		return false
	}
	if value == "" {
		for _, r := range InstrumentFieldRules {
			if r.Name == name {
				return !r.Required
			}
		}
	}
	return len([]rune(value)) <= limit
}
func StatusOptions() []Status {
	return []Status{StatusDraft, StatusPending, StatusApproved, StatusArchived}
}
func IsKnownStatus(s Status) bool {
	for _, v := range StatusOptions() {
		if v == s {
			return true
		}
	}
	return false
}
func TransitionAllowed(from, to Status) bool {
	switch from {
	case StatusDraft:
		return to == StatusPending
	case StatusPending:
		return to == StatusApproved || to == StatusDraft
	case StatusApproved:
		return to == StatusArchived
	}
	return false
}
func NextStatus(from Status, approve bool) Status {
	if from == StatusDraft {
		return StatusPending
	}
	if from == StatusPending && approve {
		return StatusApproved
	}
	if from == StatusPending {
		return StatusDraft
	}
	if from == StatusApproved {
		return StatusArchived
	}
	return from
}
func StatusRank(s Status) int {
	switch s {
	case StatusDraft:
		return 1
	case StatusPending:
		return 2
	case StatusApproved:
		return 3
	case StatusArchived:
		return 4
	}
	return 0
}
