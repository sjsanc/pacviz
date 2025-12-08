package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/command"
	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// Update handles incoming messages and updates the model (Bubble Tea interface).
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case packagesLoadedMsg:
		return m.handlePackagesLoaded(msg)
	case command.CommandResultMsg:
		return m.handleCommandResult(msg)
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

	// Apply the initial preset filter (Explicit)
	m.applyCurrentPreset()

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Handle mode-specific input
	switch m.Mode {
	case ModeCommand:
		return m.handleCommandModeInput(key)
	case ModeFilter:
		return m.handleFilterModeInput(key)
	case ModeNormal:
		return m.handleNormalModeInput(key)
	}

	return m, nil
}

func (m Model) handleNormalModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "ctrl+c", "q":
		return m, tea.Quit

	// Clear filter on escape if filter is active
	case "esc":
		if m.Viewport.Filter.Active {
			m.Viewport.ClearFilter()
		}

	// Mode switching
	case ":":
		m.EnterCommandMode()
	case "/":
		m.EnterFilterMode()

	// Vertical navigation
	case "up", "k":
		m.Viewport.SelectPrev()
	case "down", "j":
		m.Viewport.SelectNext()

	// Page navigation
	case "ctrl+u":
		m.Viewport.PageUp()
	case "ctrl+d":
		m.Viewport.PageDown()

	// Jump navigation
	case "home", "g":
		if key == "g" {
			// For vim-style gg, we'd need to track previous key
			// For now, just treat single 'g' as go to top
			m.Viewport.ScrollToTop()
		} else {
			m.Viewport.ScrollToTop()
		}
	case "G", "end":
		m.Viewport.ScrollToBottom()

	// Horizontal navigation
	case "left", "h":
		m.Viewport.PrevColumn()
	case "right", "l":
		m.Viewport.NextColumn()

	// Sort toggle
	case " ", "space":
		m.Viewport.ToggleSortCurrentColumn()

	// Preset switching
	case "tab":
		m.NextPreset()

	// TODO: Info view (enter, i)
	}

	return m, nil
}

func (m Model) handleCommandModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		// Cancel command mode
		m.ExitMode()
		return m, nil
	case "enter":
		// Execute command
		cmd := m.executeCommand()
		m.ExitMode()
		return m, cmd
	default:
		// Write to buffer
		m.WriteToBuffer(key)
		return m, nil
	}
}

func (m Model) handleFilterModeInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "esc":
		// Cancel filter mode and reset filter
		m.ExitMode()
		m.Viewport.ClearFilter()
	case "enter":
		// Accept filter and hide buffer (keep filter active)
		m.ExitMode()
	case "backspace":
		// Handle backspace
		m.WriteToBuffer(key)
		filterTerm := m.GetBufferContent()
		m.Viewport.ApplyFilter(filterTerm)
	case "left", "right":
		// Allow cursor movement keys but don't process them for now
		// (buffer editing with cursor movement can be added later if needed)
	case "up", "down", "ctrl+a", "ctrl+e", "ctrl+k", "ctrl+u", "ctrl+w":
		// Ignore control/navigation keys that shouldn't be in search
	default:
		// Only accept printable characters (exclude special/control keys)
		if isValidSearchChar(key) {
			m.WriteToBuffer(key)
			filterTerm := m.GetBufferContent()
			m.Viewport.ApplyFilter(filterTerm)
		}
	}

	return m, nil
}

func (m *Model) executeCommand() tea.Cmd {
	commandStr := m.GetBufferContent()
	return command.ExecuteCommandMsg(commandStr)
}

func (m Model) handleCommandResult(msg command.CommandResultMsg) (tea.Model, tea.Cmd) {
	result := msg.Result

	// Handle errors
	if result.Error != "" {
		m.Error = nil // Clear any previous error for now (we could show it in status)
		// TODO: Show error in status bar instead of setting m.Error
		return m, nil
	}

	// Handle quit
	if result.Quit {
		return m, tea.Quit
	}

	// Handle go to line
	if result.GoToLine >= 0 {
		m.Viewport.ScrollToLine(result.GoToLine)
	}

	// Handle scroll to top
	if result.ScrollTop {
		m.Viewport.ScrollToTop()
	}

	// Handle scroll to end
	if result.ScrollEnd {
		m.Viewport.ScrollToBottom()
	}

	return m, nil
}

func (m Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.Width = msg.Width
	m.Height = msg.Height
	// Reserve 1 line for header + 1 line for status bar
	m.Viewport.Height = msg.Height - 2
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
