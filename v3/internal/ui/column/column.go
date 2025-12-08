package column

// Type represents a column identifier.
type Type string

const (
	ColIndex       Type = "index"
	ColName        Type = "name"
	ColVersion     Type = "version"
	ColSize        Type = "size"
	ColInstallDate Type = "install_date"
	ColGroups      Type = "groups"
	ColDescription Type = "description"
)

// WidthType specifies how column width is calculated.
type WidthType int

const (
	WidthFixed WidthType = iota
	WidthPercent
	WidthAuto
)

// ColumnWidth defines width calculation parameters.
type ColumnWidth struct {
	Type    WidthType
	Min     int
	Max     int
	Size    int // pixels or percent
}

// Column represents a table column configuration.
type Column struct {
	Type       Type
	Name       string
	Width      ColumnWidth
	Sortable   bool
	Searchable bool
	Visible    bool
}

// DefaultColumns returns the default column configuration.
func DefaultColumns() []*Column {
	return []*Column{
		{
			Type:       ColIndex,
			Name:       "#",
			Width:      ColumnWidth{Type: WidthFixed, Size: 5}, // Fixed 5 char width
			Sortable:   false,
			Searchable: false,
			Visible:    true,
		},
		{
			Type:       ColName,
			Name:       "Name",
			Width:      ColumnWidth{Type: WidthFixed, Size: 40}, // 2x the size of other fixed columns
			Sortable:   true,
			Searchable: true,
			Visible:    true,
		},
		{
			Type:       ColVersion,
			Name:       "Version",
			Width:      ColumnWidth{Type: WidthFixed, Size: 20}, // Fits versions like "20250814.1-1"
			Sortable:   true,
			Searchable: false,
			Visible:    true,
		},
		{
			Type:       ColSize,
			Name:       "Size",
			Width:      ColumnWidth{Type: WidthFixed, Size: 10}, // Fits sizes like "999.9 MB"
			Sortable:   true,
			Searchable: false,
			Visible:    true,
		},
		{
			Type:       ColInstallDate,
			Name:       "Install Date",
			Width:      ColumnWidth{Type: WidthFixed, Size: 12}, // Fits "2025-11-30"
			Sortable:   true,
			Searchable: false,
			Visible:    true,
		},
		{
			Type:       ColGroups,
			Name:       "Groups",
			Width:      ColumnWidth{Type: WidthFixed, Size: 25}, // Variable but typically short
			Sortable:   true,
			Searchable: true,
			Visible:    true,
		},
		{
			Type:       ColDescription,
			Name:       "Description",
			Width:      ColumnWidth{Type: WidthAuto, Min: 20}, // Grows to fill remaining space
			Sortable:   false,
			Searchable: true,
			Visible:    true,
		},
	}
}
