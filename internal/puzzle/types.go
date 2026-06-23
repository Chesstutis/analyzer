package puzzle

import "github.com/corentings/chess/v2"

// Puzzle represents a candidate puzzle position
type Puzzle struct {
	Position       *chess.Position // Position where puzzle occurs
	PlayerMove     *chess.Move     // Move player made
	BestMove       *chess.Move     // Engine's best move
	Lines          [][]chess.Move  // Full solution variations
	AlternateMoves []AlternateMove // Similar-strength alternatives
}

// AlternateMove represents a move with similar evaluation
type AlternateMove struct {
	Move       *chess.Move // The alternate move
	Evaluation int         // Its centipawn score
	Difference int         // Gap from best move eval
}
