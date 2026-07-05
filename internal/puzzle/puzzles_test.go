package puzzle

import (
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
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

func TestTryGeneratePuzzle(t *testing.T) {
	game := chess.NewGame()
	pos := game.Position()
	bestMove := mustMove(t, game, "e2e4")
	playerMove := mustMove(t, game, "d2d4")

	tests := []struct {
		name        string
		bestScore   int
		playerScore int
		threshold   int
		wantNil     bool
	}{
		{
			name:        "below threshold returns nil",
			bestScore:   299,
			playerScore: 100,
			threshold:   200,
			wantNil:     true,
		},
		{
			name:        "exactly threshold creates puzzle",
			bestScore:   300,
			playerScore: 100,
			threshold:   200,
			wantNil:     false,
		},
		{
			name:        "above threshold creates puzzle",
			bestScore:   350,
			playerScore: 100,
			threshold:   200,
			wantNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TryGeneratePuzzle(pos, tt.bestScore, bestMove, tt.playerScore, playerMove, tt.threshold)
			if (got == nil) != tt.wantNil {
				t.Fatalf("TryGeneratePuzzle() nil = %t, want %t", got == nil, tt.wantNil)
			}
			if tt.wantNil {
				return
			}
			if got.Position != pos {
				t.Fatalf("Position = %p, want %p", got.Position, pos)
			}
			if got.BestMove != bestMove {
				t.Fatalf("BestMove = %p, want %p", got.BestMove, bestMove)
			}
			if got.PlayerMove != playerMove {
				t.Fatalf("PlayerMove = %p, want %p", got.PlayerMove, playerMove)
			}
			if got.Lines != nil {
				t.Fatalf("Lines = %#v, want nil", got.Lines)
			}
			if got.AlternateMoves != nil {
				t.Fatalf("AlternateMoves = %#v, want nil", got.AlternateMoves)
			}
		})
	}
}

func TestAlternateMovesReturnsNilUntilImplemented(t *testing.T) {
	if got := AlternameMoves(chess.NewGame().Position()); got != nil {
		t.Fatalf("AlternameMoves() = %#v, want nil", got)
	}
}

func TestPuzzleStringWithNilPositionAndMoves(t *testing.T) {
	got := (Puzzle{}).String()
	wantSubstrings := []string{
		"Puzzle\n",
		"Position: <nil>\n",
		"Player move: <nil>\n",
		"Best move: <nil>\n",
		"Solution lines: none\n",
		"Alternate moves: none\n\n",
	}

	for _, want := range wantSubstrings {
		if !strings.Contains(got, want) {
			t.Fatalf("Puzzle.String() missing %q in:\n%s", want, got)
		}
	}
}

func TestPuzzleStringWithPositionLinesAndAlternates(t *testing.T) {
	game := chess.NewGame()
	pos := game.Position()
	bestMove := mustMove(t, game, "e2e4")
	playerMove := mustMove(t, game, "d2d4")
	alternateMove := mustMove(t, game, "g1f3")

	got := Puzzle{
		Position:   pos,
		PlayerMove: playerMove,
		BestMove:   bestMove,
		Lines: [][]chess.Move{
			{*bestMove, *alternateMove},
			{},
		},
		AlternateMoves: []AlternateMove{
			{
				Move:       alternateMove,
				Evaluation: 125,
				Difference: 75,
			},
		},
	}.String()

	wantSubstrings := []string{
		"Board:",
		"FEN: ",
		"Side to move: w\n",
		"Player move: d2d4\n",
		"Best move: e2e4\n",
		"Solution lines:\n",
		"  1. e2e4 g1f3\n",
		"  2. <empty>\n",
		"Alternate moves:\n",
		"  1. g1f3 (eval: 125 cp, difference: 75 cp)\n\n",
	}

	for _, want := range wantSubstrings {
		if !strings.Contains(got, want) {
			t.Fatalf("Puzzle.String() missing %q in:\n%s", want, got)
		}
	}
}

func TestFormatMove(t *testing.T) {
	game := chess.NewGame()
	move := mustMove(t, game, "e2e4")

	tests := []struct {
		name string
		move *chess.Move
		want string
	}{
		{
			name: "nil move",
			move: nil,
			want: "<nil>",
		},
		{
			name: "move string",
			move: move,
			want: "e2e4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatMove(tt.move); got != tt.want {
				t.Fatalf("formatMove() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatMoveLine(t *testing.T) {
	game := chess.NewGame()
	e4 := mustMove(t, game, "e2e4")
	nf3 := mustMove(t, game, "g1f3")

	tests := []struct {
		name string
		line []chess.Move
		want string
	}{
		{
			name: "empty line",
			line: nil,
			want: "<empty>",
		},
		{
			name: "single move",
			line: []chess.Move{*e4},
			want: "e2e4",
		},
		{
			name: "multiple moves",
			line: []chess.Move{*e4, *nf3},
			want: "e2e4 g1f3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatMoveLine(tt.line); got != tt.want {
				t.Fatalf("formatMoveLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func mustMove(t *testing.T, game *chess.Game, want string) *chess.Move {
	t.Helper()

	moves := game.ValidMoves()
	for i := range moves {
		if moves[i].String() == want {
			return &moves[i]
		}
	}

	t.Fatalf("could not find move %q", want)
	return nil
}
