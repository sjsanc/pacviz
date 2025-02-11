package main

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetDimensions(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "right", "left", "tab":
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			m.vp.ScrollUp(1)
		case "down":
			m.vp.ScrollDown(1)
		case "esc":
			m.HandleEsc()
		default:
			m.HandleInput(msg.String())
		}
	}

	return m, cmd
}

func (m *Model) HandleEsc() {
	if m.CommandMode {
		m.ExitCommandMode()
	}

	if m.FilterMode {
		m.HideFilterBuffer()
		m.ResetFilter()
	}
}

func (m *Model) HandleInput(msg string) {
	if !m.CommandMode {
		if msg == ":" {
			m.ShowCommandBuffer()
		}
	}

	if !m.FilterMode {
		if msg == "/" {
			m.ShowFilterBuffer()
		}
	}

	if m.CommandMode {
		if msg == "enter" {
			m.ExitCommandMode()
		}

		m.WriteToBuffer(msg)

		cmd := matchCommand(m.buffer)

		cmd(m)
	}

	if m.FilterMode {
		if msg == "enter" {
			m.HideFilterBuffer()
		}

		m.WriteToBuffer(msg)

		rows := convertPkgsToRows(m.pacman.GetExplicitPkgs())
		m.vp.SetRows(filter(m.buffer[1:], rows))
	}
}

func (m *Model) WriteToBuffer(msg string) {
	if msg == "backspace" {
		if len(m.buffer) > 1 {
			m.buffer = m.buffer[:len(m.buffer)-1]
		}
		return
	}

	if msg == "enter" {
		return
	}

	m.buffer += msg
}

func (m *Model) ClearBuffer() {
	m.buffer = ""
}

func (m *Model) ShowCommandBuffer() {
	m.CommandMode = true
}

func (m *Model) HideCommandBuffer() {
	m.CommandMode = false
}

func (m *Model) ExitCommandMode() {
	m.HideCommandBuffer()
	m.ClearBuffer()
}

func (m *Model) ShowFilterBuffer() {
	m.FilterMode = true
}

func (m *Model) HideFilterBuffer() {
	m.FilterMode = false
}

// This should reset to the current View
// The current view can either be Explicit, Orphans or Sync
func (m *Model) ResetFilter() {
	rows := convertPkgsToRows(m.pacman.GetExplicitPkgs())
	m.vp.SetRows(rows)
}

func matchCommand(buffer string) func(*Model) {
	if strings.HasPrefix(buffer, ":g") {
		arg := buffer[2:]
		if arg != "" {
			line, err := strconv.Atoi(arg)
			if err == nil {
				return func(m *Model) {
					m.vp.GoTo(line)
				}
			}
		} else {
			return func(m *Model) {
				m.vp.GoTo(m.vp.offset)
			}
		}
	}

	if strings.HasPrefix(buffer, ":e") {
		return func(m *Model) {
			m.vp.GoTo(len(m.vp.rows) - 1)
		}
	}

	if strings.HasPrefix(buffer, ":t") {
		return func(m *Model) {
			m.vp.GoTo(0)
		}
	}

	return func(m *Model) {}
}

func filter(term string, rows []*Row) []*Row {
	var filtered []*Row
	for _, row := range rows {
		if strings.Contains(row.cells[1], term) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}
