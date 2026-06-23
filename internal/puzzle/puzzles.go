package puzzle

import "github.com/corentings/chess/v2"

// IsPuzzle determines whether the gap between the best move and player move is large enough.
func IsPuzzle(bestEval, playerEval, threshold int) bool {
	return bestEval-playerEval >= threshold
}

// SolveLines finds full puzzle solution variations from a position.
func SolveLines(pos *chess.Position, depth int, evalThreshold int) [][]chess.Move {
	// TODO: Implement iterative deepening to find solution lines.
	// Returns multiple variations (main line + alternatives).
	return nil
}
