package models

import (
	"strings"
	"testing"
)

func TestModelsGetNeo(t *testing.T) {
	t.Run("prints the full entry", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.krea]
checkpoint = "hyphoria_v002.safetensors"
encoders = ["qwen_3_06b_base.safetensors", "qwen_3_4b_base.safetensors"]
preset = "krea"
sizes = ["1024x1024", "832x1216"]

[neo.krea.samplers]
default = "Euler"
options = ["Euler", "Euler a"]

[neo.krea.schedulers]
default = "Automatic"
options = ["Automatic", "Karras"]

[neo.krea.steps]
default = 20
min = 10
max = 40

[neo.krea.cfg]
default = 5.0
min = 1.0
max = 15.0

[neo.krea.hires]
enabled = true

[neo.krea.instructions]
usage = "Anime and illustration work."
style = "Soft shading with crisp line art."
`)

		out, err := runGet(t, neoDeps(path), "neo", "krea")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{
			"id: krea",
			"checkpoint: hyphoria_v002.safetensors",
			"encoders: qwen_3_06b_base.safetensors, qwen_3_4b_base.safetensors",
			"preset: krea",
			"samplers:\n  default: Euler\n  options: Euler, Euler a",
			"schedulers:\n  default: Automatic\n  options: Automatic, Karras",
			"steps:\n  default: 20\n  min: 10\n  max: 40",
			"cfg:\n  default: 5\n  min: 1\n  max: 15",
			"hires:\n  enabled: true",
			"sizes: 1024x1024, 832x1216",
			"instructions:\n  usage: Anime and illustration work.\n  style: Soft shading",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("renders multi-line instructions", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.notes.instructions]
usage = """
First line.
Second line.
"""
style = """
Anime: foo
Chibi: bar
"""
`)

		out, err := runGet(t, neoDeps(path), "neo", "notes")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		want := "instructions:\n" +
			"  usage: |\n    First line.\n    Second line.\n" +
			"  style: |\n    Anime: foo\n    Chibi: bar\n"

		if !strings.Contains(out, want) {
			t.Errorf("stdout does not contain the multi-line instructions:\n%s", out)
		}
	})

	t.Run("prints only the configured fields", func(t *testing.T) {
		path := writeModelsFile(t, "[neo.bare]\n")

		out, err := runGet(t, neoDeps(path), "neo", "bare")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if out != "id: bare\n" {
			t.Errorf("stdout = %q, want just the id", out)
		}
	})

	t.Run("reports an unconfigured model", func(t *testing.T) {
		out, err := runGet(t, neoDeps(writeModelsFile(t, "")), "neo", "absent")
		if err == nil || !strings.Contains(err.Error(), "is not configured for neo") {
			t.Errorf("Execute() error = %v, want the unconfigured model", err)
		}

		if out != "" {
			t.Errorf("stdout = %q, want nothing printed", out)
		}
	})

	t.Run("rejects video", func(t *testing.T) {
		_, err := runGet(t, neoDeps(writeModelsFile(t, "")), "neo", "krea", "video")
		if err == nil || !strings.Contains(err.Error(), "does not serve video") {
			t.Errorf("Execute() error = %v, want a video rejection", err)
		}
	})

	t.Run("reports an invalid models file", func(t *testing.T) {
		path := writeModelsFile(t, `
[neo.broken.samplers]
options = ["Euler"]
`)

		_, err := runGet(t, neoDeps(path), "neo", "broken")
		if err == nil || !strings.Contains(err.Error(), "samplers.default is required") {
			t.Errorf("Execute() error = %v, want the invalid sampler reported", err)
		}
	})
}
