package puzzle

import (
	"testing"
)

func TestIsPuzzle(t *testing.T) {
	tests := []struct {
		name       string
		bestEval   int
		playerEval int
		threshold  int
		want       bool
	}{
		{
			name:       "above threshold",
			bestEval:   350,
			playerEval: 100,
			threshold:  200,
			want:       true,
		},
		{
			name:       "exactly threshold",
			bestEval:   300,
			playerEval: 100,
			threshold:  200,
			want:       true,
		},
		{
			name:       "below threshold",
			bestEval:   299,
			playerEval: 100,
			threshold:  200,
			want:       false,
		},
		{
			name:       "player move is better",
			bestEval:   100,
			playerEval: 150,
			threshold:  200,
			want:       false,
		},
		{
			name:       "works with negative evaluations",
			bestEval:   -50,
			playerEval: -300,
			threshold:  200,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPuzzle(tt.bestEval, tt.playerEval, tt.threshold); got != tt.want {
				t.Fatalf("IsPuzzle(%d, %d, %d) = %t, want %t", tt.bestEval, tt.playerEval, tt.threshold, got, tt.want)
			}
		})
	}
}

func TestSolveLinesReturnsNilUntilImplemented(t *testing.T) {
	if got := SolveLines(nil, 3, 50); got != nil {
		t.Fatalf("SolveLines() = %#v, want nil", got)
	}
}
