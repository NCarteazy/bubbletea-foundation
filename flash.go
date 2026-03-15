package foundation

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// FlashMsg is sent when a flash message should be displayed.
type FlashMsg struct {
	Message  string
	Duration int // ticks before auto-clear
}

// Flash returns a tea.Cmd that emits a FlashMsg.
func Flash(msg string, duration int) tea.Cmd {
	return func() tea.Msg { return FlashMsg{Message: msg, Duration: duration} }
}

type flashTickMsg struct{}

func flashTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg { return flashTickMsg{} })
}

type flashState struct {
	message string
	timer   int
}

func (f *flashState) set(msg string, duration int) {
	f.message = msg
	f.timer = duration
}

func (f *flashState) tick() {
	if f.timer > 0 {
		f.timer--
		if f.timer == 0 {
			f.message = ""
		}
	}
}

func (f *flashState) active() bool {
	return f.timer > 0
}
