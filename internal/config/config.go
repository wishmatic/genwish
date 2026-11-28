package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	defaultConfigDir     = ".genwish"
	defaultConfigFile    = "config.toml"
	defaultOutputPath    = "~/.genwish/out"
	defaultNeoURL        = "http://127.0.0.1:7860"
	defaultNeoModelsFile = "~/.genwish/models.toml"
)

type Config struct {
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

func Defaults() Config {
	return Config{
		Output: OutputConfig{
			Format: "webp",
			Path:   defaultOutputPath,
		},
		Neo: NeoConfig{
			URL:        defaultNeoURL,
			ModelsFile: defaultNeoModelsFile,
		},
		Booru: BooruConfig{
			RateLimit: RateLimitConfig{RPS: 1, Burst: 1, Timeout: 20},
		},
	}
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, defaultConfigDir, defaultConfigFile)
}

func EnsureDefaultFile() error {
	path := DefaultPath()
	if path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, fs.ErrExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("create config file %s: %w", path, err)
	}

	return file.Close()
}

func DefaultNeoModelsFile() string {
	path, err := expandHome(defaultNeoModelsFile)
	if err != nil {
		return ""
	}

	return path
}

func Load(path string) (Config, error) {
	cfg := Defaults()

	file := strings.TrimSpace(path)
	explicit := file != ""

	if !explicit {
		file = DefaultPath()
		if file == "" {
			return cfg.finish()
		}
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

func (c Config) Validate() error {
	if c.Booru.RateLimit.RPS <= 0 {
		return fmt.Errorf(
			"booru.rate_limit.rps must be greater than 0, got %v",
			c.Booru.RateLimit.RPS,
		)
	}

	if c.Booru.RateLimit.Burst < 1 {
		return fmt.Errorf(
			"booru.rate_limit.burst must be at least 1, got %d",
			c.Booru.RateLimit.Burst,
		)
	}

	if c.Booru.RateLimit.Timeout < 1 {
		return fmt.Errorf(
			"booru.rate_limit.timeout must be at least 1, got %d",
			c.Booru.RateLimit.Timeout,
		)
	}

	return nil
}

func expandHome(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") &&
		!strings.HasPrefix(path, fmt.Sprintf("~%s", string(os.PathSeparator))) {
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
