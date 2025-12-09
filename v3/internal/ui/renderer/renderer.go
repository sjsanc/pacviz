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
	detailPanel := RenderDetailPanel(selectedPackage, columns, colWidths, width, isSmallScreen, remoteMode)

	if isSmallScreen {
		detailLines := strings.Count(detailPanel, "\n") + 1
		header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)
		table := RenderTableWithMode(rows, columns, colWidths, selectedRow, offset, remoteMode)

		headerLines := strings.Count(header, "\n") + 1
		statusBarLines := 1
		availableTableLines := height - headerLines - statusBarLines - detailLines

		tableLines := strings.Split(table, "\n")

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

		if len(tableLines) < availableTableLines {
			fillerLines := availableTableLines - len(tableLines)
			filler := strings.Repeat("\n", fillerLines)
			table = table + filler
		}

		table = lipgloss.JoinVertical(lipgloss.Left, table, detailPanel)
		content := lipgloss.JoinVertical(lipgloss.Left, header, table)

		contentLines := strings.Count(content, "\n") + 1
		targetLines := height - statusBarLines
		if contentLines < targetLines {
			additionalFiller := strings.Repeat("\n", targetLines-contentLines)
			content = content + additionalFiller
		}

		style := lipgloss.NewStyle().Width(width)
		return style.Render(content)
	} else {
		panelWidth := CalculateDetailPanelWidth(width)

		const spacing = 2
		tableWidth := width - panelWidth - spacing

		tableColWidths := column.CalculateWidths(columns, tableWidth)

		header := RenderHeader(columns, tableColWidths, selectedCol, sortCol, sortReverse)
		table := RenderTableWithMode(rows, columns, tableColWidths, selectedRow, offset, remoteMode)

		tableContent := lipgloss.JoinVertical(lipgloss.Left, header, table)

		tableContentLines := strings.Count(tableContent, "\n") + 1
		targetLines := height - 1
		if tableContentLines < targetLines {
			additionalFiller := strings.Repeat("\n", targetLines-tableContentLines)
			tableContent = tableContent + additionalFiller
		}

		content := lipgloss.JoinHorizontal(lipgloss.Top, tableContent, "  ", detailPanel)

		style := lipgloss.NewStyle().Width(width)
		return style.Render(content)
	}
}
