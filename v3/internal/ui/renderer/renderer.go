package renderer

import (
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
) string {
	// Render header
	header := RenderHeader(columns, colWidths, selectedCol, sortCol, sortReverse)

	// Render table
	table := RenderTable(rows, columns, colWidths, selectedRow)

	// Compose and ensure it fills the width
	content := lipgloss.JoinVertical(lipgloss.Left, header, table)

	// Apply width constraint to ensure full terminal width
	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}
