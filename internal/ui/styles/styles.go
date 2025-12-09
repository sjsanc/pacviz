package styles

// This file maintains backward compatibility by providing package-level access to styles.
// All styles are now generated from themes defined in theme.go.

// Global default styles instance.
// This can be replaced to change the theme at runtime.
var Current = Default()
