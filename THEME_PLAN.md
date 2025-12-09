# Theme System Implementation Plan

## Overview

Implement a comprehensive theming system that allows users to:
- Select themes via config file (`theme = "<name>"`)
- Switch themes at runtime via `:theme <name>` command
- Load themes from multiple sources (built-in, system, user)
- Override individual colors on top of selected themes

## Current State

- **Default theme**: `styles.DarkTheme` hardcoded in memory (Tokyo Night colors)
- **Config overrides**: Individual color values can be overridden in `config.toml`
- **Global theme**: `styles.Current` is the active `*Styles` instance used throughout the app
- **Theme application**: Config loader directly modifies `styles.Current` at startup

## Design

### 1. Theme File Locations

Themes are loaded from the first match found in this priority order:

```
1. ~/.config/pacviz/themes/<name>.toml          (user themes)
2. /usr/share/pacviz/themes/<name>.toml         (system/distro themes)
3. Built-in registry in code                    (shipped themes)
```

### 2. Theme File Structure

```toml
# themes/gruvbox.toml
name = "gruvbox"
accent1 = "#83a598"      # Blue
accent2 = "#d3869b"      # Purple
accent3 = "#b8bb26"      # Green
accent4 = "#fabd2f"      # Yellow
accent5 = "#fb4934"      # Red
background = "#282828"
background_alt = "#1d2021"
foreground = "#ebdbb2"
selected = "#3c3836"
dimmed = "#928374"
remote_accent = "#fe8019"
warning_accent = "#fb4934"
```

### 3. Updated Config Structure

```toml
# config.toml
theme = "gruvbox"  # Theme name to load

[theme.overrides]
# Optional: override specific colors on top of selected theme
accent1 = "#custom_color"
```

**Backwards compatibility**: If no `theme =` line exists, treat `[theme]` section as overrides on default theme (current behavior).

### 4. Default Theme

The default theme should use terminal colors for maximum compatibility:

```go
var DefaultTheme = Theme{
    Name:          "default",
    Accent1:       "#FFFFFF",  // White
    Accent2:       "#FFFFFF",  // White
    Accent3:       "#FFFFFF",  // White
    Accent4:       "#FFFFFF",  // White
    Accent5:       "#FFFFFF",  // White
    Background:    "#000000",  // Black
    BackgroundAlt: "#000000",  // Black
    Foreground:    "#FFFFFF",  // White
    Selected:      "#FFFFFF",  // White (lipgloss should invert)
    Dimmed:        "#808080",  // Gray
    RemoteAccent:  "#FFFFFF",  // White
    WarningAccent: "#FFFFFF",  // White
}
```

### 5. Built-in Theme Registry

Create `internal/ui/styles/registry.go`:

```go
package styles

var BuiltinThemes = map[string]Theme{
    "default": DefaultTheme,
    "tokyo-night": TokyoNightTheme,  // Rename current DarkTheme
    "gruvbox": GruvboxTheme,
    // ... more themes as needed
}

func GetBuiltinTheme(name string) (Theme, bool) {
    theme, ok := BuiltinThemes[name]
    return theme, ok
}
```

### 6. Theme Loader

Create `internal/ui/styles/loader.go`:

```go
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

func mergeWithDefaults(theme Theme) Theme {
    base := DefaultTheme

    if theme.Name != "" { base.Name = theme.Name }
    if theme.Accent1 != "" { base.Accent1 = theme.Accent1 }
    if theme.Accent2 != "" { base.Accent2 = theme.Accent2 }
    if theme.Accent3 != "" { base.Accent3 = theme.Accent3 }
    if theme.Accent4 != "" { base.Accent4 = theme.Accent4 }
    if theme.Accent5 != "" { base.Accent5 = theme.Accent5 }
    if theme.Background != "" { base.Background = theme.Background }
    if theme.BackgroundAlt != "" { base.BackgroundAlt = theme.BackgroundAlt }
    if theme.Foreground != "" { base.Foreground = theme.Foreground }
    if theme.Selected != "" { base.Selected = theme.Selected }
    if theme.Dimmed != "" { base.Dimmed = theme.Dimmed }
    if theme.RemoteAccent != "" { base.RemoteAccent = theme.RemoteAccent }
    if theme.WarningAccent != "" { base.WarningAccent = theme.WarningAccent }

    return base
}

func loadUserTheme(name string) (Theme, error) {
    path, err := getUserThemePath(name)
    if err != nil {
        return Theme{}, err
    }
    return loadThemeFile(path)
}

func loadSystemTheme(name string) (Theme, error) {
    path := filepath.Join("/usr/share/pacviz/themes", name+".toml")
    return loadThemeFile(path)
}

func loadThemeFile(path string) (Theme, error) {
    var theme Theme
    _, err := toml.DecodeFile(path, &theme)
    return theme, err
}

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
```

