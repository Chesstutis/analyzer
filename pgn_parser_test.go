package analyzer

import (
	"fmt"
	"testing"
)

const ()

func TestScholars(t *testing.T) {
	g, err := ParsePGN(scholarsPgn)
	if err != nil {
		t.Errorf("error parsing PGN %s\n", scholarsPgn)
	}

	fmt.Printf("parsed game: %s\n", g)
	fmt.Printf("%s", g.Position().Board().Draw())
}

func TestChessComPgn(t *testing.T) {
	g, err := ParsePGN(chessComPgn)
	if err != nil {
		t.Errorf("error parsing PGN %s\n", chessComPgn)
	}
	fmt.Printf("parsed game: %s\n", g)
	fmt.Printf("%s", g.Position().Board().Draw())
}
