package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".nlshrc")
	
	configContent := `[openai]
model = "gpt-4o-2024-08-06"
temperature = 0.7

[safety]
confirm_execution = true
allowed_commands = ["ls", "git"]
denied_commands = ["rm"]`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Set HOME to temp dir for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.OpenAI.Model != "gpt-4o-2024-08-06" {
		t.Errorf("Expected model 'gpt-4', got %s", cfg.OpenAI.Model)
	}

	if !cfg.Safety.ConfirmExecution {
		t.Error("Expected confirm_execution to be true")
	}
}