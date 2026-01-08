package markdown

import (
	"fmt"
	"os"
	"smart-code-audit/internal/models"
)

func Write(findings []models.EnrichedFinding, path string) error {
	var out string
	for _, f := range findings {
		out += fmt.Sprintf(
			"- %s:%d %s (%.2f)\n",
			f.Finding.File,
			f.Finding.Line,
			f.Review.Impact,
			f.Review.Probability,
		)
	}
	return os.WriteFile(path, []byte(out), 0644)
}
