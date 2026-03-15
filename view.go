package foundation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

type View interface {
	ID() string
	Update(msg tea.Msg, ctx ViewContext) (View, tea.Cmd)
	Render(ctx ViewContext) string
}

type ViewContext struct {
	Width  int
	Height int
	Theme  theme.Theme
	Data   any
}

// Navigation messages
type NavigateMsg struct {
	ViewID string
	Data   any
}

type BackMsg struct{}

type ReplaceMsg struct {
	ViewID string
	Data   any
}

// Helper functions that return tea.Cmd
func Navigate(viewID string, data any) tea.Cmd {
	return func() tea.Msg { return NavigateMsg{ViewID: viewID, Data: data} }
}

func Back() tea.Cmd {
	return func() tea.Msg { return BackMsg{} }
}

func Replace(viewID string, data any) tea.Cmd {
	return func() tea.Msg { return ReplaceMsg{ViewID: viewID, Data: data} }
}
