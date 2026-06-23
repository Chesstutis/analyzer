package analyzer

import (
	"fmt"
	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
	"log"
	"strings"
	"testing"
	"time"
)

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

	gameAnalysis, err := a.AnalyzeGame(game, chess.White)
	fmt.Println("Puzzles: ")
	fmt.Println(gameAnalysis.Puzzles)
}

const chesstutisVsAlex = "1. e4 c6 2. Nc3 d5 3. d4 dxe4 4. Nxe4 Bf5 5. Bd3 Qxd4 6. Nf3 Qb6 7. O-O Nd7 8. b3 Ngf6 9. Nxf6+ Nxf6 10. Bxf5 e6 11. Bd3 Bc5 12. Bb2 Ng4 13. h3 Bxf2+ 14. Kh1 Ne3 15. Qe2 Nxf1 16. Qxf2 Qxf2 17. Bxf1 O-O 18. Nd4 Rad8 19. Rd1 c5 20. Nb5 Rxd1 0-1"

func TestChesstutisVsAlex(t *testing.T) {
	pgn, err := chess.PGN(strings.NewReader(chesstutisVsAlex))

	if err != nil {
		t.Fatalf("PGN parsing failed: %v", err)
	}

	game := chess.NewGame(pgn)

	fmt.Println(game.String())

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	a, err := NewAnalyzer(eng)

	if err != nil {
		panic(err)
	}

	gameAnalysis, err := a.AnalyzeGame(game, chess.White)
	fmt.Println("Puzzles: ")
	fmt.Println(gameAnalysis.Puzzles)
}

func TestSkipMoves(t *testing.T) {
	pgn, err := chess.PGN(strings.NewReader(chesstutisVsAlex))

	if err != nil {
		t.Fatalf("PGN parsing failed: %v", err)
	}

	game := chess.NewGame(pgn)

	fmt.Println(game.String())

	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	a, err := NewAnalyzer(eng)

	if err != nil {
		panic(err)
	}

	gameAnalysis, err := a.AnalyzeGame(game, chess.White)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range gameAnalysis.Puzzles {
		if p.Position.Turn() != chess.White {
			log.Fatal("puzzle for wrong player returned")
		}
	}

	gameAnalysis, err = a.AnalyzeGame(game, chess.Black)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range gameAnalysis.Puzzles {
		if p.Position.Turn() != chess.Black {
			log.Fatal("puzzle for wrong player returned")
		}
	}
}

func TestNilGame(t *testing.T) {
	eng, err := uci.New("stockfish")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	a, err := NewAnalyzer(eng)
	if err != nil {
		panic(err)
	}

	gameAnalysis, err := a.AnalyzeGame(nil, chess.NoColor)

	if gameAnalysis != nil {
		t.Errorf("nil game didnt return nil game analysis")
	}
}