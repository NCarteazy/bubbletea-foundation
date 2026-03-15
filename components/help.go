package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// HelpKey describes a single key binding.
type HelpKey struct {
	Key  string
	Desc string
}

// HelpSection groups related key bindings under a title.
type HelpSection struct {
	Title string
	Keys  []HelpKey
}

// Help is a dismissible overlay that displays grouped keyboard shortcuts.
type Help struct {
	Sections []HelpSection
	width    int
	height   int
	th       theme.Theme
	done     bool
}

// NewHelp creates a new Help component with the given sections, dimensions, and theme.
func NewHelp(sections []HelpSection, width, height int, th theme.Theme) *Help {
	return &Help{
		Sections: sections,
		width:    width,
		height:   height,
		th:       th,
	}
}

// Update handles key messages: ?, esc, or q dismisses the help overlay.
func (h *Help) Update(msg tea.Msg) (*Help, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyEsc:
			h.done = true
		case msg.Type == tea.KeyRunes && string(msg.Runes) == "?":
			h.done = true
		case msg.Type == tea.KeyRunes && string(msg.Runes) == "q":
			h.done = true
		}
	}
	return h, nil
}

// View renders the help sections inside a panel.
func (h *Help) View() string {
	panelWidth := h.width - 10
	if panelWidth > 60 {
		panelWidth = 60
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Primary())).
		Bold(true)

	sectionTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Secondary())).
		Bold(true)

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Primary())).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Text()))

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(h.th.Muted()))

	var b strings.Builder

	b.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	b.WriteString("\n\n")

	for i, section := range h.Sections {
		b.WriteString(sectionTitleStyle.Render(section.Title))
		b.WriteByte('\n')
		for _, key := range section.Keys {
			b.WriteString("  ")
			b.WriteString(keyStyle.Render(key.Key))
			b.WriteString("  ")
			b.WriteString(descStyle.Render(key.Desc))
			b.WriteByte('\n')
		}
		if i < len(h.Sections)-1 {
			b.WriteByte('\n')
		}
	}

	b.WriteByte('\n')
	b.WriteString(footerStyle.Render("Press ? or esc to close"))

	panel := h.th.Panel().Width(panelWidth)
	return panel.Render(b.String())
}

// Done returns whether the help overlay has been dismissed.
func (h *Help) Done() bool {
	return h.done
}
