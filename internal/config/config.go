package config

import (
	"runtime"
	"time"
)

// Config holds analyzer settings
type Config struct {
	Threads            int
	HashMB             int
	BestMoveDepth      int
	VerifyMoveTime     time.Duration
	BlunderThresholdCP int
	SkipOpeningPlies   int
}

// DefaultConfig returns analyzer defaults
func DefaultConfig() Config {
	return Config{
		Threads:            max(1, runtime.NumCPU()/2),
		HashMB:             1024,
		BestMoveDepth:      12,
		VerifyMoveTime:     100 * time.Millisecond,
		BlunderThresholdCP: 200,
		SkipOpeningPlies:   8,
	}
}

// Normalize resets invalid values to defaults
func Normalize(cfg Config) Config {
	if cfg.Threads <= 0 || cfg.Threads >= runtime.NumCPU()*2 {
		cfg.Threads = DefaultConfig().Threads
	}
	if cfg.HashMB <= 0 || cfg.HashMB >= 16384 {
		cfg.HashMB = DefaultConfig().HashMB
	}
	if cfg.BestMoveDepth <= 0 || cfg.BestMoveDepth >= 50 {
		cfg.BestMoveDepth = DefaultConfig().BestMoveDepth
	}
	if cfg.VerifyMoveTime <= 0 || cfg.VerifyMoveTime >= time.Second*10 {
		cfg.VerifyMoveTime = DefaultConfig().VerifyMoveTime
	}
	if cfg.BlunderThresholdCP <= 0 || cfg.BlunderThresholdCP >= 10000 {
		cfg.BlunderThresholdCP = DefaultConfig().BlunderThresholdCP
	}
	if cfg.SkipOpeningPlies < 0 || cfg.SkipOpeningPlies >= 20 {
		cfg.SkipOpeningPlies = 0
	}
	return cfg
}
