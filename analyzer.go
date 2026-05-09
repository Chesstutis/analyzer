package analyzer

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
)

type Config struct {
	Threads            int
	HashMB             int
	BestMoveDepth      int
	VerifyMoveTime     time.Duration
	BlunderThresholdCP int
	SkipOpeningPlies   int
}

func DefaultConfig() Config {
	return Config{
		Threads:            max(1, runtime.NumCPU()/2),
		HashMB:             1024,
		BestMoveDepth:      12,
		VerifyMoveTime:     100 * time.Millisecond,
		BlunderThresholdCP: 200,
		SkipOpeningPlies:   8,
	}
}

type (
	Analyzer struct {
		engine *uci.Engine
		cfg    Config
		mu     sync.Mutex
	}

	GameAnalysis struct {
		Puzzles []Puzzle
	}

	Puzzle struct {
		Position   *chess.Position
		PlayerMove *chess.Move
		BestMove   *chess.Move
	}
)

func NewAnalyzer(eng *uci.Engine, cfgs ...Config) (*Analyzer, error) {
	if eng == nil {
		return nil, fmt.Errorf("nil engine")
	}

	cfg := DefaultConfig()
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	if cfg.Threads <= 0 || cfg.Threads >= runtime.NumCPU()*2 {
		cfg.Threads = DefaultConfig().Threads
	}
	if cfg.HashMB <= 0 || cfg.HashMB >= 16384 {
		cfg.HashMB = DefaultConfig().HashMB
	}
	if cfg.BestMoveDepth <= 0 || cfg.BestMoveDepth >= 50 {
		cfg.BestMoveDepth = DefaultConfig().BestMoveDepth
	}
	if cfg.VerifyMoveTime <= 0 || cfg.VerifyMoveTime >= time.Second*10 {
		cfg.VerifyMoveTime = DefaultConfig().VerifyMoveTime
	}
	if cfg.BlunderThresholdCP <= 0 || cfg.BlunderThresholdCP >= 10 {
		cfg.BlunderThresholdCP = DefaultConfig().BlunderThresholdCP
	}
	if cfg.SkipOpeningPlies < 0 || cfg.BlunderThresholdCP >= 20 {
		cfg.SkipOpeningPlies = 0
	}

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

func (a *Analyzer) AnalyzeGame(game *chess.Game) (*GameAnalysis, error) {
	if game == nil {
		return nil, fmt.Errorf("nil game")
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

	for ply, move := range moves {
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

		if bestScore-playerScore >= a.cfg.BlunderThresholdCP {
			analysis.Puzzles = append(analysis.Puzzles, Puzzle{
				Position:   pos,
				PlayerMove: move,
				BestMove:   bestMove,
			})
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
