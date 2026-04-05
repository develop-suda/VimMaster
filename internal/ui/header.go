package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("229")).
	Background(lipgloss.Color("57")).
	Padding(0, 1)

var progressBarFull = lipgloss.NewStyle().
	Foreground(lipgloss.Color("42"))

var progressBarEmpty = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

// RenderHeader renders the header with stage name and progress.
func RenderHeader(stageName string, current, total int, width int) string {
	progress := renderProgress(current, total)
	title := fmt.Sprintf(" %s", stageName)
	progressStr := fmt.Sprintf("進捗: %s (%d/%d) ", progress, current, total)

	titleWidth := lipgloss.Width(title)
	progressWidth := lipgloss.Width(progressStr)
	padding := width - titleWidth - progressWidth
	if padding < 0 {
		padding = 0
	}

	paddingStr := ""
	for i := 0; i < padding; i++ {
		paddingStr += " "
	}

	line := title + paddingStr + progressStr
	return headerStyle.Width(width).Render(line)
}

func renderProgress(current, total int) string {
	if total == 0 {
		return "[-]"
	}
	barWidth := 10
	filled := barWidth * current / total
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += progressBarFull.Render("#")
		} else {
			bar += progressBarEmpty.Render("-")
		}
	}
	return "[" + bar + "]"
}
