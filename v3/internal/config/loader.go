package config

import (
	"os"
	"path/filepath"
)

// Load loads configuration from file.
func Load(path string) (*Config, error) {
	// TODO: Check if file exists
	// TODO: Parse TOML file
	// TODO: Merge with default config
	// TODO: Validate configuration

	// For now, return default config
	return DefaultConfig(), nil
}

// LoadDefault loads configuration from the default location.
func LoadDefault() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return DefaultConfig(), nil
	}

	configPath := filepath.Join(configDir, "pacviz", "config.toml")
	return Load(configPath)
}

// Save saves configuration to file.
func Save(config *Config, path string) error {
	// TODO: Create directory if needed
	// TODO: Marshal config to TOML
	// TODO: Write to file
	return nil
}
