package opening

import (
	"strings"
	"sync"

	"github.com/corentings/chess/v2"
	ecobook "github.com/corentings/chess/v2/opening"
)

type trieNode struct {
	children map[string]*trieNode
}

var (
	once sync.Once
	root *trieNode
)

func load() {
	root = &trieNode{children: make(map[string]*trieNode)}
	book := ecobook.NewBookECO()
	for _, o := range book.Possible(nil) {
		node := root
		for _, tok := range strings.Fields(o.PGN()) {
			next, ok := node.children[tok]
			if !ok {
				next = &trieNode{children: make(map[string]*trieNode)}
				node.children[tok] = next
			}
			node = next
		}
	}
}

// returns the number of BookPlies
func BookPlies(moves []*chess.Move) int {
	once.Do(load)
	if len(moves) == 0 {
		return 0
	}

	node := root
	for i, m := range moves {
		if m == nil {
			return i
		}
		next, ok := node.children[m.String()]
		if !ok {
			return i
		}
		node = next
	}
	return len(moves)
}