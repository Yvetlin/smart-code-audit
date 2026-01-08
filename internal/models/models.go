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
	IsRealIssue bool
	Probability float64
	Impact      string
	Fix         string
	Explanation string
}

type EnrichedFinding struct {
	Finding Finding
	Review  AIReview
}
