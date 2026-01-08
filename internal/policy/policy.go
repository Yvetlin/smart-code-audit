package policy

import "smart-code-audit/internal/models"

type Decision struct {
	FailPipeline bool
}

func Evaluate(cfg any, findings []models.EnrichedFinding) Decision {
	for _, f := range findings {
		if f.Review.Probability >= 0.7 {
			return Decision{FailPipeline: true}
		}
	}
	return Decision{FailPipeline: false}
}
