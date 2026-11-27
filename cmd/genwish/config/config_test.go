package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/cmd/genwish/models"
	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestAddConfigFlag(t *testing.T) {
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
			modelsFile := filepath.Join(dir, "models.toml")

			if err := os.WriteFile(
				modelsFile, []byte("[neo.sd]\ncheckpoint = \"sd.safetensors\"\n"), 0o600,
			); err != nil {
				t.Fatalf("WriteFile() error: %v", err)
			}

			configPath := filepath.Join(dir, "config.toml")
			body := "[neo]\nmodels_file = " + strconv.Quote(modelsFile) + "\n"

			if err := os.WriteFile(configPath, []byte(body), 0o600); err != nil {
				t.Fatalf("WriteFile() error: %v", err)
			}

			out := &bytes.Buffer{}
			errOut := &bytes.Buffer{}
			deps := cli.WithDefaults(cli.Deps{Out: out, Err: errOut})

			root := &cobra.Command{Use: "root"}
			root.SetOut(out)
			root.SetErr(errOut)

			AddConfigFlag(root, deps)
			root.AddCommand(models.NewCommand(deps))

			if err := testutil.Run(t, root, tt.args(configPath)...); err != nil {
				t.Fatalf("Execute() error: %v", err)
			}

			if !strings.Contains(out.String(), "sd.safetensors") {
				t.Errorf("stdout = %q, want the configured checkpoint", out)
			}
		})
	}

	t.Run("creates the default config", func(t *testing.T) {
		home := testutil.SetHome(t)

		out := &bytes.Buffer{}
		errOut := &bytes.Buffer{}
		deps := cli.WithDefaults(cli.Deps{Out: out, Err: errOut})

		root := &cobra.Command{Use: "root"}
		root.SetOut(out)
		root.SetErr(errOut)

		AddConfigFlag(root, deps)
		root.AddCommand(models.NewCommand(deps))

		if err := testutil.Run(t, root, "models", "list", "neo"); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		path := filepath.Join(home, ".genwish", "config.toml")

		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat() error: %v", err)
		}

		if info.Size() != 0 {
			t.Errorf("config size = %d, want an empty file", info.Size())
		}
	})
}
