package internal

import (
	tea "github.com/charmbracelet/bubbletea"
)

var PM = NewPacman()

type Model struct {
	Height int // Terminal height
	Width  int // Terminal width
	VP     *Viewport
}

func NewModel() *Model {
	m := &Model{}

	m.VP = NewViewport()

	m.VP.ApplyPreset("explicit")

	// m.VP.SearchSyncDBs("linux")

	// m.VP.ApplyFilter("dbus")

	return m
}

func (m *Model) SetHeight(h int) *Model {
	m.Height = h
	return m
}

func (m *Model) SetWidth(w int) *Model {
	m.Width = w
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}
