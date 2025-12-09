package styles

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// LoadTheme attempts to load a theme by name from:
// 1. User config dir (~/.config/pacviz/themes/)
// 2. System dir (/usr/share/pacviz/themes/)
// 3. Built-in registry
// Returns error if theme not found in any location.
func LoadTheme(name string) (Theme, error) {
	// Try user themes
	if userTheme, err := loadUserTheme(name); err == nil {
		return userTheme, nil
	}

	// Try system themes
	if sysTheme, err := loadSystemTheme(name); err == nil {
		return sysTheme, nil
	}

	// Try built-in registry
	if theme, ok := GetBuiltinTheme(name); ok {
		return theme, nil
	}

	return Theme{}, fmt.Errorf("theme not found: %s", name)
}

// ApplyTheme updates styles.Current with a new theme.
// Missing fields are merged from DefaultTheme.
func ApplyTheme(theme Theme) {
	// Merge with defaults for any missing fields
	merged := mergeWithDefaults(theme)
	Current = NewStyles(merged)
}

// mergeWithDefaults fills in any missing fields in the theme with DefaultTheme values.
func mergeWithDefaults(theme Theme) Theme {
	base := DefaultTheme

	if theme.Name != "" {
		base.Name = theme.Name
	}
	if theme.Accent1 != "" {
		base.Accent1 = theme.Accent1
	}
	if theme.Accent2 != "" {
		base.Accent2 = theme.Accent2
	}
	if theme.Accent3 != "" {
		base.Accent3 = theme.Accent3
	}
	if theme.Accent4 != "" {
		base.Accent4 = theme.Accent4
	}
	if theme.Accent5 != "" {
		base.Accent5 = theme.Accent5
	}
	if theme.Background != "" {
		base.Background = theme.Background
	}
	if theme.BackgroundAlt != "" {
		base.BackgroundAlt = theme.BackgroundAlt
	}
	if theme.Foreground != "" {
		base.Foreground = theme.Foreground
	}
	if theme.Selected != "" {
		base.Selected = theme.Selected
	}
	if theme.Dimmed != "" {
		base.Dimmed = theme.Dimmed
	}
	if theme.RemoteAccent != "" {
		base.RemoteAccent = theme.RemoteAccent
	}
	if theme.WarningAccent != "" {
		base.WarningAccent = theme.WarningAccent
	}

	return base
}

// loadUserTheme loads a theme from the user's config directory.
func loadUserTheme(name string) (Theme, error) {
	path, err := getUserThemePath(name)
	if err != nil {
		return Theme{}, err
	}
	return loadThemeFile(path)
}

// loadSystemTheme loads a theme from the system themes directory.
func loadSystemTheme(name string) (Theme, error) {
	path := filepath.Join("/usr/share/pacviz/themes", name+".toml")
	return loadThemeFile(path)
}

// loadThemeFile loads a theme from a TOML file.
func loadThemeFile(path string) (Theme, error) {
	var theme Theme
	_, err := toml.DecodeFile(path, &theme)
	return theme, err
}

// getUserThemePath returns the path to a theme file in the user's config directory.
func getUserThemePath(name string) (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "pacviz", "themes", name+".toml"), nil
}
