package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/config"
)

func writeModelsFile(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "models.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	return path
}

func neoDeps(path string) cli.Deps {
	return cli.Deps{Config: config.Config{Neo: config.NeoConfig{ModelsFile: path}}}
}

func TestModelsListNeo(t *testing.T) {
	t.Run("lists the configured models", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.sd]
checkpoint = "sd.safetensors"
preset = "default"

[neo.xl]
checkpoint = "xl.safetensors"
`)

		out, err := runList(t, neoDeps(path), "neo")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{
			"NEO MODELS", "ID", "CHECKPOINT", "PRESET",
			"sd", "sd.safetensors", "default",
			"xl", "xl.safetensors", "2 of 2 models shown",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("filters by id, checkpoint, and preset", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.sd]
checkpoint = "sd.safetensors"
preset = "default"

[neo.xl]
checkpoint = "xl.safetensors"
preset = "krea"
`)

		tests := map[string]string{
			"id":         "xl",
			"checkpoint": "sd.safet",
			"preset":     "krea",
		}

		for name, query := range tests {
			out, err := runList(t, neoDeps(path), "neo", "--query", query)
			if err != nil {
				t.Fatalf("Execute() error: %v", err)
			}

			if !strings.Contains(out, "1 of 2 models shown (1 matched)") {
				t.Errorf("%s query stdout = %q, want one match", name, out)
			}
		}
	})

	t.Run("defaults to empty when the default file is missing", func(t *testing.T) {
		out, err := runList(t, cli.Deps{}, "neo")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "0 of 0 models shown") {
			t.Errorf("stdout = %q, want an empty list", out)
		}
	})

	t.Run("reports a missing configured file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing.toml")

		_, err := runList(t, neoDeps(path), "neo")
		if err == nil || !strings.Contains(err.Error(), "neo.models_file") {
			t.Errorf("Execute() error = %v, want the missing file reported", err)
		}
	})

	t.Run("reports an invalid models file", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.broken]
sizes = ["wide"]
`)

		_, err := runList(t, neoDeps(path), "neo")
		if err == nil || !strings.Contains(err.Error(), `size "wide" is not WxH`) {
			t.Errorf("Execute() error = %v, want the invalid size reported", err)
		}
	})

	t.Run("rejects video", func(t *testing.T) {
		_, err := runList(t, cli.Deps{}, "neo", "video")
		if err == nil || !strings.Contains(err.Error(), "does not serve video") {
			t.Errorf("Execute() error = %v, want a video rejection", err)
		}
	})
}
