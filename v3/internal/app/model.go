package app

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/repository"
	"github.com/sjsanc/pacviz/v3/internal/ui/viewport"
)

// Model is the main Bubble Tea model for the application.
type Model struct {
	Width    int
	Height   int
	Viewport *viewport.Viewport
	Repo     repository.Repository
	Error    error
	Ready    bool
}

// packagesLoadedMsg is sent when packages are loaded.
type packagesLoadedMsg struct {
	packages []*domain.Package
	err      error
}

// NewModel creates a new application model.
func NewModel() Model {
	repo, err := repository.NewAlpmRepository()
	if err != nil {
		log.Printf("Failed to initialize repository: %v", err)
		return Model{
			Viewport: viewport.New(),
			Error:    fmt.Errorf("failed to initialize repository: %w", err),
			Ready:    false,
		}
	}

	return Model{
		Viewport: viewport.New(),
		Repo:     repo,
		Ready:    false,
	}
}

// Init initializes the model (Bubble Tea interface).
func (m Model) Init() tea.Cmd {
	return m.loadPackages
}

// loadPackages loads installed packages from the repository.
func (m Model) loadPackages() tea.Msg {
	packages, err := m.Repo.GetInstalled()
	return packagesLoadedMsg{packages: packages, err: err}
}
