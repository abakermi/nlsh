package assistant

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/abakermi/nlsh/pkg/backend"
	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/safety"
	"github.com/abakermi/nlsh/pkg/session"
	"github.com/abakermi/nlsh/pkg/color"
)

// Add method to get command
func (sa *ShellAssistant) GetCommand(input string) (string, error) {
	return sa.Backend.GenerateCommand(input, *sa.Session)
}

type ShellAssistant struct {
	Backend  backend.LLMBackend
	Checker  *safety.Checker
	Config   *config.Config
	Session  *session.Context
}

func New(backend backend.LLMBackend, cfg *config.Config, checker *safety.Checker) *ShellAssistant {
	return &ShellAssistant{
		Backend: backend,
		Config:  cfg,
		Checker: checker,
		Session: &session.Context{},
	}
}

func (sa *ShellAssistant) ExecuteCommand(command string) error {
	if !sa.Checker.IsAllowed(command) {
		return fmt.Errorf("%scommand blocked by safety rules%s", color.Red, color.Reset)
	}

	if sa.Config.Safety.ConfirmExecution && !confirmExecution(command) {
		fmt.Println("Execution cancelled")
		return nil
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err == nil {
		sa.Session.AddCommand(command)
	}

	return err
}

func confirmExecution(command string) bool {
	fmt.Printf("%sCommand: %s%s\nExecute? [y/N]: ", color.Yellow, command, color.Reset)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(scanner.Text()) == "y"
}