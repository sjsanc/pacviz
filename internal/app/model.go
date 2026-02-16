package app

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/aur"
	"github.com/sjsanc/pacviz/v3/internal/command"
	"github.com/sjsanc/pacviz/v3/internal/config"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/repository"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/viewport"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type InputMode int

const (
	ModeNormal InputMode = iota
	ModeCommand
	ModeFilter
	ModePassword
)

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

	Mode   InputMode
	Buffer string

	Presets       []domain.Preset
	CurrentPreset int

	ShowDetailPanel bool

	ViewMode      ViewMode
	RemoteQuery   string
	RemoteLoading bool
	RemoteError   string
	LocalRows     []*domain.Row
	SpinnerFrame  int

	PendingInstall bool
	Installing     bool
	InstallingPkg  string
	InstallError   string
	InstallOutput  string

	PendingRemoval bool
	RemovingPkg    string
	Removing       bool
	RemoveError    string
	RemoveOutput   string

	PasswordBuffer string
	NeedsPassword  bool

	AURClient        *aur.Client
	AURHelper        *aur.HelperConfig
	AUREnabled       bool
	syncSearchResult []*domain.Package
	syncSearchDone   bool
	aurSearchResult  []*domain.Package
	aurSearchDone    bool
}

type packagesLoadedMsg struct {
	packages []*domain.Package
	err      error
}

type remoteSearchResultMsg struct {
	packages []*domain.Package
	query    string
	err      error
}

type spinnerTickMsg struct{}

type installCompleteMsg struct {
	pkgName string
	output  string
	err     error
}

type removeCompleteMsg struct {
	pkgName string
	output  string
	err     error
}

type repositoryRefreshedMsg struct {
	shouldReload bool
}

type aurSearchResultMsg struct {
	packages []*domain.Package
	query    string
	err      error
}

type aurInfoResultMsg struct {
	found map[string]bool
	err   error
}

type aurInstallCompleteMsg struct {
	err error
}

type commandResultMsg struct {
	Result command.ExecuteResult
}

func executeCommandMsg(commandStr string) tea.Cmd {
	return func() tea.Msg {
		result := command.Execute(commandStr)
		return commandResultMsg{Result: result}
	}
}

// NewModel creates a new application model.
func NewModel(cfg *config.Config) *Model {
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

	m := &Model{
		Viewport:      viewport.New(),
		Repo:          repo,
		Ready:         false,
		Presets:       domain.DefaultPresets(),
		CurrentPreset: 0,
	}

	if !cfg.AUR.Disabled {
		timeout := time.Duration(cfg.AUR.Timeout) * time.Second
		cacheTTL := time.Duration(cfg.AUR.CacheTTL) * time.Second
		m.AURClient = aur.NewClient(timeout, cacheTTL)
		m.AURHelper = aur.DetectHelper(cfg.AUR.Helper)
		m.AUREnabled = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return m.loadPackages
}

func (m Model) loadPackages() tea.Msg {
	packages, err := m.Repo.GetInstalled()
	return packagesLoadedMsg{packages: packages, err: err}
}

func (m Model) refreshRepository(shouldReload bool) tea.Cmd {
	return func() tea.Msg {
		err := m.Repo.Refresh()
		if err != nil {
			log.Printf("Failed to refresh repository: %v", err)
		}
		return repositoryRefreshedMsg{shouldReload: shouldReload}
	}
}

func (m *Model) EnterCommandMode() {
	m.Mode = ModeCommand
	m.Buffer = ":"
}

func (m *Model) EnterFilterMode() {
	m.Mode = ModeFilter
	m.Buffer = "/"
}

func (m *Model) ExitMode() {
	m.Mode = ModeNormal
	m.Buffer = ""
}

func (m *Model) WriteToBuffer(key string) {
	if key == "backspace" {
		if len(m.Buffer) > 1 {
			m.Buffer = m.Buffer[:len(m.Buffer)-1]
		}
		return
	}

	if key == "enter" {
		return
	}

	m.Buffer += key
}

func (m *Model) ClearBuffer() {
	m.Buffer = ""
}

func (m *Model) GetBufferContent() string {
	if len(m.Buffer) > 1 {
		return m.Buffer[1:]
	}
	return ""
}

func (m *Model) EnterPasswordMode() {
	m.Mode = ModePassword
	m.PasswordBuffer = ""
	m.NeedsPassword = true
}

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

func (m *Model) ClearPasswordBuffer() {
	m.PasswordBuffer = ""
	m.NeedsPassword = false
}

// NextPreset cycles to the next preset, resetting sort and filters.
func (m *Model) NextPreset() tea.Cmd {
	m.CurrentPreset = (m.CurrentPreset + 1) % len(m.Presets)
	return m.applyCurrentPreset()
}

func (m *Model) SetPreset(presetName string) (bool, tea.Cmd) {
	for i, preset := range m.Presets {
		if string(preset.Type) == presetName {
			m.CurrentPreset = i
			cmd := m.applyCurrentPreset()
			return true, cmd
		}
	}
	return false, nil
}

// applyCurrentPreset applies the current preset filter and resets sort/filters.
func (m *Model) applyCurrentPreset() tea.Cmd {
	m.Viewport.SortColumn = column.ColName
	m.Viewport.SortReverse = false
	m.Viewport.ClearFilter()

	preset := m.Presets[m.CurrentPreset]
	m.Viewport.ApplyPresetFilter(preset.Filter)

	// If switching to AUR preset, do a lazy Info() lookup
	if preset.Type == domain.PresetAUR && m.AUREnabled && m.AURClient != nil {
		return m.doAURInfoLookup()
	}

	return nil
}

func (m *Model) EnterRemoteMode(query string) tea.Cmd {
	m.LocalRows = m.Viewport.AllRows
	m.ViewMode = ViewRemote
	m.RemoteQuery = query
	m.RemoteLoading = true
	m.RemoteError = ""
	m.SpinnerFrame = 0

	m.syncSearchResult = nil
	m.syncSearchDone = false
	m.aurSearchResult = nil
	m.aurSearchDone = false

	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = false
		}
		if col.Type == column.ColInstalled {
			col.Visible = true
		}
	}

	cmds := []tea.Cmd{
		m.doRemoteSearch(query),
		tickSpinner(),
	}

	if m.AUREnabled {
		cmds = append(cmds, m.doAURSearch(query))
	}

	return tea.Batch(cmds...)
}

