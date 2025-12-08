package config

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Config represents the application configuration.
type Config struct {
	UI          UIConfig
	Columns     ColumnConfig
	Performance PerformanceConfig
	Keybindings KeybindingsConfig
	Theme       ThemeConfig
	Pacman      PacmanConfig
}

// UIConfig contains UI-related settings.
type UIConfig struct {
	Theme         string
	ShowScrollbar bool
	MouseSupport  bool
	VimBindings   bool
}

// ColumnConfig contains column display settings.
type ColumnConfig struct {
	DefaultVisible []column.Type
	Widths         map[column.Type]ColumnWidthConfig
}

// ColumnWidthConfig defines width settings for a column.
type ColumnWidthConfig struct {
	Type     string // "auto", "fixed", "percent"
	MinWidth int
	MaxWidth int
	Width    int // for fixed, or percent value
}

// PerformanceConfig contains performance-related settings.
type PerformanceConfig struct {
	CacheTTL        int  // seconds
	DebounceFilter  int  // milliseconds
	AsyncSearch     bool
}

// KeybindingsConfig contains keybinding settings.
type KeybindingsConfig struct {
	Quit   []string
	Filter []string
	Command []string
	Info   []string
	Help   []string
}

// ThemeConfig contains theme color settings.
type ThemeConfig struct {
	Dark  ThemeColors
	Light ThemeColors
}

// ThemeColors defines a color scheme.
type ThemeColors struct {
	Accent1    string
	Accent2    string
	Background string
	Foreground string
	Selected   string
}

// PacmanConfig contains pacman-specific settings.
type PacmanConfig struct {
	DBPath string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		UI: UIConfig{
			Theme:         "dark",
			ShowScrollbar: true,
			MouseSupport:  false,
			VimBindings:   true,
		},
		Columns: ColumnConfig{
			DefaultVisible: []column.Type{
				column.ColName,
				column.ColVersion,
				column.ColSize,
				column.ColInstallDate,
				column.ColDescription,
			},
			Widths: map[column.Type]ColumnWidthConfig{
				column.ColName: {
					Type:     "fixed",
					Width: 10,
				},
				column.ColVersion: {
					Type:  "fixed",
					Width: 10,
				},
				column.ColSize: {
					Type:  "fixed",
					Width: 10,
				},
				column.ColInstallDate: {
					Type:  "fixed",
					Width: 10,
				},
				column.ColDescription: {
					Type:     "percent",
					Width:    40,
					MinWidth: 20,
				},
			},
		},
		Performance: PerformanceConfig{
			CacheTTL:       300,
			DebounceFilter: 150,
			AsyncSearch:    true,
		},
		Keybindings: KeybindingsConfig{
			Quit:    []string{"q", "ctrl+c"},
			Filter:  []string{"/"},
			Command: []string{":"},
			Info:    []string{"enter", "i"},
			Help:    []string{"?"},
		},
		Theme: ThemeConfig{
			Dark: ThemeColors{
				Accent1:    "#7aa2f7",
				Accent2:    "#bb9af7",
				Background: "#1a1b26",
				Foreground: "#c0caf5",
				Selected:   "#283457",
			},
		},
		Pacman: PacmanConfig{
			DBPath: "/var/lib/pacman",
		},
	}
}
