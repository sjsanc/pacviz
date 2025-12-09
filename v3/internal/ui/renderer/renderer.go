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
	// Render header
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

	// Render table
	table := RenderTableWithMode(rows, columns, colWidths, selectedRow, offset, remoteMode)

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

// RenderWithDetailPanel renders the UI with the detail panel either as overlay or side panel.
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
	isSmallScreen := IsSmallScreen(width)
	detailPanel := RenderDetailPanel(selectedPackage, columns, colWidths, width, isSmallScreen)

	if isSmallScreen {
		// Small screen: render detail panel as overlay above status bar
		// Count lines in detail panel
		detailLines := strings.Count(detailPanel, "\n") + 1

		// Render header
		header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

		// Render table with remote mode flag
		table := RenderTableWithMode(rows, columns, colWidths, selectedRow, offset, remoteMode)

		// Calculate available space for table (height - header - detail panel - status bar)
		// Note: status bar line is reserved but not included in this render
		headerLines := strings.Count(header, "\n") + 1
		statusBarLines := 1
		availableTableLines := height - headerLines - statusBarLines - detailLines

		// Split table into lines
		tableLines := strings.Split(table, "\n")

		// Remove rows that would collide with detail panel
		rowsInCollisionZone := 0
		if len(tableLines) > availableTableLines {
			rowsInCollisionZone = len(tableLines) - availableTableLines
			if rowsInCollisionZone > detailLines {
				rowsInCollisionZone = detailLines
			}
		}

		if rowsInCollisionZone > 0 {
			tableLines = tableLines[:len(tableLines)-rowsInCollisionZone]
		}

		table = strings.Join(tableLines, "\n")

		// Add filler lines to push detail panel to bottom if table is shorter
		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}

		// Add detail panel below
		table = lipgloss.JoinVertical(lipgloss.Left, table, detailPanel)

		// Compose sections
		content := lipgloss.JoinVertical(lipgloss.Left, header, table)

		// Ensure content fills exactly (height - statusBarLines) so status bar stays at bottom
		contentLines := strings.Count(content, "\n") + 1
		targetLines := height - statusBarLines
		if contentLines < targetLines {
			// Add filler to push content to fill the full height minus status bar
			additionalFiller := strings.Repeat("\n", targetLines-contentLines)
			content = content + additionalFiller
		}

		// Apply width constraint
		style := lipgloss.NewStyle().Width(width)
		return style.Render(content)
	} else {
		// Large screen: render detail panel on the right side
		// Calculate panel width using the same logic as the detail panel renderer
		panelWidth := CalculateDetailPanelWidth(width)

		// Adjust table width
		const spacing = 2
		tableWidth := width - panelWidth - spacing

		// Recalculate column widths for reduced table width
		tableColWidths := column.CalculateWidths(columns, tableWidth)

		// Render header with adjusted widths
		header := RenderHeader(columns, tableColWidths, selectedCol, sortCol, sortReverse)

		// Render table with adjusted widths and remote mode flag
		table := RenderTableWithMode(rows, columns, tableColWidths, selectedRow, offset, remoteMode)

		// Compose header and table vertically
		tableContent := lipgloss.JoinVertical(lipgloss.Left, header, table)

		// Add filler to table content to fill height minus status bar
		tableContentLines := strings.Count(tableContent, "\n") + 1
		targetLines := height - 1 // Reserve 1 line for status bar
		if tableContentLines < targetLines {
			additionalFiller := strings.Repeat("\n", targetLines-tableContentLines)
			tableContent = tableContent + additionalFiller
		}

		// Place table and detail panel side by side
		content := lipgloss.JoinHorizontal(lipgloss.Top, tableContent, "  ", detailPanel)

		// Apply width constraint
		style := lipgloss.NewStyle().Width(width)
		return style.Render(content)
	}
}
