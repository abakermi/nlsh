package session

import (
	"fmt"
	"strings"
)

type Context struct {
	PreviousCommands []string
	CurrentDir       string
	SystemInfo       string
}

func (sc *Context) GetContextPrompt() string {
	return fmt.Sprintf(`Previous commands:
%s
Current directory: %s`, 
		strings.Join(sc.PreviousCommands, "\n"), 
		sc.CurrentDir)
}

func (sc *Context) AddCommand(command string) {
	sc.PreviousCommands = append(sc.PreviousCommands, command)
	if len(sc.PreviousCommands) > 10 {
		sc.PreviousCommands = sc.PreviousCommands[1:]
	}
}