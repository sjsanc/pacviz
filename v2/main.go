package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sjsanc/pacviz/internal"
)

func main() {
	m := internal.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
