package main

import (
	"strings"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/styles"
)

func (m *Model) BuildCenter() string {
	width := m.width - 1
	block := lg.NewStyle().
		Width(width).
		Height(m.height)

	sb := strings.Builder{}

	widths := calcColWidths(COLUMNS, width)

	for i, h := range COLUMNS {
		w := widths[i]
		c := trunc(h.Name, w)
		s := styles.Header.Width(w).Render(c)
		sb.WriteString(s)
	}

	rows := m.vp.VisibleRows()

	for i, r := range rows {
		rsb := strings.Builder{}
		s := styles.Row
		if i%2 == 0 {
			s = styles.RowAlt
		}
		if m.vp.offset+i == m.vp.selected {
			s = styles.RowSelected
		}
		for j, c := range r.cells {
			w := widths[j]
			c := trunc(c, w)
			rsb.WriteString(s.Width(w).Render(c))
		}
		sb.WriteString(rsb.String())
	}

	for i := len(m.vp.rows); i < m.vp.height; i++ {
		sb.WriteString("\n")
	}

	return block.Render(sb.String())
}

func (m *Model) BuildScrollbar() string {
	block := lg.NewStyle().
		Width(1).
		Height(m.height)

	sb := strings.Builder{}

	for i := 0; i < 20; i++ {
		sb.WriteString("█")
	}

	return block.Render(sb.String())
}

func (m *Model) View() string {
	block := lg.JoinHorizontal(
		lg.Top,
		m.BuildCenter(),
		m.BuildScrollbar(),
	)

	return block
}
