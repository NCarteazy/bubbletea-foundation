package layout

import "github.com/charmbracelet/lipgloss"

// TwoPane renders left and right content in a 50/50 horizontal split
// with a 1-character gap between panes.
func TwoPane(left, right string, width, height int) string {
	return TwoPaneWithRatio(left, right, 0.5, width, height)
}

// TwoPaneWithRatio renders left and right content in a horizontal split
// where ratio controls the fraction of width allocated to the left pane.
// A 1-character gap separates the two panes.
func TwoPaneWithRatio(left, right string, ratio float64, width, height int) string {
	const gap = 1
	usable := width - gap
	leftWidth := int(float64(usable) * ratio)
	rightWidth := usable - leftWidth

	leftPane := lipgloss.NewStyle().
		Width(leftWidth).
		Height(height).
		Render(left)

	rightPane := lipgloss.NewStyle().
		Width(rightWidth).
		Height(height).
		Render(right)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftPane,
		lipgloss.NewStyle().Width(gap).Render(" "),
		rightPane,
	)
}
