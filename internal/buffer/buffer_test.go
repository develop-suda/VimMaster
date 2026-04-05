package buffer

import (
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer("hello\nworld")
	if len(buf.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(buf.Lines))
	}
	if buf.Lines[0] != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.Lines[0])
	}
	if buf.Lines[1] != "world" {
		t.Errorf("expected 'world', got '%s'", buf.Lines[1])
	}
}

func TestBufferMovement(t *testing.T) {
	buf := NewBuffer("abc\ndef\nghi")

	buf.MoveRight()
	if buf.CursorCol != 1 {
		t.Errorf("expected col 1, got %d", buf.CursorCol)
	}

	buf.MoveDown()
	if buf.CursorRow != 1 {
		t.Errorf("expected row 1, got %d", buf.CursorRow)
	}

	buf.MoveLeft()
	if buf.CursorCol != 0 {
		t.Errorf("expected col 0, got %d", buf.CursorCol)
	}

	buf.MoveUp()
	if buf.CursorRow != 0 {
		t.Errorf("expected row 0, got %d", buf.CursorRow)
	}
}

func TestDeleteChar(t *testing.T) {
	buf := NewBuffer("abc")
	buf.CursorCol = 1
	buf.DeleteChar()
	if buf.CurrentLine() != "ac" {
		t.Errorf("expected 'ac', got '%s'", buf.CurrentLine())
	}
}

func TestDeleteWord(t *testing.T) {
	buf := NewBuffer("hello world")
	buf.DeleteWord()
	if buf.CurrentLine() != "world" {
		t.Errorf("expected 'world', got '%s'", buf.CurrentLine())
	}
}

func TestDeleteLine(t *testing.T) {
	buf := NewBuffer("line1\nline2\nline3")
	buf.CursorRow = 1
	buf.DeleteLine()
	if len(buf.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(buf.Lines))
	}
	if buf.Lines[1] != "line3" {
		t.Errorf("expected 'line3', got '%s'", buf.Lines[1])
	}
}

func TestDeleteToEnd(t *testing.T) {
	buf := NewBuffer("hello world")
	buf.CursorCol = 5
	buf.DeleteToEnd()
	if buf.CurrentLine() != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.CurrentLine())
	}
}

func TestInsertChar(t *testing.T) {
	buf := NewBuffer("hllo")
	buf.CursorCol = 1
	buf.InsertChar('e')
	if buf.CurrentLine() != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.CurrentLine())
	}
}

func TestInsertNewline(t *testing.T) {
	buf := NewBuffer("helloworld")
	buf.CursorCol = 5
	buf.InsertNewline()
	if len(buf.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(buf.Lines))
	}
	if buf.Lines[0] != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.Lines[0])
	}
	if buf.Lines[1] != "world" {
		t.Errorf("expected 'world', got '%s'", buf.Lines[1])
	}
}

func TestBackspace(t *testing.T) {
	buf := NewBuffer("hello")
	buf.CursorCol = 3
	buf.Backspace()
	if buf.CurrentLine() != "helo" {
		t.Errorf("expected 'helo', got '%s'", buf.CurrentLine())
	}
}

func TestOpenLineBelow(t *testing.T) {
	buf := NewBuffer("line1\nline2")
	buf.OpenLineBelow()
	if len(buf.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(buf.Lines))
	}
	if buf.CursorRow != 1 {
		t.Errorf("expected row 1, got %d", buf.CursorRow)
	}
	if buf.Lines[1] != "" {
		t.Errorf("expected empty line, got '%s'", buf.Lines[1])
	}
}

func TestMoveWordForward(t *testing.T) {
	buf := NewBuffer("hello world vim")
	buf.MoveWordForward()
	if buf.CursorCol != 6 {
		t.Errorf("expected col 6, got %d", buf.CursorCol)
	}
}

func TestMoveWordBackward(t *testing.T) {
	buf := NewBuffer("hello world")
	buf.CursorCol = 6
	buf.MoveWordBackward()
	if buf.CursorCol != 0 {
		t.Errorf("expected col 0, got %d", buf.CursorCol)
	}
}

func TestFindChar(t *testing.T) {
	buf := NewBuffer("hello world")
	found := buf.FindChar('w')
	if !found {
		t.Error("expected to find 'w'")
	}
	if buf.CursorCol != 6 {
		t.Errorf("expected col 6, got %d", buf.CursorCol)
	}
}

func TestMoveToLineStartEnd(t *testing.T) {
	buf := NewBuffer("hello")
	buf.CursorCol = 3
	buf.MoveToLineStart()
	if buf.CursorCol != 0 {
		t.Errorf("expected col 0, got %d", buf.CursorCol)
	}
	buf.MoveToLineEnd()
	if buf.CursorCol != 4 {
		t.Errorf("expected col 4, got %d", buf.CursorCol)
	}
}

func TestMoveFirstLastLine(t *testing.T) {
	buf := NewBuffer("a\nb\nc")
	buf.MoveToLastLine()
	if buf.CursorRow != 2 {
		t.Errorf("expected row 2, got %d", buf.CursorRow)
	}
	buf.MoveToFirstLine()
	if buf.CursorRow != 0 {
		t.Errorf("expected row 0, got %d", buf.CursorRow)
	}
}

func TestText(t *testing.T) {
	text := "hello\nworld"
	buf := NewBuffer(text)
	if buf.Text() != text {
		t.Errorf("expected '%s', got '%s'", text, buf.Text())
	}
}
