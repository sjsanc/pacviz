package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

func RenderHeader(columns []*column.Column, colWidths []int, selectedCol int, sortCol column.Type, sortReverse bool) string {
	headers := make([]string, 0, len(columns))

	for i, col := range columns {
		if !col.Visible {
			continue
		}

		header := col.Name
		if col.Type == sortCol && col.Sortable {
			if sortReverse {
				header += " ↓"
			} else {
				header += " ↑"
			}
		} else if col.Sortable {
			header += "  "
		}

		contentWidth := colWidths[i] - (CellPadding * 2)
		if contentWidth < 1 {
			contentWidth = 1
		}

		headerWidth := lipgloss.Width(header)

		if col.Type == column.ColIndex {
			if headerWidth < contentWidth {
				header = strings.Repeat(" ", contentWidth-headerWidth) + header
			}
		} else {
			if headerWidth > contentWidth {
				if contentWidth > 3 {
					header = header[:contentWidth-3] + "..."
				} else {
					header = header[:contentWidth]
				}
			} else {
				header = header + strings.Repeat(" ", contentWidth-headerWidth)
			}
		}

		header = strings.Repeat(" ", CellPadding) + header + strings.Repeat(" ", CellPadding)

		style := styles.Current.Header
		if col.Type == column.ColIndex {
			style = styles.Current.Header.Copy().Foreground(styles.Current.Dimmed)
		}
		if i == selectedCol {
			style = style.Background(lipgloss.Color("#3d59a1"))
		}

		headers = append(headers, style.Render(header))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}
