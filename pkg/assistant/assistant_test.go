package assistant

import (
	"testing"
	"github.com/abakermi/nlsh/pkg/config"
	"github.com/abakermi/nlsh/pkg/safety"
	"github.com/abakermi/nlsh/pkg/session"
)

type mockBackend struct{}

func (m *mockBackend) GenerateCommand(prompt string, ctx session.Context) (string, error) {
	return "ls -la", nil
}

func TestShellAssistant(t *testing.T) {
	cfg := &config.Config{}
	checker := safety.NewChecker([]string{"ls"}, []string{"rm"})
	assistant := New(&mockBackend{}, cfg, checker)

	command, err := assistant.GetCommand("list files")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if command != "ls -la" {
		t.Errorf("Expected 'ls -la', got %s", command)
	}
}