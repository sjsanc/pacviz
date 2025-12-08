package app

import (
	"fmt"

	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/renderer"
)

// View renders the application (Bubble Tea interface).
func (m Model) View() string {
	// Show error if initialization failed
	if m.Error != nil {
		return fmt.Sprintf("Error: %v\nPress Ctrl+C to quit", m.Error)
	}

	// Show loading message if not ready
	if !m.Ready {
		return "Loading packages..."
	}

	// Use actual width or default to 120 if not yet initialized
	width := m.Width
	if width == 0 {
		width = 120
	}

	// Calculate column widths
	colWidths := column.CalculateWidths(m.Viewport.Columns, width)

	// Get visible rows
	visibleRows := m.Viewport.GetVisibleRows()

	// Render the UI
	return renderer.Render(
		width,
		m.Height,
		m.Viewport.Columns,
		colWidths,
		visibleRows,
		m.Viewport.SelectedRow,
		m.Viewport.SelectedCol,
		m.Viewport.SortColumn,
		m.Viewport.SortReverse,
	)
}
