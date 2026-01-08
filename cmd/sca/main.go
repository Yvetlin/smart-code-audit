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
	startTime := time.Now()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("=== Smart Code Audit started ===")

	cfgPath := flag.String("config", "configs/config.yaml", "path to config")
	flag.Parse()

	log.Println("Loading configuration...")
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Running static analyzers...")
	rawFindings := analyzers.RunAll(cfg)
	log.Printf("Static analyzers finished, raw findings: %d\n", len(rawFindings))

	log.Println("Normalizing findings...")
	findings := normalizer.Normalize(rawFindings)
	log.Printf("Normalized findings: %d\n", len(findings))

	log.Println("Initializing LLM client...")
	llmClient, err := llm.New(cfg.LLM)
	if err != nil {
		log.Fatalf("Failed to init LLM client: %v", err)
	}

	ctx := context.Background()

	log.Println("Running AI-based analysis...")
	enriched := llmClient.EnrichFindings(ctx, findings)
	log.Printf("AI analysis completed, enriched findings: %d\n", len(enriched))

	log.Println("Evaluating security policy...")
	decision := policy.Evaluate(cfg.Policy, enriched)

	log.Println("Generating reports...")
	if err := sarif.Write(enriched, "results.sarif"); err != nil {
		log.Fatalf("Failed to write SARIF report: %v", err)
	}
	if err := markdown.Write(enriched, "results.md"); err != nil {
		log.Fatalf("Failed to write Markdown report: %v", err)
	}

	duration := time.Since(startTime)
	log.Printf("Audit finished in %s\n", duration)

	if decision.FailPipeline {
		log.Println("=== Audit завершен: ОБНАРУЖЕНЫ КРИТИЧЕСКИЕ ПРОБЛЕМЫ ===")
		os.Exit(1)
	}

	log.Println("=== Audit завершен: проблем не обнаружено ===")
}
