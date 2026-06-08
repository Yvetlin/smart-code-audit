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
	cmd := exec.Command("gosec", "-fmt=json", "./...")
	cmd.Dir = target

	out, err := cmd.CombinedOutput()

	log.Println("RAW GOSEC OUTPUT:")
	log.Println(string(out))

	if err != nil {
		log.Println("gosec finished (non-zero exit code is OK)")
	}

	output := string(out)

	// 🔥 ВЫРЕЗАЕМ JSON
	idx := strings.Index(output, "{")
	if idx == -1 {
		log.Println("no JSON found in gosec output")
		return nil
	}

	clean := output[idx:]

	var report gosecReport
	if err := json.Unmarshal([]byte(clean), &report); err != nil {
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

	log.Printf("gosec findings parsed: %d", len(findings))

	return findings
}
