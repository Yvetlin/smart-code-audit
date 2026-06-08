package markdown

import (
	"fmt"
	"os"
	"strings"

	"smart-code-audit/internal/models"
)

func Write(path string, findings []models.EnrichedFinding) error {
	var b strings.Builder

	b.WriteString("# Smart Code Audit\n\n")

	if len(findings) == 0 {
		b.WriteString("Проблем не обнаружено.\n")
		return os.WriteFile(path, []byte(b.String()), 0644)
	}

	for i, f := range findings {
		verdict := "ложное срабатывание"
		if f.Review.IsRealIssue {
			verdict = "подтверждено"
		}

		b.WriteString(fmt.Sprintf("## %d. %s:%d\n\n", i+1, f.Finding.File, f.Finding.Line))
		b.WriteString(fmt.Sprintf("- **Сканер:** %s / %s\n", f.Finding.Tool, f.Finding.RuleID))
		b.WriteString(fmt.Sprintf("- **Вердикт ИИ:** %s (%.0f%%)\n", verdict, f.Review.Probability*100))
		b.WriteString(fmt.Sprintf("- **Что не так:** %s\n", f.Review.Brief()))
		b.WriteString(fmt.Sprintf("- **Как исправить:** %s\n", f.Review.Fix))
		if f.Review.Impact != "" && f.Review.Impact != f.Review.Brief() {
			b.WriteString(fmt.Sprintf("- **Риск:** %s\n", f.Review.Impact))
		}
		b.WriteString("\n")
	}

	return os.WriteFile(path, []byte(b.String()), 0644)
}
