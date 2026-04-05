package buffer

import "strings"

// Buffer represents the virtual text buffer for the editor.
type Buffer struct {
	Lines    []string
	CursorRow int
	CursorCol int
}

// NewBuffer creates a new Buffer from the given text.
func NewBuffer(text string) *Buffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &Buffer{
		Lines:    lines,
		CursorRow: 0,
		CursorCol: 0,
	}
}

// Text returns the full text content of the buffer.
func (b *Buffer) Text() string {
	return strings.Join(b.Lines, "\n")
}

// CurrentLine returns the current line.
func (b *Buffer) CurrentLine() string {
	if b.CursorRow >= 0 && b.CursorRow < len(b.Lines) {
		return b.Lines[b.CursorRow]
	}
	return ""
}

// SetCurrentLine sets the current line text.
func (b *Buffer) SetCurrentLine(s string) {
	if b.CursorRow >= 0 && b.CursorRow < len(b.Lines) {
		b.Lines[b.CursorRow] = s
	}
}

// RuneAtCursor returns the rune at the current cursor position.
func (b *Buffer) RuneAtCursor() (rune, bool) {
	line := b.CurrentLine()
	runes := []rune(line)
	if b.CursorCol >= 0 && b.CursorCol < len(runes) {
		return runes[b.CursorCol], true
	}
	return 0, false
}

// ClampCursor ensures the cursor is within valid bounds.
func (b *Buffer) ClampCursor() {
	if b.CursorRow < 0 {
		b.CursorRow = 0
	}
	if b.CursorRow >= len(b.Lines) {
		b.CursorRow = len(b.Lines) - 1
	}
	if b.CursorRow < 0 {
		b.CursorRow = 0
	}
	lineLen := len([]rune(b.CurrentLine()))
	if lineLen == 0 {
		b.CursorCol = 0
		return
	}
	if b.CursorCol < 0 {
		b.CursorCol = 0
	}
	if b.CursorCol >= lineLen {
		b.CursorCol = lineLen - 1
	}
}

// ClampCursorInsert is like ClampCursor but allows cursor at end of line (for insert mode).
func (b *Buffer) ClampCursorInsert() {
	if b.CursorRow < 0 {
		b.CursorRow = 0
	}
	if b.CursorRow >= len(b.Lines) {
		b.CursorRow = len(b.Lines) - 1
	}
	if b.CursorRow < 0 {
		b.CursorRow = 0
	}
	lineLen := len([]rune(b.CurrentLine()))
	if b.CursorCol < 0 {
		b.CursorCol = 0
	}
	if b.CursorCol > lineLen {
		b.CursorCol = lineLen
	}
}

// MoveLeft moves the cursor left.
func (b *Buffer) MoveLeft() {
	if b.CursorCol > 0 {
		b.CursorCol--
	}
}

// MoveRight moves the cursor right.
func (b *Buffer) MoveRight() {
	lineLen := len([]rune(b.CurrentLine()))
	if b.CursorCol < lineLen-1 {
		b.CursorCol++
	}
}

// MoveUp moves the cursor up.
func (b *Buffer) MoveUp() {
	if b.CursorRow > 0 {
		b.CursorRow--
		b.ClampCursor()
	}
}

// MoveDown moves the cursor down.
func (b *Buffer) MoveDown() {
	if b.CursorRow < len(b.Lines)-1 {
		b.CursorRow++
		b.ClampCursor()
	}
}

// MoveToLineStart moves cursor to the beginning of the line.
func (b *Buffer) MoveToLineStart() {
	b.CursorCol = 0
}

// MoveToLineEnd moves cursor to the end of the line.
func (b *Buffer) MoveToLineEnd() {
	lineLen := len([]rune(b.CurrentLine()))
	if lineLen > 0 {
		b.CursorCol = lineLen - 1
	} else {
		b.CursorCol = 0
	}
}

// MoveToFirstLine moves cursor to the first line.
func (b *Buffer) MoveToFirstLine() {
	b.CursorRow = 0
	b.ClampCursor()
}

