package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Jguer/go-alpm/v2"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Scrollbar
// TODO: Mouse events (clicking on rows and headers)
// TODO: Status Line and Command Buffer are separate components.
// - Status Line contains Current View (Explicit, Remote etc) and Total Row Count / Viewport
// TODO: Rename local search to filter
// TODO: Remote Search applies on enter.
// TODO: Hide InstalledDate on remote search.
// TODO: Add 1 px padding

type View struct {
	Name    string
	Rows    []*Row
	ResetFn func() []*Row
}

func (v *View) Reset() {
	v.Rows = v.ResetFn()
}

type Model struct {
	pacman           *Pacman
	height           int
	width            int
	vp               *Viewport
	buffer           string
	prefix           string
	infoboxIsVisible bool
	infoboxHeight    int
	mode             string // cmd | search

	Buffer       string
	CommandMode  bool
	FilterMode   bool
	ColumnWidths []int
	Views        []View
}

func NewModel() *Model {
	pacman := NewPacman()
	rows := convertPkgsToRows(pacman.GetExplicitPkgs())
	return &Model{
		pacman: pacman,
		vp:     &Viewport{rows: rows},
		Views: []View{
			{
				Name: "Explicit",
				Rows: convertPkgsToRows(pacman.GetExplicitPkgs()),
				ResetFn: func() []*Row {
					return convertPkgsToRows(pacman.GetExplicitPkgs())
				},
			},
			{
				Name: "Orphans",
				Rows: nil,
				ResetFn: func() []*Row {
					return convertPkgsToRows(pacman.GetOrphanPkgs())
				},
			},
			{
				Name: "Sync",
				Rows: nil,
				ResetFn: func() []*Row {
					return []*Row{}
				},
			},
		},
	}
}

func (m *Model) SetHeight(h int) *Model {
	m.height = h
	m.vp.SetHeight(h)
	return m
}
func (m *Model) SetWidth(w int) *Model {
	m.width = w
	m.ColumnWidths = calcColWidths(COLUMNS, w)
	return m
}
func (m *Model) SetDimensions(w, h int) *Model {
	m.SetWidth(w)
	m.SetHeight(h)
	return m
}

// func (m *Model) EnterCommandMode() {
// 	m.prefix = ":"
// 	m.mode = "cmd"
// 	m.CloseInfobox()
// }
// func (m *Model) EnterLocalSearchMode() {
// 	m.prefix = "search:"
// 	m.mode = "search"
// 	m.CloseInfobox()
// }

// func (m *Model) OpenInfobox() {
// 	m.infoboxIsVisible = true
// 	m.infoboxHeight = m.vp.height / 2
// 	m.vp.height = m.vp.height - m.infoboxHeight
// 	if m.vp.selected >= m.vp.height {
// 		m.vp.offset = m.vp.selected - m.vp.height + 2
// 	}
// }
// func (m *Model) CloseInfobox() {
// 	m.infoboxIsVisible = false
// 	m.vp.height = m.height - 2
// 	m.infoboxHeight = 0
// }

func (m Model) Init() tea.Cmd {
	return nil
}

func convertPkgsToRows(pkgs []alpm.IPackage) []*Row {
	rows := []*Row{}
	for i, pkg := range pkgs {
		rows = append(rows, &Row{
			cells: []string{
				strconv.Itoa(i),
				pkg.Name(),
				pkg.Version(),
				strings.Join(pkg.Groups().Slice(), ", "),
				pkg.InstallDate().Format("2006-01-02"),
				pkg.Description(),
			},
		})
	}
	return rows
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
