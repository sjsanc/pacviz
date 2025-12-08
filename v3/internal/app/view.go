package app

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/command"
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

	// Calculate the selected row index relative to the visible window
	relativeSelectedRow := m.Viewport.SelectedRow - m.Viewport.Offset

	// Generate status bar based on mode
	var statusBar string
	var commandPalette string
	var paletteRows int
	var tableUI string

	switch m.Mode {
	case ModeCommand:
		// Show command palette and buffer
		commandPalette, paletteRows = command.RenderCommandPalette(m.GetBufferContent(), width)
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModeFilter:
		// Show buffer only for filter mode
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModeNormal:
		// Show normal status line
		filterText := ""
		if m.Viewport.Filter.Active && len(m.Viewport.Filter.Terms) > 0 {
			filterText = m.Viewport.Filter.Terms[0]
		}
		presetName := m.Presets[m.CurrentPreset].Name
		statusBar = renderer.RenderStatus(
			presetName,
			len(m.Viewport.VisibleRows),
			m.Viewport.Height,
			m.Viewport.Offset,
			filterText,
			width,
		)
	}

	// Render the table with detail panel or palette overlay
	if m.ShowDetailPanel {
		// Show detail panel for selected package
		selectedPackage := m.Viewport.GetSelectedPackage()
		tableUI = renderer.RenderWithDetailPanel(
			width,
			m.Height,
			m.Viewport.Columns,
			colWidths,
			visibleRows,
			relativeSelectedRow,
			m.Viewport.SelectedCol,
			m.Viewport.SortColumn,
			m.Viewport.SortReverse,
			selectedPackage,
			m.Viewport.Offset,
		)
	} else if commandPalette != "" {
		// Show command palette overlay
		tableUI = renderer.RenderWithPaletteOverlay(
			width,
			m.Height,
			m.Viewport.Columns,
			colWidths,
			visibleRows,
			relativeSelectedRow,
			m.Viewport.SelectedCol,
			m.Viewport.SortColumn,
			m.Viewport.SortReverse,
			commandPalette,
			paletteRows,
			m.Viewport.Offset,
		)
	} else {
		// No overlay, render normally
		tableUI = renderer.RenderWithPaletteOverlay(
			width,
			m.Height,
			m.Viewport.Columns,
			colWidths,
			visibleRows,
			relativeSelectedRow,
			m.Viewport.SelectedCol,
			m.Viewport.SortColumn,
			m.Viewport.SortReverse,
			"",
			0,
			m.Viewport.Offset,
		)
	}

	// Add status bar at the bottom
	if statusBar != "" {
		return lipgloss.JoinVertical(lipgloss.Left, tableUI, statusBar)
	}

	return tableUI
}
