package viewport

import (
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Viewport manages the visible portion of the package list.
type Viewport struct {
	AllRows     []*domain.Row
	VisibleRows []*domain.Row

	Offset int
	Height int

	SelectedRow int
	SelectedCol int

	SortColumn  column.Type
	SortReverse bool

	Filter domain.FilterState

	Columns []*column.Column
}

func New() *Viewport {
	return &Viewport{
		AllRows:     make([]*domain.Row, 0),
		VisibleRows: make([]*domain.Row, 0),
		Columns:     column.DefaultColumns(),
		SortColumn:  column.ColName,
		SortReverse: false,
		SelectedCol: 1,
	}
}

func (v *Viewport) SetRows(rows []*domain.Row) {
	v.AllRows = rows
	v.sortRows()
	v.updateVisibleRows()
}

func (v *Viewport) GetVisibleRows() []*domain.Row {
	if v.Height == 0 {
		return v.VisibleRows
	}

	start := v.Offset
	end := v.Offset + v.Height

	if start >= len(v.VisibleRows) {
		return []*domain.Row{}
	}

	if end > len(v.VisibleRows) {
		end = len(v.VisibleRows)
	}

	return v.VisibleRows[start:end]
}

func (v *Viewport) GetSelectedPackage() *domain.Package {
	if v.SelectedRow < 0 || v.SelectedRow >= len(v.VisibleRows) {
		return nil
	}
	return v.VisibleRows[v.SelectedRow].Package
}

// Filtering is handled separately via ApplyFilter and ApplyPresetFilter.
func (v *Viewport) updateVisibleRows() {
	v.VisibleRows = v.AllRows
}
