package app

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/repository"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/viewport"
)

// SpinnerFrames contains the ASCII spinner animation frames.
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// InputMode represents the current input mode of the application.
type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCommand
	ModeFilter
)

// ViewMode represents whether viewing local or remote packages.
type ViewMode int

const (
	ViewLocal ViewMode = iota
	ViewRemote
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

	// Detail panel state
	ShowDetailPanel bool

	// Remote mode state
	ViewMode      ViewMode
	RemoteQuery   string        // The search query
	RemoteLoading bool          // Whether a remote query is in progress
	RemoteError   string        // Error message from remote query
	LocalRows     []*domain.Row // Cached local rows when in remote mode
	SpinnerFrame  int           // Current spinner animation frame
}

// packagesLoadedMsg is sent when packages are loaded.
type packagesLoadedMsg struct {
	packages []*domain.Package
	err      error
}

// remoteSearchResultMsg is sent when a remote search completes.
type remoteSearchResultMsg struct {
	packages []*domain.Package
	query    string
	err      error
}

// spinnerTickMsg is sent to animate the spinner.
type spinnerTickMsg struct{}

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

// SetPreset sets a specific preset by name.
func (m *Model) SetPreset(presetName string) bool {
	for i, preset := range m.Presets {
		if string(preset.Type) == presetName {
			m.CurrentPreset = i
			m.applyCurrentPreset()
			return true
		}
	}
	return false
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

// EnterRemoteMode starts a remote search.
func (m *Model) EnterRemoteMode(query string) tea.Cmd {
	// Cache current local rows
	m.LocalRows = m.Viewport.AllRows

	m.ViewMode = ViewRemote
	m.RemoteQuery = query
	m.RemoteLoading = true
	m.RemoteError = ""
	m.SpinnerFrame = 0

	// Hide install date column for remote mode
	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = false
			break
		}
	}

	// Start spinner and search concurrently
	return tea.Batch(
		m.doRemoteSearch(query),
		tickSpinner(),
	)
}

// ExitRemoteMode returns to local mode and restores state.
func (m *Model) ExitRemoteMode() {
	m.ViewMode = ViewLocal
	m.RemoteQuery = ""
	m.RemoteLoading = false
	m.RemoteError = ""

	// Show install date column again
	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = true
			break
		}
	}

	// Restore local rows
	if m.LocalRows != nil {
		m.Viewport.SetRows(m.LocalRows)
		m.LocalRows = nil
	}

	// Reset to default view
	m.Viewport.SortColumn = column.ColName
	m.Viewport.SortReverse = false
	m.Viewport.ClearFilter()
	m.Viewport.ScrollToTop()
	m.CurrentPreset = 0
	m.applyCurrentPreset()
}

// doRemoteSearch performs the remote search query.
func (m Model) doRemoteSearch(query string) tea.Cmd {
	return func() tea.Msg {
		packages, err := m.Repo.Search(query)
		return remoteSearchResultMsg{
			packages: packages,
			query:    query,
			err:      err,
		}
	}
}

// tickSpinner returns a command that sends spinner tick messages.
func tickSpinner() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

// GetSpinner returns the current spinner frame character.
func (m Model) GetSpinner() string {
	return SpinnerFrames[m.SpinnerFrame%len(SpinnerFrames)]
}
