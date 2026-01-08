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

	var result []models.EnrichedFinding

	for _, f := range findings {
		review, err := c.reviewFinding(ctx, f)
		if err != nil {
			// 🔥 FALLBACK
			result = append(result, models.EnrichedFinding{
				Finding: f,
				Review: models.AIReview{
					IsRealIssue: true,
					Probability: 0.8,
					Impact:      "Potential security vulnerability detected by static analysis",
					Fix:         "Review the code and apply secure coding practices",
					Explanation: "Fallback: LLM response could not be parsed or returned an error",
				},
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

func (c *Client) reviewFinding(
	ctx context.Context,
	f models.Finding,
) (*models.AIReview, error) {

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	userPrompt := BuildUserPrompt(f)

	resp, err := c.client.ChatCompletions(
		ctx,
		llmhub.ChatCompletionRequest{
			Model: c.model,
			Messages: []llmhub.ChatMessage{
				{Role: "system", Content: SystemPrompt},
				{Role: "user", Content: userPrompt},
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

	// 🔒 защита от текста до JSON
	if idx := strings.Index(content, "{"); idx != -1 {
		content = content[idx:]
	}

	var review models.AIReview
	if err := json.Unmarshal([]byte(content), &review); err != nil {
		log.Printf("LLM JSON parse error: %v\n", err)
		return nil, err
	}

	return &review, nil
}
