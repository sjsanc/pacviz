package main

type ColType int

const (
	Fixed ColType = iota
	Percent
)

type ColDef struct {
	Name string
	Type ColType
	Size int // px | %
}

func calcColWidths(cols []ColDef, maxWidth int) []int {
	var totalFixed int
	var totalPercent int

	for _, col := range cols {
		switch col.Type {
		case Fixed:
			totalFixed += col.Size
		case Percent:
			totalPercent += col.Size
		}
	}

	totalAvailable := maxWidth - totalFixed
	widths := make([]int, len(cols))

	for i, col := range cols {
		switch col.Type {
		case Fixed:
			widths[i] = col.Size
		case Percent:
			widths[i] = (col.Size * totalAvailable) / 100
		}
	}

	combinedWidth := 0
	for _, w := range widths {
		combinedWidth += w
	}

	// Distribute remaining width to last column
	if combinedWidth < maxWidth {
		widths[len(widths)-1] += maxWidth - combinedWidth
	}

	return widths
}
