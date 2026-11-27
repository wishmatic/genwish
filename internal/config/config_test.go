package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("applies the defaults", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("OUTPUT_FORMAT", "")
		t.Setenv("SD_URL", "")
		t.Setenv("ERROR_DETAIL", "")
		t.Setenv("RATE_LIMIT_RPS", "")
		t.Setenv("RATE_LIMIT_BURST", "")
		t.Setenv("REQUEST_TIMEOUT_SECONDS", "")
		t.Setenv("MAX_LIMIT", "")
		t.Setenv("MAX_OFFSET", "")
		t.Setenv("NEO_MODELS_FILE", "")
		t.Setenv("BLOCKED_TAGS", "")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.LogLevel != "info" || cfg.DefaultOutput != "webp" || cfg.SDURL != "http://127.0.0.1:7860" {
			t.Errorf("Load() = %+v, want the documented defaults", cfg)
		}

		if cfg.RateLimitRPS != 1 || cfg.RateLimitBurst != 1 || cfg.RequestTimeoutSeconds != 20 {
			t.Errorf("Load() = %+v, want the booru defaults", cfg)
		}

		if cfg.MaxLimit != 100 || cfg.MaxOffset != 1000 {
			t.Errorf("Load() = %+v, want the tag limits", cfg)
		}

		if cfg.NEOModelsFile != "" {
			t.Errorf("NEOModelsFile = %q, want empty", cfg.NEOModelsFile)
		}
	})

	t.Run("reads the environment", func(t *testing.T) {
		t.Setenv("OUTPUT_FORMAT", "png")
		t.Setenv("SD_URL", "http://forge:7860")
		t.Setenv("ERROR_DETAIL", "verbose")
		t.Setenv("RATE_LIMIT_RPS", "2.5")
		t.Setenv("RATE_LIMIT_BURST", "3")
		t.Setenv("REQUEST_TIMEOUT_SECONDS", "30")
		t.Setenv("MAX_LIMIT", "50")
		t.Setenv("MAX_OFFSET", "500")
		t.Setenv("NEO_MODELS_FILE", "/tmp/neo-models.txt")
		t.Setenv("BLOCKED_TAGS", "gore, blood")

		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.DefaultOutput != "png" || cfg.SDURL != "http://forge:7860" {
			t.Errorf("Load() = %+v, want the output and Forge settings", cfg)
		}

		if cfg.NEOModelsFile != "/tmp/neo-models.txt" {
			t.Errorf("NEOModelsFile = %q, want the configured path", cfg.NEOModelsFile)
		}

		if !cfg.VerboseErrors() {
			t.Errorf("Load() = %+v, want verbose errors", cfg)
		}
	})
}

func TestBlockedTags(t *testing.T) {
	t.Run("normalizes and deduplicates", func(t *testing.T) {
		cfg := Config{BlockedTagsRaw: "Blue Hair, blue_hair, , gore"}

		got := cfg.BlockedTags()
		want := []string{"blue_hair", "gore"}

		if len(got) != len(want) {
			t.Fatalf("BlockedTags() = %v, want %v", got, want)
		}

		for i := range want {
			if got[i] != want[i] {
				t.Errorf("BlockedTags()[%d] = %q, want %q", i, got[i], want[i])
			}
		}
	})

	t.Run("is empty when unset", func(t *testing.T) {
		if got := (Config{}).BlockedTags(); len(got) != 0 {
			t.Errorf("BlockedTags() = %v, want none", got)
		}
	})
}

func TestValidate(t *testing.T) {
	t.Run("accepts the defaults", func(t *testing.T) {
		if err := (Config{RateLimitRPS: 1, RateLimitBurst: 1, MaxLimit: 100, MaxOffset: 1000}).Validate(); err != nil {
			t.Fatalf("Validate() error: %v", err)
		}
	})

	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name:    "rejects a zero rate limit",
			cfg:     Config{RateLimitRPS: 0, RateLimitBurst: 1, MaxLimit: 100, MaxOffset: 1000},
			wantErr: "RATE_LIMIT_RPS",
		},
		{
			name:    "rejects a zero burst",
			cfg:     Config{RateLimitRPS: 1, RateLimitBurst: 0, MaxLimit: 100, MaxOffset: 1000},
			wantErr: "RATE_LIMIT_BURST",
		},
		{
			name:    "rejects a large max limit",
			cfg:     Config{RateLimitRPS: 1, RateLimitBurst: 1, MaxLimit: 101, MaxOffset: 1000},
			wantErr: "MAX_LIMIT",
		},
		{
			name:    "rejects a zero max offset",
			cfg:     Config{RateLimitRPS: 1, RateLimitBurst: 1, MaxLimit: 100, MaxOffset: 0},
			wantErr: "MAX_OFFSET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %v, want it to name %s", err, tt.wantErr)
			}
		})
	}
}
