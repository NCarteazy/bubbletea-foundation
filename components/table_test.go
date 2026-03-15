package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/NCarteazy/bubbletea-foundation/theme"
)

func cols() []Column {
	return []Column{
		{Name: "Name", Width: 20, Align: lipgloss.Left},
		{Name: "Size", Width: 10, Align: lipgloss.Right},
	}
}

func TestTableRender(t *testing.T) {
	rows := [][]string{
		{"alpha", "1.5 GB"},
		{"beta", "300 MB"},
	}
	tbl := NewTable(cols(), rows, theme.TokyoNight)
	tbl.SetSize(80, 10)

	out := tbl.View()
	if out == "" {
		t.Fatal("View() returned empty string")
	}
}

func TestTableNavigation(t *testing.T) {
	rows := [][]string{
		{"a", "1"},
		{"b", "2"},
		{"c", "3"},
	}
	tbl := NewTable(cols(), rows, theme.TokyoNight)
	tbl.SetSize(80, 10)

	// Initially cursor should be at 0.
	if tbl.SelectedRow() != 0 {
		t.Fatalf("expected cursor at 0, got %d", tbl.SelectedRow())
	}

	// Press j to move down.
	tbl, _ = tbl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if tbl.SelectedRow() != 1 {
		t.Fatalf("expected cursor at 1 after j, got %d", tbl.SelectedRow())
	}

	// Press k to move up.
	tbl, _ = tbl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if tbl.SelectedRow() != 0 {
		t.Fatalf("expected cursor at 0 after k, got %d", tbl.SelectedRow())
	}

	// Press k again — should not go below 0.
	tbl, _ = tbl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if tbl.SelectedRow() != 0 {
		t.Fatalf("expected cursor at 0 (clamped), got %d", tbl.SelectedRow())
	}
}

func TestTableSort(t *testing.T) {
	rows := [][]string{
		{"alpha", "1.5 GB"},
		{"beta", "300 MB"},
		{"gamma", "2 GB"},
	}
	tbl := NewTable(cols(), rows, theme.TokyoNight)
	tbl.SetSize(80, 10)

	// Press s to sort.
	tbl, _ = tbl.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

	if tbl.RowCount() != 3 {
		t.Fatalf("expected 3 rows after sort, got %d", tbl.RowCount())
	}
}

func TestTableEmpty(t *testing.T) {
	tbl := NewTable(cols(), nil, theme.TokyoNight)
	tbl.SetSize(80, 10)

	out := tbl.View()
	if out == "" {
		t.Fatal("View() returned empty string for nil rows")
	}
}

func TestTableSetRows(t *testing.T) {
	rows := [][]string{
		{"a", "1"},
		{"b", "2"},
	}
	tbl := NewTable(cols(), rows, theme.TokyoNight)
	tbl.SetSize(80, 10)

	if tbl.RowCount() != 2 {
		t.Fatalf("expected 2 rows, got %d", tbl.RowCount())
	}

	newRows := [][]string{
		{"x", "10"},
		{"y", "20"},
		{"z", "30"},
	}
	tbl.SetRows(newRows)

	if tbl.RowCount() != 3 {
		t.Fatalf("expected 3 rows after SetRows, got %d", tbl.RowCount())
	}
}
