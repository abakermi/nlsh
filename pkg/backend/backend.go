package backend

import (
	"context"
	"fmt"
	"strings"
	"github.com/sashabaranov/go-openai"
	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/session"
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
	return &OpenAIBackend{
		client:    openai.NewClient(apiKey),
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