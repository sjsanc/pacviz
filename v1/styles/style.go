package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Fg      = lipgloss.Color("#E6E1CF")
	FgMuted = lipgloss.Color("#6C7891")
	Bg1     = lipgloss.Color("#0A0E14")
	Bg2     = lipgloss.Color("#1F2430")
	BgAlt   = lipgloss.Color("#292D38")
	Accent1 = lipgloss.Color("#D95757")
	Accent2 = lipgloss.Color("#39BAE6")
	Accent3 = lipgloss.Color("#D2A6FF")
	Accent4 = lipgloss.Color("#7FD962")
	Accent5 = lipgloss.Color("#FF8F40")

	// Table
	Cell        = lipgloss.NewStyle().Background(Bg1).Foreground(Fg)
	CellVersion = lipgloss.NewStyle().Foreground(Accent4).Bold(true)
	Row         = lipgloss.NewStyle().Background(Bg1).Foreground(Fg)
	RowAlt      = lipgloss.NewStyle().Background(Bg2).Foreground(Fg)
	RowSelected = lipgloss.NewStyle().Background(Accent1).Foreground(Bg1).Bold(true)
	Header      = lipgloss.NewStyle().Background(BgAlt).Foreground(Accent2).Bold(true)
	Footer      = lipgloss.NewStyle().Background(BgAlt)

	// Infobox
	Infobox = lipgloss.NewStyle().Background(Bg1).Foreground(Fg).Padding(1, 2)
	PkgInfo = lipgloss.NewStyle().Foreground(Accent2).Bold(true)
	Label   = lipgloss.NewStyle().Foreground(Accent3).Bold(true)
	Value   = lipgloss.NewStyle().Foreground(Fg)
)
