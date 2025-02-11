package internal

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		m.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "left":
			m.VP.PrevColumn()
		case "right":
			m.VP.NextColumn()
		case " ":
			col := m.VP.GetCurrentColumn()
			rev := !m.VP.SortReversed
			m.VP.ApplySort(col, rev)
		case "up":
			m.VP.ScrollUp(1)
		case "down":
			m.VP.ScrollDown(1)
		}
	}

	return m, cmd
}
