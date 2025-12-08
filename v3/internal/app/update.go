package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// Update handles incoming messages and updates the model (Bubble Tea interface).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case packagesLoadedMsg:
		return m.handlePackagesLoaded(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	}

	return m, nil
}

func (m Model) handlePackagesLoaded(msg packagesLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.Error = msg.err
		return m, nil
	}

	// Convert packages to rows
	rows := domain.PackagesToRows(msg.packages)
	m.Viewport.SetRows(rows)
	m.Ready = true

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	// TODO: Arrow key navigation
	// TODO: Command mode (`:`)
	// TODO: Filter mode (`/`)
	// TODO: Sort toggle (space)
	// TODO: Info view (enter, i)
	// TODO: Tab for preset switching
	}

	return m, nil
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height
	// Reserve 1 line for header
	m.Viewport.Height = msg.Height - 1
	return m, nil
}
