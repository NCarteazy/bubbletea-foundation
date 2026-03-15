package layout

import "github.com/charmbracelet/lipgloss"

// OverlayCenter places overlay content centered within the given dimensions.
// If width or height is <= 0, the overlay string is returned as-is.
func OverlayCenter(width, height int, background, overlay string) string {
	if width <= 0 || height <= 0 {
		return overlay
	}
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		overlay,
		lipgloss.WithWhitespaceChars(" "),
	)
}
