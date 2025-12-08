package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Render renders the entire UI.
func Render(
	width, height int,
	columns []*column.Column,
	colWidths []int,
	rows []*domain.Row,
	selectedRow, selectedCol int,
	sortCol column.Type,
	sortReverse bool,
	statusBar string,
	offset int,
) string {
	// Render header
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

	// Render table
	table := RenderTable(rows, columns, colWidths, selectedRow, offset)

	// Compose sections
	content := lipgloss.JoinVertical(lipgloss.Left, header, table)

	// Add status bar if provided
	if statusBar != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, statusBar)
	}

	// Apply width constraint to ensure full terminal width
	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}

// RenderWithPaletteOverlay renders the UI with command palette overlaying the bottom table rows.
func RenderWithPaletteOverlay(
	width, height int,
	columns []*column.Column,
	colWidths []int,
	rows []*domain.Row,
	selectedRow, selectedCol int,
	sortCol column.Type,
	sortReverse bool,
	paletteContent string,
	paletteRows int,
	offset int,
) string {
	// Render header
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

	// Render table
	table := RenderTable(rows, columns, colWidths, selectedRow, offset)

	// Calculate available space for table (height - header - status bar - palette)
	headerLines := strings.Count(header, "\n") + 1
	statusBarLines := 1
	availableTableLines := height - headerLines - statusBarLines - paletteRows

	// Split table into lines
	tableLines := strings.Split(table, "\n")

	// If palette is active, overlay it on the bottom rows
	if paletteContent != "" && paletteRows > 0 {
		// Calculate how many table rows would actually collide with the palette
		// Only rows that extend into the palette area need to be removed
		rowsInCollisionZone := 0
		if len(tableLines) > availableTableLines {
			// Table extends past available space, so it would collide with palette
			rowsInCollisionZone = len(tableLines) - availableTableLines
			// But don't remove more than palette rows
			if rowsInCollisionZone > paletteRows {
				rowsInCollisionZone = paletteRows
			}
		}

		// Only remove rows that would actually be obscured
		if rowsInCollisionZone > 0 {
			tableLines = tableLines[:len(tableLines)-rowsInCollisionZone]
		}

		table = strings.Join(tableLines, "\n")

		// Add filler lines to push palette to bottom if table is shorter than available space
		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}

		// Add palette below
		table = lipgloss.JoinVertical(lipgloss.Left, table, paletteContent)
	} else {
		// No palette: add filler lines to fill the viewport
		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}
	}

	// Compose sections
	content := lipgloss.JoinVertical(lipgloss.Left, header, table)

	// Apply width constraint to ensure full terminal width
	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}
