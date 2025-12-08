package viewport

import (
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// Viewport manages the visible portion of the package list.
type Viewport struct {
	AllRows     []*domain.Row
	VisibleRows []*domain.Row

	// Pagination
	Offset int
	Height int

	// Selection
	SelectedRow int
	SelectedCol int

	// Sorting
	SortColumn  column.Type
	SortReverse bool

	// Filtering
	Filter domain.FilterState

	// Columns
	Columns []*column.Column
}

// New creates a new viewport.
func New() *Viewport {
	return &Viewport{
		AllRows:     make([]*domain.Row, 0),
		VisibleRows: make([]*domain.Row, 0),
		Columns:     column.DefaultColumns(),
		SortColumn:  column.ColName,
		SortReverse: false,
		SelectedCol: 1, // Start at first selectable column (skip index column at 0)
	}
}

// SetRows sets the rows and updates the visible rows.
func (v *Viewport) SetRows(rows []*domain.Row) {
	v.AllRows = rows
	v.sortRows()
	v.updateVisibleRows()
}

// GetVisibleRows returns the currently visible rows based on offset and height.
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

// GetSelectedPackage returns the currently selected package.
func (v *Viewport) GetSelectedPackage() *domain.Package {
	if v.SelectedRow < 0 || v.SelectedRow >= len(v.VisibleRows) {
		return nil
	}
	return v.VisibleRows[v.SelectedRow].Package
}

// updateVisibleRows updates the visible rows based on current filter.
func (v *Viewport) updateVisibleRows() {
	// For now, just copy all rows
	// TODO: Apply filtering
	v.VisibleRows = v.AllRows
}
