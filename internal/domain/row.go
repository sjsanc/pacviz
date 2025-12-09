package domain

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Row is the display representation of a package.
type Row struct {
	Package  *Package
	Cells    map[column.Type]string
	Selected bool
	Filtered bool
}

// NewRow creates a new row from a package.
func NewRow(pkg *Package) *Row {
	return &Row{
		Package: pkg,
		Cells:   make(map[column.Type]string),
	}
}
