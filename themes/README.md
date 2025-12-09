# Pacviz Themes

This directory contains example theme files for pacviz. These themes are also available as built-in themes in the application.

## Available Themes

- **catppuccin-mocha.toml** - Soothing pastel theme from the Catppuccin palette
- **gruvbox.toml** - Warm, retro groove color scheme
- **monokai.toml** - Colorful theme inspired by the classic Monokai editor theme
- **darcula.toml** - Dark theme from IntelliJ IDEA

## Using Themes

### 1. Via Config File

Edit your `~/.config/pacviz/config.toml`:

```toml
theme = "catppuccin-mocha"
```

### 2. Runtime Theme Switching

Use the `:theme` command while running pacviz:

```
:theme gruvbox
:th monokai
```

### 3. Custom Themes

You can create your own themes by:

1. Copying one of these theme files to `~/.config/pacviz/themes/mytheme.toml`
2. Modifying the colors to your liking
3. Loading it with `theme = "mytheme"` in your config or `:theme mytheme` at runtime

## Theme File Structure

```toml
name = "mytheme"
accent1 = "#89b4fa"      # Primary accent color
accent2 = "#cba6f7"      # Secondary accent color
accent3 = "#a6e3a1"      # Tertiary accent color
accent4 = "#f9e2af"      # Quaternary accent color
accent5 = "#f38ba8"      # Quinary accent color
background = "#1e1e2e"   # Main background
background_alt = "#181825"  # Alternate background for striped rows
foreground = "#cdd6f4"   # Main text color
selected = "#45475a"     # Selection background
dimmed = "#6c7086"       # Dimmed/disabled text
remote_accent = "#fab387"  # Accent for remote mode
warning_accent = "#f38ba8"  # Accent for warnings/destructive actions
```

## Theme Loading Priority

Themes are loaded in this order:

1. `~/.config/pacviz/themes/<name>.toml` (user themes)
2. `/usr/share/pacviz/themes/<name>.toml` (system themes)
3. Built-in registry (hardcoded themes)

## Color Overrides

You can override specific colors on top of any theme:

```toml
theme = "gruvbox"

[theme.overrides]
accent1 = "#custom_color"
background = "#282828"
```
