package opening

import (
	"strings"
	"sync"

	"github.com/corentings/chess/v2"
	ecobook "github.com/corentings/chess/v2/opening"
)

var (
	once  sync.Once
	lines [][]string
)

func load() {
	book := ecobook.NewBookECO()
	for _, o := range book.Possible(nil) {
		lines = append(lines, strings.Fields(o.PGN()))
	}
}

//returns the number of BookPlies
func BookPlies(moves []*chess.Move) int {
	once.Do(load)
	if len(moves) == 0 {
		return 0
	}

	played := make([]string, len(moves))
	for i, m := range moves {
		if m == nil {
			return i
		}
		played[i] = m.String()
	}

	best := 0
	for _, line := range lines {
		n := min(len(line), len(played))
		i := 0
		for i < n && line[i] == played[i] {
			i++
		}
		if i > best {
			best = i
		}
		if best == len(played) {
			break
		}
	}
	return best
}
