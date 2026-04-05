package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/develop-suda/VimMaster/internal/vim"
)

var (
	statusLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("252")).
		Padding(0, 1)

	modeNormalStyle = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("63")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)

	modeInsertStyle = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("34")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)

	modeVisualStyle = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("208")).
		Foreground(lipgloss.Color("255")).
		Padding(0, 1)
)

// RenderStatusLine renders the status line showing mode, position, and pending keys.
func RenderStatusLine(mode vim.Mode, row, col int, pendingOp vim.PendingOp, strokes int, width int) string {
	// Mode badge
	var modeBadge string
	switch mode {
	case vim.InsertMode:
		modeBadge = modeInsertStyle.Render(fmt.Sprintf(" %s ", mode.String()))
	case vim.VisualMode:
		modeBadge = modeVisualStyle.Render(fmt.Sprintf(" %s ", mode.String()))
	default:
		modeBadge = modeNormalStyle.Render(fmt.Sprintf(" %s ", mode.String()))
	}

	// Position
	pos := fmt.Sprintf(" %d:%d", row+1, col+1)

	// Pending key display
	pendingStr := ""
	if pendingOp != vim.OpNone {
		pendingStr = fmt.Sprintf(" | キー入力: %s_", pendingOp.String())
	}

	// Strokes count
	strokesStr := fmt.Sprintf(" | 手数: %d", strokes)

	rightPart := pos + pendingStr + strokesStr

	content := modeBadge + statusLineStyle.Render(rightPart)

	lineWidth := lipgloss.Width(content)
	if lineWidth < width {
		padding := width - lineWidth
		paddingStr := ""
		for i := 0; i < padding; i++ {
			paddingStr += " "
		}
		content = content + statusLineStyle.Render(paddingStr)
	}

	return content
}
