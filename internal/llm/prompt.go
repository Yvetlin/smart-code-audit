package llm

import (
	"fmt"

	"smart-code-audit/internal/models"
)

const SystemPrompt = `
You are an expert secure code reviewer working in a CI/CD pipeline.

Your task:
- determine whether the reported issue is a real security problem
- estimate the probability that the issue is real (0.0 - 1.0)
- explain the impact
- suggest a concrete fix

Return ONLY valid JSON.
Do NOT add explanations outside JSON.

JSON schema:
{
  "is_real_issue": true,
  "probability": 0.0,
  "impact": "string",
  "fix": "string",
  "explanation": "string"
}
`

func BuildUserPrompt(f models.Finding) string {
	return fmt.Sprintf(
		`Tool: %s
Rule ID: %s
Severity: %s
Message: %s
File: %s
Line: %d

Code snippet:
%s
`,
		f.Tool,
		f.RuleID,
		f.Severity,
		f.Message,
		f.File,
		f.Line,
		f.Snippet,
	)
}
