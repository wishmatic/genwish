package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/config"
)

func TestModelsListNeo(t *testing.T) {
	t.Run("reads the configured file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "neo-models.txt")
		body := "# checkpoints\nsd.safetensors\n\nxl.safetensors\n"

		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		out, err := runList(t, cli.Deps{Config: config.Config{Neo: config.NeoConfig{ModelsFile: path}}}, "neo")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"NEO MODELS", "sd.safetensors", "xl.safetensors", "2 of 2 models shown"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}

		if strings.Contains(out, "# checkpoints") {
			t.Errorf("stdout contains the comment line:\n%s", out)
		}
	})

	t.Run("defaults to empty when unset", func(t *testing.T) {
		out, err := runList(t, cli.Deps{}, "neo")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "0 of 0 models shown") {
			t.Errorf("stdout = %q, want an empty list", out)
		}
	})

	t.Run("reports a missing configured file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.txt")

		_, err := runList(t, cli.Deps{Config: config.Config{Neo: config.NeoConfig{ModelsFile: path}}}, "neo")
		if err == nil || !strings.Contains(err.Error(), "neo.models_file") {
			t.Errorf("Execute() error = %v, want the missing file reported", err)
		}
	})

	t.Run("rejects video", func(t *testing.T) {
		_, err := runList(t, cli.Deps{}, "neo", "video")
		if err == nil || !strings.Contains(err.Error(), "does not serve video") {
			t.Errorf("Execute() error = %v, want a video rejection", err)
		}
	})
}
