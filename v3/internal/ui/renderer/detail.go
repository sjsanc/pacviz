package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// RenderDetailPanel renders the package detail panel.
// For small screens (width < 120), it's rendered above the status bar as an overlay.
// For large screens (width >= 120), it's rendered on the right side.
func RenderDetailPanel(pkg *domain.Package, columns []*column.Column, colWidths []int, width int, isSmallScreen bool) string {
	if pkg == nil {
		return ""
	}

	// Style for labels (same as index column)
	labelStyle := styles.Index
	valueStyle := lipgloss.NewStyle().Foreground(styles.Foreground)

	// Create a temporary row to get formatted cell values
	row := domain.PackageToRow(pkg, 0)

	// Build the detail lines
	var lines []string
	for _, col := range columns {
		// Skip the index column
		if col.Type == column.ColIndex {
			continue
		}

		label := col.Name
		value := row.Cells[col.Type]

		// Format: "Label: value"
		line := labelStyle.Render(label+":") + " " + valueStyle.Render(value)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")

	if isSmallScreen {
		// Small screen: render as overlay above status bar
		// Add border and padding
		panelStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Accent1).
			Padding(0, 1).
			Width(width - 4)

		return panelStyle.Render(content)
	} else {
		// Large screen: render as side panel on the right
		panelWidth := CalculateDetailPanelWidth(width)

		panelStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Accent1).
			Padding(0, 1).
			Width(panelWidth - 4)

		return panelStyle.Render(content)
	}
}

// IsSmallScreen determines if the screen is small based on width.
// Below this threshold, the detail panel renders as an overlay above the status bar.
// At or above this threshold, it renders side-by-side with the table.
func IsSmallScreen(width int) bool {
	// Minimum width for side-by-side rendering:
	// - Table needs at least 80 chars to be readable
	// - Detail panel needs at least 50 chars
	// - 2 chars for spacing
	// Total: 132 chars minimum
	return width < 140
}

// CalculateDetailPanelWidth calculates the width of the detail panel for large screens.
// Returns the total width of the panel including borders.
func CalculateDetailPanelWidth(width int) int {
	const maxPanelWidth = 60
	const minTableWidth = 80
	const spacing = 2

	// Calculate panel width ensuring table has minimum width
	panelWidth := maxPanelWidth
	if width-spacing < minTableWidth+panelWidth {
		// Not enough space for max panel width, reduce it
		panelWidth = width - minTableWidth - spacing
	}

	// Ensure panel has minimum width of 45 chars
	if panelWidth < 45 {
		panelWidth = 45
	}

	return panelWidth
}
