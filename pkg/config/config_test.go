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
model = "gpt-4-turbo-preview"
temperature = 0.7

[safety]
confirm_execution = true
allowed_commands = [
    "ls *",
    "touch *",
    "mkdir *",
    "git *"
]
denied_commands = [
    "rm -rf /*",
    "dd if=/dev/*",
    "mkfs.*"
]`

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

	// Test OpenAI configuration
	if cfg.OpenAI.Model != "gpt-4-turbo-preview" {
		t.Errorf("Expected model 'gpt-4-turbo-preview', got %s", cfg.OpenAI.Model)
	}
	if cfg.OpenAI.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", cfg.OpenAI.Temperature)
	}

	// Test Safety configuration
	if !cfg.Safety.ConfirmExecution {
		t.Error("Expected confirm_execution to be true")
	}

	// Test allowed commands
	expectedAllowed := []string{"ls *", "touch *", "mkdir *", "git *"}
	if len(cfg.Safety.AllowedCommands) != len(expectedAllowed) {
		t.Errorf("Expected %d allowed commands, got %d", len(expectedAllowed), len(cfg.Safety.AllowedCommands))
	}
	for i, cmd := range expectedAllowed {
		if cfg.Safety.AllowedCommands[i] != cmd {
			t.Errorf("Expected allowed command '%s', got '%s'", cmd, cfg.Safety.AllowedCommands[i])
		}
	}

	// Test denied commands
	expectedDenied := []string{"rm -rf /*", "dd if=/dev/*", "mkfs.*"}
	if len(cfg.Safety.DeniedCommands) != len(expectedDenied) {
		t.Errorf("Expected %d denied commands, got %d", len(expectedDenied), len(cfg.Safety.DeniedCommands))
	}
	for i, cmd := range expectedDenied {
		if cfg.Safety.DeniedCommands[i] != cmd {
			t.Errorf("Expected denied command '%s', got '%s'", cmd, cfg.Safety.DeniedCommands[i])
		}
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test default values
	if cfg.OpenAI.Model == "" {
		t.Error("Expected default model to be set")
	}
	if cfg.OpenAI.Temperature == 0 {
		t.Error("Expected default temperature to be set")
	}
	if len(cfg.Safety.AllowedCommands) == 0 {
		t.Error("Expected default allowed commands to be set")
	}
	if len(cfg.Safety.DeniedCommands) == 0 {
		t.Error("Expected default denied commands to be set")
	}
}