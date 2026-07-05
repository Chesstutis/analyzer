package analyzer

import (
	"fmt"
	"sync"

	"github.com/chesstutis/analyzer/internal/config"
	"github.com/chesstutis/analyzer/internal/puzzle"
	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
)

// Re-export internal types for public API
type Puzzle = puzzle.Puzzle
type AlternateMove = puzzle.AlternateMove
type Config = config.Config

type (
	Analyzer struct {
		engine *uci.Engine
		cfg    Config
		mu     sync.Mutex
	}

	GameAnalysis struct {
		Puzzles []Puzzle
	}
)

func NewAnalyzer(eng *uci.Engine, cfgs ...Config) (*Analyzer, error) {
	if eng == nil {
		return nil, fmt.Errorf("nil engine")
	}

	cfg := config.DefaultConfig()
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	cfg = config.Normalize(cfg)

	if err := eng.Run(
		uci.CmdUCI,
		uci.CmdSetOption{Name: "Threads", Value: fmt.Sprintf("%d", cfg.Threads)},
		uci.CmdSetOption{Name: "Hash", Value: fmt.Sprintf("%d", cfg.HashMB)},
		uci.CmdSetOption{Name: "Ponder", Value: "false"},
		uci.CmdIsReady,
	); err != nil {
		return nil, err
	}

	return &Analyzer{
		engine: eng,
		cfg:    cfg,
	}, nil
}

func (a *Analyzer) Close() {
	if a.engine != nil {
		a.engine.Close()
	}
}

func (a *Analyzer) AnalyzeGame(game *chess.Game, color chess.Color) (*GameAnalysis, error) {
	if game == nil {
		return nil, fmt.Errorf("error nil game")
	}
	if color == chess.NoColor {
		return nil, fmt.Errorf("error no color")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	analysis := &GameAnalysis{
		Puzzles: make([]Puzzle, 0, max(4, len(game.Moves())/12)),
	}

	if err := a.engine.Run(uci.CmdUCINewGame, uci.CmdIsReady); err != nil {
		return nil, err
	}

	notation := chess.UCINotation{}
	moves := game.Moves()

	parity := 0
	if color == chess.Black {
		parity = 1
	}

	for ply, move := range moves {
		// skip other players moves
		if ply%2 != parity {
			continue
		}
		if move == nil || move.Parent() == nil {
			continue
		}

		if ply < a.cfg.SkipOpeningPlies {
			continue
		}

		pos := move.Parent().Position()
		if pos == nil {
			continue
		}

		cmdPos := uci.CmdPosition{Position: pos}

		// Main search.
		if err := a.engine.Run(cmdPos, uci.CmdGo{Depth: a.cfg.BestMoveDepth}); err != nil {
			return nil, err
		}

		bestRes := a.engine.SearchResults()
		bestMove := bestRes.BestMove
		if bestMove == nil {
			continue
		}

		bestScore := bestRes.Info.Score.CP

		// Skip expensive second search if player already played engine move.
		if sameMove(pos, notation, move, bestMove) {
			continue
		}

		// Verify only candidate mistakes.
		if err := a.engine.Run(cmdPos, uci.CmdGo{
			MoveTime:    a.cfg.VerifyMoveTime,
			SearchMoves: []*chess.Move{move},
		}); err != nil {
			return nil, err
		}

		playerRes := a.engine.SearchResults()
		playerScore := playerRes.Info.Score.CP

		puz := puzzle.TryGeneratePuzzle(pos, bestScore, bestMove, playerScore, move, a.cfg.BlunderThresholdCP)
		if puz != nil {
			analysis.Puzzles = append(analysis.Puzzles, *puz)
		}
	}
	return analysis, nil
}

func sameMove(pos *chess.Position, notation chess.UCINotation, a, b *chess.Move) bool {
	if a == nil || b == nil || pos == nil {
		return false
	}
	return notation.Encode(pos, a) == notation.Encode(pos, b)
}