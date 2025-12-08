package viewport

// Scroll moves the viewport by delta rows.
func (v *Viewport) Scroll(delta int) {
	// TODO: Update offset and selectedRow
	// TODO: Bounds checking
}

// ScrollToTop jumps to the first row.
func (v *Viewport) ScrollToTop() {
	// TODO: Set offset=0, selectedRow=0
}

// ScrollToBottom jumps to the last row.
func (v *Viewport) ScrollToBottom() {
	// TODO: Set to last row
}

// ScrollToLine jumps to a specific line number.
func (v *Viewport) ScrollToLine(line int) {
	// TODO: Set selectedRow, adjust offset if needed
}
