package foundation

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// KeyHint represents a keyboard shortcut hint shown in the status bar.
type KeyHint struct {
	Key  string
	Desc string
}

func renderStatusBar(hints []KeyHint, flash *flashState, width int, th theme.Theme) string {
	barStyle := th.StatusBar().Width(width)

	if flash != nil && flash.active() {
		flashStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(th.Success())).
			Bold(true)
		content := flashStyle.Render(flash.message)
		return barStyle.Align(lipgloss.Center).Render(content)
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(th.Primary())).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(th.Muted()))

	sepStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(th.Muted()))

	var parts []string
	for _, h := range hints {
		parts = append(parts, keyStyle.Render(h.Key)+" "+descStyle.Render(h.Desc))
	}

	content := strings.Join(parts, sepStyle.Render(" │ "))

	return barStyle.Align(lipgloss.Center).Render(content)
}
