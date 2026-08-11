package opening

import (
	"strings"
	"testing"

	"github.com/corentings/chess/v2"
)

func gameMoves(t *testing.T, pgn string) []*chess.Move {
	t.Helper()
	game, err := chess.PGN(strings.NewReader(pgn))
	if err != nil {
		t.Fatalf("PGN parsing failed: %v", err)
	}
	return chess.NewGame(game).Moves()
}

func TestBookPliesEmpty(t *testing.T) {
	if got := BookPlies(nil); got != 0 {
		t.Errorf("BookPlies(nil) = %d, want 0", got)
	}
	if got := BookPlies([]*chess.Move{}); got != 0 {
		t.Errorf("BookPlies([]) = %d, want 0", got)
	}
}

// Berlin Defense main line is a well-known ECO entry (C65); a game that
// stops exactly there should be recognized as book for every ply played.
func TestBookPliesKnownLine(t *testing.T) {
	moves := gameMoves(t, "1. e4 e5 2. Nf3 Nc6 3. Bb5 Nf6 4. O-O")
	if got, want := BookPlies(moves), len(moves); got != want {
		t.Errorf("BookPlies(Berlin Defense) = %d, want %d (len(moves))", got, want)
	}
}

// A bizarre-but-legal early king move isn't part of any known opening
// line, so book detection should stop at the point of divergence rather
// than crediting the whole game.
func TestBookPliesDiverges(t *testing.T) {
	moves := gameMoves(t, "1. e4 e5 2. Nf3 Nc6 3. Bb5 Nf6 4. Ke2 Bc5")
	got := BookPlies(moves)
	// The shared prefix with the Berlin Defense (e4 e5 Nf3 Nc6 Bb5 Nf6) is
	// 6 plies; White's 4th move (Ke2 instead of O-O) is where play leaves
	// known theory.
	if got != 6 {
		t.Errorf("BookPlies(off-book after Ke2) = %d, want 6", got)
	}
	if got >= len(moves) {
		t.Errorf("BookPlies(%d) should be less than total plies played (%d)", got, len(moves))
	}
}
