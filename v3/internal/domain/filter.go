package domain

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// FilterState represents the current filter configuration.
type FilterState struct {
	Active bool
	Terms  []string
	Column column.Type // empty = all columns
	Regex  bool
}

// NewFilterState creates a new filter state.
func NewFilterState() FilterState {
	return FilterState{
		Active: false,
		Terms:  make([]string, 0),
	}
}
