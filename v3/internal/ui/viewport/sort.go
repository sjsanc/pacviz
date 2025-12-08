package viewport

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// ApplySort sorts the rows by the specified column.
func (v *Viewport) ApplySort(col column.Type, reverse bool) {
	v.SortColumn = col
	v.SortReverse = reverse
	// TODO: Sort AllRows by column type
	// TODO: Update VisibleRows
}

// ToggleSort toggles the sort direction on the current column.
func (v *Viewport) ToggleSort() {
	v.SortReverse = !v.SortReverse
	// TODO: Re-sort with new direction
}

// SortByColumn sorts by a specific column (used when selecting columns).
func (v *Viewport) SortByColumn(col column.Type) {
	if v.SortColumn == col {
		v.ToggleSort()
	} else {
		v.ApplySort(col, false)
	}
}
