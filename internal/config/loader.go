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
		Theme struct {
			Accent1        string `toml:"accent1"`
			Accent2        string `toml:"accent2"`
			Accent3        string `toml:"accent3"`
			Accent4        string `toml:"accent4"`
			Accent5        string `toml:"accent5"`
			Background     string `toml:"background"`
			BackgroundAlt  string `toml:"background_alt"`
			Foreground     string `toml:"foreground"`
			Selected       string `toml:"selected"`
			Dimmed         string `toml:"dimmed"`
			RemoteAccent   string `toml:"remote_accent"`
			WarningAccent  string `toml:"warning_accent"`
		} `toml:"theme"`
	}{}

	_, err := toml.DecodeFile(path, &tomlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply theme colors if any are specified
	if tomlConfig.Theme.Accent1 != "" || tomlConfig.Theme.Accent2 != "" ||
		tomlConfig.Theme.Accent3 != "" || tomlConfig.Theme.Accent4 != "" ||
		tomlConfig.Theme.Accent5 != "" || tomlConfig.Theme.Background != "" ||
		tomlConfig.Theme.BackgroundAlt != "" || tomlConfig.Theme.Foreground != "" ||
		tomlConfig.Theme.Selected != "" || tomlConfig.Theme.Dimmed != "" ||
		tomlConfig.Theme.RemoteAccent != "" || tomlConfig.Theme.WarningAccent != "" {

		theme := styles.DarkTheme

		if tomlConfig.Theme.Accent1 != "" {
			theme.Accent1 = tomlConfig.Theme.Accent1
		}
		if tomlConfig.Theme.Accent2 != "" {
			theme.Accent2 = tomlConfig.Theme.Accent2
		}
		if tomlConfig.Theme.Accent3 != "" {
			theme.Accent3 = tomlConfig.Theme.Accent3
		}
		if tomlConfig.Theme.Accent4 != "" {
			theme.Accent4 = tomlConfig.Theme.Accent4
		}
		if tomlConfig.Theme.Accent5 != "" {
			theme.Accent5 = tomlConfig.Theme.Accent5
		}
		if tomlConfig.Theme.Background != "" {
			theme.Background = tomlConfig.Theme.Background
		}
		if tomlConfig.Theme.BackgroundAlt != "" {
			theme.BackgroundAlt = tomlConfig.Theme.BackgroundAlt
		}
		if tomlConfig.Theme.Foreground != "" {
			theme.Foreground = tomlConfig.Theme.Foreground
		}
		if tomlConfig.Theme.Selected != "" {
			theme.Selected = tomlConfig.Theme.Selected
		}
		if tomlConfig.Theme.Dimmed != "" {
			theme.Dimmed = tomlConfig.Theme.Dimmed
		}
		if tomlConfig.Theme.RemoteAccent != "" {
			theme.RemoteAccent = tomlConfig.Theme.RemoteAccent
		}
		if tomlConfig.Theme.WarningAccent != "" {
			theme.WarningAccent = tomlConfig.Theme.WarningAccent
		}

		styles.Current = styles.NewStyles(theme)
	}

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