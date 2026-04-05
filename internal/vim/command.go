package vim

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/develop-suda/VimMaster/internal/buffer"
)

// PendingOp represents a pending operator (like 'd' waiting for a motion).
type PendingOp int

const (
	OpNone PendingOp = iota
	OpDelete   // d
	OpChange   // c
	OpYank     // y
	OpFind     // f
	OpFindBefore // t
	OpGo       // g (waiting for second g)
)

// String returns the display string for pending operations.
func (p PendingOp) String() string {
	switch p {
	case OpDelete:
		return "d"
	case OpChange:
		return "c"
	case OpYank:
		return "y"
	case OpFind:
		return "f"
	case OpFindBefore:
		return "t"
	case OpGo:
		return "g"
	default:
		return ""
	}
}

// HandleResult represents the result of handling a key.
type HandleResult struct {
	ModeChanged bool
	NewMode     Mode
	Strokes     int
}

// HandleNormalMode processes a key press in normal mode.
func HandleNormalMode(msg tea.KeyMsg, buf *buffer.Buffer, pending *PendingOp) HandleResult {
	result := HandleResult{Strokes: 1}

	// Handle pending operations first
	if *pending != OpNone {
		return handlePendingOp(msg, buf, pending)
	}

	switch msg.String() {
	// Movement
	case "h", "left":
		buf.MoveLeft()
	case "j", "down":
		buf.MoveDown()
	case "k", "up":
		buf.MoveUp()
	case "l", "right":
		buf.MoveRight()
	case "w":
		buf.MoveWordForward()
	case "e":
		buf.MoveWordEnd()
	case "b":
		buf.MoveWordBackward()
	case "0":
		buf.MoveToLineStart()
	case "$":
		buf.MoveToLineEnd()
	case "G":
		buf.MoveToLastLine()

	// Editing
	case "x":
		buf.DeleteChar()
	case "D":
		buf.DeleteToEnd()

	// Mode switching
	case "i":
		result.ModeChanged = true
		result.NewMode = InsertMode
	case "a":
		buf.MoveRight()
		result.ModeChanged = true
		result.NewMode = InsertMode
	case "A":
		buf.MoveToLineEnd()
		lineLen := len([]rune(buf.CurrentLine()))
		if lineLen > 0 {
			buf.CursorCol = lineLen
		}
		result.ModeChanged = true
		result.NewMode = InsertMode
	case "o":
		buf.OpenLineBelow()
		result.ModeChanged = true
		result.NewMode = InsertMode
	case "I":
		buf.MoveToLineStart()
		result.ModeChanged = true
		result.NewMode = InsertMode

	// Pending operators
	case "d":
		*pending = OpDelete
		result.Strokes = 0 // Don't count until completed
	case "c":
		*pending = OpChange
		result.Strokes = 0
	case "y":
		*pending = OpYank
		result.Strokes = 0
	case "f":
		*pending = OpFind
		result.Strokes = 0
	case "t":
		*pending = OpFindBefore
		result.Strokes = 0
	case "g":
		*pending = OpGo
		result.Strokes = 0

	default:
		result.Strokes = 0 // Unknown key, don't count
	}

	return result
}

func handlePendingOp(msg tea.KeyMsg, buf *buffer.Buffer, pending *PendingOp) HandleResult {
	result := HandleResult{Strokes: 1}
	key := msg.String()

	switch *pending {
	case OpDelete:
		switch key {
		case "w":
			buf.DeleteWord()
		case "d":
			buf.DeleteLine()
		case "$":
			buf.DeleteToEnd()
		default:
			result.Strokes = 0
		}
	case OpChange:
		switch key {
		case "w":
			buf.ChangeWord()
			result.ModeChanged = true
			result.NewMode = InsertMode
		default:
			result.Strokes = 0
		}
	case OpFind:
		if len(key) == 1 {
			runes := []rune(key)
			buf.FindChar(runes[0])
		} else {
			result.Strokes = 0
		}
	case OpFindBefore:
		if len(key) == 1 {
			runes := []rune(key)
			buf.FindCharBefore(runes[0])
		} else {
			result.Strokes = 0
		}
	case OpGo:
		switch key {
		case "g":
			buf.MoveToFirstLine()
		default:
			result.Strokes = 0
		}
	default:
		result.Strokes = 0
	}

	*pending = OpNone
	return result
}

// HandleInsertMode processes a key press in insert mode.
func HandleInsertMode(msg tea.KeyMsg, buf *buffer.Buffer) HandleResult {
	result := HandleResult{Strokes: 1}

	switch msg.Type {
	case tea.KeyEsc:
		result.ModeChanged = true
		result.NewMode = NormalMode
		if buf.CursorCol > 0 {
			buf.CursorCol--
		}
		return result
	case tea.KeyBackspace:
		buf.Backspace()
	case tea.KeyEnter:
		buf.InsertNewline()
	case tea.KeyLeft:
		buf.MoveLeft()
	case tea.KeyRight:
		lineLen := len([]rune(buf.CurrentLine()))
		if buf.CursorCol < lineLen {
			buf.CursorCol++
		}
	case tea.KeyUp:
		buf.MoveUp()
		buf.ClampCursorInsert()
	case tea.KeyDown:
		buf.MoveDown()
		buf.ClampCursorInsert()
	default:
		if msg.Type == tea.KeyRunes {
			for _, r := range msg.Runes {
				buf.InsertChar(r)
			}
		}
	}

	return result
}
