package version

import "strings"

// Compare compares two semantic version strings.
// Returns:
//   -1 if a < b
//    0 if a == b
//    1 if a > b
func Compare(a, b string) int {
	// TODO: Parse epoch (e.g., "1:2.0.0")
	// TODO: Parse version segments (e.g., "1.2.3")
	// TODO: Parse release suffix (e.g., "1.2.3-1")
	// TODO: Compare numerically, not lexicographically

	// Stub: simple string comparison
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

// Parse parses a version string into components.
func Parse(version string) (epoch int, segments []int, release string) {
	// TODO: Implement proper version parsing
	// Format: [epoch:]version[-release]
	// Example: "1:2.0.3-1" -> epoch=1, segments=[2,0,3], release="1"
	return 0, nil, ""
}

// Less returns true if version a is less than version b.
func Less(a, b string) bool {
	return Compare(a, b) < 0
}

// Equal returns true if versions a and b are equal.
func Equal(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}
