package main

import (
	"os"
	"testing"
	"github.com/abakermi/nlsh/pkg/assistant"
	"github.com/abakermi/nlsh/pkg/backend"
	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/safety"
)

func TestIntegration(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping integration test: OPENAI_API_KEY not set")
	}

	cfg := &config.Config{}
	cfg.OpenAI.Model = "gpt-4o-2024-08-06"
	cfg.OpenAI.Temperature = 0.7

	systemCtx := "You are a shell assistant"
	llmBackend := backend.NewOpenAIBackend(os.Getenv("OPENAI_API_KEY"), cfg, systemCtx)
	
	safetyChecker := safety.NewChecker(
		[]string{"ls *"},
		[]string{"rm *"},
	)

	shellAssistant := assistant.New(llmBackend, cfg, safetyChecker)

	command, err := shellAssistant.GetCommand("list all files")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if command == "" {
		t.Error("Expected non-empty command")
	}
}