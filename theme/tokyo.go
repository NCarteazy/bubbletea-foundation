package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
)

// TokyoNight is the default theme using the Tokyo Night color palette.
var TokyoNight Theme = tokyoNight{}

type tokyoNight struct{}

// Colors

func (t tokyoNight) Bg() string        { return "#1a1b26" }
func (t tokyoNight) Primary() string    { return "#7dcfff" }
func (t tokyoNight) Secondary() string  { return "#bb9af7" }
func (t tokyoNight) Success() string    { return "#9ece6a" }
func (t tokyoNight) Warning() string    { return "#e0af68" }
func (t tokyoNight) Danger() string     { return "#f7768e" }
func (t tokyoNight) Text() string       { return "#c0caf5" }
func (t tokyoNight) Muted() string      { return "#565f89" }
func (t tokyoNight) Surface() string    { return "#24283b" }
func (t tokyoNight) Overlay() string    { return "#414868" }

// Pre-built styles

func (t tokyoNight) Title() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Primary())).
		Bold(true)
}

func (t tokyoNight) Panel() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.Surface())).
		Foreground(lipgloss.Color(t.Text())).
		Padding(0, 1)
}

func (t tokyoNight) ActivePanel() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.Surface())).
		Foreground(lipgloss.Color(t.Text())).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.Primary())).
		Padding(0, 1)
}

func (t tokyoNight) StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.Overlay())).
		Foreground(lipgloss.Color(t.Text())).
		Padding(0, 1)
}

func (t tokyoNight) Breadcrumb() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted()))
}

func (t tokyoNight) TableHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Primary())).
		Bold(true)
}

func (t tokyoNight) SelectedRow() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(t.Overlay())).
		Foreground(lipgloss.Color(t.Text()))
}

// Utilities

func (t tokyoNight) FormatSize(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	return humanize.Bytes(uint64(bytes))
}

func (t tokyoNight) SizeBar(ratio float64, maxWidth int) string {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	filled := int(ratio * float64(maxWidth))
	empty := maxWidth - filled

	var color string
	switch {
	case ratio >= 0.7:
		color = t.Danger()
	case ratio >= 0.4:
		color = t.Warning()
	default:
		color = t.Success()
	}

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted()))

	return filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))
}
