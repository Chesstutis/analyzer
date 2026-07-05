package engine

import "github.com/corentings/chess/v2"

// SearchResult represents an engine search outcome
type SearchResult struct {
	BestMove   *chess.Move
	Evaluation int // Centipawns
	Depth      int
	// TODO: Add Info field from uci.Info for richer context
}
