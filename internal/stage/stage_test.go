package stage

import (
	"testing"

	"github.com/develop-suda/VimMaster/internal/buffer"
)

func TestCheckClear(t *testing.T) {
	buf := buffer.NewBuffer("hello world")
	if !CheckClear(buf, "hello world") {
		t.Error("expected clear")
	}
	if CheckClear(buf, "hello") {
		t.Error("expected not clear")
	}
}

func TestCheckClearTrimming(t *testing.T) {
	buf := buffer.NewBuffer("hello world\n")
	if !CheckClear(buf, "hello world") {
		t.Error("expected clear with trimmed newline")
	}
}

func TestCalculateRating(t *testing.T) {
	tests := []struct {
		strokes    int
		minStrokes int
		expected   Rating
	}{
		{5, 5, RatingS},
		{3, 5, RatingS},
		{7, 5, RatingA},
		{9, 5, RatingB},
		{15, 5, RatingC},
	}

	for _, tt := range tests {
		rating := CalculateRating(tt.strokes, tt.minStrokes)
		if rating != tt.expected {
			t.Errorf("strokes=%d, minStrokes=%d: expected %s, got %s",
				tt.strokes, tt.minStrokes, tt.expected.String(), rating.String())
		}
	}
}

func TestLoadAllStages(t *testing.T) {
	stages, err := LoadAllStages()
	if err != nil {
		t.Fatalf("failed to load stages: %v", err)
	}
	if len(stages) == 0 {
		t.Error("expected at least one stage")
	}
	// Check stages are sorted by ID
	for i := 1; i < len(stages); i++ {
		if stages[i].ID <= stages[i-1].ID {
			t.Errorf("stages not sorted: %d <= %d", stages[i].ID, stages[i-1].ID)
		}
	}
}
