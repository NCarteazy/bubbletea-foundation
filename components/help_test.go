package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

func testSections() []HelpSection {
	return []HelpSection{
		{
			Title: "Navigation",
			Keys: []HelpKey{
				{Key: "j/k", Desc: "Move up/down"},
				{Key: "enter", Desc: "Select item"},
			},
		},
		{
			Title: "Actions",
			Keys: []HelpKey{
				{Key: "d", Desc: "Delete"},
				{Key: "q", Desc: "Quit"},
			},
		},
	}
}

func TestHelpRender(t *testing.T) {
	h := NewHelp(testSections(), 80, 40, theme.TokyoNight)
	out := h.View()

	if !strings.Contains(out, "Navigation") {
		t.Fatal("expected section title 'Navigation' in View()")
	}
	if !strings.Contains(out, "Actions") {
		t.Fatal("expected section title 'Actions' in View()")
	}
	if !strings.Contains(out, "j/k") {
		t.Fatal("expected key 'j/k' in View()")
	}
	if !strings.Contains(out, "Move up/down") {
		t.Fatal("expected desc 'Move up/down' in View()")
	}
	if !strings.Contains(out, "Keyboard Shortcuts") {
		t.Fatal("expected title 'Keyboard Shortcuts' in View()")
	}
}

func TestHelpDismissEsc(t *testing.T) {
	h := NewHelp(testSections(), 80, 40, theme.TokyoNight)

	h, _ = h.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !h.Done() {
		t.Fatal("expected Done() to be true after esc")
	}
}

func TestHelpDismissQuestion(t *testing.T) {
	h := NewHelp(testSections(), 80, 40, theme.TokyoNight)

	h, _ = h.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !h.Done() {
		t.Fatal("expected Done() to be true after ?")
	}
}
