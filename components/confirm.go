package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// Confirm is a y/n confirmation dialog component.
type Confirm struct {
	Prompt    string
	OnConfirm func() tea.Cmd
	th        theme.Theme
	done      bool
	confirmed bool
}

// NewConfirm creates a new Confirm with the given prompt and theme.
func NewConfirm(prompt string, th theme.Theme) *Confirm {
	return &Confirm{
		Prompt: prompt,
		th:     th,
	}
}

// Update handles key messages: y/Y confirms, n/N/esc cancels.
func (c *Confirm) Update(msg tea.Msg) (*Confirm, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.Type == tea.KeyRunes && (string(msg.Runes) == "y" || string(msg.Runes) == "Y"):
			c.done = true
			c.confirmed = true
			if c.OnConfirm != nil {
				return c, c.OnConfirm()
			}
		case msg.Type == tea.KeyRunes && (string(msg.Runes) == "n" || string(msg.Runes) == "N"):
			c.done = true
			c.confirmed = false
		case msg.Type == tea.KeyEsc:
			c.done = true
			c.confirmed = false
		}
	}
	return c, nil
}

// View renders the confirmation prompt in a panel with hint text.
func (c *Confirm) View() string {
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.th.Warning())).
		Bold(true)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.th.Muted()))

	prompt := promptStyle.Render(c.Prompt)
	hint := hintStyle.Render("[y] confirm  [n/esc] cancel")

	content := prompt + "\n\n" + hint

	return c.th.Panel().Render(content)
}

// Done returns whether the user has made a choice.
func (c *Confirm) Done() bool {
	return c.done
}

// Confirmed returns whether the user confirmed (pressed y/Y).
func (c *Confirm) Confirmed() bool {
	return c.confirmed
}
