package analyzers

import (
	"log"
	"smart-code-audit/internal/models"
)

type RawFinding = models.Finding

func RunAll(target string) []RawFinding {
	var findings []RawFinding

	findings = append(findings, runGosec(target)...)
	findings = append(findings, runSemgrep(target)...)
	findings = append(findings, runGolangci(target)...)

	log.Printf("TOTAL findings: %d", len(findings))

	return findings
}
