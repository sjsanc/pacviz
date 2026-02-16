package column

// CalculateWidths computes actual pixel widths for columns based on terminal width.
func CalculateWidths(columns []*Column, terminalWidth int) []int {
	widths := make([]int, len(columns))
	remainingWidth := terminalWidth
	autoColumns := make([]int, 0)

	for i, col := range columns {
		if !col.Visible {
			widths[i] = 0
			continue
		}

		if col.Width.Type == WidthFixed {
			widths[i] = col.Width.Size
			remainingWidth -= col.Width.Size
		}
	}

	for i, col := range columns {
		if !col.Visible || col.Width.Type != WidthPercent {
			continue
		}

		width := (terminalWidth * col.Width.Size) / 100

		if col.Width.Min > 0 && width < col.Width.Min {
			width = col.Width.Min
		}
		if col.Width.Max > 0 && width > col.Width.Max {
			width = col.Width.Max
		}

		widths[i] = width
		remainingWidth -= width
	}

	for i, col := range columns {
		if !col.Visible {
			continue
		}
		if col.Width.Type == WidthAuto {
			autoColumns = append(autoColumns, i)
		}
	}

	if len(autoColumns) > 0 && remainingWidth > 0 {
		widthPerAuto := remainingWidth / len(autoColumns)
		remainder := remainingWidth % len(autoColumns)

		for i, colIdx := range autoColumns {
			width := widthPerAuto
			if i == len(autoColumns)-1 {
				width += remainder
			}

			col := columns[colIdx]

			if col.Width.Min > 0 && width < col.Width.Min {
				width = col.Width.Min
			}
			if col.Width.Max > 0 && width > col.Width.Max {
				width = col.Width.Max
			}

			widths[colIdx] = width
		}
	}

	return widths
}
