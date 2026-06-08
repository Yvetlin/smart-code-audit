package policy

import (
	"smart-code-audit/internal/config"
	"smart-code-audit/internal/models"
)

func Evaluate(findings []models.EnrichedFinding, cfg config.PolicyConfig) bool {
	for _, f := range findings {
		if !f.Review.IsRealIssue {
			continue
		}
		if f.Review.Probability >= cfg.MinProbability {
			return true
		}
	}
	return false
}

func Confirmed(findings []models.EnrichedFinding, cfg config.PolicyConfig) []models.EnrichedFinding {
	var out []models.EnrichedFinding
	for _, f := range findings {
		if f.Review.IsRealIssue && f.Review.Probability >= cfg.MinProbability {
			out = append(out, f)
		}
	}
	return out
}
