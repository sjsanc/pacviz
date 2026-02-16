package domain

// PresetType identifies different package view presets.
type PresetType string

const (
	PresetExplicit   PresetType = "explicit"
	PresetDependency PresetType = "dependency"
	PresetOrphans    PresetType = "orphans"
	PresetForeign    PresetType = "foreign"
	PresetAUR        PresetType = "aur"
	PresetUpdatable  PresetType = "updatable"
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
			Description: "Non-repo, non-AUR packages",
			Filter: func(p *Package) bool {
				return p.IsForeign && !p.IsAUR
			},
		},
		{
			Type:        PresetAUR,
			Name:        "AUR",
			Description: "AUR packages",
			Filter: func(p *Package) bool {
				return p.IsAUR
			},
		},
		{
			Type:        PresetUpdatable,
			Name:        "Updatable",
			Description: "Packages with available updates",
			Filter: func(p *Package) bool {
				return p.HasUpdate
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
