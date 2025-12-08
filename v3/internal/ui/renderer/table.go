package renderer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// RenderTable renders the package table rows.
func RenderTable(rows []*domain.Row, columns []*column.Column, colWidths []int, selectedRow int, offset int) string {
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
			var content string
			if col.Type == column.ColIndex {
				// Index should reflect absolute position in the full list (1-based)
				content = fmt.Sprintf("%d", offset+rowIdx+1)
			} else {
				content = row.Cells[col.Type]
			}

			// Account for padding in width calculation
			contentWidth := colWidths[colIdx] - (CellPadding * 2)
			if contentWidth < 1 {
				contentWidth = 1
			}

			// Handle text alignment for index column (right-aligned)
			if col.Type == column.ColIndex {
				// Right-align for index column
				if len(content) < contentWidth {
					content = strings.Repeat(" ", contentWidth-len(content)) + content
				}
			} else {
				// Default left-align
				if len(content) > contentWidth {
					if contentWidth > 3 {
						content = content[:contentWidth-3] + "..."
					} else {
						content = content[:contentWidth]
					}
				} else {
					content = content + strings.Repeat(" ", contentWidth-len(content))
				}
			}

			// Add horizontal padding
			content = strings.Repeat(" ", CellPadding) + content + strings.Repeat(" ", CellPadding)

			// Choose style
			var style lipgloss.Style
			if col.Type == column.ColIndex {
				// Index column uses dimmed style (dark-ish foreground)
				style = styles.Index
				if rowIdx == selectedRow {
					style = styles.RowSelected.Copy().Foreground(styles.Dimmed)
				} else if rowIdx%2 == 0 {
					style = styles.Index.Copy().Background(styles.Background)
				} else {
					style = styles.Index.Copy().Background(lipgloss.Color("#16161e"))
				}
			} else {
				// Regular row styling
				if rowIdx == selectedRow {
					style = styles.RowSelected
				} else if rowIdx%2 == 0 {
					style = styles.Row
				} else {
					style = styles.RowAlt
				}
			}

			cells = append(cells, style.Render(content))
		}

		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}
