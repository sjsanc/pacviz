package internal

import "sort"

// Viewport manages state relating to table's rows.
type Viewport struct {
	Rows            []*Row
	Offset          int
	Height          int
	CurrentRowIndex int
	CurrentColIndex int
	SortingBy       ColType
	SortReversed    bool
	Filter          []string
}

func NewViewport() *Viewport {
	return &Viewport{}
}

// ================================================================================
// ### SCROLLING
// ================================================================================

func (vp *Viewport) GoTo(idx int) {
	if idx < 0 {
		idx = 0
	} else if idx >= len(vp.Rows) {
		idx = len(vp.Rows) - 1
	}

	if idx < vp.Offset {
		vp.Offset = idx
	} else if idx >= vp.Offset+vp.Height {
		vp.Offset = idx - vp.Height + 1
	}

	vp.CurrentRowIndex = idx
}
func (vp *Viewport) ScrollUp(rows int) {
	vp.GoTo(vp.CurrentRowIndex - rows)
}
func (vp *Viewport) ScrollDown(rows int) {
	vp.GoTo(vp.CurrentRowIndex + rows)
}

// ================================================================================
// ### ROWS
// ================================================================================

func (vp *Viewport) SetRows(rows []*Row) {
	vp.Offset = 0 // Reset view
	vp.Rows = rows
}
func (vp *Viewport) GetRow(idx int) *Row {
	if idx >= len(vp.Rows) {
		return nil
	}
	return vp.Rows[idx]
}
func (vp *Viewport) GetSelectedRow() *Row {
	return vp.GetRow(vp.CurrentRowIndex)
}
func (vp *Viewport) GetVisibleRows() []*Row {
	if vp.Height >= len(vp.Rows) {
		return vp.Rows
	}

	return vp.Rows[vp.Offset : vp.Offset+vp.Height]
}

// ================================================================================
// ### COLUMNS
// ================================================================================

func (vp *Viewport) NextColumn() {
	vp.CurrentColIndex++
	if vp.CurrentColIndex >= len(COLUMNS) {
		vp.CurrentColIndex = 0
	}
}

func (vp *Viewport) PrevColumn() {
	vp.CurrentColIndex--
	if vp.CurrentColIndex < 0 {
		vp.CurrentColIndex = len(COLUMNS) - 1
	}
}

func (vp *Viewport) GetCurrentColumn() ColType {
	return COLUMNS[vp.CurrentColIndex]
}

// ================================================================================
// ### SORTING
// ================================================================================

func (vp *Viewport) ApplySort(sortBy ColType, reverse bool) {
	vp.SortingBy = sortBy
	vp.SortReversed = reverse

	// TODO: This is a naive implementation. Some columns won't sort properly.
	// For example, sorting by version number will put 9 after 24 (9 > 2).
	sort.Slice(vp.Rows, func(i, j int) bool {
		if reverse {
			return vp.Rows[i].Cells[sortBy] > vp.Rows[j].Cells[sortBy]
		}
		return vp.Rows[i].Cells[sortBy] < vp.Rows[j].Cells[sortBy]
	})
}

// ================================================================================
// ### FILTERING
// ================================================================================

func (vp *Viewport) ApplyFilter(terms ...string) {
	vp.Filter = terms
	var filtered []*Row
	for _, row := range vp.Rows {
		if row.Matches(terms...) {
			filtered = append(filtered, row)
		}
	}
	vp.SetRows(filtered)
}

// ================================================================================
// ### PRESETS
// ================================================================================

func (vp *Viewport) ApplyPreset(preset string) {
	switch preset {
	case "explicit":
		vp.SetRows(pkgsToRows(PM.GetExplicitPkgs()))
	case "orphans":
		vp.SetRows(pkgsToRows(PM.GetOrphanPkgs()))
	}
}

// ================================================================================
// ### SYNC SEARCH
// ================================================================================

func (vp *Viewport) SearchSyncDBs(terms ...string) {
	vp.SetRows(pkgsToRows(PM.SearchSyncDBs(terms...)))
}
