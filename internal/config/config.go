package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	defaultConfigFile = "config.toml"
	defaultOutputPath = "~/.genwish/out"
	defaultNeoURL     = "http://127.0.0.1:7860"
)

type Config struct {
	LogLevel    string `toml:"-"`
	ErrorDetail string `toml:"-"`

	Output OutputConfig
	Neo    NeoConfig
	Booru  BooruConfig
}

type OutputConfig struct {
	Format string `toml:"format"`
	Path   string `toml:"path"`
}

type NeoConfig struct {
	URL        string `toml:"url"`
	ModelsFile string `toml:"models_file"`
}

type BooruConfig struct {
	RateLimit RateLimitConfig `toml:"rate_limit"`
}

type RateLimitConfig struct {
	RPS     float64 `toml:"rps"`
	Burst   int     `toml:"burst"`
	Timeout int     `toml:"timeout"`
}

type envConfig struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	ErrorDetail string `env:"ERROR_DETAIL" envDefault:"useful"`
}

func Defaults() Config {
	return Config{
		LogLevel:    "info",
		ErrorDetail: "useful",
		Output: OutputConfig{
			Format: "webp",
			Path:   defaultOutputPath,
		},
		Neo: NeoConfig{URL: defaultNeoURL},
		Booru: BooruConfig{
			RateLimit: RateLimitConfig{RPS: 1, Burst: 1, Timeout: 20},
		},
	}
}

func Load(path string) (Config, error) {
	_ = godotenv.Load()

	var fromEnv envConfig
	if err := env.Parse(&fromEnv); err != nil {
		return Config{}, fmt.Errorf("parse environment: %w", err)
	}

	cfg := Defaults()
	cfg.LogLevel = fromEnv.LogLevel
	cfg.ErrorDetail = fromEnv.ErrorDetail

	file := strings.TrimSpace(path)
	explicit := file != ""

	if !explicit {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg.finish()
		}

		file = filepath.Join(home, ".genwish", defaultConfigFile)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		if !explicit && errors.Is(err, os.ErrNotExist) {
			return cfg.finish()
		}

		return Config{}, fmt.Errorf("config file %s: %w", file, err)
	}

	meta, err := toml.Decode(string(data), &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("config file %s: %w", file, err)
	}

	if unknown := meta.Undecoded(); len(unknown) > 0 {
		return Config{}, fmt.Errorf("config file %s: unknown key %s", file, formatKeys(unknown))
	}

	return cfg.finish()
}

func (c Config) finish() (Config, error) {
	if err := c.normalize(); err != nil {
		return Config{}, err
	}

	if err := c.Validate(); err != nil {
		return Config{}, err
	}

	return c, nil
}

func (c *Config) normalize() error {
	if strings.TrimSpace(c.Output.Format) == "" {
		c.Output.Format = "webp"
	}

	if strings.TrimSpace(c.Output.Path) == "" {
		c.Output.Path = defaultOutputPath
	}

	outputPath, err := expandHome(c.Output.Path)
	if err != nil {
		return err
	}
	c.Output.Path = outputPath

	if strings.TrimSpace(c.Neo.URL) == "" {
		c.Neo.URL = defaultNeoURL
	}

	if strings.TrimSpace(c.Neo.ModelsFile) != "" {
		modelsFile, err := expandHome(c.Neo.ModelsFile)
		if err != nil {
			return err
		}
		c.Neo.ModelsFile = modelsFile
	}

	return nil
}

func (c Config) VerboseErrors() bool {
	return c.ErrorDetail == "verbose"
}

func (c Config) Validate() error {
	if c.Booru.RateLimit.RPS <= 0 {
		return fmt.Errorf("booru.rate_limit.rps must be greater than 0, got %v", c.Booru.RateLimit.RPS)
	}

	if c.Booru.RateLimit.Burst < 1 {
		return fmt.Errorf("booru.rate_limit.burst must be at least 1, got %d", c.Booru.RateLimit.Burst)
	}

	if c.Booru.RateLimit.Timeout < 1 {
		return fmt.Errorf("booru.rate_limit.timeout must be at least 1, got %d", c.Booru.RateLimit.Timeout)
	}

	return nil
}

func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("expand %q: %w", path, err)
	}

	return filepath.Join(home, path[1:]), nil
}

func formatKeys(keys []toml.Key) string {
	names := make([]string, 0, len(keys))

	for _, key := range keys {
		names = append(names, fmt.Sprintf("%q", key.String()))
	}

	return strings.Join(names, ", ")
}
