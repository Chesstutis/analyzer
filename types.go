package analyzer

import (
	"time"

	chess "github.com/corentings/chess/v2"
)

// --------------------------
//   Move classifications
// --------------------------

type MoveClassification int

const (
	//
	// numbers could probably use some tweaking
	//
	// ClassBrilliant - difficult to find move
	ClassBest       = iota // top engine move
	ClassExcellent         // <= 10 cpl -
	ClassGood              // <= 25 cpl -
	ClassInaccuracy        // <= 50 cpl -
	ClassMistake           // <= 100 cpl -
	ClassBlunder           // > 100 cpl -
	ClassMiss              // missed forced mate, or missed winning sequence
	ClassForced            // only legal move
)

// String returns the human-readable label.
func (c MoveClassification) String() string {
	switch c {
	case ClassBest:
		return "Best"
	case ClassExcellent:
		return "Excellent"
	case ClassGood:
		return "Good"
	case ClassInaccuracy:
		return "Inaccuracy"
	case ClassMistake:
		return "Mistake"
	case ClassBlunder:
		return "Blunder"
	case ClassMiss:
		return "Miss"
	case ClassForced:
		return "Forced"
	default:
		return "Unknown"
	}
}

// ------------------------
//   Move Analysis
// ------------------------

type Eval struct {
	IsMate bool
}

type AlternativeMove struct {
	Move     *chess.Move
	MoveEval Eval
}

type MoveAnalysis struct {
	Color chess.Color // color who makes the move
	Move  chess.Move  // the move that is played

	EvalBefore Eval // engine eval before the move was made
	EvalAfter  Eval // engine eval after the move is made

	BestMove *chess.Move // Best move in the position
	BestEval Eval        // evaluation of the best move

	Class MoveClassification

	MakePuzzle bool // move is mistake, blunder, or miss. -> make puzzle.
}

// -----------------------
//   Puzzles
// -----------------------

type PuzzleTheme int

const (
	PuzzleBlunder = iota
	PuzzleMistake
	PuzzleMiss
)

type Puzzle struct {
	Theme PuzzleTheme
	For   chess.Color

	StartPos   chess.Position
	BestMove   *chess.Move
	PlayerMove *chess.Move

	// list of N best moves in position
	// set N in the config
	Alternatives []AlternativeMove
}

func (AnalyzedGame) PuzzlesForColor(c chess.Color) []Puzzle {

}

// ---------------------
//    Analysis Result
// ---------------------

type AnalyzedGame struct {
	Puzzles []Puzzle

	WhiteAccuracy float64 // [0, 100]
	BlackAccuracy float64

	// list of move analyses
	// used for game analysis
	// has board eval and move class
	MoveStats []MoveAnalysis
}

// -----------
//   Config
// -----------

type Config struct {
	// only one of these two can be set or else config is invalid
	EngineDepth  int
	AnalysisTime time.Duration
}
