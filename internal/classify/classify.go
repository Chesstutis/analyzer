package classify

// import "github.com/corentings/chess/v2"

// Grade represents move quality classification
type Grade int

const (
	Brilliant   Grade = iota // Significantly better than engine
	Best                      // Engine's best move
	Good                      // Minor improvement available
	Inaccuracy                // Moderate error
	Mistake                   // Significant error
	Blunder                   // Critical error
)

// Classify determines move quality from position and play evaluations
func Classify(bestEval, playerEval int) Grade {
	// TODO: Implement classification logic based on eval gap
	return Blunder
}

// ThresholdCP defines eval gaps for each grade
// TODO: Define grade boundaries
