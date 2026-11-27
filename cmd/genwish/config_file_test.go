package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestConfigFileFlag(t *testing.T) {
	tests := []struct {
		name string
		args func(configPath string) []string
	}{
		{
			name: "before the subcommand",
			args: func(path string) []string {
				return []string{"--config-file", path, "models", "list", "neo"}
			},
		},
		{
			name: "after the subcommand",
			args: func(path string) []string {
				return []string{"models", "list", "neo", "--config-file", path}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			models := filepath.Join(dir, "neo-models.txt")

			if err := os.WriteFile(models, []byte("sd.safetensors\n"), 0o600); err != nil {
				t.Fatalf("WriteFile() error: %v", err)
			}

			configPath := filepath.Join(dir, "config.toml")
			body := "[neo]\nmodels_file = " + strconv.Quote(models) + "\n"

			if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
				t.Fatalf("WriteFile() error: %v", err)
			}

			out, _, err := runCommand(t, cli.Deps{}, tt.args(configPath)...)
			if err != nil {
				t.Fatalf("Execute() error: %v", err)
			}

			if !strings.Contains(out, "sd.safetensors") {
				t.Errorf("stdout = %q, want the configured checkpoint", out)
			}
		})
	}
}
