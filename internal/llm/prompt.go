package llm

import (
	"fmt"

	"smart-code-audit/internal/models"
)

const SystemPrompt = `Ты эксперт по безопасности кода в CI/CD.

Статический анализатор сообщил о потенциальной проблеме. Твоя задача:
1. Перепроверить finding: это реальная проблема или ложное срабатывание?
2. Кратко (1-2 предложения) объяснить, что именно не так с указанной строкой кода.
3. Дать конкретную рекомендацию по исправлению — что заменить или как переписать.

Пиши по-русски. Отвечай ТОЛЬКО валидным JSON без markdown и комментариев.

Схема ответа:
{
  "is_real_issue": true,
  "probability": 0.85,
  "summary": "кратко что не так со строкой",
  "fix": "конкретная рекомендация по исправлению",
  "impact": "чем опасно, если не исправить"
}`

func BuildUserPrompt(f models.Finding) string {
	snippet := f.Snippet
	if snippet == "" {
		snippet = "(фрагмент кода недоступен)"
	}

	return fmt.Sprintf(`Сообщение сканера (%s / %s):
%s

Файл: %s
Строка: %d
Severity: %s

Контекст кода:
%s`,
		f.Tool,
		f.RuleID,
		f.Message,
		f.File,
		f.Line,
		f.Severity,
		snippet,
	)
}
