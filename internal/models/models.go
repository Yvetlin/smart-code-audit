package models

type Finding struct {
	Tool     string
	RuleID   string
	Severity string
	Message  string
	File     string
	Line     int
	Snippet  string
}

type AIReview struct {
	IsRealIssue bool    `json:"is_real_issue"`
	Probability float64 `json:"probability"`
	Summary     string  `json:"summary"`
	Impact      string  `json:"impact"`
	Fix         string  `json:"fix"`
	Explanation string  `json:"explanation"`
}

func (r AIReview) Brief() string {
	if r.Summary != "" {
		return r.Summary
	}
	if r.Impact != "" {
		return r.Impact
	}
	return r.Explanation
}

type EnrichedFinding struct {
	Finding Finding
	Review  AIReview
}
