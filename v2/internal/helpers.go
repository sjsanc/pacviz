package internal

import (
	"github.com/Jguer/go-alpm/v2"
)

func pkgsToRows(pkgs []alpm.IPackage) []*Row {
	rows := make([]*Row, len(pkgs))
	for i, pkg := range pkgs {
		rows[i] = &Row{
			Cells: map[ColType]string{
				ColName: pkg.Name(),
				ColVer:  pkg.Version(),
				ColDate: pkg.InstallDate().Format("2006-01-02"),
				ColDesc: pkg.Description(),
			},
		}
	}
	return rows
}
