package analyzer

import (
	"fmt"
	"log"
	"os"

	chess "github.com/corentings/chess/v2"
)

func ParsePGN(filePath string) (*chess.Game, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, nil
}

func ParsePGNString(pgn string) (*chess.Game, error) {

	return nil, nil
}

func ParseMultiPGN(filePath string) ([]*chess.Game, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, nil
}
