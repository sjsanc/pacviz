package domain

import (
	"fmt"
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/ui/column"
)

func PackagesToRows(packages []*Package) []*Row {
	rows := make([]*Row, 0, len(packages))
	for idx, pkg := range packages {
		rows = append(rows, PackageToRow(pkg, idx+1))
	}
	return rows
}

func PackageToRow(pkg *Package, index int) *Row {
	row := NewRow(pkg)

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
	row.Cells[column.ColDeps] = fmt.Sprintf("%d", pkg.DependencyCount)
	row.Cells[column.ColGroups] = strings.Join(pkg.Groups, ", ")
	row.Cells[column.ColDescription] = pkg.Description

	row.Cells[column.ColURL] = pkg.URL
	row.Cells[column.ColLicenses] = strings.Join(pkg.Licenses, ", ")
	row.Cells[column.ColArchitecture] = pkg.Architecture
	row.Cells[column.ColPackager] = pkg.Packager
	row.Cells[column.ColBuildDate] = pkg.BuildDate.Format("2006-01-02")
	row.Cells[column.ColDependencies] = strings.Join(pkg.Dependencies, ", ")
	row.Cells[column.ColOptDepends] = formatOptDepends(pkg.OptDepends)
	row.Cells[column.ColConflicts] = strings.Join(pkg.Conflicts, ", ")
	row.Cells[column.ColProvides] = strings.Join(pkg.Provides, ", ")
	row.Cells[column.ColReplaces] = strings.Join(pkg.Replaces, ", ")
	row.Cells[column.ColRequired] = strings.Join(pkg.Required, ", ")
	row.Cells[column.ColInstallReason] = formatInstallReason(pkg.InstallReason)
	row.Cells[column.ColIsOrphan] = formatBool(pkg.IsOrphan)
	row.Cells[column.ColIsForeign] = formatBool(pkg.IsForeign)
	row.Cells[column.ColHasUpdate] = formatBool(pkg.HasUpdate)
	row.Cells[column.ColNewVersion] = pkg.NewVersion

	return row
}

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

func formatOptDepends(optDepends map[string]string) string {
	if len(optDepends) == 0 {
		return ""
	}
	var deps []string
	for name, desc := range optDepends {
		if desc != "" {
			deps = append(deps, fmt.Sprintf("%s: %s", name, desc))
		} else {
			deps = append(deps, name)
		}
	}
	return strings.Join(deps, ", ")
}

func formatInstallReason(reason InstallReason) string {
	switch reason {
	case ReasonExplicit:
		return "Explicit"
	case ReasonDependency:
		return "Dependency"
	default:
		return "Unknown"
	}
}

func formatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
