package analyzer

import (
	"fmt"
	"time"

	"github.com/corentings/chess/v2"
	"github.com/corentings/chess/v2/uci"
)

type Analyzer struct {
	engine *uci.Engine
}

type GameAnalysis struct {
	Puzzles []Puzzle
}

type Puzzle struct {
	Position   *chess.Position
	PlayerMove *chess.Move
	BestMove   *chess.Move
}

func NewAnalyzer(eng *uci.Engine) (*Analyzer, error) {
	return &Analyzer{
		engine: eng,
	}, nil
}

func (a *Analyzer) Close() {
	a.engine.Close()
}

func (a *Analyzer) AnalyzeGame(game *chess.Game) (*GameAnalysis, error) {
	gameAnalysis := GameAnalysis{}
	// initialize uci engine with new game
	if err := a.engine.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	for _, move := range game.Moves() {
		pos := move.Parent().Position()

		// find best engine move
		cmdPos := uci.CmdPosition{Position: pos}
		cmdGo := uci.CmdGo{MoveTime: time.Second} // should replace time with config value

		if err := a.engine.Run(cmdPos, cmdGo); err != nil {
			panic(err)
		}

		bestMove := a.engine.SearchResults().BestMove
		bestScore := a.engine.SearchResults().Info.Score

		// analyer player move
		cmdGo = uci.CmdGo{
			MoveTime:    time.Second,
			SearchMoves: []*chess.Move{move},
		}

		if err := a.engine.Run(cmdPos, cmdGo); err != nil {
			return nil, err
		}

		playerScore := a.engine.SearchResults().Info.Score

		// make puzzles if blunder

		scoreDiff := bestScore.CP - playerScore.CP
		if scoreDiff >= 200 { // replace with config value
			newPuzzle := Puzzle{
				Position:   pos,
				BestMove:   bestMove,
				PlayerMove: move,
			}
			gameAnalysis.Puzzles = append(gameAnalysis.Puzzles, newPuzzle)
		}
	}

	return &gameAnalysis, nil
}
