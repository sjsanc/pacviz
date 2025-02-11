package main

type Row struct {
	cells []string
}

type Viewport struct {
	rows     []*Row
	offset   int // index of the first visible row relative to all rows
	height   int // number of visible rows
	selected int // index of the selected row relative to all rows
}

func NewViewport(rows []*Row) *Viewport {
	return &Viewport{
		rows: rows,
	}
}

func (vp *Viewport) GoTo(line int) {
	// Ensure the line is within bounds
	if line < 0 {
		line = 0
	} else if line >= len(vp.rows) {
		line = len(vp.rows) - 1
	}

	if line < vp.offset {
		vp.offset = line
	} else if line >= vp.offset+vp.height {
		vp.offset = line - vp.height + 1
	}

	vp.selected = line
}

func (vp *Viewport) ScrollUp(rows int) {
	vp.GoTo(vp.selected - rows)
}
func (vp *Viewport) ScrollDown(rows int) {
	vp.GoTo(vp.selected + rows)
}

func (vp *Viewport) GetRow(index int) *Row {
	if index >= len(vp.rows) {
		return nil
	}
	return vp.rows[index]
}
func (vp *Viewport) GetSelectedRow() *Row {
	return vp.GetRow(vp.selected)
}

func (vp *Viewport) SetHeight(h int) {
	vp.height = h - 2 // Account for header and footer
}
func (vp *Viewport) SetRows(rows []*Row) {
	vp.offset = 0
	vp.rows = rows
}

func (vp *Viewport) VisibleRows() []*Row {
	if vp.height >= len(vp.rows) {
		return vp.rows
	}
	return vp.rows[vp.offset : vp.offset+vp.height]
}
