package analyzers

import (
	"smart-code-audit/internal/config"
	"smart-code-audit/internal/models"
)

type RawFinding = models.Finding

func RunAll(cfg *config.Config) []RawFinding {
	var all []RawFinding
	all = append(all, runGosec(cfg.Project.Target)...)
	return all
}
