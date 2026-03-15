package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color palette, pre-built styles, and utility methods
// used throughout the TUI.
type Theme interface {
	// Colors
	Bg() string
	Primary() string
	Secondary() string
	Success() string
	Warning() string
	Danger() string
	Text() string
	Muted() string
	Surface() string

	// Pre-built styles
	Title() lipgloss.Style
	Panel() lipgloss.Style
	ActivePanel() lipgloss.Style
	StatusBar() lipgloss.Style
	Breadcrumb() lipgloss.Style
	TableHeader() lipgloss.Style
	SelectedRow() lipgloss.Style

	// Utilities
	FormatSize(bytes int64) string
	SizeBar(ratio float64, maxWidth int) string
}
