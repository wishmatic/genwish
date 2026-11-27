package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/wishmatic/genwish/internal/booru"
)

const maxTagLimit = 100

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	DefaultOutput string `env:"OUTPUT_FORMAT" envDefault:"webp"`
	SDURL         string `env:"SD_URL" envDefault:"http://127.0.0.1:7860"`

	ErrorDetail    string  `env:"ERROR_DETAIL" envDefault:"useful"`
	RateLimitRPS   float64 `env:"RATE_LIMIT_RPS" envDefault:"1"`
	RateLimitBurst int     `env:"RATE_LIMIT_BURST" envDefault:"1"`

	RequestTimeoutSeconds int `env:"REQUEST_TIMEOUT_SECONDS" envDefault:"20"`

	MaxLimit  int `env:"MAX_LIMIT" envDefault:"100"`
	MaxOffset int `env:"MAX_OFFSET" envDefault:"1000"`

	NEOModelsFile string `env:"NEO_MODELS_FILE"`

	BlockedTagsRaw string `env:"BLOCKED_TAGS"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	return cfg, nil
}

func (c Config) VerboseErrors() bool {
	return c.ErrorDetail == "verbose"
}

func (c Config) BlockedTags() []string {
	seen := make(map[string]bool)
	out := make([]string, 0)

	for _, raw := range splitList(c.BlockedTagsRaw) {
		tag := booru.NormalizeTag(raw)
		if tag == "" || seen[tag] {
			continue
		}

		seen[tag] = true
		out = append(out, tag)
	}

	return out
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}

	return out
}

func (c Config) Validate() error {
	if c.RateLimitRPS <= 0 {
		return fmt.Errorf("RATE_LIMIT_RPS must be greater than 0, got %v", c.RateLimitRPS)
	}

	if c.RateLimitBurst < 1 {
		return fmt.Errorf("RATE_LIMIT_BURST must be at least 1, got %d", c.RateLimitBurst)
	}

	if c.MaxLimit < 1 || c.MaxLimit > maxTagLimit {
		return fmt.Errorf("MAX_LIMIT must be between 1 and %d, got %d", maxTagLimit, c.MaxLimit)
	}

	if c.MaxOffset < 1 {
		return fmt.Errorf("MAX_OFFSET must be at least 1, got %d", c.MaxOffset)
	}

	return nil
}
