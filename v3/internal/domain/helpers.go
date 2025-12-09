package domain

import (
	"fmt"
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

// PackagesToRows converts packages to table rows.
func PackagesToRows(packages []*Package) []*Row {
	rows := make([]*Row, 0, len(packages))
	for idx, pkg := range packages {
		rows = append(rows, PackageToRow(pkg, idx+1))
	}
	return rows
}

// PackageToRow converts a single package to a row.
func PackageToRow(pkg *Package, index int) *Row {
	row := NewRow(pkg)

	// Format cells
	row.Cells[column.ColIndex] = fmt.Sprintf("%d", index)
	row.Cells[column.ColRepo] = pkg.Repository
	row.Cells[column.ColName] = pkg.Name
	row.Cells[column.ColVersion] = pkg.Version
	row.Cells[column.ColSize] = formatSize(pkg.InstalledSize)
	row.Cells[column.ColInstallDate] = pkg.InstallDate.Format("2006-01-02")
	if pkg.Installed {
		row.Cells[column.ColInstalled] = "Yes"
	} else {
		row.Cells[column.ColInstalled] = "No"
	}
	row.Cells[column.ColGroups] = strings.Join(pkg.Groups, ", ")
	row.Cells[column.ColDescription] = pkg.Description

	return row
}

// formatSize formats bytes into human-readable size.
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
