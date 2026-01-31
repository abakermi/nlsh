package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Backend string `toml:"backend"` // "openai" or "gemini"
	OpenAI  struct {
		Model       string  `toml:"model"`
		Temperature float64 `toml:"temperature"`
		BaseURL     string  `toml:"base_url"`
	} `toml:"openai"`
	Gemini struct {
		Model       string  `toml:"model"`
		Temperature float64 `toml:"temperature"`
	} `toml:"gemini"`
	Safety struct {
		ConfirmExecution bool     `toml:"confirm_execution"`
		AllowedCommands  []string `toml:"allowed_commands"`
		DeniedCommands   []string `toml:"denied_commands"`
	} `toml:"safety"`
}

func Load() (*Config, error) {
	config := &Config{}
	configPath := filepath.Join(os.Getenv("HOME"), ".nlshrc")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := `# Backend to use: "openai" or "gemini"
backend = "openai"

[openai]
model = "gpt-4-turbo-preview"
temperature = 0.7
# base_url = "http://localhost:11434/v1" # Uncomment for local models (Ollama, etc)

[gemini]
model = "gemini-2.0-flash"
temperature = 0.7

[safety]
confirm_execution = true
allowed_commands = [
    "ls *",
    "touch *",
    "mkdir *",
    "echo *",
    "cat *",
    "cp *",
    "mv *",
    "git *",
    "docker *",
    "code *",
    "vim *",
    "nano *"
]
denied_commands = [
    "rm -rf /*",
    "rm -rf /",
    "dd if=/dev/*",
    "mkfs.*",
    "> /dev/*",
    "shutdown *",
    "reboot *",
    "halt *",
    "*--no-preserve-root*"
]`

		err := os.WriteFile(configPath, []byte(defaultConfig), 0644)
		if err != nil {
			return nil, fmt.Errorf("error creating default config: %v", err)
		}
	}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	setDefaults(config)
	return config, nil
}

func setDefaults(config *Config) {
	if config.Backend == "" {
		config.Backend = "openai"
	}
	if config.OpenAI.Model == "" {
		config.OpenAI.Model = "gpt-4o-2024-08-06"
	}
	if config.OpenAI.Temperature == 0 {
		config.OpenAI.Temperature = 0.7
	}
	if config.Gemini.Model == "" {
		config.Gemini.Model = "gemini-2.0-flash"
	}
	if config.Gemini.Temperature == 0 {
		config.Gemini.Temperature = 0.7
	}
}