package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/v3/internal/app"
	"github.com/sjsanc/pacviz/v3/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "", "Path to config file (TOML format)")
	flag.StringVar(&configPath, "config", "", "Path to config file (TOML format)")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadWithOverride(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Create model with loaded config
	model := app.NewModel()
	_ = cfg // Config is used via styles.Current which is set by config loader

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
