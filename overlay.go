package foundation

import tea "github.com/charmbracelet/bubbletea"

type Overlay interface {
	Update(msg tea.Msg) (Overlay, tea.Cmd)
	View() string
	Done() bool
}

type ShowOverlayMsg struct {
	Overlay Overlay
}

func ShowOverlay(o Overlay) tea.Cmd {
	return func() tea.Msg { return ShowOverlayMsg{Overlay: o} }
}

type overlayStack struct {
	entries []Overlay
}

func (s *overlayStack) push(o Overlay) {
	s.entries = append(s.entries, o)
}

func (s *overlayStack) pop() {
	if len(s.entries) == 0 {
		return
	}
	s.entries = s.entries[:len(s.entries)-1]
}

func (s *overlayStack) top() Overlay {
	if len(s.entries) == 0 {
		return nil
	}
	return s.entries[len(s.entries)-1]
}

func (s *overlayStack) active() bool {
	return len(s.entries) > 0
}

func (s *overlayStack) len() int {
	return len(s.entries)
}
