package suggestion

import (
	"sync"
	"testing"
)

func TestAddCommand(t *testing.T) {
	cs := New()
	cs.AddCommand("ls -l")
	cs.AddCommand("ls -l")

	if cs.commandFreq["ls -l"] != 2 {
		t.Errorf("Expected frequency 2 for 'ls -l', got %d", cs.commandFreq["ls -l"])
	}
}

func TestAddContextualCommand(t *testing.T) {
	cs := New()
	cs.AddContextualCommand("cd dir", "ls -l")

	commands := cs.contextCommands["cd dir"]
	if len(commands) != 1 || commands[0] != "ls -l" {
		t.Errorf("Expected ['ls -l'] for context 'cd dir', got %v", commands)
	}

	// Test duplicate prevention
	cs.AddContextualCommand("cd dir", "ls -l")
	if len(cs.contextCommands["cd dir"]) != 1 {
		t.Error("Duplicate command was added to context")
	}
}

func TestGetSuggestions(t *testing.T) {
	cs := New()

	// Add some command history
	cs.AddCommand("ls -l")
	cs.AddCommand("ls -l")
	cs.AddCommand("ls -la")
	cs.AddCommand("cd /tmp")

	// Add contextual relationship
	cs.AddContextualCommand("cd /tmp", "ls -l")

	// Test frequency-based suggestions
	suggestions := cs.GetSuggestions("ls", "")
	if len(suggestions) == 0 {
		t.Error("Expected suggestions for 'ls', got none")
	}
	if suggestions[0] != "ls -l" {
		t.Errorf("Expected most frequent command 'ls -l', got %s", suggestions[0])
	}

	// Test contextual suggestions
	suggestions = cs.GetSuggestions("", "cd /tmp")
	if len(suggestions) == 0 || suggestions[0] != "ls -l" {
		t.Error("Expected contextual suggestion 'ls -l' after 'cd /tmp'")
	}
}

func TestConcurrentAccess(t *testing.T) {
	cs := New()
	var wg sync.WaitGroup

	// Test concurrent command additions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.AddCommand("ls -l")
		}()
	}

	// Test concurrent contextual additions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.AddContextualCommand("cd dir", "ls -l")
		}()
	}

	// Test concurrent suggestions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.GetSuggestions("ls", "cd dir")
		}()
	}

	wg.Wait()

	// Verify final state
	if cs.commandFreq["ls -l"] != 100 {
		t.Errorf("Expected frequency 100 for 'ls -l', got %d", cs.commandFreq["ls -l"])
	}
	if len(cs.contextCommands["cd dir"]) != 1 {
		t.Error("Concurrent contextual additions resulted in duplicate entries")
	}
}
