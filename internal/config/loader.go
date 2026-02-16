package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// Load loads configuration from file.
// Returns default config if file doesn't exist (only when optional is true).
// Fails if file doesn't exist and optional is false.
func Load(path string, optional bool) (*Config, error) {
	config := DefaultConfig()

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if optional {
				return config, nil
			}
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat config file: %w", err)
	}

	// Parse TOML file
	tomlConfig := struct {
		SelectedTheme string `toml:"selected_theme"` // Theme name to load
		AUR           struct {
			Helper   string `toml:"helper"`
			Disabled bool   `toml:"disabled"`
			Timeout  int    `toml:"timeout"`
			CacheTTL int    `toml:"cache_ttl"`
		} `toml:"aur"`
		Theme struct {
			// Nested overrides (new format)
			Overrides struct {
				Accent1       string `toml:"accent1"`
				Accent2       string `toml:"accent2"`
				Accent3       string `toml:"accent3"`
				Accent4       string `toml:"accent4"`
				Accent5       string `toml:"accent5"`
				Background    string `toml:"background"`
				BackgroundAlt string `toml:"background_alt"`
				Foreground    string `toml:"foreground"`
				Selected      string `toml:"selected"`
				Dimmed        string `toml:"dimmed"`
				RemoteAccent  string `toml:"remote_accent"`
				WarningAccent string `toml:"warning_accent"`
			} `toml:"overrides"`

			// Flat structure (backwards compatibility)
			Accent1       string `toml:"accent1"`
			Accent2       string `toml:"accent2"`
			Accent3       string `toml:"accent3"`
			Accent4       string `toml:"accent4"`
			Accent5       string `toml:"accent5"`
			Background    string `toml:"background"`
			BackgroundAlt string `toml:"background_alt"`
			Foreground    string `toml:"foreground"`
			Selected      string `toml:"selected"`
			Dimmed        string `toml:"dimmed"`
			RemoteAccent  string `toml:"remote_accent"`
			WarningAccent string `toml:"warning_accent"`
		} `toml:"theme"`
	}{}

	_, err := toml.DecodeFile(path, &tomlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply AUR config overrides
	if tomlConfig.AUR.Helper != "" {
		config.AUR.Helper = tomlConfig.AUR.Helper
	}
	if tomlConfig.AUR.Disabled {
		config.AUR.Disabled = true
	}
	if tomlConfig.AUR.Timeout > 0 {
		config.AUR.Timeout = tomlConfig.AUR.Timeout
	}
	if tomlConfig.AUR.CacheTTL > 0 {
		config.AUR.CacheTTL = tomlConfig.AUR.CacheTTL
	}

	// Load base theme by name (defaults to "default" if not specified)
	themeName := tomlConfig.SelectedTheme
	if themeName == "" {
		themeName = "default"
	}

	baseTheme, err := styles.LoadTheme(themeName)
	if err != nil {
		return nil, fmt.Errorf("failed to load theme '%s': %w", themeName, err)
	}

	// Apply overrides from both nested and flat structures (flat takes precedence for backwards compat)

	// Check nested overrides first
	if tomlConfig.Theme.Overrides.Accent1 != "" {
		baseTheme.Accent1 = tomlConfig.Theme.Overrides.Accent1
	}
	if tomlConfig.Theme.Overrides.Accent2 != "" {
		baseTheme.Accent2 = tomlConfig.Theme.Overrides.Accent2
	}
	if tomlConfig.Theme.Overrides.Accent3 != "" {
		baseTheme.Accent3 = tomlConfig.Theme.Overrides.Accent3
	}
	if tomlConfig.Theme.Overrides.Accent4 != "" {
		baseTheme.Accent4 = tomlConfig.Theme.Overrides.Accent4
	}
	if tomlConfig.Theme.Overrides.Accent5 != "" {
		baseTheme.Accent5 = tomlConfig.Theme.Overrides.Accent5
	}
	if tomlConfig.Theme.Overrides.Background != "" {
		baseTheme.Background = tomlConfig.Theme.Overrides.Background
	}
	if tomlConfig.Theme.Overrides.BackgroundAlt != "" {
		baseTheme.BackgroundAlt = tomlConfig.Theme.Overrides.BackgroundAlt
	}
	if tomlConfig.Theme.Overrides.Foreground != "" {
		baseTheme.Foreground = tomlConfig.Theme.Overrides.Foreground
	}
	if tomlConfig.Theme.Overrides.Selected != "" {
		baseTheme.Selected = tomlConfig.Theme.Overrides.Selected
	}
	if tomlConfig.Theme.Overrides.Dimmed != "" {
		baseTheme.Dimmed = tomlConfig.Theme.Overrides.Dimmed
	}
	if tomlConfig.Theme.Overrides.RemoteAccent != "" {
		baseTheme.RemoteAccent = tomlConfig.Theme.Overrides.RemoteAccent
	}
	if tomlConfig.Theme.Overrides.WarningAccent != "" {
		baseTheme.WarningAccent = tomlConfig.Theme.Overrides.WarningAccent
	}

	// Apply flat overrides (backwards compatibility - these take precedence)
	if tomlConfig.Theme.Accent1 != "" {
		baseTheme.Accent1 = tomlConfig.Theme.Accent1
	}
	if tomlConfig.Theme.Accent2 != "" {
		baseTheme.Accent2 = tomlConfig.Theme.Accent2
	}
	if tomlConfig.Theme.Accent3 != "" {
		baseTheme.Accent3 = tomlConfig.Theme.Accent3
	}
	if tomlConfig.Theme.Accent4 != "" {
		baseTheme.Accent4 = tomlConfig.Theme.Accent4
	}
	if tomlConfig.Theme.Accent5 != "" {
		baseTheme.Accent5 = tomlConfig.Theme.Accent5
	}
	if tomlConfig.Theme.Background != "" {
		baseTheme.Background = tomlConfig.Theme.Background
	}
	if tomlConfig.Theme.BackgroundAlt != "" {
		baseTheme.BackgroundAlt = tomlConfig.Theme.BackgroundAlt
	}
	if tomlConfig.Theme.Foreground != "" {
		baseTheme.Foreground = tomlConfig.Theme.Foreground
	}
	if tomlConfig.Theme.Selected != "" {
		baseTheme.Selected = tomlConfig.Theme.Selected
	}
	if tomlConfig.Theme.Dimmed != "" {
		baseTheme.Dimmed = tomlConfig.Theme.Dimmed
	}
	if tomlConfig.Theme.RemoteAccent != "" {
		baseTheme.RemoteAccent = tomlConfig.Theme.RemoteAccent
	}
	if tomlConfig.Theme.WarningAccent != "" {
		baseTheme.WarningAccent = tomlConfig.Theme.WarningAccent
	}

	// Always apply the loaded theme (with any overrides)
	styles.ApplyTheme(baseTheme)

	return config, nil
}

// LoadDefault loads configuration from the default location.
// It follows XDG conventions with fallbacks:
// 1. $XDG_CONFIG_HOME/pacviz/config.toml (or ~/.config/pacviz/config.toml)
// 2. Returns default config if no file found
func LoadDefault() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}
	return Load(configPath, true)
}

// LoadWithOverride loads configuration from a file, or the default location if path is empty.
// If path is explicitly provided, the file must exist or an error is returned.
func LoadWithOverride(path string) (*Config, error) {
	if path != "" {
		return Load(path, false)
	}
	return LoadDefault()
}

// getConfigPath returns the config file path following XDG conventions.
func getConfigPath() (string, error) {
	// Check XDG_CONFIG_HOME first
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "pacviz", "config.toml"), nil
	}

	// Fall back to ~/.config/pacviz/config.toml
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "pacviz", "config.toml"), nil
}