package suggestion

import (
	"sort"
	"strings"
	"sync"
)

// CommandSuggester provides intelligent command suggestions based on history
type CommandSuggester struct {
	commandFreq     map[string]int
	contextCommands map[string][]string
	mutex           sync.RWMutex
}

// New creates a new CommandSuggester instance
func New() *CommandSuggester {
	return &CommandSuggester{
		commandFreq:     make(map[string]int),
		contextCommands: make(map[string][]string),
	}
}

// AddCommand records a command and updates frequency statistics
func (cs *CommandSuggester) AddCommand(command string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.commandFreq[command]++
}

// AddContextualCommand records relationships between commands
func (cs *CommandSuggester) AddContextualCommand(prevCommand, newCommand string) {
	if prevCommand == "" {
		return
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	commands := cs.contextCommands[prevCommand]
	for _, cmd := range commands {
		if cmd == newCommand {
			return
		}
	}
	cs.contextCommands[prevCommand] = append(cs.contextCommands[prevCommand], newCommand)
}

// GetSuggestions returns command suggestions based on history and context
func (cs *CommandSuggester) GetSuggestions(currentInput string, lastCommand string) []string {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	suggestions := make(map[string]int)

	// Add contextual suggestions based on last command
	if lastCommand != "" {
		for _, cmd := range cs.contextCommands[lastCommand] {
			suggestions[cmd] += 2 // Give more weight to contextual suggestions
		}
	}

	// Add frequency-based suggestions
	for cmd, freq := range cs.commandFreq {
		if strings.Contains(strings.ToLower(cmd), strings.ToLower(currentInput)) {
			suggestions[cmd] += freq
		}
	}

	// Sort suggestions by score
	var result []string
	for cmd := range suggestions {
		result = append(result, cmd)
	}

	sort.Slice(result, func(i, j int) bool {
		return suggestions[result[i]] > suggestions[result[j]]
	})

	// Return top 5 suggestions
	if len(result) > 5 {
		result = result[:5]
	}

	return result
}
