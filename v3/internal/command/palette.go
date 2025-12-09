package command

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// CommandDef defines a command with its metadata.
type CommandDef struct {
	Name        string
	Aliases     []string
	Args        string
	Description string
}

// GetAllCommands returns all available commands, optionally filtered by mode.
// If isRemoteMode is true, only show commands available in remote mode.
// If isRemoteMode is false, only show commands available in local mode.
func GetAllCommands(isRemoteMode bool) []CommandDef {
	baseCommands := []CommandDef{
		{
			Name:        "g",
			Aliases:     []string{"goto"},
			Args:        "<line>",
			Description: "Go to line number",
		},
		{
			Name:        "t",
			Aliases:     []string{"top"},
			Args:        "",
			Description: "Jump to top (first row)",
		},
		{
			Name:        "e",
			Aliases:     []string{"end"},
			Args:        "",
			Description: "Jump to end (last row)",
		},
		{
			Name:        "q",
			Aliases:     []string{"quit"},
			Args:        "",
			Description: "Quit application",
		},
		{
			Name:        "sort",
			Aliases:     []string{},
			Args:        "<col> <asc|desc>",
			Description: "Sort by column",
		},
		{
			Name:        "preset",
			Aliases:     []string{},
			Args:        "<name>",
			Description: "Switch to preset view",
		},
		{
			Name:        "filter",
			Aliases:     []string{},
			Args:        "<term>",
			Description: "Apply filter",
		},
		{
			Name:        "clear",
			Aliases:     []string{},
			Args:        "",
			Description: "Clear current filter",
		},
		{
			Name:        "help",
			Aliases:     []string{"?"},
			Args:        "",
			Description: "Show help screen",
		},
	}

	if isRemoteMode {
		baseCommands = append(baseCommands, CommandDef{
			Name:        "i",
			Aliases:     []string{"install"},
			Args:        "",
			Description: "Install selected package",
		})
	} else {
		baseCommands = append(baseCommands, CommandDef{
			Name:        "r",
			Aliases:     []string{"remove"},
			Args:        "",
			Description: "Remove selected package",
		})
	}

	return baseCommands
}

// FilterCommands filters commands based on the current buffer input.
func FilterCommands(buffer string, commands []CommandDef) []CommandDef {
	if buffer == "" {
		return commands
	}

	lowerBuffer := strings.ToLower(buffer)
	filtered := make([]CommandDef, 0)

	for _, cmd := range commands {
		// Match against command name
		if strings.HasPrefix(strings.ToLower(cmd.Name), lowerBuffer) {
			filtered = append(filtered, cmd)
			continue
		}

		// Match against aliases
		for _, alias := range cmd.Aliases {
			if strings.HasPrefix(strings.ToLower(alias), lowerBuffer) {
				filtered = append(filtered, cmd)
				break
			}
		}
	}

	return filtered
}

// RenderCommandPalette renders the command palette as table-style rows.
// Returns the rendered palette and the number of rows it occupies.
func RenderCommandPalette(buffer string, width int, isRemoteMode bool) (string, int) {
	commands := GetAllCommands(isRemoteMode)
	filtered := FilterCommands(buffer, commands)

	if len(filtered) == 0 {
		return "", 0
	}

	// Limit to first 6 commands to avoid taking too much space
	maxDisplay := 6
	if len(filtered) > maxDisplay {
		filtered = filtered[:maxDisplay]
	}

	// Row style matching table rows
	rowStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Foreground).
		Background(styles.Current.Selected).
		Width(width)

	commandStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Accent1).
		Bold(true)

	argsStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Accent4).
		Italic(true)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Dimmed)

	// Build command list
	var lines []string
	for _, cmd := range filtered {
		cmdText := commandStyle.Render(":" + cmd.Name)

		var argsText string
		if cmd.Args != "" {
			argsText = " " + argsStyle.Render(cmd.Args)
		}

		descText := " - " + descStyle.Render(cmd.Description)

		content := "  " + cmdText + argsText + descText
		line := rowStyle.Render(content)
		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...), len(lines)
}

// RenderOutputPalette renders the output palette showing command output.
// Returns the rendered palette and the number of rows it occupies.
func RenderOutputPalette(output string, width int) (string, int) {
	if output == "" {
		return "", 0
	}

	// Split output into lines
	outputLines := strings.Split(output, "\n")

	// Limit to first 10 lines to avoid taking too much space
	maxDisplay := 10
	if len(outputLines) > maxDisplay {
		outputLines = outputLines[:maxDisplay]
	}

	// Row style matching command palette
	rowStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Foreground).
		Background(styles.Current.Selected).
		Width(width)

	// Build output list
	var lines []string
	for _, line := range outputLines {
		if line == "" {
			continue
		}
		content := "  " + line
		lines = append(lines, rowStyle.Render(content))
	}

	if len(lines) == 0 {
		return "", 0
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...), len(lines)
}
