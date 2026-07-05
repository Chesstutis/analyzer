package config

import (
	"runtime"
	"testing"
	"time"
)

// TestDefaultConfig validates default values
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Threads != max(1, runtime.NumCPU()/2) {
		t.Errorf("Threads = %d, want %d", cfg.Threads, max(1, runtime.NumCPU()/2))
	}
	if cfg.HashMB != 1024 {
		t.Errorf("HashMB = %d, want 1024", cfg.HashMB)
	}
	if cfg.BestMoveDepth != 12 {
		t.Errorf("BestMoveDepth = %d, want 12", cfg.BestMoveDepth)
	}
	if cfg.VerifyMoveTime != 100*time.Millisecond {
		t.Errorf("VerifyMoveTime = %s, want %s", cfg.VerifyMoveTime, 100*time.Millisecond)
	}
	if cfg.BlunderThresholdCP != 200 {
		t.Errorf("BlunderThresholdCP = %d, want 200", cfg.BlunderThresholdCP)
	}
	if cfg.SkipOpeningPlies != 8 {
		t.Errorf("SkipOpeningPlies = %d, want 8", cfg.SkipOpeningPlies)
	}
}

// TestNormalize validates config normalization
func TestNormalize(t *testing.T) {
	defaults := DefaultConfig()
	valid := Config{
		Threads:            max(1, runtime.NumCPU()*2-1),
		HashMB:             16383,
		BestMoveDepth:      49,
		VerifyMoveTime:     10*time.Second - time.Nanosecond,
		BlunderThresholdCP: 9999,
		SkipOpeningPlies:   19,
	}

	tests := []struct {
		name string
		cfg  Config
		want Config
	}{
		{
			name: "valid values are preserved",
			cfg:  valid,
			want: valid,
		},
		{
			name: "zero threads uses default",
			cfg: Config{
				Threads:            0,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            defaults.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "too many threads uses default",
			cfg: Config{
				Threads:            runtime.NumCPU() * 2,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            defaults.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "zero hash uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             0,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             defaults.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "too much hash uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             16384,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             defaults.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "zero best move depth uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      0,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      defaults.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "too high best move depth uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      50,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      defaults.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "zero verify move time uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     0,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     defaults.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "too high verify move time uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     10 * time.Second,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     defaults.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "zero blunder threshold uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: 0,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: defaults.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "too high blunder threshold uses default",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: 10000,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: defaults.BlunderThresholdCP,
				SkipOpeningPlies:   valid.SkipOpeningPlies,
			},
		},
		{
			name: "negative opening plies uses zero",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   -1,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   0,
			},
		},
		{
			name: "too high opening plies uses zero",
			cfg: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   20,
			},
			want: Config{
				Threads:            valid.Threads,
				HashMB:             valid.HashMB,
				BestMoveDepth:      valid.BestMoveDepth,
				VerifyMoveTime:     valid.VerifyMoveTime,
				BlunderThresholdCP: valid.BlunderThresholdCP,
				SkipOpeningPlies:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.cfg); got != tt.want {
				t.Errorf("Normalize(%+v) = %+v, want %+v", tt.cfg, got, tt.want)
			}
		})
	}
}
