package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// Update handles incoming messages and updates the model (Bubble Tea interface).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case packagesLoadedMsg:
		return m.handlePackagesLoaded(msg)
	case remoteSearchResultMsg:
		return m.handleRemoteSearchResult(msg)
	case aurSearchResultMsg:
		return m.handleAURSearchResult(msg)
	case aurInfoResultMsg:
		return m.handleAURInfoResult(msg)
	case aurInstallCompleteMsg:
		return m.handleAURInstallComplete(msg)
	case repositoryRefreshedMsg:
		return m.handleRepositoryRefreshed(msg)
	case spinnerTickMsg:
		return m.handleSpinnerTick()
	case installCompleteMsg:
		return m.handleInstallComplete(msg)
	case removeCompleteMsg:
		return m.handleRemoveComplete(msg)
	case commandResultMsg:
		return m.handleCommandResult(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.MouseMsg:
		return m.handleMouseEvent(msg)
	}

	return m, nil
}

func (m Model) handlePackagesLoaded(msg packagesLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.Error = msg.err.Error()
		return m, nil
	}

	rows := domain.PackagesToRows(msg.packages)

	// In remote mode, update cached local rows and re-run search to refresh install status
	if m.ViewMode == ViewRemote {
		m.LocalRows = rows
		m.Ready = true

		var cmds []tea.Cmd
		cmds = append(cmds, m.doRemoteSearch(m.RemoteQuery))
		if m.AUREnabled {
			cmds = append(cmds, m.doAURSearch(m.RemoteQuery))
		}
		m.syncSearchDone = false
		m.aurSearchDone = false
		m.syncSearchResult = nil
		m.aurSearchResult = nil
		return m, tea.Batch(cmds...)
	}

	m.Viewport.SetRows(rows)
	m.Ready = true

	presetCmd := m.applyCurrentPreset()

	if m.Viewport.SelectedRow >= len(m.Viewport.VisibleRows) {
		if len(m.Viewport.VisibleRows) > 0 {
			m.Viewport.SelectedRow = len(m.Viewport.VisibleRows) - 1
		} else {
			m.Viewport.SelectedRow = 0
		}
	}
	if m.Viewport.SelectedRow < 0 && len(m.Viewport.VisibleRows) > 0 {
		m.Viewport.SelectedRow = 0
	}

	// Look up which foreign packages are AUR packages so Repo column shows "aur"
	var cmds []tea.Cmd
	if presetCmd != nil {
		cmds = append(cmds, presetCmd)
	}
	if m.AUREnabled && m.AURClient != nil {
		cmds = append(cmds, m.doAURInfoLookup())
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m Model) handleRemoteSearchResult(msg remoteSearchResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.syncSearchResult = nil
	} else {
		m.syncSearchResult = msg.packages
	}
	m.syncSearchDone = true

	if !m.AUREnabled || m.aurSearchDone {
		return m.finalizeSearchResults(msg.query, msg.err)
	}

	return m, nil
}

func (m Model) handleAURSearchResult(msg aurSearchResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.aurSearchResult = nil
	} else {
		m.aurSearchResult = msg.packages
	}
	m.aurSearchDone = true

	if m.syncSearchDone {
		var syncErr error
		if m.syncSearchResult == nil && m.syncSearchDone {
		}
		return m.finalizeSearchResults(msg.query, syncErr)
	}

	return m, nil
}

func (m Model) finalizeSearchResults(query string, syncErr error) (tea.Model, tea.Cmd) {
	m.RemoteLoading = false

	syncPkgs := m.syncSearchResult
	aurPkgs := m.aurSearchResult

	var merged []*domain.Package
	if len(syncPkgs) > 0 && len(aurPkgs) > 0 {
		merged = mergeSearchResults(syncPkgs, aurPkgs)
	} else if len(syncPkgs) > 0 {
		merged = syncPkgs
	} else if len(aurPkgs) > 0 {
		merged = aurPkgs
	}

	if len(merged) == 0 {
		if syncErr != nil {
			m.RemoteError = syncErr.Error()
		} else {
			m.RemoteError = "No packages found for: " + query
		}
		return m, nil
	}

	m.RemoteError = ""
	rows := domain.PackagesToRows(merged)
	m.Viewport.SetRows(rows)
	m.Viewport.ScrollToTop()

	return m, nil
}

