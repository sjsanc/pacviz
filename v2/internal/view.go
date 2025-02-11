package internal

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	sb := strings.Builder{}

	for i, c := range COLUMNS {
		if i == m.VP.CurrentColIndex {
			sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("21")).Render(string(c)))
		} else {
			sb.WriteString(string(c))
		}
	}

	sb.WriteString("\n")

	return sb.String()
}
