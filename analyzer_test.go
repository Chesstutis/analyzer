package analyzer

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
)

func TestAnalyzer(t *testing.T) {

	// create random game
	game := chess.NewGame()

	for game.Outcome() == chess.NoOutcome {
		moves := game.ValidMoves()
		move := moves[rand.Intn(len(moves))]
		if err := game.Move(&move, nil); err != nil {
			panic(err)
		}
	}

	// game is made time to analyze

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	a, err := NewAnalyzer(eng)

	if err != nil {
		panic(err)
	}

	gameAnalysis, err := a.AnalyzeGame(game)
	fmt.Println(gameAnalysis.Puzzles)
}
