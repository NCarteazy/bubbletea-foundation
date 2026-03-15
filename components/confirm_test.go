package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

func TestConfirmRender(t *testing.T) {
	c := NewConfirm("Delete all files?", theme.TokyoNight)
	out := c.View()
	if !strings.Contains(out, "Delete all files?") {
		t.Fatalf("expected prompt text in View(), got: %s", out)
	}
	if !strings.Contains(out, "y") {
		t.Fatal("expected [y] hint in View()")
	}
	if !strings.Contains(out, "n") {
		t.Fatal("expected [n] hint in View()")
	}
}

func TestConfirmYes(t *testing.T) {
	called := false
	c := NewConfirm("Proceed?", theme.TokyoNight)
	c.OnConfirm = func() tea.Cmd { called = true; return nil }

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	if !c.Done() {
		t.Fatal("expected Done() to be true after y")
	}
	if !c.Confirmed() {
		t.Fatal("expected Confirmed() to be true after y")
	}
	if !called {
		t.Fatal("expected OnConfirm to be called")
	}
}

func TestConfirmNo(t *testing.T) {
	c := NewConfirm("Proceed?", theme.TokyoNight)

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	if !c.Done() {
		t.Fatal("expected Done() to be true after n")
	}
	if c.Confirmed() {
		t.Fatal("expected Confirmed() to be false after n")
	}
}

func TestConfirmEsc(t *testing.T) {
	c := NewConfirm("Proceed?", theme.TokyoNight)

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !c.Done() {
		t.Fatal("expected Done() to be true after esc")
	}
	if c.Confirmed() {
		t.Fatal("expected Confirmed() to be false after esc")
	}
}
