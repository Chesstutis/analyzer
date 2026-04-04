package analyzer

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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

func TestStockfishGame(t *testing.T) {
	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()
	// initialize uci with new game
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}
	// have stockfish play speed chess against itself (10 msec per move)
	game := chess.NewGame()
	for game.Outcome() == chess.NoOutcome {
		cmdPos := uci.CmdPosition{Position: game.Position()}
		cmdGo := uci.CmdGo{MoveTime: time.Second / 1000}
		if err := eng.Run(cmdPos, cmdGo); err != nil {
			panic(err)
		}
		move := eng.SearchResults().BestMove
		if err := game.Move(move, nil); err != nil {
			panic(err)
		}
	}
	fmt.Println(game.String())

	a, err := NewAnalyzer(eng)

	if err != nil {
		panic(err)
	}

	gameAnalysis, err := a.AnalyzeGame(game)
	fmt.Println("Puzzles: ")
	fmt.Println(gameAnalysis.Puzzles)
}
