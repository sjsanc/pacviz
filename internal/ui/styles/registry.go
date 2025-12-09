package styles

// GruvboxTheme is a warm, retro groove color scheme.
var GruvboxTheme = Theme{
	Name:          "gruvbox",
	Accent1:       "#83a598", // Blue
	Accent2:       "#d3869b", // Purple
	Accent3:       "#b8bb26", // Green
	Accent4:       "#fabd2f", // Yellow
	Accent5:       "#fb4934", // Red
	Background:    "#282828", // bg0
	BackgroundAlt: "#1d2021", // bg0_h
	Foreground:    "#ebdbb2", // fg1
	Selected:      "#3c3836", // bg1
	Dimmed:        "#928374", // gray
	RemoteAccent:  "#fe8019", // Orange
	WarningAccent: "#fb4934", // Red
}

// MonokaiTheme is a colorful theme inspired by the Monokai color scheme.
var MonokaiTheme = Theme{
	Name:          "monokai",
	Accent1:       "#66d9ef", // Blue
	Accent2:       "#ae81ff", // Purple
	Accent3:       "#a6e22e", // Green
	Accent4:       "#e6db74", // Yellow
	Accent5:       "#f92672", // Red
	Background:    "#272822", // Background
	BackgroundAlt: "#1e1f1c", // Darker background
	Foreground:    "#f8f8f2", // Foreground
	Selected:      "#49483e", // Selection
	Dimmed:        "#75715e", // Comment
	RemoteAccent:  "#fd971f", // Orange
	WarningAccent: "#f92672", // Red
}

// DarculaTheme is the dark theme from IntelliJ IDEA.
var DarculaTheme = Theme{
	Name:          "darcula",
	Accent1:       "#6897bb", // Cyan/Blue
	Accent2:       "#9876aa", // Purple
	Accent3:       "#6a8759", // Green
	Accent4:       "#ffc66d", // Yellow
	Accent5:       "#ff6b68", // Red
	Background:    "#2b2b2b", // Background
	BackgroundAlt: "#323232", // Lighter background
	Foreground:    "#a9b7c6", // Foreground
	Selected:      "#214283", // Selection
	Dimmed:        "#808080", // Comment
	RemoteAccent:  "#cc7832", // Orange
	WarningAccent: "#ff6b68", // Red
}

// CatppuccinMochaTheme is a soothing pastel theme.
var CatppuccinMochaTheme = Theme{
	Name:          "catppuccin-mocha",
	Accent1:       "#89b4fa", // Blue
	Accent2:       "#cba6f7", // Mauve
	Accent3:       "#a6e3a1", // Green
	Accent4:       "#f9e2af", // Yellow
	Accent5:       "#f38ba8", // Red
	Background:    "#1e1e2e", // Base
	BackgroundAlt: "#181825", // Mantle
	Foreground:    "#cdd6f4", // Text
	Selected:      "#45475a", // Surface1
	Dimmed:        "#6c7086", // Overlay0
	RemoteAccent:  "#fab387", // Peach
	WarningAccent: "#f38ba8", // Red
}

// BuiltinThemes is a registry of all built-in themes.
var BuiltinThemes = map[string]Theme{
	"default":           DefaultTheme,
	"tokyo-night":       TokyoNightTheme,
	"gruvbox":           GruvboxTheme,
	"monokai":           MonokaiTheme,
	"darcula":           DarculaTheme,
	"catppuccin-mocha":  CatppuccinMochaTheme,
}

// GetBuiltinTheme retrieves a built-in theme by name.
// Returns the theme and true if found, or an empty theme and false if not found.
func GetBuiltinTheme(name string) (Theme, bool) {
	theme, ok := BuiltinThemes[name]
	return theme, ok
}
