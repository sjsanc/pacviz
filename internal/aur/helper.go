package aur

import (
	"os/exec"
)

// HelperConfig represents a detected AUR helper.
type HelperConfig struct {
	Name string
	Path string
}

// defaultHelpers is the priority-ordered list of AUR helpers to try.
var defaultHelpers = []string{"yay", "paru", "pikaur", "trizen"}

// DetectHelper finds an installed AUR helper.
// If configOverride is set, only that helper is checked.
// Otherwise, tries yay, paru, pikaur, trizen in order.
func DetectHelper(configOverride string) *HelperConfig {
	if configOverride != "" {
		path, err := exec.LookPath(configOverride)
		if err != nil {
			return nil
		}
		return &HelperConfig{Name: configOverride, Path: path}
	}

	for _, name := range defaultHelpers {
		path, err := exec.LookPath(name)
		if err == nil {
			return &HelperConfig{Name: name, Path: path}
		}
	}

	return nil
}

// InstallCmd builds an exec.Cmd for installing packages via the AUR helper.
// No --noconfirm is used — the helper runs interactively for PKGBUILD review.
func InstallCmd(helper *HelperConfig, names []string) *exec.Cmd {
	args := append([]string{"-S"}, names...)
	return exec.Command(helper.Path, args...)
}