func (m *Model) ExitRemoteMode() {
	m.ViewMode = ViewLocal
	m.RemoteQuery = ""
	m.RemoteLoading = false
	m.RemoteError = ""

	for _, col := range m.Viewport.Columns {
		if col.Type == column.ColInstallDate {
			col.Visible = true
		}
		if col.Type == column.ColInstalled {
			col.Visible = false
		}
	}

	if m.LocalRows != nil {
		m.Viewport.SetRows(m.LocalRows)
		m.LocalRows = nil
	}

	m.Viewport.SortColumn = column.ColName
	m.Viewport.SortReverse = false
	m.Viewport.ClearFilter()
	m.Viewport.ScrollToTop()
	m.CurrentPreset = 0
	_ = m.applyCurrentPreset()
}

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

func (m Model) doAURSearch(query string) tea.Cmd {
	return func() tea.Msg {
		packages, err := m.AURClient.Search(query)
		return aurSearchResultMsg{
			packages: packages,
			query:    query,
			err:      err,
		}
	}
}

// mergeSearchResults merges sync and AUR results. Sync wins on name collision.
func mergeSearchResults(syncPkgs, aurPkgs []*domain.Package) []*domain.Package {
	seen := make(map[string]bool, len(syncPkgs))
	for _, pkg := range syncPkgs {
		seen[pkg.Name] = true
	}

	merged := make([]*domain.Package, len(syncPkgs))
	copy(merged, syncPkgs)

	for _, pkg := range aurPkgs {
		if !seen[pkg.Name] {
			merged = append(merged, pkg)
		}
	}

	return merged
}

func (m Model) doAURInfoLookup() tea.Cmd {
	return func() tea.Msg {
		var foreignNames []string
		for _, row := range m.Viewport.AllRows {
			if row.Package != nil && row.Package.IsForeign {
				foreignNames = append(foreignNames, row.Package.Name)
			}
		}

		if len(foreignNames) == 0 {
			return aurInfoResultMsg{found: map[string]bool{}}
		}

		found, err := m.AURClient.Info(foreignNames)
		return aurInfoResultMsg{found: found, err: err}
	}
}

func tickSpinner() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

func (m Model) GetSpinner() string {
	return spinnerFrames[m.SpinnerFrame%len(spinnerFrames)]
}

func (m *Model) InitiateInstall(pkgName string) {
	m.PendingInstall = true
	m.InstallingPkg = pkgName
	m.InstallError = ""
}

func (m *Model) CancelInstall() {
	m.PendingInstall = false
	m.InstallingPkg = ""
}

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

func (m *Model) InitiateRemoval(pkgName string) {
	m.PendingRemoval = true
	m.RemovingPkg = pkgName
	m.RemoveError = ""
}

func (m Model) isSelectedPackageAUR() bool {
	if m.Viewport.SelectedRow < 0 || m.Viewport.SelectedRow >= len(m.Viewport.VisibleRows) {
		return false
	}
	row := m.Viewport.VisibleRows[m.Viewport.SelectedRow]
	return row.Package != nil && row.Package.IsAUR
}

func IsRunningAsRoot() bool {
	return os.Geteuid() == 0
}

func (m *Model) CancelRemoval() {
	m.PendingRemoval = false
	m.RemovingPkg = ""
}

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

// doAURInstall suspends the TUI and runs the AUR helper interactively.
func (m Model) doAURInstall(pkgName string) tea.Cmd {
	if m.AURHelper == nil {
		return func() tea.Msg {
			return aurInstallCompleteMsg{err: fmt.Errorf("no AUR helper found (install yay, paru, pikaur, or trizen)")}
		}
	}

	cmd := aur.InstallCmd(m.AURHelper, []string{pkgName})
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return aurInstallCompleteMsg{err: err}
	})
}

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
