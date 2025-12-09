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
	if m.Error != "" {
		return fmt.Sprintf("Error: %s\nPress Ctrl+C to quit", m.Error)
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
	var outputPalette string
	var paletteRows int
	var tableUI string
	isRemoteMode := m.ViewMode == ViewRemote

	switch m.Mode {
	case ModeCommand:
		// Show command palette and buffer
		commandPalette, paletteRows = command.RenderCommandPalette(m.GetBufferContent(), width, isRemoteMode)
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModeFilter:
		// Show buffer only for filter mode
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModePassword:
		// Show password prompt with masked input
		prompt := "[sudo] password: "
		masked := ""
		for range m.PasswordBuffer {
			masked += "*"
		}
		statusBar = renderer.RenderWarningStatus(prompt+masked, width)
	case ModeNormal:
		// Show output palette if there's removal or installation output
		if m.RemoveOutput != "" {
			var rows int
			outputPalette, rows = command.RenderOutputPalette(m.RemoveOutput, width)
			paletteRows = rows
		} else if m.InstallOutput != "" {
			var rows int
			outputPalette, rows = command.RenderOutputPalette(m.InstallOutput, width)
			paletteRows = rows
		}

		filterText := ""
		if m.Viewport.Filter.Active && len(m.Viewport.Filter.Terms) > 0 {
			filterText = m.Viewport.Filter.Terms[0]
		}

		// Check for warning states first
		if m.PendingInstall {
			// Show warning status for pending installation
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("⚠ Press Enter to install %s or Esc to cancel", m.InstallingPkg),
				width,
			)
		} else if m.PendingRemoval {
			// Show warning status for pending removal
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("⚠ Press Enter to remove %s or Esc to cancel", m.RemovingPkg),
				width,
			)
		} else if m.Installing {
			// Show status for active installation
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("%s Installing %s...", m.GetSpinner(), m.InstallingPkg),
				width,
			)
		} else if m.Removing {
			// Show status for active removal
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("%s Removing %s...", m.GetSpinner(), m.RemovingPkg),
				width,
			)
		} else if m.InstallOutput != "" && m.InstallError == "" {
			// Show success message for installation
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("✓ Package installed successfully. Press Enter to dismiss."),
				width,
			)
		} else if m.InstallError != "" {
			// Show installation error
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("Error installing package: %s", m.InstallError),
				width,
			)
		} else if m.RemoveError != "" {
			// Show removal error
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("Error removing package: %s", m.RemoveError),
				width,
			)
		} else if isRemoteMode {
			// Show remote status line
			errorMsg := m.RemoteError
			statusBar = renderer.RenderRemoteStatus(
				m.RemoteQuery,
				len(m.Viewport.VisibleRows),
				m.Viewport.Height,
				m.Viewport.Offset,
				filterText,
				m.RemoteLoading,
				m.GetSpinner(),
				errorMsg,
				m.Installing,
				m.InstallingPkg,
				width,
			)
		} else {
			// Show normal status line
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
			isRemoteMode,
		)
		// Add status bar at the bottom for detail panel
		if statusBar != "" {
			return lipgloss.JoinVertical(lipgloss.Left, tableUI, statusBar)
		}
		return tableUI
	} else if commandPalette != "" || outputPalette != "" {
		// Show command palette or output palette overlay
		palette := commandPalette
		if outputPalette != "" {
			palette = outputPalette
		}
		tableUI = renderer.RenderWithPaletteOverlayAndMode(
			width,
			m.Height,
			m.Viewport.Columns,
			colWidths,
			visibleRows,
			relativeSelectedRow,
			m.Viewport.SelectedCol,
			m.Viewport.SortColumn,
			m.Viewport.SortReverse,
			palette,
			paletteRows,
			m.Viewport.Offset,
			isRemoteMode,
		)
	} else {
		// No overlay, render normally
		tableUI = renderer.RenderWithPaletteOverlayAndMode(
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
			isRemoteMode,
		)
	}

	// Add status bar at the bottom
	if statusBar != "" {
		return lipgloss.JoinVertical(lipgloss.Left, tableUI, statusBar)
	}

	return tableUI
}
