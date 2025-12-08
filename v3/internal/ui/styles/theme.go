package styles

// Theme represents a color theme configuration.
type Theme struct {
	Name       string
	Accent1    string
	Accent2    string
	Accent3    string
	Accent4    string
	Accent5    string
	Background string
	Foreground string
	Selected   string
	Dimmed     string
}

// DarkTheme is the default dark theme.
var DarkTheme = Theme{
	Name:       "dark",
	Accent1:    "#7aa2f7",
	Accent2:    "#bb9af7",
	Accent3:    "#9ece6a",
	Accent4:    "#e0af68",
	Accent5:    "#f7768e",
	Background: "#1a1b26",
	Foreground: "#c0caf5",
	Selected:   "#283457",
	Dimmed:     "#565f89",
}

// ApplyTheme applies a theme to the global styles.
func ApplyTheme(theme Theme) {
	// TODO: Update global style variables based on theme
}
