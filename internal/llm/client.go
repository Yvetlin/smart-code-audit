package llm

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gotoailab/llmhub"

	"smart-code-audit/internal/config"
	"smart-code-audit/internal/models"
)

type Client struct {
	client *llmhub.Client
	model  string
}

func New(cfg config.LLMConfig) (*Client, error) {
	apiKey := os.Getenv(cfg.APIKeyEnv)

	clientCfg := llmhub.ClientConfig{
		Model: cfg.Model,
	}

	switch cfg.Provider {
	case "openai":
		clientCfg.Provider = llmhub.ProviderOpenAI
		clientCfg.APIKey = apiKey

	case "ollama":
		clientCfg.Provider = llmhub.ProviderOllama

	case "openai_compat":
		clientCfg.Provider = llmhub.ProviderOpenAI
		clientCfg.APIKey = apiKey

		if baseURL := os.Getenv(cfg.BaseURLEnv); baseURL != "" {
			clientCfg.BaseURL = baseURL
		}

	default:
		return nil, errors.New("unknown LLM provider")
	}

	c, err := llmhub.NewClient(clientCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: c,
		model:  cfg.Model,
	}, nil
}

func (c *Client) EnrichFindings(
	ctx context.Context,
	findings []models.Finding,
) []models.EnrichedFinding {
	result := make([]models.EnrichedFinding, 0, len(findings))

	for i, f := range findings {
		log.Printf("AI review [%d/%d]: %s:%d (%s)", i+1, len(findings), f.File, f.Line, f.RuleID)

		review, err := c.reviewFinding(ctx, f)
		if err != nil {
			result = append(result, models.EnrichedFinding{
				Finding: f,
				Review:  fallbackReview(f),
			})
			continue
		}

		result = append(result, models.EnrichedFinding{
			Finding: f,
			Review:  *review,
		})
	}

	return result
}

func fallbackReview(f models.Finding) models.AIReview {
	return models.AIReview{
		IsRealIssue: true,
		Probability: 0.7,
		Summary:     f.Message,
		Impact:      "Требует ручной проверки — LLM не смог перепроверить автоматически",
		Fix:         "Проверьте указанную строку и устраните проблему согласно рекомендации сканера",
		Explanation: "Fallback: ответ LLM не распознан",
	}
}

func (c *Client) reviewFinding(
	ctx context.Context,
	f models.Finding,
) (*models.AIReview, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := c.client.ChatCompletions(
		ctx,
		llmhub.ChatCompletionRequest{
			Model: c.model,
			Messages: []llmhub.ChatMessage{
				{Role: "system", Content: SystemPrompt},
				{Role: "user", Content: BuildUserPrompt(f)},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("empty LLM response")
	}

	content, ok := resp.Choices[0].Message.Content.(string)
	if !ok {
		return nil, errors.New("LLM response content is not a string")
	}

	log.Printf("RAW LLM RESPONSE:\n%s\n", content)

	if idx := strings.Index(content, "{"); idx != -1 {
		content = content[idx:]
	}
	if idx := strings.LastIndex(content, "}"); idx != -1 {
		content = content[:idx+1]
	}

	var review models.AIReview
	if err := json.Unmarshal([]byte(content), &review); err != nil {
		log.Printf("LLM JSON parse error: %v\n", err)
		return nil, err
	}

	if review.Summary == "" {
		review.Summary = review.Impact
	}
	if review.Fix == "" {
		review.Fix = review.Explanation
	}

	return &review, nil
}
