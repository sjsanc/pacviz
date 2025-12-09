package viewport

import (
	"testing"

	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

func TestSelectNext(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 20)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	// Initial selection should be 0
	if v.SelectedRow != 0 {
		t.Errorf("Expected initial SelectedRow to be 0, got %d", v.SelectedRow)
	}

	// Move down
	v.SelectNext()
	if v.SelectedRow != 1 {
		t.Errorf("Expected SelectedRow to be 1, got %d", v.SelectedRow)
	}

	// Move to bottom of first viewport
	for i := 0; i < 8; i++ {
		v.SelectNext()
	}
	if v.SelectedRow != 9 {
		t.Errorf("Expected SelectedRow to be 9, got %d", v.SelectedRow)
	}

	// Moving beyond viewport should scroll
	v.SelectNext()
	if v.SelectedRow != 10 {
		t.Errorf("Expected SelectedRow to be 10, got %d", v.SelectedRow)
	}
	if v.Offset != 1 {
		t.Errorf("Expected Offset to be 1 after scrolling, got %d", v.Offset)
	}
}

func TestSelectPrev(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 20)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	// Start at row 10
	v.SelectRow(10)
	if v.SelectedRow != 10 {
		t.Errorf("Expected SelectedRow to be 10, got %d", v.SelectedRow)
	}

	// Move up
	v.SelectPrev()
	if v.SelectedRow != 9 {
		t.Errorf("Expected SelectedRow to be 9, got %d", v.SelectedRow)
	}

	// Move to top of viewport
	for i := 0; i < 8; i++ {
		v.SelectPrev()
	}
	if v.SelectedRow != 1 {
		t.Errorf("Expected SelectedRow to be 1, got %d", v.SelectedRow)
	}

	// Moving before viewport should scroll
	v.SelectPrev()
	if v.SelectedRow != 0 {
		t.Errorf("Expected SelectedRow to be 0, got %d", v.SelectedRow)
	}
	if v.Offset != 0 {
		t.Errorf("Expected Offset to be 0 after scrolling to top, got %d", v.Offset)
	}
}

func TestScrollToTop(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 20)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	// Start somewhere in the middle
	v.SelectedRow = 15
	v.Offset = 10

	v.ScrollToTop()

	if v.SelectedRow != 0 {
		t.Errorf("Expected SelectedRow to be 0, got %d", v.SelectedRow)
	}
	if v.Offset != 0 {
		t.Errorf("Expected Offset to be 0, got %d", v.Offset)
	}
}

func TestScrollToBottom(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 20)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	v.ScrollToBottom()

	if v.SelectedRow != 19 {
		t.Errorf("Expected SelectedRow to be 19, got %d", v.SelectedRow)
	}
	if v.Offset != 10 {
		t.Errorf("Expected Offset to be 10 (20-10), got %d", v.Offset)
	}
}

func TestEnsureSelectionVisible(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 50)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	// Test: Selection before viewport
	v.Offset = 20
	v.SelectedRow = 15
	v.EnsureSelectionVisible()
	if v.Offset != 15 {
		t.Errorf("Expected Offset to be 15 when selection is before viewport, got %d", v.Offset)
	}

	// Test: Selection after viewport
	v.Offset = 0
	v.SelectedRow = 15
	v.EnsureSelectionVisible()
	if v.Offset != 6 {
		t.Errorf("Expected Offset to be 6 (15-10+1) when selection is after viewport, got %d", v.Offset)
	}

	// Test: Selection within viewport (should not change offset)
	v.Offset = 10
	v.SelectedRow = 15
	v.EnsureSelectionVisible()
	if v.Offset != 10 {
		t.Errorf("Expected Offset to remain 10 when selection is within viewport, got %d", v.Offset)
	}
}

func TestPageUpDown(t *testing.T) {
	v := New()
	v.Height = 10

	// Create test rows
	rows := make([]*domain.Row, 50)
	for i := range rows {
		rows[i] = &domain.Row{
			Cells: make(map[column.Type]string),
		}
	}
	v.SetRows(rows)

	// Start at row 20
	v.SelectedRow = 20
	v.Offset = 15

	// Page down (half page = 5 rows)
	v.PageDown()
	if v.SelectedRow != 25 {
		t.Errorf("Expected SelectedRow to be 25 after PageDown, got %d", v.SelectedRow)
	}

	// Page up (half page = 5 rows)
	v.PageUp()
	if v.SelectedRow != 20 {
		t.Errorf("Expected SelectedRow to be 20 after PageUp, got %d", v.SelectedRow)
	}
}
