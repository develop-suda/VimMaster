package vim

// Mode represents the current Vim mode.
type Mode int

const (
	NormalMode Mode = iota
	InsertMode
	VisualMode
)

// String returns the display string for the mode.
func (m Mode) String() string {
	switch m {
	case NormalMode:
		return "NORMAL"
	case InsertMode:
		return "INSERT"
	case VisualMode:
		return "VISUAL"
	default:
		return "UNKNOWN"
	}
}
