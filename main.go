package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	systemPrompt = `You are a helpful shell command assistant. Convert natural language requests into appropriate shell commands.
Provide ONLY the command without any explanation. If you're unsure or the request is unclear, respond with "UNCLEAR".
Ensure the command is safe and won't harm the system.`
)

type ShellAssistant struct {
	client *openai.Client
}

func NewShellAssistant(apiKey string) *ShellAssistant {
	return &ShellAssistant{
		client: openai.NewClient(apiKey),
	}
}

func (sa *ShellAssistant) GetCommand(input string) (string, error) {
	resp, err := sa.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: input,
				},
			},
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

func (sa *ShellAssistant) ExecuteCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	assistant := NewShellAssistant(apiKey)

	// Check if command-line argument is provided
	if len(os.Args) > 1 {
		// Join all arguments as a single instruction
		input := strings.Join(os.Args[1:], " ")
		command, err := assistant.GetCommand(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Executing: %s\n", command)
		if err := assistant.ExecuteCommand(command); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Interactive mode
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
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Executing: %s\n", command)
		if err := assistant.ExecuteCommand(command); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
