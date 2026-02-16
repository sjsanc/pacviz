package command

import (
	"strconv"
	"strings"
)

// ExecuteResult represents the result of executing a command.
type ExecuteResult struct {
	Quit           bool
	GoToLine       int
	ScrollTop      bool
	ScrollEnd      bool
	Error          string
	PresetChange   string
	RemoteSearch   string
	InstallPackage bool
	RemovePackage  bool
	ThemeName      string
}

// Execute parses and executes a command string.
func Execute(commandStr string) ExecuteResult {
	commandStr = strings.TrimSpace(commandStr)
	if commandStr == "" {
		return ExecuteResult{GoToLine: -1}
	}

	// Allow :g22 syntax (no space between command and number)
	if strings.HasPrefix(commandStr, "g") && len(commandStr) > 1 {
		if commandStr[1] >= '0' && commandStr[1] <= '9' {
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
	case "s", "search":
		return executeSearch(args)
	case "i", "install":
		return ExecuteResult{InstallPackage: true, GoToLine: -1}
	case "r", "remove":
		return ExecuteResult{RemovePackage: true, GoToLine: -1}
	case "theme", "th":
		return executeTheme(args)
	default:
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Unknown command: " + command,
		}
	}
}

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

	return ExecuteResult{
		GoToLine: line - 1,
	}
}

func executePreset(args []string) ExecuteResult {
	if len(args) == 0 {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Usage: :p <preset> (explicit, dependency, orphans, foreign, aur, updatable, all)",
		}
	}

	preset := args[0]
	validPresets := map[string]bool{
		"explicit":   true,
		"dependency": true,
		"orphans":    true,
		"foreign":    true,
		"aur":        true,
		"updatable":  true,
		"all":        true,
	}

	if !validPresets[preset] {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Invalid preset: " + preset + " (valid: explicit, dependency, orphans, foreign, aur, updatable, all)",
		}
	}

	return ExecuteResult{
		GoToLine:     -1,
		PresetChange: preset,
	}
}

func executeSearch(args []string) ExecuteResult {
	if len(args) == 0 {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Usage: :search <query>",
		}
	}

	query := strings.Join(args, " ")

	return ExecuteResult{
		GoToLine:     -1,
		RemoteSearch: query,
	}
}

func executeTheme(args []string) ExecuteResult {
	if len(args) == 0 {
		return ExecuteResult{
			GoToLine: -1,
			Error:    "Usage: :theme <name>",
		}
	}

	themeName := args[0]

	return ExecuteResult{
		GoToLine:  -1,
		ThemeName: themeName,
	}
}
