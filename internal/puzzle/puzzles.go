package puzzle

import (
	"fmt"
	"strings"

	"github.com/corentings/chess/v2"
)

// TODO: implement move classifications
//   the blunder threshold value should be removes from the config and IsPuzzle should no longer take that as a parameter.
//   all puzzle logic should be handled in this file
//   move classes: Brilliant (might be hard to determine: should be a difficult to find move, skip for now), Best, Great, Good, Inaccuraccy, Mistake, Blunder, Forced

// IsPuzzle determines whether the gap between the best move and player move is large enough.
func IsPuzzle(bestEval, playerEval, threshold int) bool {
	return bestEval-playerEval >= threshold
}

func TryGeneratePuzzle(pos *chess.Position, bestScore int, bestMove *chess.Move, playerScore int, playerMove *chess.Move, threshold int) *Puzzle {
	if bestScore-playerScore < threshold {
		return nil
	}

	return &Puzzle{
		Position:       pos,
		PlayerMove:     playerMove,
		BestMove:       bestMove,
		Lines:          nil,
		AlternateMoves: nil,
	}

}

// SolveLines finds full puzzle solution variations from a position.
func SolveLines(pos *chess.Position, depth int, evalThreshold int) [][]chess.Move {
	// TODO: Implement iterative deepening to find solution lines.
	// Returns multiple variations (main line + alternatives).
	return nil
}

func AlternameMoves(pos *chess.Position) []chess.Move {

	return nil
}

func (p Puzzle) String() string {
	var sb strings.Builder

	sb.WriteString("Puzzle\n")
	if p.Position == nil {
		sb.WriteString("Position: <nil>\n")
	} else {
		sb.WriteString("Board:")
		sb.WriteString(p.Position.Board().Draw())
		sb.WriteString(fmt.Sprintf("FEN: %s\n", p.Position.String()))
		sb.WriteString(fmt.Sprintf("Side to move: %s\n", p.Position.Turn()))
	}

	sb.WriteString(fmt.Sprintf("Player move: %s\n", formatMove(p.PlayerMove)))
	sb.WriteString(fmt.Sprintf("Best move: %s\n", formatMove(p.BestMove)))

	if len(p.Lines) == 0 {
		sb.WriteString("Solution lines: none\n")
	} else {
		sb.WriteString("Solution lines:\n")
		for i, line := range p.Lines {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, formatMoveLine(line)))
		}
	}

	if len(p.AlternateMoves) == 0 {
		sb.WriteString("Alternate moves: none")
	} else {
		sb.WriteString("Alternate moves:\n")
		for i, alternate := range p.AlternateMoves {
			sb.WriteString(fmt.Sprintf(
				"  %d. %s (eval: %d cp, difference: %d cp)",
				i+1,
				formatMove(alternate.Move),
				alternate.Evaluation,
				alternate.Difference,
			))
			if i < len(p.AlternateMoves)-1 {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String() + "\n\n"
}

func formatMove(move *chess.Move) string {
	if move == nil {
		return "<nil>"
	}
	return move.String()
}

func formatMoveLine(line []chess.Move) string {
	if len(line) == 0 {
		return "<empty>"
	}

	moves := make([]string, 0, len(line))
	for i := range line {
		moves = append(moves, line[i].String())
	}
	return strings.Join(moves, " ")
}
