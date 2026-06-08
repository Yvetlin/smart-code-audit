package analyzers

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type semgrepReport struct {
	Results []struct {
		CheckID string `json:"check_id"`
		Path    string `json:"path"`
		Start   struct {
			Line int `json:"line"`
		} `json:"start"`
		Extra struct {
			Message  string `json:"message"`
			Severity string `json:"severity"`
		} `json:"extra"`
	} `json:"results"`
}

func runSemgrep(target string) []RawFinding {
	cmd := exec.Command("semgrep", "--config=p/ci", "--json", ".")
	cmd.Dir = target

	out, err := cmd.CombinedOutput()

	log.Println("RAW SEMGREP OUTPUT:")
	log.Println(string(out))

	if err != nil {
		log.Println("semgrep finished (non-zero exit is OK)")
	}

	output := string(out)

	// 🔥 ВЫРЕЗАЕМ JSON
	idx := strings.Index(output, "{")
	if idx == -1 {
		log.Println("no JSON found in semgrep output")
		return nil
	}

	clean := output[idx:]

	var report semgrepReport
	if err := json.Unmarshal([]byte(clean), &report); err != nil {
		log.Printf("failed to parse semgrep JSON: %v", err)
		return nil
	}

	var findings []RawFinding

	for _, r := range report.Results {
		findings = append(findings, RawFinding{
			Tool:     "semgrep",
			RuleID:   r.CheckID,
			Severity: r.Extra.Severity,
			Message:  r.Extra.Message,
			File:     r.Path,
			Line:     r.Start.Line,
		})
	}

	log.Printf("semgrep findings parsed: %d", len(findings))

	return findings
}
