package policy

import "smart-code-audit/internal/models"

func Evaluate(findings []models.EnrichedFinding) bool {
	for _, f := range findings {
		if f.Review.Probability >= 0.7 {
			return true
		}
	}
	return false
}
