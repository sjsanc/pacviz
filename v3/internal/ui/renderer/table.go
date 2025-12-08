package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// RenderTable renders the package table rows.
func RenderTable(rows []*domain.Row, columns []*column.Column, colWidths []int, selectedRow int) string {
	if len(rows) == 0 {
		return styles.Row.Render("No packages to display")
	}

	renderedRows := make([]string, 0, len(rows))

	for rowIdx, row := range rows {
		cells := make([]string, 0, len(columns))

		for colIdx, col := range columns {
			if !col.Visible {
				continue
			}

			// Get cell content
			content := row.Cells[col.Type]

			// Account for padding in width calculation
			contentWidth := colWidths[colIdx] - (CellPadding * 2)
			if contentWidth < 1 {
				contentWidth = 1
			}

			// Truncate or pad to content width
			if len(content) > contentWidth {
				if contentWidth > 3 {
					content = content[:contentWidth-3] + "..."
				} else {
					content = content[:contentWidth]
				}
			} else {
				content = content + strings.Repeat(" ", contentWidth-len(content))
			}

			// Add horizontal padding
			content = strings.Repeat(" ", CellPadding) + content + strings.Repeat(" ", CellPadding)

			// Choose style
			var style lipgloss.Style
			if rowIdx == selectedRow {
				style = styles.RowSelected
			} else if rowIdx%2 == 0 {
				style = styles.Row
			} else {
				style = styles.RowAlt
			}

			cells = append(cells, style.Render(content))
		}

		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}
