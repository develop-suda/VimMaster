package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/develop-suda/VimMaster/internal/buffer"
)

var (
	lineNumberStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(4).
		Align(lipgloss.Right)

	cursorStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("229")).
		Foreground(lipgloss.Color("0"))

	editorBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)

	descriptionStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("117")).
		Bold(true).
		Padding(0, 1)
)

// RenderEditor renders the main editor area with line numbers and cursor.
func RenderEditor(buf *buffer.Buffer, description string, width int) string {
	var sb strings.Builder

	// Description line
	sb.WriteString(descriptionStyle.Render(fmt.Sprintf("課題: %s", description)))
	sb.WriteString("\n\n")

	for i, line := range buf.Lines {
		lineNum := lineNumberStyle.Render(fmt.Sprintf("%d ", i+1))
		sb.WriteString(lineNum)
		sb.WriteString(" ")

		if i == buf.CursorRow {
			sb.WriteString(renderLineWithCursor(line, buf.CursorCol))
		} else {
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}

	innerWidth := width - 4
	if innerWidth < 20 {
		innerWidth = 20
	}

	return editorBorderStyle.Width(innerWidth).Render(sb.String())
}

func renderLineWithCursor(line string, cursorCol int) string {
	runes := []rune(line)
	if len(runes) == 0 {
		return cursorStyle.Render(" ")
	}

	var sb strings.Builder
	for i, r := range runes {
		if i == cursorCol {
			sb.WriteString(cursorStyle.Render(string(r)))
		} else {
			sb.WriteRune(r)
		}
	}

	// If cursor is past end of line, show block cursor at end
	if cursorCol >= len(runes) {
		sb.WriteString(cursorStyle.Render(" "))
	}

	return sb.String()
}
