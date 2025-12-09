package renderer

import (
	"fmt"
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// InputMode represents the application's input mode (for status display).
type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCommand
	ModeFilter
)

// RenderStatus renders the status/footer bar.
func RenderStatus(preset string, totalRows, visibleRows, offset int, filter string, width int) string {
	// Default status line for normal mode
	start := offset + 1
	end := min(offset+visibleRows, totalRows)
	status := fmt.Sprintf("Preset: %s | Showing %d-%d of %d",
		preset,
		start,
		end,
		totalRows)

	if filter != "" {
		status += fmt.Sprintf(" | Filter: \"%s\"", filter)
	}

	return styles.StatusBar.Width(width).Render(status)
}

// RenderStatusWithBuffer renders the status bar with an input buffer (for command/filter mode).
func RenderStatusWithBuffer(buffer string, width int) string {
	return styles.StatusBar.Width(width).Render(buffer)
}

// RenderRemoteStatus renders the status bar for remote mode.
func RenderRemoteStatus(query string, totalRows, visibleRows, offset int, filter string, loading bool, spinner string, errorMsg string, width int) string {
	var status string

	if loading {
		status = fmt.Sprintf("%s Searching: %s", spinner, query)
	} else {
		start := offset + 1
		end := min(offset+visibleRows, totalRows)
		status = fmt.Sprintf("SEARCH: %s | Showing %d-%d of %d",
			query,
			start,
			end,
			totalRows)

		if filter != "" {
			status += fmt.Sprintf(" | Filter: \"%s\"", filter)
		}
	}

	// If there's an error, show it right-aligned
	if errorMsg != "" {
		// Calculate available space for error
		statusLen := len(status)
		errorLen := len(errorMsg)
		padding := width - statusLen - errorLen - 4 // 4 for padding and separator

		if padding > 0 {
			status = status + strings.Repeat(" ", padding) + "| " + errorMsg
		} else {
			// If not enough space, just show error
			status = errorMsg
		}
	}

	return styles.RemoteStatusBar.Width(width).Render(status)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
