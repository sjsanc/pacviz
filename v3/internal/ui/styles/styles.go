package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	Accent1      = lipgloss.Color("#7aa2f7")
	Accent2      = lipgloss.Color("#bb9af7")
	Accent3      = lipgloss.Color("#9ece6a")
	Accent4      = lipgloss.Color("#e0af68")
	Accent5      = lipgloss.Color("#f7768e")
	Background   = lipgloss.Color("#1a1b26")
	Foreground   = lipgloss.Color("#c0caf5")
	Selected     = lipgloss.Color("#283457")
	Dimmed       = lipgloss.Color("#565f89")
	RemoteAccent = lipgloss.Color("#e0af68") // Yellow accent for remote mode
)

// Styles
var (
	Header = lipgloss.NewStyle().
		Bold(true).
		Foreground(Accent1).
		Background(Background)

	Row = lipgloss.NewStyle().
		Foreground(Foreground).
		Background(Background)

	RowAlt = lipgloss.NewStyle().
		Foreground(Foreground).
		Background(lipgloss.Color("#16161e"))

	RowSelected = lipgloss.NewStyle().
		Foreground(Foreground).
		Background(Selected).
		Bold(true)

	Status = lipgloss.NewStyle().
		Foreground(Dimmed).
		Background(Background)

	Footer = lipgloss.NewStyle().
		Foreground(Accent2).
		Background(Background)

	StatusBar = lipgloss.NewStyle().
		Foreground(Foreground).
		Background(Background).
		Padding(0, 1)

	Index = lipgloss.NewStyle().
		Foreground(Dimmed)

	// Remote mode styles
	RemoteStatusBar = lipgloss.NewStyle().
			Foreground(Background).
			Background(RemoteAccent).
			Padding(0, 1)

	RemoteRowSelected = lipgloss.NewStyle().
				Foreground(Background).
				Background(RemoteAccent).
				Bold(true)
)
