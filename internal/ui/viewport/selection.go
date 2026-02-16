package viewport

func (v *Viewport) SelectRow(index int) {
	if len(v.VisibleRows) == 0 {
		v.SelectedRow = 0
		return
	}

	if index < 0 {
		index = 0
	} else if index >= len(v.VisibleRows) {
		index = len(v.VisibleRows) - 1
	}

	v.SelectedRow = index
	v.EnsureSelectionVisible()
}

func (v *Viewport) SelectNext() {
	if len(v.VisibleRows) == 0 {
		return
	}

	if v.SelectedRow < len(v.VisibleRows)-1 {
		v.SelectedRow++
		v.EnsureSelectionVisible()
	}
}

func (v *Viewport) SelectPrev() {
	if len(v.VisibleRows) == 0 {
		return
	}

	if v.SelectedRow > 0 {
		v.SelectedRow--
		v.EnsureSelectionVisible()
	}
}

func (v *Viewport) EnsureSelectionVisible() {
	if v.Height == 0 {
		return
	}

	if v.SelectedRow < v.Offset {
		v.Offset = v.SelectedRow
	}

	if v.SelectedRow >= v.Offset+v.Height {
		v.Offset = v.SelectedRow - v.Height + 1
	}

	if v.Offset < 0 {
		v.Offset = 0
	}
}

func (v *Viewport) centerSelection() {
	if v.Height == 0 {
		return
	}

	targetOffset := v.SelectedRow - (v.Height / 2)

	if targetOffset < 0 {
		targetOffset = 0
	}

	maxOffset := len(v.VisibleRows) - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if targetOffset > maxOffset {
		targetOffset = maxOffset
	}

	v.Offset = targetOffset
}

func (v *Viewport) SelectColumn(index int) {
	if len(v.Columns) == 0 {
		v.SelectedCol = 1
		return
	}

	if index < 1 {
		index = 1
	} else if index >= len(v.Columns) {
		index = len(v.Columns) - 1
	}

	v.SelectedCol = index
}

func (v *Viewport) NextColumn() {
	if len(v.Columns) == 0 {
		return
	}

	v.SelectedCol++
	if v.SelectedCol >= len(v.Columns) {
		v.SelectedCol = 1
	}
}

func (v *Viewport) PrevColumn() {
	if len(v.Columns) == 0 {
		return
	}

	v.SelectedCol--
	if v.SelectedCol < 1 {
		v.SelectedCol = len(v.Columns) - 1
	}
}
