package viewport

import (
	"testing"

	"github.com/sjsanc/pacviz/v3/internal/domain"
)

func TestScrollToLine_Centering(t *testing.T) {
	tests := []struct {
		name           string
		totalRows      int
		viewportHeight int
		targetLine     int
		expectedRow    int
		expectedOffset int
	}{
		{
			name:           "center line in middle of large list",
			totalRows:      100,
			viewportHeight: 20,
			targetLine:     50,
			expectedRow:    50,
			expectedOffset: 40, // 50 - 20/2 = 40
		},
		{
			name:           "line near start",
			totalRows:      100,
			viewportHeight: 20,
			targetLine:     5,
			expectedRow:    5,
			expectedOffset: 0, // Can't go negative
		},
		{
			name:           "line near end",
			totalRows:      100,
			viewportHeight: 20,
			targetLine:     95,
			expectedRow:    95,
			expectedOffset: 80, // maxOffset = 100 - 20 = 80
		},
		{
			name:           "line at start",
			totalRows:      100,
			viewportHeight: 20,
			targetLine:     0,
			expectedRow:    0,
			expectedOffset: 0,
		},
		{
			name:           "line at end",
			totalRows:      100,
			viewportHeight: 20,
			targetLine:     99,
			expectedRow:    99,
			expectedOffset: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			v.Height = tt.viewportHeight

			// Create dummy rows
			rows := make([]*domain.Row, tt.totalRows)
			for i := 0; i < tt.totalRows; i++ {
				rows[i] = &domain.Row{}
			}
			v.SetRows(rows)

			// Execute
			v.ScrollToLine(tt.targetLine)

			// Verify
			if v.SelectedRow != tt.expectedRow {
				t.Errorf("SelectedRow = %d, want %d", v.SelectedRow, tt.expectedRow)
			}
			if v.Offset != tt.expectedOffset {
				t.Errorf("Offset = %d, want %d", v.Offset, tt.expectedOffset)
			}
		})
	}
}
