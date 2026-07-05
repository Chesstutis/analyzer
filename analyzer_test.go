package analyzer

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
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

// black blunders in both
const evergreenGame = "1. e4 e5 2. Nf3 Nc6 3. Bc4 Bc5 4. b4 Bxb4 5. c3 Ba5 6. d4 exd4 7. O-O d3 8. Qb3 Qf6 9. e5 Qg6 10. Re1 Nge7 11. Ba3 b5 12. Qxb5 Rb8 13. Qa4 Bb6 14. Nbd2 Bb7 15. Ne4 Qf5 16. Bxd3 Qh5 17. Nf6+ gxf6 18. exf6 Rg8 19. Rad1 Qxf3 20. Rxe7+ Nxe7 {White gives mate in 4 moves.} 21.Qxd7+ Kxd7 22.Bf5+ Ke8 23.Bd7+ Kf8 24.Bxe7# 1-0"
const operaGame = "1.e4 e5 2.Nf3 d6 3.d4 Bg4 4.dxe5 Bxf3 5.Qxf3 dxe5 6.Bc4 Nf6 7.Qb3 Qe7 8.Nc3 c6 9.Bg5 b5 10.Nxb5 cxb5 11.Bxb5+ Nbd7 12.O-O-O Rd8 13.Rxd7 Rxd7 14.Rd1 Qe6 15.Bxd7+ Nxd7 16.Qb8+ Nxb8 17.Rd8+ 1-0"

func TestAnalyzeConfigs(t *testing.T) {
	referenceCfg := Config{
		Threads:            4,
		HashMB:             1024,
		BestMoveDepth:      14,
		VerifyMoveTime:     500 * time.Millisecond,
		BlunderThresholdCP: 200,
		SkipOpeningPlies:   8,
	}

	candidateCfg := []struct {
		name   string
		Config Config
	}{
		{"verify25ms", Config{Threads: 2, HashMB: 256, BestMoveDepth: 10, VerifyMoveTime: 25 * time.Millisecond, BlunderThresholdCP: 200, SkipOpeningPlies: 8}},
		{"verify50ms", Config{Threads: 2, HashMB: 256, BestMoveDepth: 10, VerifyMoveTime: 50 * time.Millisecond, BlunderThresholdCP: 200, SkipOpeningPlies: 8}},
		{"verify75ms", Config{Threads: 2, HashMB: 256, BestMoveDepth: 10, VerifyMoveTime: 75 * time.Millisecond, BlunderThresholdCP: 200, SkipOpeningPlies: 8}},
		{"verify100ms", Config{Threads: 2, HashMB: 256, BestMoveDepth: 10, VerifyMoveTime: 100 * time.Millisecond, BlunderThresholdCP: 200, SkipOpeningPlies: 8}},
	}

	// TODO: get more game data...
	games := [...]string{chesstutisVsAlex, evergreenGame, operaGame}

	for i, g := range games {
		// reference
		pgn, err := chess.PGN(strings.NewReader(g))
		if err != nil {
			t.Fatalf("PGN parsing failed: %v", err)
		}
		game := chess.NewGame(pgn)

		eng, err := uci.New("stockfish")
		if err != nil {
			t.Fatalf("stockfish startup failed: %v", err)
		}
		defer eng.Close()

		a, err := NewAnalyzer(eng, referenceCfg)
		if err != nil {
			t.Fatalf("analyzer creation failed: %v", err)
		}

		referenceAnalysis, err := a.AnalyzeGame(game, chess.Black)
		if err != nil {
			t.Fatalf("analysis failed")
		}
		fmt.Printf("Reference Analysis %d: # puzzles (%d)\n", i+1, len(referenceAnalysis.Puzzles))
		for _, puz := range referenceAnalysis.Puzzles {
			fmt.Printf("\tbest move: %s, player move: %s\n", puz.BestMove, puz.PlayerMove)
		}
		// candidate configs
		for j, cfg := range candidateCfg {
			a, err := NewAnalyzer(eng, cfg.Config)
			if err != nil {
				t.Fatal(err)
			}
			candidateAnalysis, err := a.AnalyzeGame(game, chess.Black)
			fmt.Printf("Candidate Analysis %d: # puzzles (%d)\n", j+1, len(candidateAnalysis.Puzzles))
			for _, puz := range candidateAnalysis.Puzzles {
				fmt.Printf("\tbest move: %s, player move: %s\n", puz.BestMove, puz.PlayerMove)
			}
		}
		fmt.Println()
		fmt.Println("==================================================")
		fmt.Println("==================================================")
		fmt.Println("==================================================")
		fmt.Println()
	}
}
