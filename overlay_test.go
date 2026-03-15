package foundation

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type mockOverlay struct {
	done    bool
	content string
}

func (m mockOverlay) Update(msg tea.Msg) (Overlay, tea.Cmd) { return m, nil }
func (m mockOverlay) View() string                          { return m.content }
func (m mockOverlay) Done() bool                            { return m.done }

func TestOverlayStack_PushAndTop(t *testing.T) {
	s := &overlayStack{}

	first := mockOverlay{content: "first"}
	second := mockOverlay{content: "second"}

	s.push(first)
	s.push(second)

	if s.len() != 2 {
		t.Fatalf("expected len=2, got %d", s.len())
	}
	if s.top().View() != "second" {
		t.Fatalf("expected top='second', got %q", s.top().View())
	}

	s.pop()

	if s.len() != 1 {
		t.Fatalf("expected len=1 after pop, got %d", s.len())
	}
	if s.top().View() != "first" {
		t.Fatalf("expected top='first' after pop, got %q", s.top().View())
	}
}

func TestOverlayStack_Empty(t *testing.T) {
	s := &overlayStack{}

	if s.active() {
		t.Fatal("expected active()=false on empty stack")
	}
	if s.top() != nil {
		t.Fatal("expected top()=nil on empty stack")
	}

	// pop on empty stack should not panic
	s.pop()
}