func (m Model) handleAURInfoResult(msg aurInfoResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		return m, nil
	}

	for _, row := range m.Viewport.AllRows {
		if row.Package != nil && msg.found[row.Package.Name] {
			row.Package.IsAUR = true
			row.Package.Repository = "aur"
			row.Cells[column.ColRepo] = "aur"
		}
	}

	_ = m.applyCurrentPreset()

	return m, nil
}

func (m Model) handleAURInstallComplete(msg aurInstallCompleteMsg) (tea.Model, tea.Cmd) {
	m.Installing = false
	m.InstallingPkg = ""

	m.InstallError = ""
	if msg.err != nil {
		m.InstallError = msg.err.Error()
		m.InstallOutput = "AUR installation failed"
	} else {
		m.InstallOutput = "AUR package installed successfully"
	}

	cmds := []tea.Cmd{m.refreshRepository(true)}

	if m.ViewMode == ViewRemote && m.RemoteQuery != "" {
		cmds = append(cmds, m.doRemoteSearch(m.RemoteQuery))
		if m.AUREnabled {
			cmds = append(cmds, m.doAURSearch(m.RemoteQuery))
		}
		m.syncSearchDone = false
		m.aurSearchDone = false
		m.syncSearchResult = nil
		m.aurSearchResult = nil
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleRepositoryRefreshed(msg repositoryRefreshedMsg) (tea.Model, tea.Cmd) {
	if msg.shouldReload {
		return m, m.loadPackages
	}
	return m, nil
}

func (m Model) handleSpinnerTick() (tea.Model, tea.Cmd) {
	if !m.RemoteLoading && !m.Installing && !m.Removing {
		return m, nil
	}

	m.SpinnerFrame++
	return m, tickSpinner()
}

func (m Model) handleInstallComplete(msg installCompleteMsg) (tea.Model, tea.Cmd) {
	m.Installing = false
	m.InstallingPkg = ""

	m.InstallError = ""
	m.InstallOutput = msg.output

	if msg.err != nil {
		m.InstallError = msg.err.Error()
	}

	return m, m.refreshRepository(true)
}

func (m Model) handleRemoveComplete(msg removeCompleteMsg) (tea.Model, tea.Cmd) {
	m.Removing = false
	m.RemovingPkg = ""

	m.RemoveError = ""
	m.RemoveOutput = msg.output

	if msg.err != nil {
		m.RemoveError = msg.err.Error()
	}

	return m, m.refreshRepository(true)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch m.Mode {
	case ModeCommand:
		return m.handleCommandModeInput(key)
	case ModeFilter:
		return m.handleFilterModeInput(key)
	case ModePassword:
		return m.handlePasswordModeInput(key)
	case ModeNormal:
		return m.handleNormalModeInput(key)
	}

	return m, nil
}

func (m Model) handleNormalModeInput(key string) (tea.Model, tea.Cmd) {
	if m.PendingInstall {
		switch key {
		case "enter":
			pkgName := m.InstallingPkg

			if m.isSelectedPackageAUR() {
				m.PendingInstall = false
				m.Installing = true
				return m, m.doAURInstall(pkgName)
			}

			if !IsRunningAsRoot() {
				m.EnterPasswordMode()
				return m, nil
			}
			return m, m.InstallPackage(pkgName, "")
		case "esc", "ctrl+c":
			m.CancelInstall()
			return m, nil
		default:
			return m, nil
		}
	}

	if m.PendingRemoval {
		switch key {
		case "enter":
			if !IsRunningAsRoot() {
				m.EnterPasswordMode()
				return m, nil
			}
			pkgName := m.RemovingPkg
			return m, m.RemovePackage(pkgName, "")
		case "esc", "ctrl+c":
			m.CancelRemoval()
			return m, nil
		default:
			return m, nil
		}
	}

	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		if m.RemoveOutput != "" {
			m.RemoveOutput = ""
			m.RemoveError = ""
			return m, nil
		}
		if m.InstallOutput != "" {
			m.InstallOutput = ""
			m.InstallError = ""
			return m, nil
		}
		m.ShowDetailPanel = !m.ShowDetailPanel
	case "esc":
		if m.ShowDetailPanel {
			m.ShowDetailPanel = false
			return m, nil
		}
		if m.RemoveOutput != "" {
			m.RemoveOutput = ""
			return m, nil
		}
		if m.InstallOutput != "" {
			m.InstallOutput = ""
			return m, nil
		}
		if m.ViewMode == ViewRemote {
			m.ExitRemoteMode()
			return m, nil
		}
		if m.Viewport.Filter.Active {
			m.Viewport.ClearFilter()
		}
	case ":":
		m.EnterCommandMode()
	case "/":
		m.EnterFilterMode()
	case "up", "k":
		m.Viewport.SelectPrev()
	case "down", "j":
		m.Viewport.SelectNext()
	case "ctrl+u":
		m.Viewport.PageUp()
	case "ctrl+d":
		m.Viewport.PageDown()
	case "home", "g":
		if key == "g" {
			m.Viewport.ScrollToTop()
		} else {
			m.Viewport.ScrollToTop()
		}
	case "G", "end":
		m.Viewport.ScrollToBottom()
	case "left", "h":
		m.Viewport.PrevColumn()
	case "right", "l":
		m.Viewport.NextColumn()
	case " ", "space":
		m.Viewport.ToggleSortCurrentColumn()
	case "tab":
		cmd := m.NextPreset()
		return m, cmd
	case "i":
		if m.ViewMode == ViewRemote && m.ShowDetailPanel && m.Viewport.SelectedRow >= 0 && m.Viewport.SelectedRow < len(m.Viewport.VisibleRows) {
			selectedRow := m.Viewport.VisibleRows[m.Viewport.SelectedRow]
			pkgName := selectedRow.Cells[column.ColName]
			m.InitiateInstall(pkgName)
		}
	}

	return m, nil
}

func (m Model) handleCommandModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.ExitMode()
		return m, nil
	case "enter":
		cmd := m.executeCommand()
		m.ExitMode()
		return m, cmd
	default:
		m.WriteToBuffer(key)
		return m, nil
	}
}

