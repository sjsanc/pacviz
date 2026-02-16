package domain

import (
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

type FilterState struct {
	Active bool
	Terms  []string
	Column column.Type
	Regex  bool
}

func NewFilterState() FilterState {
	return FilterState{
		Active: false,
		Terms:  make([]string, 0),
	}
}
