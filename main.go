package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/abakermi/nlsh/pkg/assistant"
	"github.com/abakermi/nlsh/pkg/backend"
	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/color"
	"github.com/abakermi/nlsh/pkg/safety"
)

const systemPromptTemplate = `You are a expert system shell assistant. Convert natural language requests into appropriate shell commands for %s.
Current Directory: %s
Shell: %s
OS: %s

Follow these rules:
1. Prefer standard GNU coreutils over platform-specific tools
2. Use pipes for complex operations
3. Avoid destructive commands (rm -rf, dd, etc)
4. Add comments for complex commands
5. Consider directory context

Provide ONLY the command without explanation. If unclear, respond with "UNCLEAR".`

func getSystemContext() string {
	dir, _ := os.Getwd()
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "unknown"
	}
	return fmt.Sprintf(
		"OS: %s\nShell: %s\nCurrent Directory: %s",
		runtime.GOOS,
		shell,
		dir,
	)
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	systemCtx := fmt.Sprintf(systemPromptTemplate, runtime.GOOS, getSystemContext(), os.Getenv("SHELL"), runtime.GOOS)
	
	llmBackend := backend.NewOpenAIBackend(apiKey, cfg, systemCtx)
	
	safetyChecker := safety.NewChecker(
		cfg.Safety.AllowedCommands,
		append([]string{
			"rm * -rf*",
			"dd *",
			"mkfs*",
			"*--no-preserve-root*",
		}, cfg.Safety.DeniedCommands...),
	)

	shellAssistant := assistant.New(llmBackend, cfg, safetyChecker)

	fmt.Printf("%s[System]%s Natural Language Shell initialized\n", color.Green, color.Reset)
	
	if len(os.Args) > 1 {
		handleSingleCommand(shellAssistant, os.Args[1:])
		return
	}

	runInteractiveMode(shellAssistant)
}

func handleSingleCommand(assistant *assistant.ShellAssistant, args []string) {
	input := strings.Join(args, " ")
	command, err := assistant.GetCommand(input)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", color.Red, err, color.Reset)
		os.Exit(1)
	}

	if err := assistant.ExecuteCommand(command); err != nil {
		fmt.Printf("%sError executing command: %v%s\n", color.Red, err, color.Reset)
		os.Exit(1)
	}
}

func runInteractiveMode(assistant *assistant.ShellAssistant) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Natural Language Shell (nlsh) - Type 'exit' to quit")
	fmt.Println("Enter your request in natural language:")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		command, err := assistant.GetCommand(input)
		if err != nil {
			fmt.Printf("%sError: %v%s\n", color.Red, err, color.Reset)
			continue
		}

		if err := assistant.ExecuteCommand(command); err != nil {
			fmt.Printf("%sError executing command: %v%s\n", color.Red, err, color.Reset)
		}
	}
}