func (m Model) handleFilterModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.ExitMode()
		m.Viewport.ClearFilter()
	case "enter":
		m.ExitMode()
	case "backspace":
		m.WriteToBuffer(key)
		filterTerm := m.GetBufferContent()
		m.Viewport.ApplyFilter(filterTerm)
	case "left", "right":
	case "up", "down", "ctrl+a", "ctrl+e", "ctrl+k", "ctrl+u", "ctrl+w":
	default:
		if isValidSearchChar(key) {
			m.WriteToBuffer(key)
			filterTerm := m.GetBufferContent()
			m.Viewport.ApplyFilter(filterTerm)
		}
	}

	return m, nil
}

func (m Model) handlePasswordModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		m.ExitMode()
		m.ClearPasswordBuffer()
		if m.PendingInstall {
			m.CancelInstall()
		}
		if m.PendingRemoval {
			m.CancelRemoval()
		}
		return m, nil
	case "enter":
		password := m.PasswordBuffer
		m.ExitMode()
		m.ClearPasswordBuffer()

		if m.PendingInstall {
			pkgName := m.InstallingPkg
			return m, m.InstallPackage(pkgName, password)
		}
		if m.PendingRemoval {
			pkgName := m.RemovingPkg
			return m, m.RemovePackage(pkgName, password)
		}
		return m, nil
	default:
		m.WriteToPasswordBuffer(key)
		return m, nil
	}
}

func (m *Model) executeCommand() tea.Cmd {
	commandStr := m.GetBufferContent()
	return executeCommandMsg(commandStr)
}

