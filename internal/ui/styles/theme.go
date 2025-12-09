package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme configuration.
// This is the single source of truth for all colors in the application.
type Theme struct {
	Name          string
	Accent1       string
	Accent2       string
	Accent3       string
	Accent4       string
	Accent5       string
	Background    string
	BackgroundAlt string // Alternate background for striped rows
	Foreground    string
	Selected      string
	Dimmed        string
	RemoteAccent  string // Accent color for remote mode
	WarningAccent string // Accent color for warning/destructive actions
}

// Styles contains all lipgloss styles used in the application.
// Styles are generated from a Theme using NewStyles().
type Styles struct {
	// Color values
	Accent1       lipgloss.Color
	Accent2       lipgloss.Color
	Accent3       lipgloss.Color
	Accent4       lipgloss.Color
	Accent5       lipgloss.Color
	Background    lipgloss.Color
	BackgroundAlt lipgloss.Color
	Foreground    lipgloss.Color
	Selected      lipgloss.Color
	Dimmed        lipgloss.Color
	RemoteAccent  lipgloss.Color
	WarningAccent lipgloss.Color

	// Component styles
	Header            lipgloss.Style
	Row               lipgloss.Style
	RowAlt            lipgloss.Style
	RowSelected       lipgloss.Style
	Status            lipgloss.Style
	Footer            lipgloss.Style
	StatusBar         lipgloss.Style
	Index             lipgloss.Style
	RemoteStatusBar   lipgloss.Style
	RemoteRowSelected lipgloss.Style
	WarningStatusBar  lipgloss.Style
}

// DefaultTheme uses simple black/white terminal colors for maximum compatibility.
var DefaultTheme = Theme{
	Name:          "default",
	Accent1:       "#FFFFFF",
	Accent2:       "#FFFFFF",
	Accent3:       "#FFFFFF",
	Accent4:       "#FFFFFF",
	Accent5:       "#FFFFFF",
	Background:    "#000000",
	BackgroundAlt: "#000000",
	Foreground:    "#FFFFFF",
	Selected:      "#FFFFFF",
	Dimmed:        "#808080",
	RemoteAccent:  "#FFFFFF",
	WarningAccent: "#FFFFFF",
}

// TokyoNightTheme is a dark theme with Tokyo Night colors.
var TokyoNightTheme = Theme{
	Name:          "tokyo-night",
	Accent1:       "#7aa2f7",
	Accent2:       "#bb9af7",
	Accent3:       "#9ece6a",
	Accent4:       "#e0af68",
	Accent5:       "#f7768e",
	Background:    "#1a1b26",
	BackgroundAlt: "#16161e",
	Foreground:    "#c0caf5",
	Selected:      "#283457",
	Dimmed:        "#565f89",
	RemoteAccent:  "#e0af68",
	WarningAccent: "#ff9e64", // Bright orange for warnings
}

// NewStyles creates a Styles instance from a Theme.
func NewStyles(theme Theme) *Styles {
	s := &Styles{
		// Store color values
		Accent1:       lipgloss.Color(theme.Accent1),
		Accent2:       lipgloss.Color(theme.Accent2),
		Accent3:       lipgloss.Color(theme.Accent3),
		Accent4:       lipgloss.Color(theme.Accent4),
		Accent5:       lipgloss.Color(theme.Accent5),
		Background:    lipgloss.Color(theme.Background),
		BackgroundAlt: lipgloss.Color(theme.BackgroundAlt),
		Foreground:    lipgloss.Color(theme.Foreground),
		Selected:      lipgloss.Color(theme.Selected),
		Dimmed:        lipgloss.Color(theme.Dimmed),
		RemoteAccent:  lipgloss.Color(theme.RemoteAccent),
		WarningAccent: lipgloss.Color(theme.WarningAccent),
	}

	// Build component styles from colors
	s.Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(s.Accent1).
		Background(s.Background)

	s.Row = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Background)

	s.RowAlt = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.BackgroundAlt)

	s.RowSelected = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Selected).
		Bold(true)

	s.Status = lipgloss.NewStyle().
		Foreground(s.Dimmed).
		Background(s.Background)

	s.Footer = lipgloss.NewStyle().
		Foreground(s.Accent2).
		Background(s.Background)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(s.Foreground).
		Background(s.Background).
		Padding(0, 1)

	s.Index = lipgloss.NewStyle().
		Foreground(s.Dimmed)

	s.RemoteStatusBar = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.RemoteAccent).
		Padding(0, 1)

	s.RemoteRowSelected = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.RemoteAccent).
		Bold(true)

	s.WarningStatusBar = lipgloss.NewStyle().
		Foreground(s.Background).
		Background(s.WarningAccent).
		Padding(0, 1).
		Bold(true)

	return s
}

// Default returns the default styles (using DefaultTheme).
func Default() *Styles {
	return NewStyles(DefaultTheme)
}
