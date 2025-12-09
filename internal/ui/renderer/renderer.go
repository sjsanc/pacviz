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
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)
	table := RenderTable(rows, columns, colWidths, selectedRow, offset)
	content := lipgloss.JoinVertical(lipgloss.Left, header, table)

	if statusBar != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, statusBar)
	}

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
	return RenderWithPaletteOverlayAndMode(width, height, columns, colWidths, rows, selectedRow, selectedCol, sortCol, sortReverse, paletteContent, paletteRows, offset, false)
}

// RenderWithPaletteOverlayAndMode renders the UI with optional remote mode styling.
func RenderWithPaletteOverlayAndMode(
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
	remoteMode bool,
) string {
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)
	table := RenderTableWithMode(rows, columns, colWidths, selectedRow, offset, remoteMode)

	headerLines := strings.Count(header, "\n") + 1
	statusBarLines := 1
	availableTableLines := height - headerLines - statusBarLines - paletteRows

	tableLines := strings.Split(table, "\n")

	if paletteContent != "" && paletteRows > 0 {
		rowsInCollisionZone := 0
		if len(tableLines) > availableTableLines {
			rowsInCollisionZone = len(tableLines) - availableTableLines
			if rowsInCollisionZone > paletteRows {
				rowsInCollisionZone = paletteRows
			}
		}

		if rowsInCollisionZone > 0 {
			tableLines = tableLines[:len(tableLines)-rowsInCollisionZone]
		}

		table = strings.Join(tableLines, "\n")

		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}

		table = lipgloss.JoinVertical(lipgloss.Left, table, paletteContent)
	} else {
		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, header, table)
	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}

// RenderWithDetailPanel renders the UI with the detail panel inserted between the table and status bar.
// The detail panel shortens the table height by its own height, but the status bar remains at the bottom.
func RenderWithDetailPanel(
	width, height int,
	columns []*column.Column,
	colWidths []int,
	rows []*domain.Row,
	selectedRow, selectedCol int,
	sortCol column.Type,
	sortReverse bool,
	selectedPackage *domain.Package,
	offset int,
	remoteMode bool,
) string {
	// Render the detail panel
	detailPanel := RenderDetailPanel(selectedPackage, columns, colWidths, width, true, remoteMode)

	detailLines := strings.Count(detailPanel, "\n") + 1
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

	headerLines := strings.Count(header, "\n") + 1
	statusBarLines := 1
	// Available table lines = total height - header - detail panel - status bar
	availableTableLines := height - headerLines - statusBarLines - detailLines

	// Render table with available height
	table := RenderTableWithMode(rows, columns, colWidths, selectedRow, offset, remoteMode)
	tableLines := strings.Split(table, "\n")

	// Clip table rows to fit available space
	if len(tableLines) > availableTableLines {
		tableLines = tableLines[:availableTableLines]
	}

	// Pad table to fill available space
	if len(tableLines) < availableTableLines {
		fillerLines := availableTableLines - len(tableLines)
		filler := strings.Repeat("\n", fillerLines)
		tableLines = append(tableLines, strings.Split(filler, "\n")...)
	}

	// Build the layout: header -> table -> detail panel (will be followed by status bar in view.go)
	table = strings.Join(tableLines, "\n")
	content := lipgloss.JoinVertical(lipgloss.Left, header, table, detailPanel)

	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}
