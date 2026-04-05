package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var hintStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("243")).
	Italic(true).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

// RenderHint renders the hint area.
func RenderHint(hint string, showHint bool, width int) string {
	if !showHint || hint == "" {
		return hintStyle.Width(width - 4).Render("ヒント: '?' を押すとヒントが表示されます")
	}
	return hintStyle.Width(width - 4).Render("ヒント: " + hint)
}
