package stage

import (
	"strings"

	"github.com/develop-suda/VimMaster/internal/buffer"
)

// CheckClear checks if the buffer matches the expected text for the stage.
func CheckClear(buf *buffer.Buffer, expected string) bool {
	actual := strings.TrimRight(buf.Text(), "\n ")
	expected = strings.TrimRight(expected, "\n ")
	return actual == expected
}
