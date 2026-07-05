package engine

import (
	"sync"

	"github.com/corentings/chess/v2/uci"
)

// Manager wraps UCI engine with thread safety
type Manager struct {
	engine *uci.Engine
	mu     sync.Mutex
}

// NewManager creates a new engine manager
func NewManager(eng *uci.Engine) *Manager {
	return &Manager{
		engine: eng,
	}
}

// Close shuts down the engine
func (m *Manager) Close() error {
	if m.engine != nil {
		m.engine.Close()
	}
	return nil
}

// TODO: Add orchestration methods (e.g., SearchBest, SearchMove)
