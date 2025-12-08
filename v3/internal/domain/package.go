package domain

import "time"

// InstallReason indicates why a package was installed.
type InstallReason int

const (
	ReasonExplicit InstallReason = iota
	ReasonDependency
)

// Package represents a complete package entry.
type Package struct {
	Name         string
	Version      string
	Description  string
	Architecture string
	URL          string
	Licenses     []string
	Groups       []string
	Dependencies []string
	OptDepends   map[string]string // name -> description
	Conflicts    []string
	Provides     []string
	Replaces     []string

	// Installation metadata
	Installed     bool
	InstallDate   time.Time
	InstallReason InstallReason
	InstalledSize int64
	Packager      string
	BuildDate     time.Time

	// Computed fields
	Required  []string // packages depending on this
	IsOrphan  bool
	IsForeign bool
	HasUpdate bool
	NewVersion string
}
