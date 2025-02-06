package main

import (
	"fmt"
	"strings"

	"github.com/sjsanc/pacviz/styles"
)

var COLUMNS = []ColDef{
	{"ID", Fixed, 4},
	{"Name", Percent, 15},
	{"Version", Percent, 15},
	{"Groups", Percent, 10},
	{"Installed", Percent, 10},
	{"Description", Percent, 50},
}

func (m *Model) View() string {
	lines := make([]string, m.height)

	m.ColumnWidths = calcColWidths(COLUMNS, m.width)

	rows := m.vp.VisibleRows()

	m.RenderHeader(&lines, COLUMNS)
	m.RenderRows(&lines, rows)
	m.RenderStatus(&lines)

	var body strings.Builder
	for i, l := range lines {
		body.WriteString(l)
		if i != len(lines)-1 {
			body.WriteString("\n")
		}
	}
	return body.String()
}

func trunc(s string, w int) string {
	if w < 3 {
		// Handle startup edge case
		return ""
	}
	if len(s) > w && w > 3 {
		return s[:w-3] + "..."
	}
	if len(s) > w {
		return s[:w]
	}
	return s
}

func (m *Model) RenderHeader(lines *[]string, cols []ColDef) {
	line := strings.Builder{}
	for i, h := range cols {
		w := m.ColumnWidths[i]
		c := trunc(h.Name, w)
		s := styles.Header.Width(w).Render(c)
		line.WriteString(s)
	}
	*lines = append(*lines, line.String())
}

func (m *Model) RenderRows(lines *[]string, rows []*Row) {
	for i, r := range rows {
		row := strings.Builder{}
		s := styles.Row
		if i%2 == 0 {
			s = styles.RowAlt
		}
		if m.vp.offset+i == m.vp.selected {
			s = styles.RowSelected
		}
		for j, c := range r.cells {
			w := m.ColumnWidths[j]
			c := trunc(c, w)
			row.WriteString(s.Width(w).Render(c))
		}
		*lines = append(*lines, row.String())
	}
	for i := len(rows); i < m.vp.height; i++ {
		*lines = append(*lines, "")
	}
}

func (m *Model) RenderStatus(lines *[]string) {
	if m.CommandMode {
		text := m.buffer
		s := styles.Footer.Width(m.width).Render(text)
		*lines = append(*lines, s)
		return
	}

	if m.FilterMode {
		text := m.buffer
		s := styles.Footer.Width(m.width).Render(text)
		*lines = append(*lines, s)
		return
	}

	// Rows: 1-100 of 100
	text := fmt.Sprintf("Rows: %d-%d of %d", m.vp.offset, m.vp.offset+m.vp.height-1, len(m.vp.rows))
	s := styles.Footer.Width(m.width).Render(text)
	*lines = append(*lines, s)
}

func (m *Model) RenderEmptySpace(lines *[]string, count int) {
	for i := 0; i < count; i++ {
		*lines = append(*lines, "")
	}
}
