package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/config"
)

func TestOutputDir(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		requested string
		want      string
	}{
		{
			name:      "uses the requested directory",
			path:      "/from/config",
			requested: "/from/flag",
			want:      "/from/flag",
		},
		{
			name: "falls back to the configured path",
			path: "/from/config",
			want: "/from/config",
		},
		{
			name: "falls back to the working directory",
			want: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{Output: config.OutputConfig{Path: tt.path}}

			if got := OutputDir(cfg, tt.requested); got != tt.want {
				t.Errorf("OutputDir() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWriteImage(t *testing.T) {
	t.Run("writes under the kind and date", func(t *testing.T) {
		base := t.TempDir()

		path, err := WriteImage(base, "images", []byte("image bytes"), "image/png")
		if err != nil {
			t.Fatalf("WriteImage() error: %v", err)
		}

		if want := filepath.Join(base, "images"); filepath.Dir(filepath.Dir(path)) != want {
			t.Errorf("path = %s, want it under %s", path, want)
		}

		if ext := filepath.Ext(path); ext != ".png" {
			t.Errorf("extension = %s, want .png", ext)
		}

		name := strings.TrimSuffix(filepath.Base(path), ".png")
		if len(name) != 36 || name[14] != '4' {
			t.Errorf("name = %q, want a version 4 UUID", name)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile() error: %v", err)
		}

		if string(data) != "image bytes" {
			t.Errorf("data = %q, want the written bytes", data)
		}
	})
}
