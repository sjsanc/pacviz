package viewport

import (
	"sort"
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// ApplySort sorts the rows by the specified column.
func (v *Viewport) ApplySort(col column.Type, reverse bool) {
	v.SortColumn = col
	v.SortReverse = reverse
	v.sortRows()
	v.updateVisibleRows()
}

// ToggleSort toggles the sort direction on the current column.
func (v *Viewport) ToggleSort() {
	v.SortReverse = !v.SortReverse
	v.sortRows()
	v.updateVisibleRows()
}

// sortRows sorts AllRows based on the current sort column and direction.
func (v *Viewport) sortRows() {
	sort.Slice(v.AllRows, func(i, j int) bool {
		less := v.compareRows(v.AllRows[i], v.AllRows[j])
		if v.SortReverse {
			return !less
		}
		return less
	})
}

// compareRows compares two rows based on the current sort column.
func (v *Viewport) compareRows(a, b *domain.Row) bool {
	// Handle nil packages (for tests)
	if a.Package == nil || b.Package == nil {
		return false
	}

	switch v.SortColumn {
	case column.ColName:
		return strings.ToLower(a.Package.Name) < strings.ToLower(b.Package.Name)
	case column.ColVersion:
		return a.Package.Version < b.Package.Version
	case column.ColSize:
		return a.Package.InstalledSize < b.Package.InstalledSize
	case column.ColInstallDate:
		return a.Package.InstallDate.Before(b.Package.InstallDate)
	case column.ColGroups:
		aGroups := strings.Join(a.Package.Groups, ", ")
		bGroups := strings.Join(b.Package.Groups, ", ")
		return strings.ToLower(aGroups) < strings.ToLower(bGroups)
	default:
		// Default to name sorting
		return strings.ToLower(a.Package.Name) < strings.ToLower(b.Package.Name)
	}
}

// SortByColumn sorts by a specific column (used when selecting columns).
func (v *Viewport) SortByColumn(col column.Type) {
	if v.SortColumn == col {
		v.ToggleSort()
	} else {
		v.ApplySort(col, false)
	}
}

// ToggleSortCurrentColumn toggles sorting on the currently selected column.
func (v *Viewport) ToggleSortCurrentColumn() {
	if v.SelectedCol < 0 || v.SelectedCol >= len(v.Columns) {
		return
	}

	selectedColumn := v.Columns[v.SelectedCol]

	// Don't sort if column is not sortable
	if !selectedColumn.Sortable {
		return
	}

	v.SortByColumn(selectedColumn.Type)
}
