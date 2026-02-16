package renderer

import (
	"fmt"
	"strings"

	"github.com/sjsanc/pacviz/v3/internal/ui/styles"
)

func RenderStatus(preset string, totalRows, visibleRows, offset int, filter string, width int) string {
	start := offset + 1
	end := min(offset+visibleRows, totalRows)
	status := fmt.Sprintf("Preset: %s | Showing %d-%d of %d",
		preset,
		start,
		end,
		totalRows)

	if filter != "" {
		status += fmt.Sprintf(" | Filter: \"%s\"", filter)
	}

	return styles.Current.StatusBar.Width(width).Render(status)
}

func RenderStatusWithBuffer(buffer string, width int) string {
	return styles.Current.StatusBar.Width(width).Render(buffer)
}

func RenderRemoteStatus(query string, totalRows, visibleRows, offset int, filter string, loading bool, spinner string, errorMsg string, installing bool, installingPkg string, width int) string {
	var status string

	if installing {
		status = fmt.Sprintf("%s Installing %s...", spinner, installingPkg)
	} else if loading {
		status = fmt.Sprintf("%s Searching: %s", spinner, query)
	} else {
		start := offset + 1
		end := min(offset+visibleRows, totalRows)
		status = fmt.Sprintf("SEARCH: %s | Showing %d-%d of %d",
			query,
			start,
			end,
			totalRows)

		if filter != "" {
			status += fmt.Sprintf(" | Filter: \"%s\"", filter)
		}
	}

	if errorMsg != "" {
		statusLen := len(status)
		errorLen := len(errorMsg)
		padding := width - statusLen - errorLen - 4

		if padding > 0 {
			status = status + strings.Repeat(" ", padding) + "| " + errorMsg
		} else {
			status = errorMsg
		}
	}

	return styles.Current.RemoteStatusBar.Width(width).Render(status)
}

func RenderWarningStatus(message string, width int) string {
	return styles.Current.WarningStatusBar.Width(width).Render(message)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
