package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/session"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/genai"
)

type LLMBackend interface {
	GenerateCommand(prompt string, context session.Context) (string, error)
}

type OpenAIBackend struct {
	client    *openai.Client
	config    *config.Config
	systemCtx string
}

func NewOpenAIBackend(apiKey string, cfg *config.Config, systemCtx string) *OpenAIBackend {
	clientConfig := openai.DefaultConfig(apiKey)
	if cfg.OpenAI.BaseURL != "" {
		clientConfig.BaseURL = cfg.OpenAI.BaseURL
	}
	
	return &OpenAIBackend{
		client:    openai.NewClientWithConfig(clientConfig),
		config:    cfg,
		systemCtx: systemCtx,
	}
}

func (b *OpenAIBackend) GenerateCommand(prompt string, ctx session.Context) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: b.systemCtx,
		},
	}

	if len(ctx.PreviousCommands) > 0 {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: ctx.GetContextPrompt(),
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := b.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       b.config.OpenAI.Model,
			Messages:    messages,
			Temperature: float32(b.config.OpenAI.Temperature),
		},
	)
	if err != nil {
		return "", fmt.Errorf("error getting completion: %v", err)
	}

	command := strings.TrimSpace(resp.Choices[0].Message.Content)
	if command == "UNCLEAR" {
		return "", fmt.Errorf("unclear request, please be more specific")
	}

	return command, nil
}


// GeminiBackend implements LLMBackend using Google's Gemini API
type GeminiBackend struct {
	client    *genai.Client
	config    *config.Config
	systemCtx string
}

func NewGeminiBackend(apiKey string, cfg *config.Config, systemCtx string) (*GeminiBackend, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	return &GeminiBackend{
		client:    client,
		config:    cfg,
		systemCtx: systemCtx,
	}, nil
}

func (b *GeminiBackend) GenerateCommand(prompt string, ctx session.Context) (string, error) {
	// Build the full prompt with system context
	fullPrompt := b.systemCtx + "\n\n"

	if len(ctx.PreviousCommands) > 0 {
		fullPrompt += ctx.GetContextPrompt() + "\n\n"
	}

	fullPrompt += "User request: " + prompt

	temperature := float32(b.config.Gemini.Temperature)
	result, err := b.client.Models.GenerateContent(
		context.Background(),
		b.config.Gemini.Model,
		genai.Text(fullPrompt),
		&genai.GenerateContentConfig{
			Temperature: &temperature,
		},
	)
	if err != nil {
		return "", fmt.Errorf("error getting Gemini completion: %v", err)
	}

	command := strings.TrimSpace(result.Text())
	if command == "UNCLEAR" {
		return "", fmt.Errorf("unclear request, please be more specific")
	}

	return command, nil
}
