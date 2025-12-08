package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// RenderHeader renders the table header with column names.
func RenderHeader(columns []*column.Column, colWidths []int, selectedCol int, sortCol column.Type, sortReverse bool) string {
	headers := make([]string, 0, len(columns))

	for i, col := range columns {
		if !col.Visible {
			continue
		}

		// Build header text with sort indicator
		header := col.Name
		if col.Type == sortCol && col.Sortable {
			if sortReverse {
				header += " ↓"
			} else {
				header += " ↑"
			}
		}

		// Account for padding in width calculation
		contentWidth := colWidths[i] - (CellPadding * 2)
		if contentWidth < 1 {
			contentWidth = 1
		}

		// Handle text alignment for index column (right-aligned)
		if col.Type == column.ColIndex {
			// Right-align for index column
			if len(header) < contentWidth {
				header = strings.Repeat(" ", contentWidth-len(header)) + header
			}
		} else {
			// Default left-align
			if len(header) > contentWidth {
				if contentWidth > 3 {
					header = header[:contentWidth-3] + "..."
				} else {
					header = header[:contentWidth]
				}
			} else {
				header = header + strings.Repeat(" ", contentWidth-len(header))
			}
		}

		// Add horizontal padding
		header = strings.Repeat(" ", CellPadding) + header + strings.Repeat(" ", CellPadding)

		// Apply style
		style := styles.Header
		if col.Type == column.ColIndex {
			// Index column header uses dimmed style
			style = styles.Header.Copy().Foreground(styles.Dimmed)
		}
		if i == selectedCol {
			style = style.Background(lipgloss.Color("#3d59a1"))
		}

		headers = append(headers, style.Render(header))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}
