package viewport

import (
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// ApplyFilter filters rows by a simple text term (case-insensitive substring match).
// This is used for live filtering in filter mode.
func (v *Viewport) ApplyFilter(term string) {
	if term == "" {
		// Empty filter shows all rows
		v.VisibleRows = v.AllRows
		v.Filter = domain.FilterState{Active: false}
		return
	}

	// Filter rows by term
	filtered := make([]*domain.Row, 0, len(v.AllRows))
	lowerTerm := strings.ToLower(term)

	for _, row := range v.AllRows {
		// Search in name and description columns
		name := strings.ToLower(row.Cells[column.ColName])
		desc := strings.ToLower(row.Cells[column.ColDescription])

		if strings.Contains(name, lowerTerm) || strings.Contains(desc, lowerTerm) {
			filtered = append(filtered, row)
		}
	}

	v.VisibleRows = filtered
	v.Filter = domain.FilterState{
		Active: true,
		Terms:  []string{term},
	}

	// Reset selection to top of filtered results
	v.SelectedRow = 0
	v.Offset = 0
}

// ClearFilter removes all filters and restores all rows.
func (v *Viewport) ClearFilter() {
	v.Filter = domain.FilterState{Active: false}
	v.VisibleRows = v.AllRows
	v.SelectedRow = 0
	v.Offset = 0
}

// ApplyPresetFilter applies a preset filter function to the rows.
func (v *Viewport) ApplyPresetFilter(filterFunc func(*domain.Package) bool) {
	if filterFunc == nil {
		// No filter, show all rows
		v.VisibleRows = v.AllRows
		return
	}

	// Filter rows by preset function
	filtered := make([]*domain.Row, 0, len(v.AllRows))
	for _, row := range v.AllRows {
		if row.Package != nil && filterFunc(row.Package) {
			filtered = append(filtered, row)
		}
	}

	v.VisibleRows = filtered

	// Reset selection to top of filtered results
	v.SelectedRow = 0
	v.Offset = 0

	// Re-sort the filtered rows
	v.sortRows()
}
