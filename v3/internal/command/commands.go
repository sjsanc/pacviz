package command

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ArgType specifies the type of command argument.
type ArgType int

const (
	ArgString ArgType = iota
	ArgInt
	ArgEnum
)

// Arg defines a command argument.
type Arg struct {
	Name     string
	Required bool
	Type     ArgType
	Values   []string // For enum types
}

// Command represents a user command.
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Args        []Arg
	Execute     ExecuteFunc
}

// ExecuteFunc is the function signature for command execution.
type ExecuteFunc func(args []string) tea.Cmd

// DefaultCommands returns the standard command set.
func DefaultCommands() []*Command {
	return []*Command{
		{
			Name:        "quit",
			Aliases:     []string{"q"},
			Description: "Exit application",
			Args:        nil,
			Execute: func(args []string) tea.Cmd {
				return tea.Quit
			},
		},
		{
			Name:        "goto",
			Aliases:     []string{"g"},
			Description: "Go to line number",
			Args: []Arg{
				{Name: "line", Required: true, Type: ArgInt},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement goto logic
				return nil
			},
		},
		{
			Name:        "search",
			Aliases:     []string{"s"},
			Description: "Search sync databases",
			Args: []Arg{
				{Name: "query", Required: true, Type: ArgString},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement search logic
				return nil
			},
		},
		{
			Name:        "install",
			Aliases:     []string{"i"},
			Description: "Install package",
			Args: []Arg{
				{Name: "package", Required: false, Type: ArgString},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement install logic
				return nil
			},
		},
		{
			Name:        "remove",
			Aliases:     []string{"r"},
			Description: "Remove package",
			Args: []Arg{
				{Name: "package", Required: false, Type: ArgString},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement remove logic
				return nil
			},
		},
		{
			Name:        "sort",
			Aliases:     nil,
			Description: "Sort by column",
			Args: []Arg{
				{Name: "column", Required: true, Type: ArgString},
				{Name: "direction", Required: false, Type: ArgEnum, Values: []string{"asc", "desc"}},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement sort logic
				return nil
			},
		},
		{
			Name:        "filter",
			Aliases:     []string{"f"},
			Description: "Apply filter",
			Args: []Arg{
				{Name: "term", Required: true, Type: ArgString},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement filter logic
				return nil
			},
		},
		{
			Name:        "preset",
			Aliases:     []string{"p"},
			Description: "Switch to preset view",
			Args: []Arg{
				{Name: "preset", Required: true, Type: ArgEnum, Values: []string{"explicit", "dependency", "orphans", "foreign", "all"}},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement preset switch logic
				return nil
			},
		},
		{
			Name:        "info",
			Aliases:     nil,
			Description: "Show package details",
			Args: []Arg{
				{Name: "package", Required: false, Type: ArgString},
			},
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement info view logic
				return nil
			},
		},
		{
			Name:        "help",
			Aliases:     []string{"h", "?"},
			Description: "Show help",
			Args:        nil,
			Execute: func(args []string) tea.Cmd {
				// TODO: Implement help screen logic
				return nil
			},
		},
	}
}
