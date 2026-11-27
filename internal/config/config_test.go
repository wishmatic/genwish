package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	return path
}

func setHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	return home
}

func TestLoad(t *testing.T) {
	t.Run("applies the built-in defaults", func(t *testing.T) {
		home := setHome(t)
		t.Setenv("LOG_LEVEL", "")
		t.Setenv("ERROR_DETAIL", "")

		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.LogLevel != "info" || cfg.ErrorDetail != "useful" {
			t.Errorf("Load() = %+v, want the environment defaults", cfg)
		}

		if cfg.Output.Format != "webp" || cfg.Output.Path != filepath.Join(home, ".genwish", "out") {
			t.Errorf("Output = %+v, want webp under the home directory", cfg.Output)
		}

		if cfg.Neo.URL != "http://127.0.0.1:7860" || cfg.Neo.ModelsFile != "" {
			t.Errorf("Neo = %+v, want the local Forge URL", cfg.Neo)
		}

		rate := cfg.Booru.RateLimit
		if rate.RPS != 1 || rate.Burst != 1 || rate.Timeout != 20 {
			t.Errorf("RateLimit = %+v, want 1/1/20", rate)
		}
	})

	t.Run("reads a config file", func(t *testing.T) {
		home := setHome(t)
		t.Setenv("LOG_LEVEL", "debug")
		t.Setenv("ERROR_DETAIL", "verbose")

		path := writeConfig(t, `
[output]
format = "png"
path = "~/art"

[neo]
url = "http://forge:7860"
models_file = "~/models.txt"

[booru.rate_limit]
rps = 2.5
burst = 3
timeout = 30
`)

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.LogLevel != "debug" || cfg.ErrorDetail != "verbose" {
			t.Errorf("Load() = %+v, want the environment values", cfg)
		}

		if cfg.Output.Format != "png" || cfg.Output.Path != filepath.Join(home, "art") {
			t.Errorf("Output = %+v, want a png under the expanded home path", cfg.Output)
		}

		if cfg.Neo.URL != "http://forge:7860" || cfg.Neo.ModelsFile != filepath.Join(home, "models.txt") {
			t.Errorf("Neo = %+v, want the configured URL and model list", cfg.Neo)
		}

		rate := cfg.Booru.RateLimit
		if rate.RPS != 2.5 || rate.Burst != 3 || rate.Timeout != 30 {
			t.Errorf("RateLimit = %+v, want the configured limits", rate)
		}
	})

	t.Run("keeps defaults for omitted keys", func(t *testing.T) {
		home := setHome(t)

		path := writeConfig(t, `
[output]
format = "jpeg"
`)

		cfg, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.Output.Format != "jpeg" || cfg.Output.Path != filepath.Join(home, ".genwish", "out") {
			t.Errorf("Output = %+v, want only the format overridden", cfg.Output)
		}

		if cfg.Neo.URL != "http://127.0.0.1:7860" {
			t.Errorf("Neo.URL = %q, want the default", cfg.Neo.URL)
		}

		if cfg.Booru.RateLimit.Timeout != 20 {
			t.Errorf("Timeout = %d, want the default", cfg.Booru.RateLimit.Timeout)
		}
	})

	t.Run("rejects a missing explicit file", func(t *testing.T) {
		setHome(t)
		path := filepath.Join(t.TempDir(), "missing.toml")

		_, err := Load(path)
		if err == nil || !strings.Contains(err.Error(), "config file") || !strings.Contains(err.Error(), path) {
			t.Errorf("Load() error = %v, want the missing path named", err)
		}
	})

	t.Run("rejects an unknown key", func(t *testing.T) {
		setHome(t)
		path := writeConfig(t, `
[output]
formatt = "png"
`)

		_, err := Load(path)
		if err == nil || !strings.Contains(err.Error(), "unknown key") || !strings.Contains(err.Error(), "output.formatt") {
			t.Errorf("Load() error = %v, want the unknown key named", err)
		}
	})

	t.Run("rejects an environment setting in the file", func(t *testing.T) {
		setHome(t)
		path := writeConfig(t, "LogLevel = \"debug\"\n")

		_, err := Load(path)
		if err == nil || !strings.Contains(err.Error(), "unknown key") {
			t.Errorf("Load() error = %v, want environment settings rejected from the file", err)
		}
	})

	t.Run("rejects malformed TOML", func(t *testing.T) {
		setHome(t)
		path := writeConfig(t, "[output\nformat = \"png\"\n")

		_, err := Load(path)
		if err == nil || !strings.Contains(err.Error(), "config file") {
			t.Errorf("Load() error = %v, want a parse failure naming the file", err)
		}
	})

	t.Run("ignores removed environment variables", func(t *testing.T) {
		home := setHome(t)

		for _, key := range []string{
			"OUTPUT_FORMAT", "SD_URL", "NEO_MODELS_FILE", "RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST", "REQUEST_TIMEOUT_SECONDS", "MAX_LIMIT", "MAX_OFFSET", "BLOCKED_TAGS",
		} {
			t.Setenv(key, "nonsense")
		}

		cfg, err := Load("")
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if cfg.Output.Format != "webp" || cfg.Output.Path != filepath.Join(home, ".genwish", "out") {
			t.Errorf("Output = %+v, want the defaults", cfg.Output)
		}

		if cfg.Neo.URL != "http://127.0.0.1:7860" || cfg.Neo.ModelsFile != "" {
			t.Errorf("Neo = %+v, want the defaults", cfg.Neo)
		}

		if cfg.Booru.RateLimit != (RateLimitConfig{RPS: 1, Burst: 1, Timeout: 20}) {
			t.Errorf("RateLimit = %+v, want the defaults", cfg.Booru.RateLimit)
		}
	})
}