func (m Model) handleCommandResult(msg commandResultMsg) (tea.Model, tea.Cmd) {
	result := msg.Result

	if result.Error != "" {
		m.RemoteError = result.Error
		return m, nil
	}

	if result.Quit {
		return m, tea.Quit
	}

	if result.RemoteSearch != "" {
		return m, m.EnterRemoteMode(result.RemoteSearch)
	}

	if result.PresetChange != "" {
		ok, cmd := m.SetPreset(result.PresetChange)
		if !ok {
			m.RemoteError = "Invalid preset"
		} else if cmd != nil {
			return m, cmd
		}
	}

	if result.GoToLine >= 0 {
		m.Viewport.ScrollToLine(result.GoToLine)
	}

	if result.ScrollTop {
		m.Viewport.ScrollToTop()
	}

	if result.ScrollEnd {
		m.Viewport.ScrollToBottom()
	}

	if result.InstallPackage {
		if m.ViewMode == ViewRemote && m.Viewport.SelectedRow >= 0 && m.Viewport.SelectedRow < len(m.Viewport.VisibleRows) {
			selectedRow := m.Viewport.VisibleRows[m.Viewport.SelectedRow]
			pkgName := selectedRow.Cells[column.ColName]
			m.InitiateInstall(pkgName)
		} else if m.ViewMode != ViewRemote {
			m.RemoteError = "Install command only works in search mode"
		} else {
			m.RemoteError = "No package selected"
		}
	}

	if result.RemovePackage {
		if m.ViewMode == ViewLocal && m.Viewport.SelectedRow >= 0 && m.Viewport.SelectedRow < len(m.Viewport.VisibleRows) {
			selectedRow := m.Viewport.VisibleRows[m.Viewport.SelectedRow]
			pkgName := selectedRow.Cells[column.ColName]
			m.InitiateRemoval(pkgName)
		} else if m.ViewMode != ViewLocal {
			m.RemoteError = "Remove command only works in local mode"
		} else {
			m.RemoteError = "No package selected"
		}
	}

	if result.ThemeName != "" {
		theme, err := styles.LoadTheme(result.ThemeName)
		if err != nil {
			m.RemoteError = fmt.Sprintf("Error loading theme: %v", err)
			return m, nil
		}
		styles.ApplyTheme(theme)
		m.RemoteError = fmt.Sprintf("Theme changed to: %s", result.ThemeName)
	}

	return m, nil
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height
	m.Viewport.Height = msg.Height - 2
	m.Viewport.EnsureSelectionVisible()
	return m, nil
}

func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.Mode != ModeNormal {
		return m, nil
	}

	switch msg.Type {
	case tea.MouseWheelUp:
		m.Viewport.SelectPrev()
		return m, nil
	case tea.MouseWheelDown:
		m.Viewport.SelectNext()
		return m, nil
	}

	if msg.Type != tea.MouseLeft {
		return m, nil
	}

	if msg.Y <= 0 {
		return m, nil
	}

	clickedRow := msg.Y - 1 + m.Viewport.Offset

	if clickedRow >= 0 && clickedRow < len(m.Viewport.VisibleRows) {
		m.Viewport.SelectedRow = clickedRow
	}

	return m, nil
}

// isValidSearchChar checks if a key is a valid printable character for search.
func isValidSearchChar(key string) bool {
	if len(key) == 1 {
		return true
	}

	validMultiChars := map[string]bool{
		"space": true,
		"tab":   true,
		"-":     true,
		"_":     true,
		".":     true,
		",":     true,
		"+":     true,
		"=":     true,
		"[":     true,
		"]":     true,
		"{":     true,
		"}":     true,
		"(":     true,
		")":     true,
		"@":     true,
		"#":     true,
		"$":     true,
		"%":     true,
		"^":     true,
		"&":     true,
		"*":     true,
		"/":     true,
		"|":     true,
		"\\":    true,
		"?":     true,
		"!":     true,
		"`":     true,
		"~":     true,
		"'":     true,
		"\"":    true,
		";":     true,
		":":     true,
		"<":     true,
		">":     true,
	}

	_, isValid := validMultiChars[key]
	return isValid
}
