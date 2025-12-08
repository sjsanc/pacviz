package renderer

import (
	"fmt"

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
