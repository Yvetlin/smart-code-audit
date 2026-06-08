package analyzers

import (
	"encoding/json"
	"log"
	"os/exec"
)

type golangciReport struct {
	Issues []struct {
		FromLinter string `json:"FromLinter"`
		Text       string `json:"Text"`
		Pos        struct {
			Filename string `json:"Filename"`
			Line     int    `json:"Line"`
		} `json:"Pos"`
	} `json:"Issues"`
}

func runGolangci(target string) []RawFinding {
	cmd := exec.Command("golangci-lint", "run", "--timeout", "1m", "--output.json.path=stdout", "./...")
	cmd.Dir = target

	out, err := cmd.CombinedOutput()

	log.Println("RAW GOLANGCI OUTPUT:")
	log.Println(string(out))

	if err != nil {
		log.Println("golangci finished (non-zero exit is OK)")
	}

	var report golangciReport
	if err := json.Unmarshal(out, &report); err != nil {
		log.Printf("failed to parse golangci JSON: %v", err)
		return nil
	}

	var findings []RawFinding

	for _, i := range report.Issues {
		findings = append(findings, RawFinding{
			Tool:     "golangci",
			RuleID:   i.FromLinter,
			Severity: "INFO",
			Message:  i.Text,
			File:     i.Pos.Filename,
			Line:     i.Pos.Line,
		})
	}

	log.Printf("golangci findings parsed: %d", len(findings))

	return findings
}
