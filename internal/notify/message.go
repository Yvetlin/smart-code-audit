package notify

import (
	"fmt"
	"os"
	"strings"

	"smart-code-audit/internal/config"
	"smart-code-audit/internal/models"
	"smart-code-audit/internal/policy"
)

func BuildMessage(findings []models.EnrichedFinding, critical bool, cfg config.PolicyConfig) string {
	confirmed := policy.Confirmed(findings, cfg)

	var msg strings.Builder
	msg.WriteString("🔍 *Smart Code Audit*\n\n")

	if critical {
		msg.WriteString("❌ ИИ подтвердил проблемы безопасности\n\n")
	} else {
		msg.WriteString("✅ Критических проблем не обнаружено\n\n")
	}

	msg.WriteString(fmt.Sprintf("Сканеров: %d | Подтверждено ИИ: %d\n\n", len(findings), len(confirmed)))

	limit := 5
	shown := 0
	for _, f := range confirmed {
		if shown >= limit {
			msg.WriteString("\n...и ещё\n")
			break
		}
		shown++

		msg.WriteString(fmt.Sprintf(
			"• `%s:%d` (%.0f%%)\n%s\n💡 _%s_\n\n",
			shortPath(f.Finding.File),
			f.Finding.Line,
			f.Review.Probability*100,
			f.Review.Brief(),
			f.Review.Fix,
		))
	}

	if shown == 0 && len(findings) > 0 {
		msg.WriteString("ИИ отфильтровал findings как ложные срабатывания.\n\n")
	}

	runURL := os.Getenv("GITHUB_RUN_URL")
	if runURL == "" {
		runURL = "https://github.com/Yvetlin/smart-code-audit/actions"
	}
	msg.WriteString("🔗 GitHub Actions: " + runURL)
	return msg.String()
}

func shortPath(path string) string {
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		return path[idx+1:]
	}
	return path
}
