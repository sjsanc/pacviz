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
		// Calculate width for side panel (30% of screen width, min 40 chars)
		panelWidth := width * 30 / 100
		if panelWidth < 40 {
			panelWidth = 40
		}

		panelStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Accent1).
			Padding(0, 1).
			Width(panelWidth - 4)

		return panelStyle.Render(content)
	}
}

// IsSmallScreen determines if the screen is small based on width.
func IsSmallScreen(width int) bool {
	return width < 120
}