// MoveToLastLine moves cursor to the last line.
func (b *Buffer) MoveToLastLine() {
	b.CursorRow = len(b.Lines) - 1
	b.ClampCursor()
}

// MoveWordForward moves the cursor to the start of the next word.
func (b *Buffer) MoveWordForward() {
	line := []rune(b.CurrentLine())
	col := b.CursorCol

	if col >= len(line) {
		// Move to next line
		if b.CursorRow < len(b.Lines)-1 {
			b.CursorRow++
			b.CursorCol = 0
		}
		return
	}

	// Skip current word
	for col < len(line) && !isSpace(line[col]) {
		col++
	}
	// Skip whitespace
	for col < len(line) && isSpace(line[col]) {
		col++
	}

	if col >= len(line) {
		// Move to next line
		if b.CursorRow < len(b.Lines)-1 {
			b.CursorRow++
			b.CursorCol = 0
		} else {
			b.CursorCol = len(line) - 1
			if b.CursorCol < 0 {
				b.CursorCol = 0
			}
		}
	} else {
		b.CursorCol = col
	}
}

// MoveWordEnd moves the cursor to the end of the current/next word.
func (b *Buffer) MoveWordEnd() {
	line := []rune(b.CurrentLine())
	col := b.CursorCol

	if col >= len(line)-1 {
		if b.CursorRow < len(b.Lines)-1 {
			b.CursorRow++
			line = []rune(b.CurrentLine())
			col = 0
			// Skip whitespace
			for col < len(line) && isSpace(line[col]) {
				col++
			}
			// Move to end of word
			for col < len(line)-1 && !isSpace(line[col+1]) {
				col++
			}
			b.CursorCol = col
		}
		return
	}

	col++
	// Skip whitespace
	for col < len(line) && isSpace(line[col]) {
		col++
	}
	// Move to end of word
	for col < len(line)-1 && !isSpace(line[col+1]) {
		col++
	}
	b.CursorCol = col
}

// MoveWordBackward moves the cursor to the start of the previous word.
func (b *Buffer) MoveWordBackward() {
	line := []rune(b.CurrentLine())
	col := b.CursorCol

	if col <= 0 {
		if b.CursorRow > 0 {
			b.CursorRow--
			b.MoveToLineEnd()
		}
		return
	}

	col--
	// Skip whitespace
	for col > 0 && isSpace(line[col]) {
		col--
	}
	// Move to start of word
	for col > 0 && !isSpace(line[col-1]) {
		col--
	}
	b.CursorCol = col
}

// DeleteChar deletes the character at cursor (like 'x').
func (b *Buffer) DeleteChar() {
	line := []rune(b.CurrentLine())
	if b.CursorCol >= 0 && b.CursorCol < len(line) {
		newLine := append(line[:b.CursorCol], line[b.CursorCol+1:]...)
		b.SetCurrentLine(string(newLine))
		if b.CursorCol >= len(newLine) && len(newLine) > 0 {
			b.CursorCol = len(newLine) - 1
		}
	}
}

// DeleteWord deletes from cursor to the start of the next word (like 'dw').
func (b *Buffer) DeleteWord() {
	line := []rune(b.CurrentLine())
	col := b.CursorCol

	if col >= len(line) {
		return
	}

	endCol := col
	// Delete current word chars
	for endCol < len(line) && !isSpace(line[endCol]) {
		endCol++
	}
	// Delete trailing whitespace
	for endCol < len(line) && isSpace(line[endCol]) {
		endCol++
	}

	newLine := append(line[:col], line[endCol:]...)
	b.SetCurrentLine(string(newLine))
	if b.CursorCol >= len(newLine) && len(newLine) > 0 {
		b.CursorCol = len(newLine) - 1
	}
}

