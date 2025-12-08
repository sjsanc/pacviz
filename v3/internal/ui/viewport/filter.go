package viewport

import (
	"github.com/sjsanc/pacviz/v3/internal/domain"
)

// ApplyFilter filters the rows based on the filter state.
func (v *Viewport) ApplyFilter(filter domain.FilterState) {
	v.Filter = filter
	// TODO: Filter AllRows to create VisibleRows
	// TODO: Reset selection if needed
}

// ClearFilter removes all filters.
func (v *Viewport) ClearFilter() {
	v.Filter = domain.FilterState{}
	v.VisibleRows = v.AllRows
}
