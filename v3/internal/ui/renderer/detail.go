package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sjsanc/pacviz/v3/internal/domain"
	"github.com/sjsanc/pacviz/v3/internal/ui/column"
	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

// detailField defines a field to display in the detail panel.
type detailField struct {
	label string
	colType column.Type
}

var detailPanelFields = []detailField{
	{label: "Name", colType: column.ColName},
	{label: "Version", colType: column.ColVersion},
	{label: "Repository", colType: column.ColRepo},
	{label: "Architecture", colType: column.ColArchitecture},
	{label: "Installed", colType: column.ColInstalled},
	{label: "Install Date", colType: column.ColInstallDate},
	{label: "Install Reason", colType: column.ColInstallReason},
	{label: "Is Orphan", colType: column.ColIsOrphan},
	{label: "Is Foreign", colType: column.ColIsForeign},
	{label: "Has Update", colType: column.ColHasUpdate},
	{label: "New Version", colType: column.ColNewVersion},
	{label: "Description", colType: column.ColDescription},
	{label: "URL", colType: column.ColURL},
	{label: "Packager", colType: column.ColPackager},
	{label: "Build Date", colType: column.ColBuildDate},
	{label: "Licenses", colType: column.ColLicenses},
	{label: "Size", colType: column.ColSize},
	{label: "Groups", colType: column.ColGroups},
	{label: "Dependencies", colType: column.ColDependencies},
	{label: "Optional Dependencies", colType: column.ColOptDepends},
	{label: "Required By", colType: column.ColRequired},
	{label: "Provides", colType: column.ColProvides},
	{label: "Conflicts", colType: column.ColConflicts},
	{label: "Replaces", colType: column.ColReplaces},
}

// RenderDetailPanel renders the package detail panel as an overlay above the status bar.
// If isRemote is true, show install commands at the bottom using remote colors.
func RenderDetailPanel(pkg *domain.Package, columns []*column.Column, colWidths []int, width int, _ bool, isRemote bool) string {
	if pkg == nil {
		return ""
	}

	// Style for labels (same as index column)
	labelStyle := styles.Current.Index
	valueStyle := lipgloss.NewStyle().Foreground(styles.Current.Foreground)

	// Create a temporary row to get formatted cell values
	row := domain.PackageToRow(pkg, 0)

	// Collect all fields to display
	var fields []struct {
		label string
		value string
	}

	for _, field := range detailPanelFields {
		value := row.Cells[field.colType]

		// Skip empty values for non-essential fields
		if value == "" && field.colType != column.ColHasUpdate &&
			field.colType != column.ColIsOrphan && field.colType != column.ColIsForeign {
			continue
		}

		fields = append(fields, struct {
			label string
			value string
		}{label: field.label, value: value})
	}

	// Split fields into two columns
	mid := (len(fields) + 1) / 2
	leftFields := fields[:mid]
	rightFields := fields[mid:]

	// Calculate column width (half of available width, minus padding and gap)
	colWidth := (width - 8) / 2

	// Build left and right column content
	var leftContent, rightContent strings.Builder

	for _, f := range leftFields {
		line := labelStyle.Render(f.label+":") + " " + valueStyle.Render(f.value)
		if leftContent.Len() > 0 {
			leftContent.WriteString("\n")
		}
		leftContent.WriteString(line)
	}

	for _, f := range rightFields {
		line := labelStyle.Render(f.label+":") + " " + valueStyle.Render(f.value)
		if rightContent.Len() > 0 {
			rightContent.WriteString("\n")
		}
		rightContent.WriteString(line)
	}

	// Style columns with max width to prevent overflow
	leftStyle := lipgloss.NewStyle().
		Width(colWidth).
		MaxWidth(colWidth)
	rightStyle := lipgloss.NewStyle().
		Width(colWidth).
		MaxWidth(colWidth)

	leftCol := leftStyle.Render(leftContent.String())
	rightCol := rightStyle.Render(rightContent.String())

	// Join columns horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// Add install commands at the bottom if in remote mode
	if isRemote {
		content = content + "\n\n" + renderInstallCommands(pkg.Name)
	}

	// Render as overlay above status bar
	panelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Current.Accent1).
		Padding(0, 1).
		Width(width - 4)

	return panelStyle.Render(content)
}


// renderInstallCommands renders the install command options using remote colors.
func renderInstallCommands(pkgName string) string {
	commandStyle := lipgloss.NewStyle().
		Foreground(styles.Current.Background).
		Background(styles.Current.RemoteAccent).
		Bold(true).
		Padding(0, 1)

	textStyle := lipgloss.NewStyle().
		Foreground(styles.Current.RemoteAccent)

	cmd1 := commandStyle.Render("i")
	cmd2 := commandStyle.Render(":install")
	cmd3 := commandStyle.Render("Ctrl+I")

	text := textStyle.Render(" to install ") + lipgloss.NewStyle().Bold(true).Foreground(styles.Current.Foreground).Render(pkgName)

	return cmd1 + " " + cmd2 + " " + cmd3 + text
}
