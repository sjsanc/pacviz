package domain

// PresetType identifies different package view presets.
type PresetType string

const (
	PresetExplicit   PresetType = "explicit"
	PresetDependency PresetType = "dependency"
	PresetOrphans    PresetType = "orphans"
	PresetForeign    PresetType = "foreign"
	PresetAll        PresetType = "all"
)

// Preset defines a view preset configuration.
type Preset struct {
	Type        PresetType
	Name        string
	Description string
	Filter      func(*Package) bool
}

// DefaultPresets returns the standard preset configurations.
func DefaultPresets() []Preset {
	return []Preset{
		{
			Type:        PresetExplicit,
			Name:        "Explicit",
			Description: "User-installed packages",
			Filter: func(p *Package) bool {
				return p.InstallReason == ReasonExplicit
			},
		},
		{
			Type:        PresetDependency,
			Name:        "Dependencies",
			Description: "Auto-installed dependencies",
			Filter: func(p *Package) bool {
				return p.InstallReason == ReasonDependency
			},
		},
		{
			Type:        PresetOrphans,
			Name:        "Orphans",
			Description: "Dependencies with no dependents",
			Filter: func(p *Package) bool {
				return p.IsOrphan
			},
		},
		{
			Type:        PresetForeign,
			Name:        "Foreign",
			Description: "Packages not in sync databases",
			Filter: func(p *Package) bool {
				return p.IsForeign
			},
		},
		{
			Type:        PresetAll,
			Name:        "All",
			Description: "All installed packages",
			Filter: func(p *Package) bool {
				return true
			},
		},
	}
}
