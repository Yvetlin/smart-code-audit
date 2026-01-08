package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"smart-code-audit/internal/analyzers"
	"smart-code-audit/internal/config"
	"smart-code-audit/internal/llm"
	"smart-code-audit/internal/normalizer"
	"smart-code-audit/internal/policy"
	"smart-code-audit/internal/report/markdown"
	"smart-code-audit/internal/report/sarif"
)

func main() {
	start := time.Now()

	log.Println("=== Smart Code Audit started ===")

	// ===== FLAGS =====
	configPath := flag.String("config", "configs/config.yaml", "path to config")
	targetPath := flag.String("target", ".", "path to target project")
	flag.Parse()

	// ===== CONFIG =====
	log.Println("Loading configuration...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// ===== STATIC ANALYSIS =====
	log.Println("Running static analyzers...")
	raw := analyzers.RunAll(*targetPath)
	log.Printf("Static analyzers finished, raw findings: %d", len(raw))

	// ===== NORMALIZATION =====
	log.Println("Normalizing findings...")
	findings := normalizer.Normalize(raw)
	log.Printf("Normalized findings: %d", len(findings))

	// ===== LLM =====
	log.Println("Initializing LLM client...")
	llmClient, err := llm.New(cfg.LLM)
	if err != nil {
		log.Fatalf("Failed to init LLM client: %v", err)
	}

	log.Println("Running AI-based analysis...")
	enriched := llmClient.EnrichFindings(context.Background(), findings)
	log.Printf("AI analysis completed, enriched findings: %d", len(enriched))

	// ===== POLICY =====
	log.Println("Evaluating security policy...")
	critical := policy.Evaluate(enriched)

	// ===== REPORTS =====
	log.Println("Generating reports...")
	_ = markdown.Write("results.md", enriched)
	_ = sarif.Write("results.sarif", enriched)

	log.Printf("Audit finished in %s", time.Since(start))

	if critical {
		log.Println("=== Audit завершен: ОБНАРУЖЕНЫ КРИТИЧЕСКИЕ ПРОБЛЕМЫ ===")
		os.Exit(1)
	}

	log.Println("=== Audit завершен: проблем не обнаружено ===")
}
