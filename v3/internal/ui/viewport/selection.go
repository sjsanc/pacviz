package viewport

// SelectRow changes the selected row.
func (v *Viewport) SelectRow(index int) {
	// TODO: Set selectedRow with bounds checking
	// TODO: Adjust offset if needed to keep selection visible
}

// SelectColumn changes the selected column.
func (v *Viewport) SelectColumn(index int) {
	// TODO: Set selectedCol with bounds checking
}

// NextColumn moves to the next column.
func (v *Viewport) NextColumn() {
	// TODO: Increment selectedCol (wrap around)
}

// PrevColumn moves to the previous column.
func (v *Viewport) PrevColumn() {
	// TODO: Decrement selectedCol (wrap around)
}
