package app

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/repository"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/viewport"
)

// InputMode represents the current input mode of the application.
type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCommand
	ModeFilter
)

// Model is the main Bubble Tea model for the application.
type Model struct {
	Width    int
	Height   int
	Viewport *viewport.Viewport
	Repo     repository.Repository
	Error    error
	Ready    bool

	// Input mode state
	Mode   InputMode
	Buffer string

	// Preset state
	Presets       []domain.Preset
	CurrentPreset int // Index into Presets slice
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
			Viewport:      viewport.New(),
			Error:         fmt.Errorf("failed to initialize repository: %w", err),
			Ready:         false,
			Presets:       domain.DefaultPresets(),
			CurrentPreset: 0, // Start with Explicit preset
		}
	}

	return Model{
		Viewport:      viewport.New(),
		Repo:          repo,
		Ready:         false,
		Presets:       domain.DefaultPresets(),
		CurrentPreset: 0, // Start with Explicit preset
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

// Buffer Management Methods

// EnterCommandMode switches to command mode and initializes the buffer.
func (m *Model) EnterCommandMode() {
	m.Mode = ModeCommand
	m.Buffer = ":"
}

// EnterFilterMode switches to filter mode and initializes the buffer.
func (m *Model) EnterFilterMode() {
	m.Mode = ModeFilter
	m.Buffer = "/"
}

// ExitMode returns to normal mode and clears the buffer.
func (m *Model) ExitMode() {
	m.Mode = ModeNormal
	m.Buffer = ""
}

// WriteToBuffer appends a character or handles special keys in the buffer.
func (m *Model) WriteToBuffer(key string) {
	// Handle backspace
	if key == "backspace" {
		if len(m.Buffer) > 1 { // Keep the prefix (: or /)
			m.Buffer = m.Buffer[:len(m.Buffer)-1]
		}
		return
	}

	// Ignore enter key (handled separately)
	if key == "enter" {
		return
	}

	// Append character to buffer
	m.Buffer += key
}

// ClearBuffer resets the buffer to empty.
func (m *Model) ClearBuffer() {
	m.Buffer = ""
}

// GetBufferContent returns the buffer without the prefix.
func (m *Model) GetBufferContent() string {
	if len(m.Buffer) > 1 {
		return m.Buffer[1:] // Strip the : or / prefix
	}
	return ""
}

// NextPreset cycles to the next preset, resetting sort and filters.
func (m *Model) NextPreset() {
	m.CurrentPreset = (m.CurrentPreset + 1) % len(m.Presets)
	m.applyCurrentPreset()
}

// applyCurrentPreset applies the current preset filter and resets sort/filters.
func (m *Model) applyCurrentPreset() {
	// Reset sort to default (Name, ascending)
	m.Viewport.SortColumn = column.ColName
	m.Viewport.SortReverse = false

	// Clear any active filters
	m.Viewport.ClearFilter()

	// Apply preset filter
	preset := m.Presets[m.CurrentPreset]
	m.Viewport.ApplyPresetFilter(preset.Filter)
}
