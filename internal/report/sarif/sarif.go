package sarif

import (
	"encoding/json"
	"os"

	"smart-code-audit/internal/models"
)

func Write(path string, findings []models.EnrichedFinding) error {
	b, _ := json.MarshalIndent(findings, "", "  ")
	return os.WriteFile(path, b, 0644)
}
