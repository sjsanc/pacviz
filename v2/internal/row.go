package internal

import "strings"

type Row struct {
	Cells map[ColType]string
}

func (r *Row) Matches(terms ...string) bool {
	for _, term := range terms {
		found := false
		for _, cell := range r.Cells {
			if strings.Contains(cell, term) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
