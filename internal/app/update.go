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
	m.Viewport.SetRows(rows)
	m.Ready = true

	m.applyCurrentPreset()

	// Ensure SelectedRow is within valid bounds after data reload
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

	return m, nil
}

func (m Model) handleRemoteSearchResult(msg remoteSearchResultMsg) (tea.Model, tea.Cmd) {
	m.RemoteLoading = false

	if msg.err != nil {
		m.RemoteError = msg.err.Error()
		return m, nil
	}

	if len(msg.packages) == 0 {
		m.RemoteError = "No packages found for: " + msg.query
		return m, nil
	}

	m.RemoteError = ""
	rows := domain.PackagesToRows(msg.packages)
	m.Viewport.SetRows(rows)
	m.Viewport.ScrollToTop()

	return m, nil
}

func (m Model) handleRepositoryRefreshed(msg repositoryRefreshedMsg) (tea.Model, tea.Cmd) {
	if msg.shouldReload {
		return m, m.loadPackages
	}
	return m, nil
}

func (m Model) handleSpinnerTick() (tea.Model, tea.Cmd) {
	// Only continue spinning if we're still loading, installing, or removing
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

	// Refresh repository, then reload packages
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

	// Refresh repository, then reload packages
	return m, m.refreshRepository(true)
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle mode-specific input
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
			// Check if we need a password
			if !IsRunningAsRoot() {
				// Prompt for password
				m.EnterPasswordMode()
				return m, nil
			}
			// Already root, proceed with installation
			pkgName := m.InstallingPkg
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
			// Check if we need a password
			if !IsRunningAsRoot() {
				// Prompt for password
				m.EnterPasswordMode()
				return m, nil
			}
			// Already root, proceed with removal
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
		// If output palette is active, dismiss it and clear state
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
		// Otherwise toggle detail panel
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
		m.NextPreset()
	case "i":
		if m.ViewMode == ViewRemote && m.ShowDetailPanel && m.Viewport.SelectedRow >= 0 && m.Viewport.SelectedRow < len(m.Viewport.VisibleRows) {
			// In detail view in remote mode, pressing 'i' initiates install
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
		// Cancel password entry and operation
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
		// Submit password and proceed with operation
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
		// Add character to password buffer
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
		if !m.SetPreset(result.PresetChange) {
			m.RemoteError = "Invalid preset"
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
	// Reserve 1 line for header + 1 line for status bar
	m.Viewport.Height = msg.Height - 2
	// Ensure selection is still visible after resize
	m.Viewport.EnsureSelectionVisible()
	return m, nil
}

func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.Mode != ModeNormal {
		return m, nil
	}

	// Handle mouse wheel scroll
	switch msg.Type {
	case tea.MouseWheelUp:
		m.Viewport.SelectPrev()
		return m, nil
	case tea.MouseWheelDown:
		m.Viewport.SelectNext()
		return m, nil
	}

	// Handle mouse click
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
// It filters out control keys and special navigation keys that shouldn't appear in search text.
func isValidSearchChar(key string) bool {
	// Single character printable keys
	if len(key) == 1 {
		// Allow letters, numbers, and common punctuation
		return true
	}

	// Allow specific multi-character keys that are printable
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
		"\\": true,
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
