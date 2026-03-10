package analyzer

import (
	chess "github.com/corentings/chess/v2"
	"strings"
)

func ParsePGN(pgn string) (game *chess.Game, err error) {

	reader := strings.NewReader(pgn)

	parsedPgn, err := chess.PGN(reader)
	if err != nil {
		panic(err)
	}

	game = chess.NewGame(parsedPgn)

	return nil, nil
}
