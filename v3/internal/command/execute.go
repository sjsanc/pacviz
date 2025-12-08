package command

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ExecuteResult represents the result of executing a command.
type ExecuteResult struct {
	Quit         bool   // Whether to quit the application
	GoToLine     int    // Line number to jump to (-1 = no jump)
	ScrollTop    bool   // Whether to scroll to top
	ScrollEnd    bool   // Whether to scroll to end
	Error        string // Error message if command failed
	PresetChange string // Preset to switch to (empty = no change)
}

// Execute parses and executes a command string.
func Execute(commandStr string) ExecuteResult {
	commandStr = strings.TrimSpace(commandStr)
	if commandStr == "" {
		return ExecuteResult{GoToLine: -1}
	}

	// Special handling for :g command - allow :g22 syntax
	if strings.HasPrefix(commandStr, "g") && len(commandStr) > 1 {
		// Check if character after 'g' is a digit
		if commandStr[1] >= '0' && commandStr[1] <= '9' {
			// Extract the number part
			return executeGoTo([]string{commandStr[1:]})
		}
	}

	parts := strings.Fields(commandStr)
	if len(parts) == 0 {
		return ExecuteResult{GoToLine: -1}
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "g", "goto":
		return executeGoTo(args)
	case "t", "top":
		return ExecuteResult{ScrollTop: true, GoToLine: -1}
	case "e", "end":
		return ExecuteResult{ScrollEnd: true, GoToLine: -1}
	case "q", "quit":
		return ExecuteResult{Quit: true, GoToLine: -1}
	case "p", "preset":
		return executePreset(args)
	default:
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Unknown command: " + command,
		}
	}
}

// executeGoTo handles the :g <line> command.
func executeGoTo(args []string) ExecuteResult {
	if len(args) == 0 {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Usage: :g <line>",
		}
	}

	lineStr := args[0]
	line, err := strconv.Atoi(lineStr)
	if err != nil {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Invalid line number: " + lineStr,
		}
	}

	// Convert to 0-based index (user sees 1-based line numbers)
	return ExecuteResult{
		GoToLine: line - 1,
	}
}

// executePreset handles the :p <preset> command.
func executePreset(args []string) ExecuteResult {
	if len(args) == 0 {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Usage: :p <preset> (explicit, dependency, orphans, foreign, all)",
		}
	}

	preset := args[0]
	// Validate preset name
	validPresets := map[string]bool{
		"explicit":   true,
		"dependency": true,
		"orphans":    true,
		"foreign":    true,
		"all":        true,
	}

	if !validPresets[preset] {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Invalid preset: " + preset + " (valid: explicit, dependency, orphans, foreign, all)",
		}
	}

	return ExecuteResult{
		GoToLine:     -1,
		PresetChange: preset,
	}
}

// CommandResultMsg is sent when a command needs to affect the model.
type CommandResultMsg struct {
	Result ExecuteResult
}

// ExecuteCommandMsg creates a tea.Cmd that executes a command.
func ExecuteCommandMsg(commandStr string) tea.Cmd {
	return func() tea.Msg {
		result := Execute(commandStr)
		return CommandResultMsg{Result: result}
	}
}
