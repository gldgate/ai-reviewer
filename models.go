package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	googleoption "google.golang.org/api/option"
)

var (
	envCache map[string]string
	envOnce  sync.Once
)

func getEnv(key string) string {
	envOnce.Do(func() {
		envCache = make(map[string]string)
		data, err := os.ReadFile("KEYS.env")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					envCache[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}
	})

	if val, ok := envCache[key]; ok {
		return val
	}
	return os.Getenv(key)
}

type ModelClient interface {
	Generate(ctx context.Context, prompt string, maxTokens int) (ModelResult, error)
	GenerateJSON(ctx context.Context, prompt string, maxTokens int) (ModelResult, error)
}

type ModelResult struct {
	Text      string
	TokensIn  int
	TokensOut int
	Provider  string
	Model     string
}

type ModelCategory string

const (
	FastestGood  ModelCategory = "fastest_good"
	Balanced     ModelCategory = "balanced"
	BestCode     ModelCategory = "best_code"
	FrontierBest ModelCategory = "frontier_best"
)

// OpenAI Client
type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (c *OpenAIClient) Generate(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	return c.generate(ctx, prompt, maxTokens, false)
}

func (c *OpenAIClient) GenerateJSON(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	return c.generate(ctx, prompt, maxTokens, true)
}

func (c *OpenAIClient) generate(ctx context.Context, prompt string, maxTokens int, jsonMode bool) (ModelResult, error) {
	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	if jsonMode {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}
	if maxTokens > 0 {
		req.MaxTokens = maxTokens
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return ModelResult{}, err
	}

	return ModelResult{
		Text:      resp.Choices[0].Message.Content,
		TokensIn:  resp.Usage.PromptTokens,
		TokensOut: resp.Usage.CompletionTokens,
		Provider:  "openai",
		Model:     c.model,
	}, nil
}

// Anthropic Client
type AnthropicClient struct {
	client *anthropic.Client
	model  string
}

func NewAnthropicClient(apiKey, model string) *AnthropicClient {
	c := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicClient{
		client: &c,
		model:  model,
	}
}

func (c *AnthropicClient) Generate(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	return c.generate(ctx, prompt, maxTokens, false)
}

func (c *AnthropicClient) GenerateJSON(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	// Anthropic doesn't have a simple "JSON mode" flag in the same way OpenAI does,
	// but we can ask for it in the prompt or use tool use.
	// For now, we'll just append a JSON instruction if it's not already there and use regular Generate.
	// Actually, newer Anthropic models support structured output via tools, but for simplicity
	// here we will just rely on the system prompt for now, OR we could implement tool use.
	// The prompt already asks for JSON.
	return c.generate(ctx, prompt, maxTokens, true)
}

func (c *AnthropicClient) generate(ctx context.Context, prompt string, maxTokens int, jsonMode bool) (ModelResult, error) {
	params := anthropic.MessageNewParams{
		Model: anthropic.Model(c.model),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	}
	if maxTokens > 0 {
		params.MaxTokens = int64(maxTokens)
	} else {
		params.MaxTokens = 4096
	}

	message, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return ModelResult{}, err
	}

	return ModelResult{
		Text:      message.Content[0].Text,
		TokensIn:  int(message.Usage.InputTokens),
		TokensOut: int(message.Usage.OutputTokens),
		Provider:  "anthropic",
		Model:     c.model,
	}, nil
}

// Gemini Client
type GeminiClient struct {
	client *genai.Client
	model  string
}

func NewGeminiClient(ctx context.Context, apiKey, model string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, googleoption.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (c *GeminiClient) Generate(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	return c.generate(ctx, prompt, maxTokens, false)
}

func (c *GeminiClient) GenerateJSON(ctx context.Context, prompt string, maxTokens int) (ModelResult, error) {
	return c.generate(ctx, prompt, maxTokens, true)
}

func (c *GeminiClient) generate(ctx context.Context, prompt string, maxTokens int, jsonMode bool) (ModelResult, error) {
	model := c.client.GenerativeModel(c.model)
	if jsonMode {
		model.ResponseMIMEType = "application/json"
	}
	if maxTokens > 0 {
		model.MaxOutputTokens = genai.Ptr(int32(maxTokens))
	}
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return ModelResult{}, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return ModelResult{}, fmt.Errorf("empty response from Gemini")
	}

	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			text += string(t)
		}
	}

	tokensIn := 0
	tokensOut := 0
	if resp.UsageMetadata != nil {
		tokensIn = int(resp.UsageMetadata.PromptTokenCount)
		tokensOut = int(resp.UsageMetadata.CandidatesTokenCount)
	}

	return ModelResult{
		Text:      text,
		TokensIn:  tokensIn,
		TokensOut: tokensOut,
		Provider:  "gemini",
		Model:     c.model,
	}, nil
}

func GetModelClient(ctx context.Context, provider, model string) (ModelClient, error) {
	switch provider {
	case "openai":
		apiKey := getEnv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return NewOpenAIClient(apiKey, model), nil
	case "anthropic":
		apiKey := getEnv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return NewAnthropicClient(apiKey, model), nil
	case "gemini":
		apiKey := getEnv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY not set")
		}
		return NewGeminiClient(ctx, apiKey, model)
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}
