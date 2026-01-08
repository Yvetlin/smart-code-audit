package analyzers

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type gosecReport struct {
	Issues []struct {
		RuleID   string `json:"rule_id"`
		Severity string `json:"severity"`
		Details  string `json:"details"`
		File     string `json:"file"`
		Line     string `json:"line"`
		Code     string `json:"code"`
	} `json:"Issues"`
}

func runGosec(target string) []RawFinding {
	cmd := exec.Command("gosec", "-fmt=json", target)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("gosec finished with findings (non-zero exit code)")
	}

	output := string(out)

	jsonStart := strings.Index(output, "{")
	if jsonStart == -1 {
		log.Println("gosec output does not contain JSON")
		return nil
	}

	var report gosecReport
	if err := json.Unmarshal([]byte(output[jsonStart:]), &report); err != nil {
		log.Printf("failed to parse gosec JSON: %v", err)
		return nil
	}

	var findings []RawFinding
	for _, i := range report.Issues {
		line, _ := strconv.Atoi(i.Line)

		findings = append(findings, RawFinding{
			Tool:     "gosec",
			RuleID:   i.RuleID,
			Severity: i.Severity,
			Message:  i.Details,
			File:     i.File,
			Line:     line,
			Snippet:  i.Code,
		})
	}

	return findings
}
