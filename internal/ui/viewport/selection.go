package viewport

// SelectRow changes the selected row to an absolute index.
func (v *Viewport) SelectRow(index int) {
	if len(v.VisibleRows) == 0 {
		v.SelectedRow = 0
		return
	}

	// Bounds checking
	if index < 0 {
		index = 0
	} else if index >= len(v.VisibleRows) {
		index = len(v.VisibleRows) - 1
	}

	v.SelectedRow = index
	v.EnsureSelectionVisible()
}

// SelectNext moves the selection down by one row.
func (v *Viewport) SelectNext() {
	if len(v.VisibleRows) == 0 {
		return
	}

	if v.SelectedRow < len(v.VisibleRows)-1 {
		v.SelectedRow++
		v.EnsureSelectionVisible()
	}
}

// SelectPrev moves the selection up by one row.
func (v *Viewport) SelectPrev() {
	if len(v.VisibleRows) == 0 {
		return
	}

	if v.SelectedRow > 0 {
		v.SelectedRow--
		v.EnsureSelectionVisible()
	}
}

// EnsureSelectionVisible adjusts the viewport offset to keep the selected row visible.
func (v *Viewport) EnsureSelectionVisible() {
	if v.Height == 0 {
		return
	}

	// If selection is before the viewport, scroll up
	if v.SelectedRow < v.Offset {
		v.Offset = v.SelectedRow
	}

	// If selection is after the viewport, scroll down
	if v.SelectedRow >= v.Offset+v.Height {
		v.Offset = v.SelectedRow - v.Height + 1
	}

	// Ensure offset doesn't go negative
	if v.Offset < 0 {
		v.Offset = 0
	}
}

// centerSelection adjusts the viewport offset to center the selected row.
func (v *Viewport) centerSelection() {
	if v.Height == 0 {
		return
	}

	// Try to center the selected row in the viewport
	targetOffset := v.SelectedRow - (v.Height / 2)

	// Bounds checking
	if targetOffset < 0 {
		targetOffset = 0
	}

	// Don't scroll past the end
	maxOffset := len(v.VisibleRows) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if targetOffset > maxOffset {
		targetOffset = maxOffset
	}

	v.Offset = targetOffset
}

// SelectColumn changes the selected column.
func (v *Viewport) SelectColumn(index int) {
	if len(v.Columns) == 0 {
		v.SelectedCol = 1 // Skip index column
		return
	}

	// Bounds checking - never allow index 0 (index column)
	if index < 1 {
		index = 1
	} else if index >= len(v.Columns) {
		index = len(v.Columns) - 1
	}

	v.SelectedCol = index
}

// NextColumn moves to the next column.
func (v *Viewport) NextColumn() {
	if len(v.Columns) == 0 {
		return
	}

	v.SelectedCol++
	if v.SelectedCol >= len(v.Columns) {
		v.SelectedCol = 1 // Wrap around to first selectable column (skip index)
	}
}

// PrevColumn moves to the previous column.
func (v *Viewport) PrevColumn() {
	if len(v.Columns) == 0 {
		return
	}

	v.SelectedCol--
	if v.SelectedCol < 1 { // Don't allow index 0 (index column)
		v.SelectedCol = len(v.Columns) - 1 // Wrap around
	}
}
