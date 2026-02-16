package config

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Config represents the application configuration.
// Note: Theme configuration has been moved to internal/ui/styles package.
// Use styles.Current to access the active theme, and styles.NewStyles() to create themes.
type Config struct {
	Columns     ColumnConfig
	Performance PerformanceConfig
	Keybindings KeybindingsConfig
	Pacman      PacmanConfig
	AUR         AURConfig
}

// AURConfig contains AUR-related settings.
type AURConfig struct {
	Helper   string // Explicit helper name ("yay", "paru"). Empty = auto-detect
	Disabled bool   // Disable all AUR features
	Timeout  int    // HTTP timeout seconds (default 5)
	CacheTTL int    // Cache TTL seconds (default 300)
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

// PacmanConfig contains pacman-specific settings.
type PacmanConfig struct {
	DBPath string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
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
		Pacman: PacmanConfig{
			DBPath: "/var/lib/pacman",
		},
		AUR: AURConfig{
			Timeout:  5,
			CacheTTL: 300,
		},
	}
}
