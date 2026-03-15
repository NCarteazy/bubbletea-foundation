package foundation

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/NCarteazy/bubbletea-foundation/theme"
)

func renderBreadcrumbs(crumbs []string, th theme.Theme) string {
	if len(crumbs) == 0 {
		return ""
	}

	sep := lipgloss.NewStyle().
		Foreground(lipgloss.Color(th.Muted())).
		Render(" > ")

	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(th.Muted()))

	lastStyle := th.Breadcrumb()

	var result string
	for i, crumb := range crumbs {
		if i > 0 {
			result += sep
		}
		if i == len(crumbs)-1 {
			result += lastStyle.Render(crumb)
		} else {
			result += mutedStyle.Render(crumb)
		}
	}

	return result
}
