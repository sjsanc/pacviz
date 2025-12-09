package app

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/command"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/repository"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/viewport"
)

// spinnerFrames contains the ASCII spinner animation frames.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// InputMode represents the current input mode of the application.
type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCommand
	ModeFilter
	ModePassword
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
	Error    string
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

	// Installation state
	PendingInstall bool   // Whether a package is pending installation confirmation
	Installing     bool   // Whether a package is being installed
	InstallingPkg  string // Name of package pending/being installed
	InstallError   string // Error message from installation
	InstallOutput  string // Output from installation operation

	// Removal state
	PendingRemoval bool   // Whether a package is pending removal confirmation
	RemovingPkg    string // Name of package pending/being removed
	Removing       bool   // Whether a package is currently being removed
	RemoveError    string // Error message from removal
	RemoveOutput   string // Output from removal operation

	// Password state
	PasswordBuffer string // Buffer for password input (not displayed)
	NeedsPassword  bool   // Whether we need to prompt for password
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

// installCompleteMsg is sent when installation completes.
type installCompleteMsg struct {
	pkgName string
	output  string
	err     error
}

// removeCompleteMsg is sent when removal completes.
type removeCompleteMsg struct {
	pkgName string
	output  string
	err     error
}

// repositoryRefreshedMsg is sent when repository refresh completes.
type repositoryRefreshedMsg struct {
	shouldReload bool
}

// commandResultMsg is sent when a command needs to affect the model.
type commandResultMsg struct {
	Result command.ExecuteResult
}

// executeCommandMsg creates a tea.Cmd that executes a command.
func executeCommandMsg(commandStr string) tea.Cmd {
	return func() tea.Msg {
		result := command.Execute(commandStr)
		return commandResultMsg{Result: result}
	}
}

// NewModel creates a new application model.
func NewModel() *Model {
	repo, err := repository.NewAlpmRepository()
	if err != nil {
		log.Printf("Failed to initialize repository: %v", err)
		return &Model{
			Viewport:      viewport.New(),
			Error:         fmt.Sprintf("failed to initialize repository: %v", err),
			Ready:         false,
			Presets:       domain.DefaultPresets(),
			CurrentPreset: 0,
		}
	}

	return &Model{
		Viewport:      viewport.New(),
		Repo:          repo,
		Ready:         false,
		Presets:       domain.DefaultPresets(),
		CurrentPreset: 0,
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

// refreshRepository refreshes the package database.
func (m Model) refreshRepository(shouldReload bool) tea.Cmd {
	return func() tea.Msg {
		err := m.Repo.Refresh()
		if err != nil {
			log.Printf("Failed to refresh repository: %v", err)
		}
		return repositoryRefreshedMsg{shouldReload: shouldReload}
	}
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

// EnterPasswordMode switches to password mode.
func (m *Model) EnterPasswordMode() {
	m.Mode = ModePassword
	m.PasswordBuffer = ""
	m.NeedsPassword = true
}

// WriteToPasswordBuffer appends a character to the password buffer.
func (m *Model) WriteToPasswordBuffer(key string) {
	if key == "backspace" {
		if len(m.PasswordBuffer) > 0 {
			m.PasswordBuffer = m.PasswordBuffer[:len(m.PasswordBuffer)-1]
		}
		return
	}

	if key == "enter" {
		return
	}

	m.PasswordBuffer += key
}

// ClearPasswordBuffer resets the password buffer.
func (m *Model) ClearPasswordBuffer() {
	m.PasswordBuffer = ""
	m.NeedsPassword = false
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

	// Hide install date column and show installed column for remote mode
	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = false
		}
		if col.Type == column.ColInstalled {
			col.Visible = true
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

	// Show install date column and hide installed column again
	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = true
		}
		if col.Type == column.ColInstalled {
			col.Visible = false
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
	return spinnerFrames[m.SpinnerFrame%len(spinnerFrames)]
}

// InitiateInstall sets up the pending installation state.
func (m *Model) InitiateInstall(pkgName string) {
	m.PendingInstall = true
	m.InstallingPkg = pkgName
	m.InstallError = ""
}

// CancelInstall cancels the pending installation.
func (m *Model) CancelInstall() {
	m.PendingInstall = false
	m.InstallingPkg = ""
}

// InstallPackage starts installation of a package.
func (m *Model) InstallPackage(pkgName string, password string) tea.Cmd {
	m.Installing = true
	m.PendingInstall = false
	m.InstallingPkg = pkgName
	m.InstallError = ""

	return tea.Batch(
		m.doInstall(pkgName, password),
		tickSpinner(),
	)
}

// doInstall performs the package installation.
func (m Model) doInstall(pkgName string, password string) tea.Cmd {
	return func() tea.Msg {
		output, err := m.Repo.Install([]string{pkgName}, password)
		return installCompleteMsg{
			pkgName: pkgName,
			output:  output,
			err:     err,
		}
	}
}

// InitiateRemoval sets up the pending removal state.
func (m *Model) InitiateRemoval(pkgName string) {
	m.PendingRemoval = true
	m.RemovingPkg = pkgName
	m.RemoveError = ""
}

// IsRunningAsRoot checks if the program is running with root privileges.
func IsRunningAsRoot() bool {
	return os.Geteuid() == 0
}

// CancelRemoval cancels the pending removal.
func (m *Model) CancelRemoval() {
	m.PendingRemoval = false
	m.RemovingPkg = ""
}

// RemovePackage starts removal of a package.
func (m *Model) RemovePackage(pkgName string, password string) tea.Cmd {
	m.Removing = true
	m.PendingRemoval = false
	m.RemovingPkg = pkgName
	m.RemoveError = ""

	return tea.Batch(
		m.doRemove(pkgName, password),
		tickSpinner(),
	)
}

// doRemove performs the package removal.
func (m Model) doRemove(pkgName string, password string) tea.Cmd {
	return func() tea.Msg {
		output, err := m.Repo.Remove([]string{pkgName}, false, password)
		return removeCompleteMsg{
			pkgName: pkgName,
			output:  output,
			err:     err,
		}
	}
}
