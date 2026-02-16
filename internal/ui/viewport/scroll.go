package viewport

func (v *Viewport) Scroll(delta int) {
	if len(v.VisibleRows) == 0 {
		return
	}

	newSelectedRow := v.SelectedRow + delta

	if newSelectedRow < 0 {
		newSelectedRow = 0
	} else if newSelectedRow >= len(v.VisibleRows) {
		newSelectedRow = len(v.VisibleRows) - 1
	}

	v.SelectedRow = newSelectedRow
	v.EnsureSelectionVisible()
}

func (v *Viewport) ScrollToTop() {
	if len(v.VisibleRows) == 0 {
		return
	}

	v.SelectedRow = 0
	v.Offset = 0
}

func (v *Viewport) ScrollToBottom() {
	if len(v.VisibleRows) == 0 {
		return
	}

	v.SelectedRow = len(v.VisibleRows) - 1

	if v.Height > 0 && len(v.VisibleRows) > v.Height {
		v.Offset = len(v.VisibleRows) - v.Height
	} else {
		v.Offset = 0
	}
}

// ScrollToLine jumps to a specific line number (0-indexed) and centers it in the viewport.
func (v *Viewport) ScrollToLine(line int) {
	if len(v.VisibleRows) == 0 {
		return
	}

	if line < 0 {
		line = 0
	} else if line >= len(v.VisibleRows) {
		line = len(v.VisibleRows) - 1
	}

	v.SelectedRow = line
	v.centerSelection()
}

func (v *Viewport) PageUp() {
	if v.Height == 0 {
		return
	}
	delta := v.Height / 2
	if delta < 1 {
		delta = 1
	}
	v.Scroll(-delta)
}

func (v *Viewport) PageDown() {
	if v.Height == 0 {
		return
	}
	delta := v.Height / 2
	if delta < 1 {
		delta = 1
	}
	v.Scroll(delta)
}
