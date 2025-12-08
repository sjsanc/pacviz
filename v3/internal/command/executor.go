package command

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Executor executes parsed commands.
type Executor struct {
	registry *Registry
}

// NewExecutor creates a new command executor.
func NewExecutor(registry *Registry) *Executor {
	return &Executor{
		registry: registry,
	}
}

// Execute executes a command string.
func (e *Executor) Execute(input string) tea.Cmd {
	name, args := Parse(input)
	cmd := e.registry.Get(name)

	if cmd == nil {
		// TODO: Return error message command
		return nil
	}

	// TODO: Validate arguments
	// TODO: Execute command
	return cmd.Execute(args)
}