func TestDefaults(t *testing.T) {
	t.Run("matches the documented values", func(t *testing.T) {
		cfg := Defaults()

		if cfg.LogLevel != "info" || cfg.ErrorDetail != "useful" {
			t.Errorf("Defaults() = %+v, want the environment defaults", cfg)
		}

		if cfg.Output != (OutputConfig{Format: "webp", Path: "~/.genwish/out"}) {
			t.Errorf("Output = %+v, want webp under ~/.genwish/out", cfg.Output)
		}

		if cfg.Neo != (NeoConfig{URL: "http://127.0.0.1:7860"}) {
			t.Errorf("Neo = %+v, want the local Forge URL", cfg.Neo)
		}
	})
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		rate    RateLimitConfig
		wantErr string
	}{
		{
			name:    "accepts the defaults",
			rate:    RateLimitConfig{RPS: 1, Burst: 1, Timeout: 20},
			wantErr: "",
		},
		{
			name:    "rejects a zero rate",
			rate:    RateLimitConfig{RPS: 0, Burst: 1, Timeout: 20},
			wantErr: "booru.rate_limit.rps",
		},
		{
			name:    "rejects a zero burst",
			rate:    RateLimitConfig{RPS: 1, Burst: 0, Timeout: 20},
			wantErr: "booru.rate_limit.burst",
		},
		{
			name:    "rejects a zero timeout",
			rate:    RateLimitConfig{RPS: 1, Burst: 1, Timeout: 0},
			wantErr: "booru.rate_limit.timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := (Config{Booru: BooruConfig{RateLimit: tt.rate}}).Validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error: %v", err)
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Validate() error = %v, want it to name %s", err, tt.wantErr)
			}
		})
	}
}

func TestVerboseErrors(t *testing.T) {
	tests := []struct {
		name   string
		detail string
		want   bool
	}{
		{name: "verbose", detail: "verbose", want: true},
		{name: "useful", detail: "useful", want: false},
		{name: "empty", detail: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (Config{ErrorDetail: tt.detail}).VerboseErrors(); got != tt.want {
				t.Errorf("VerboseErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}
