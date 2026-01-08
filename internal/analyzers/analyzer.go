package analyzers

import (
	"smart-code-audit/internal/models"
)

type RawFinding = models.Finding

func RunAll(target string) []RawFinding {
	var findings []RawFinding
	findings = append(findings, runGosec(target)...)
	return findings
}
