package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	foundation "github.com/NCarteazy/bubbletea-foundation"
	"github.com/NCarteazy/bubbletea-foundation/components"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

// --- homeView ---

type homeView struct{}

func (v homeView) ID() string { return "Home" }

func (v homeView) Update(msg tea.Msg, ctx foundation.ViewContext) (foundation.View, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter":
			return v, foundation.Navigate("Files", nil)
		case "q":
			return v, tea.Quit
		}
	}
	return v, nil
}

func (v homeView) Render(ctx foundation.ViewContext) string {
	title := ctx.Theme.Title().Render("Welcome to Example App")

	bar := ctx.Theme.SizeBar(0.42, 40)
	label := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ctx.Theme.Muted())).
		Render("Disk usage: 42%")

	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ctx.Theme.Muted())).
		Render("\nPress enter to browse files, q to quit.")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		bar,
		label,
		hint,
	)
}

// --- filesView ---

type filesView struct {
	table components.Table
	built bool
}

func (v filesView) ID() string { return "Files" }

func (v filesView) Update(msg tea.Msg, ctx foundation.ViewContext) (foundation.View, tea.Cmd) {
	if !v.built {
		v.table = buildFileTable(ctx.Theme)
		v.table.SetSize(ctx.Width, ctx.Height-2)
		v.built = true
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "enter":
			idx := v.table.SelectedRow()
			return v, foundation.Navigate("Detail", idx)
		case "d":
			c := components.NewConfirm("Delete this file?", ctx.Theme)
			c.OnConfirm = func() tea.Cmd {
				return foundation.Flash("File deleted!", 30)
			}
			return v, foundation.ShowOverlay(foundation.ConfirmOverlay(c))
		case "esc":
			return v, foundation.Back()
		}
	}

	updated, cmd := v.table.Update(msg)
	v.table = updated
	return v, cmd
}

func (v filesView) Render(ctx foundation.ViewContext) string {
	if !v.built {
		return ""
	}
	return v.table.View()
}

func buildFileTable(th theme.Theme) components.Table {
	columns := []components.Column{
		{Name: "Name", Width: 30, Align: lipgloss.Left},
		{Name: "Size", Width: 12, Align: lipgloss.Right},
		{Name: "Modified", Width: 20, Align: lipgloss.Left},
	}

	rows := [][]string{
		{"Documents", "2.4 GB", "2026-03-10"},
		{"Downloads", "8.1 GB", "2026-03-14"},
		{"node_modules", "1.2 GB", "2026-02-28"},
		{"Library/Caches", "4.7 GB", "2026-03-15"},
		{".docker", "12.3 GB", "2026-03-01"},
		{"Applications", "15.6 GB", "2026-01-20"},
		{"Pictures", "3.2 GB", "2026-03-12"},
		{"Music", "900 MB", "2025-12-05"},
	}

	return components.NewTable(columns, rows, th)
}

// --- detailView ---

type detailView struct{}

func (v detailView) ID() string { return "Detail" }

func (v detailView) Update(msg tea.Msg, ctx foundation.ViewContext) (foundation.View, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "esc" {
			return v, foundation.Back()
		}
	}
	return v, nil
}

func (v detailView) Render(ctx foundation.ViewContext) string {
	idx := 0
	if n, ok := ctx.Data.(int); ok {
		idx = n
	}
	title := ctx.Theme.Title().Render(fmt.Sprintf("Detail view for item %d", idx))
	hint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ctx.Theme.Muted())).
		Render("\nPress esc to go back.")
	return lipgloss.JoinVertical(lipgloss.Left, title, hint)
}

// --- main ---

func main() {
	app := foundation.New(
		foundation.WithTitle("Example App"),
		foundation.WithTheme(theme.TokyoNight),
		foundation.WithInitialView("Home"),
		foundation.WithViews(homeView{}, filesView{}, detailView{}),
		foundation.WithHelp([]foundation.KeySection{
			{Title: "Navigation", Keys: []foundation.KeyHint{
				{Key: "j/k", Desc: "move up/down"},
				{Key: "enter", Desc: "select / drill in"},
				{Key: "esc", Desc: "go back"},
			}},
			{Title: "Actions", Keys: []foundation.KeyHint{
				{Key: "d", Desc: "delete selected"},
				{Key: "s", Desc: "cycle sort"},
			}},
		}),
		foundation.WithStatusHints([]foundation.KeyHint{
			{Key: "j/k", Desc: "navigate"},
			{Key: "enter", Desc: "select"},
			{Key: "esc", Desc: "back"},
			{Key: "?", Desc: "help"},
			{Key: "q", Desc: "quit"},
		}),
	)
	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
