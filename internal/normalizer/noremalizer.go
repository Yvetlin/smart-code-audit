package normalizer

import (
	"path/filepath"

	"smart-code-audit/internal/contextloader"
	"smart-code-audit/internal/models"
)

func Normalize(target string, raw []models.Finding) []models.Finding {
	out := make([]models.Finding, 0, len(raw))
	for _, f := range raw {
		file := f.File
		if !filepath.IsAbs(file) {
			file = filepath.Join(target, file)
		}
		f.Snippet = contextloader.EnrichSnippet(file, f.Line, f.Snippet)
		out = append(out, f)
	}
	return out
}
