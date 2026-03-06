package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
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

func TestPrintOnlyNoStartupOutput(t *testing.T) {
	// When --print-only is passed, the "[System]" startup line must NOT appear in stdout.
	cmd := exec.Command("./nlsh", "--print-only", "list files")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// We expect non-zero exit (no API key in test env), but stdout must be clean
	cmd.Run()
	if strings.Contains(stdout.String(), "[System]") {
		t.Errorf("--print-only stdout must not contain '[System]', got: %q", stdout.String())
	}
}