// DeleteToEnd deletes from cursor to end of line (like 'd$' or 'D').
func (b *Buffer) DeleteToEnd() {
	line := []rune(b.CurrentLine())
	if b.CursorCol < len(line) {
		newLine := line[:b.CursorCol]
		b.SetCurrentLine(string(newLine))
		if b.CursorCol > 0 {
			b.CursorCol--
		}
	}
}

// DeleteLine deletes the entire current line (like 'dd').
func (b *Buffer) DeleteLine() {
	if len(b.Lines) <= 1 {
		b.Lines = []string{""}
		b.CursorCol = 0
		return
	}
	b.Lines = append(b.Lines[:b.CursorRow], b.Lines[b.CursorRow+1:]...)
	if b.CursorRow >= len(b.Lines) {
		b.CursorRow = len(b.Lines) - 1
	}
	b.ClampCursor()
}

// InsertChar inserts a character at the cursor position (insert mode).
func (b *Buffer) InsertChar(ch rune) {
	line := []rune(b.CurrentLine())
	col := b.CursorCol
	if col > len(line) {
		col = len(line)
	}
	newLine := make([]rune, 0, len(line)+1)
	newLine = append(newLine, line[:col]...)
	newLine = append(newLine, ch)
	newLine = append(newLine, line[col:]...)
	b.SetCurrentLine(string(newLine))
	b.CursorCol = col + 1
}

// InsertNewline inserts a newline at cursor, splitting the line.
func (b *Buffer) InsertNewline() {
	line := []rune(b.CurrentLine())
	col := b.CursorCol
	if col > len(line) {
		col = len(line)
	}

	before := string(line[:col])
	after := string(line[col:])

	b.Lines[b.CursorRow] = before
	newLines := make([]string, 0, len(b.Lines)+1)
	newLines = append(newLines, b.Lines[:b.CursorRow+1]...)
	newLines = append(newLines, after)
	newLines = append(newLines, b.Lines[b.CursorRow+1:]...)
	b.Lines = newLines
	b.CursorRow++
	b.CursorCol = 0
}

// OpenLineBelow opens a new line below cursor (like 'o').
func (b *Buffer) OpenLineBelow() {
	newLines := make([]string, 0, len(b.Lines)+1)
	newLines = append(newLines, b.Lines[:b.CursorRow+1]...)
	newLines = append(newLines, "")
	newLines = append(newLines, b.Lines[b.CursorRow+1:]...)
	b.Lines = newLines
	b.CursorRow++
	b.CursorCol = 0
}

// Backspace deletes the character before cursor (insert mode).
func (b *Buffer) Backspace() {
	if b.CursorCol > 0 {
		line := []rune(b.CurrentLine())
		newLine := append(line[:b.CursorCol-1], line[b.CursorCol:]...)
		b.SetCurrentLine(string(newLine))
		b.CursorCol--
	} else if b.CursorRow > 0 {
		// Join with previous line
		prevLine := b.Lines[b.CursorRow-1]
		curLine := b.CurrentLine()
		b.Lines[b.CursorRow-1] = prevLine + curLine
		b.Lines = append(b.Lines[:b.CursorRow], b.Lines[b.CursorRow+1:]...)
		b.CursorRow--
		b.CursorCol = len([]rune(prevLine))
	}
}

// FindChar finds the next occurrence of a character on the current line (like 'f').
func (b *Buffer) FindChar(ch rune) bool {
	line := []rune(b.CurrentLine())
	for i := b.CursorCol + 1; i < len(line); i++ {
		if line[i] == ch {
			b.CursorCol = i
			return true
		}
	}
	return false
}

// FindCharBefore finds the next occurrence of a character and moves one before it (like 't').
func (b *Buffer) FindCharBefore(ch rune) bool {
	line := []rune(b.CurrentLine())
	for i := b.CursorCol + 1; i < len(line); i++ {
		if line[i] == ch {
			b.CursorCol = i - 1
			return true
		}
	}
	return false
}

// ChangeWord deletes from cursor to next word start (like 'cw') - same as dw but used with mode switch.
func (b *Buffer) ChangeWord() {
	b.DeleteWord()
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}
