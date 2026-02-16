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
	if m.Error != "" {
		return fmt.Sprintf("Error: %s\nPress Ctrl+C to quit", m.Error)
	}

	if !m.Ready {
		return "Loading packages..."
	}

	width := m.Width
	if width == 0 {
		width = 120
	}

	colWidths := column.CalculateWidths(m.Viewport.Columns, width)
	visibleRows := m.Viewport.GetVisibleRows()
	relativeSelectedRow := m.Viewport.SelectedRow - m.Viewport.Offset

	var statusBar string
	var commandPalette string
	var outputPalette string
	var paletteRows int
	var tableUI string
	isRemoteMode := m.ViewMode == ViewRemote

	switch m.Mode {
	case ModeCommand:
		commandPalette, paletteRows = command.RenderCommandPalette(m.GetBufferContent(), width, isRemoteMode)
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModeFilter:
		statusBar = renderer.RenderStatusWithBuffer(m.Buffer, width)
	case ModePassword:
		prompt := "[sudo] password: "
		masked := ""
		for range m.PasswordBuffer {
			masked += "*"
		}
		statusBar = renderer.RenderWarningStatus(prompt+masked, width)
	case ModeNormal:
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

		if m.PendingInstall {
			installMsg := fmt.Sprintf("⚠ Press Enter to install %s or Esc to cancel", m.InstallingPkg)
			if m.isSelectedPackageAUR() && m.AURHelper != nil {
				installMsg = fmt.Sprintf("⚠ Press Enter to install %s via %s or Esc to cancel", m.InstallingPkg, m.AURHelper.Name)
			}
			statusBar = renderer.RenderWarningStatus(installMsg, width)
		} else if m.PendingRemoval {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("⚠ Press Enter to remove %s or Esc to cancel", m.RemovingPkg),
				width,
			)
		} else if m.Installing {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("%s Installing %s...", m.GetSpinner(), m.InstallingPkg),
				width,
			)
		} else if m.Removing {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("%s Removing %s...", m.GetSpinner(), m.RemovingPkg),
				width,
			)
		} else if m.InstallOutput != "" && m.InstallError == "" {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("✓ Package installed successfully. Press Enter to dismiss."),
				width,
			)
		} else if m.InstallError != "" {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("Error installing package: %s", m.InstallError),
				width,
			)
		} else if m.RemoveError != "" {
			statusBar = renderer.RenderWarningStatus(
				fmt.Sprintf("Error removing package: %s", m.RemoveError),
				width,
			)
		} else if isRemoteMode {
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

	if m.ShowDetailPanel {
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
		if statusBar != "" {
			return lipgloss.JoinVertical(lipgloss.Left, tableUI, statusBar)
		}
		return tableUI
	} else if commandPalette != "" || outputPalette != "" {
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

	if statusBar != "" {
		return lipgloss.JoinVertical(lipgloss.Left, tableUI, statusBar)
	}

	return tableUI
}