### 7. Runtime Theme Switching

**Command**: `:theme <name>`

**Implementation in `internal/command/execute.go`**:

```go
// Add to ExecuteResult struct
type ExecuteResult struct {
    // ... existing fields
    ThemeName string  // NEW: theme to apply
}

// Add to Execute function
case "theme", "t":
    if len(parts) < 2 {
        return ExecuteResult{Error: "usage: :theme <name>"}
    }
    return ExecuteResult{ThemeName: parts[1]}
```

**Handle in `internal/app/update.go`**:

```go
func (m Model) handleCommandResult(result command.ExecuteResult) (Model, tea.Cmd) {
    // ... existing handling

    if result.ThemeName != "" {
        theme, err := styles.LoadTheme(result.ThemeName)
        if err != nil {
            m.statusMessage = fmt.Sprintf("Error: %v", err)
            return m, nil
        }
        styles.ApplyTheme(theme)
        m.statusMessage = fmt.Sprintf("Theme changed to: %s", result.ThemeName)
    }

    // ... rest of function
}
```

### 8. Config Loading Changes

Modify `internal/config/loader.go`:

```go
// Update TOML structure
tomlConfig := struct {
    ThemeName string `toml:"theme"`  // NEW: theme name
    Theme struct {
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

        // Backwards compat: flat structure
        Accent1       string `toml:"accent1"`
        Accent2       string `toml:"accent2"`
        // ... etc
    } `toml:"theme"`
}{}

// Loading logic:
// 1. Load base theme by name (or use DefaultTheme)
baseTheme := styles.DefaultTheme
if tomlConfig.ThemeName != "" {
    if loaded, err := styles.LoadTheme(tomlConfig.ThemeName); err == nil {
        baseTheme = loaded
    }
}

// 2. Apply overrides (check both nested and flat for backwards compat)
// 3. Call styles.ApplyTheme(baseTheme)
```

### 9. Implementation Checklist

- [ ] Create `internal/ui/styles/registry.go` with built-in themes
- [ ] Create `internal/ui/styles/loader.go` with theme loading logic
- [ ] Update `internal/ui/styles/theme.go` to rename DarkTheme → TokyoNightTheme
- [ ] Add DefaultTheme (black/white terminal colors)
- [ ] Update `internal/config/loader.go` to support `theme = "<name>"`
- [ ] Add `:theme` command to `internal/command/execute.go`
- [ ] Handle theme changes in `internal/app/update.go`
- [ ] Test theme switching at runtime
- [ ] Test config-based theme loading
- [ ] Test theme file loading from user/system directories
- [ ] Test backwards compatibility with existing configs
- [ ] Update example `config.toml` with theme documentation

### 10. Error Handling

- **Theme not found**: Show error in status bar, keep current theme
- **Invalid TOML**: Show parse error, keep current theme
- **Missing fields**: Merge with DefaultTheme (see `mergeWithDefaults()`)
- **File I/O errors**: Graceful degradation to built-in themes
- **Color validation**: Let lipgloss handle invalid hex values

### 11. Important Implementation Notes

**⚠️ Runtime style updates**: Ensure all rendering code uses `styles.Current.*` dynamically, not cached values:

```go
// BAD: Caches the style at init time
var myStyle = styles.Current.Header

// GOOD: References current theme each render
func render() string {
    return styles.Current.Header.Render("text")
}
```

**Re-render triggering**: After theme change, Update() returns a message which automatically triggers View(), so the new theme appears immediately.

## Future Enhancements (TODO)

### TODO: Error Reporting
Need a better error reporting mechanism for theme loading failures. Currently errors just show in status bar briefly. Consider:
- Dedicated error pane
- Error log file
- More descriptive error messages with suggestions

### TODO: Command Palette Improvements
Add support for positional argument completion in command palette. When user types `:theme ` and triggers completion, show available themes:
- Built-in themes
- User themes from `~/.config/pacviz/themes/`
- System themes from `/usr/share/pacviz/themes/`

This requires:
- Command argument parsing enhancement
- Filesystem scanning for theme files
- Autocomplete UI component
- Tab completion support

## Testing Strategy

1. **Unit tests** for theme loading (`loader_test.go`)
2. **Integration tests** for config loading with themes
3. **Manual testing**:
   - Create custom theme in `~/.config/pacviz/themes/`
   - Test `:theme <name>` command
   - Test `theme = "<name>"` in config
   - Test override behavior
   - Test error cases (missing theme, invalid TOML)
   - Test backwards compatibility with existing configs

## Migration Path

Existing users with theme overrides in their config will continue to work:
- No `theme =` line → overrides apply to DefaultTheme (now black/white instead of Tokyo Night)
- If users want to keep Tokyo Night look, they should add `theme = "tokyo-night"` to their config
