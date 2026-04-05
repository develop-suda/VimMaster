package stage

// Stage represents a learning stage/challenge.
type Stage struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	InitialText string   `json:"initial_text"`
	ExpectedText string  `json:"expected_text"`
	Hint        string   `json:"hint"`
	MinStrokes  int      `json:"min_strokes"`
	Category    string   `json:"category"`
	AllowedKeys []string `json:"allowed_keys,omitempty"`
}

// StageFile represents the JSON structure of a stage file.
type StageFile struct {
	Stages []Stage `json:"stages"`
}

// Rating represents the performance rating.
type Rating int

const (
	RatingC Rating = iota // Cleared
	RatingB               // Good
	RatingA               // Great
	RatingS               // Perfect (minimum strokes)
)

// String returns the display string for the rating.
func (r Rating) String() string {
	switch r {
	case RatingS:
		return "S"
	case RatingA:
		return "A"
	case RatingB:
		return "B"
	default:
		return "C"
	}
}

// CalculateRating calculates a rating based on strokes used vs minimum.
func CalculateRating(strokes, minStrokes int) Rating {
	if minStrokes <= 0 {
		return RatingC
	}
	ratio := float64(strokes) / float64(minStrokes)
	switch {
	case ratio <= 1.0:
		return RatingS
	case ratio <= 1.5:
		return RatingA
	case ratio <= 2.0:
		return RatingB
	default:
		return RatingC
	}
}
