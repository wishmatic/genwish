package neomodels

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func writeModels(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "models.toml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	return path
}

func TestLoad(t *testing.T) {
	t.Run("returns the not-exist error", func(t *testing.T) {
		_, err := Load(filepath.Join(t.TempDir(), "models.toml"))
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Load() error = %v, want not exist", err)
		}
	})

	t.Run("reads the file", func(t *testing.T) {
		path := writeModels(t, "[neo.krea]\ncheckpoint = \"krea.safetensors\"\n")

		registry, err := Load(path)
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		model, ok := registry.Get("krea")
		if !ok || model.Checkpoint != "krea.safetensors" {
			t.Errorf("model = %+v, want the checkpoint", model)
		}
	})

	t.Run("reads the shipped example", func(t *testing.T) {
		registry, err := Load(filepath.Join("..", "..", "example", "models.toml"))
		if err != nil {
			t.Fatalf("Load(example) error: %v", err)
		}

		if !slices.Equal(registry.IDs(), []string{"illustrious", "krea"}) {
			t.Errorf("IDs() = %v, want both example models", registry.IDs())
		}

		model, ok := registry.Get("krea")
		if !ok || model.Checkpoint != "hyphoria_v002.safetensors" || !model.Hires.Enabled {
			t.Errorf("model = %+v, want the example entry", model)
		}

		if !strings.Contains(model.Instructions.Usage, "\n") {
			t.Errorf("usage = %q, want multi-line instructions", model.Instructions.Usage)
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("decodes every field", func(t *testing.T) {
		registry, err := Parse([]byte(`
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
usage = "Anime work."
style = "Soft shading."
`))
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}

		model, ok := registry.Get("krea")
		if !ok {
			t.Fatal("Parse() did not keep the krea model")
		}

		if model.Checkpoint != "hyphoria_v002.safetensors" || model.Preset != "krea" {
			t.Errorf("model = %+v, want the checkpoint and preset", model)
		}

		if !slices.Equal(model.Encoders, []string{
			"qwen_3_06b_base.safetensors", "qwen_3_4b_base.safetensors",
		}) {
			t.Errorf("encoders = %v, want both encoders", model.Encoders)
		}

		if model.Samplers.Default != "Euler" ||
			!slices.Equal(model.Samplers.Options, []string{"Euler", "Euler a"}) {
			t.Errorf("samplers = %+v, want Euler and its options", model.Samplers)
		}

		if model.Schedulers.Default != "Automatic" ||
			!slices.Equal(model.Schedulers.Options, []string{"Automatic", "Karras"}) {
			t.Errorf("schedulers = %+v, want Automatic and its options", model.Schedulers)
		}

		if model.Steps.Default == nil || *model.Steps.Default != 20 ||
			model.Steps.Min == nil || *model.Steps.Min != 10 ||
			model.Steps.Max == nil || *model.Steps.Max != 40 {
			t.Errorf("steps = %+v, want 20 within 10 and 40", model.Steps)
		}

		if model.CFG.Default == nil || *model.CFG.Default != 5 ||
			model.CFG.Min == nil || *model.CFG.Min != 1 ||
			model.CFG.Max == nil || *model.CFG.Max != 15 {
			t.Errorf("cfg = %+v, want 5 within 1 and 15", model.CFG)
		}

		if !model.Hires.Enabled {
			t.Error("hires = false, want it enabled")
		}

		if !slices.Equal(model.Sizes, []string{"1024x1024", "832x1216"}) {
			t.Errorf("sizes = %v, want both sizes", model.Sizes)
		}

		if model.Instructions.Usage != "Anime work." ||
			model.Instructions.Style != "Soft shading." {
			t.Errorf("instructions = %+v, want both strings", model.Instructions)
		}
	})

	t.Run("accepts an empty file", func(t *testing.T) {
		registry, err := Parse(nil)
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}

		if len(registry) != 0 {
			t.Errorf("registry = %v, want nothing", registry)
		}
	})

	t.Run("accepts a model with no fields", func(t *testing.T) {
		registry, err := Parse([]byte("[neo.bare]\n"))
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}

		if _, ok := registry.Get("bare"); !ok {
			t.Error("Parse() dropped the empty model")
		}
	})

	t.Run("rejects an unknown key", func(t *testing.T) {
		_, err := Parse([]byte("[neo.krea]\ncheckpont = \"typo\"\n"))
		if err == nil || !containsAll(err.Error(), "unknown key", "checkpont") {
			t.Errorf("Parse() error = %v, want the unknown key", err)
		}
	})

	t.Run("rejects malformed TOML", func(t *testing.T) {
		_, err := Parse([]byte("[neo.krea\n"))
		if err == nil {
			t.Error("Parse() error = nil, want a parse failure")
		}
	})

	t.Run("names the invalid model", func(t *testing.T) {
		_, err := Parse([]byte("[neo.broken]\nsizes = [\"wide\"]\n"))
		if err == nil || !containsAll(err.Error(), "neo.broken", "WxH") {
			t.Errorf("Parse() error = %v, want the model and size named", err)
		}
	})
}

func TestRegistryIDs(t *testing.T) {
	t.Run("sorts the ids", func(t *testing.T) {
		registry, err := Parse([]byte("[neo.zeta]\n\n[neo.alpha]\n\n[neo.mid]\n"))
		if err != nil {
			t.Fatalf("Parse() error: %v", err)
		}

		if !slices.Equal(registry.IDs(), []string{"alpha", "mid", "zeta"}) {
			t.Errorf("IDs() = %v, want sorted ids", registry.IDs())
		}
	})
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}

	return true
}
