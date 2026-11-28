package genimage

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/config"
	"github.com/wishmatic/genwish/internal/forge"
	"github.com/wishmatic/genwish/internal/testutil"
)

const neoModelsFixture = `
[neo.krea]
checkpoint = "hyphoria_v002.safetensors"
encoders = ["qwen_3_06b_base.safetensors", "qwen_3_4b_base.safetensors"]
preset = "krea"

[neo.krea.samplers]
default = "Euler"
options = ["Euler", "Euler a"]

[neo.krea.schedulers]
default = "Karras"
options = ["Automatic", "Karras"]

[neo.krea.steps]
default = 30
min = 10
max = 40

[neo.krea.cfg]
default = 5.0

[neo.krea.hires]
enabled = true
`

func TestGenImageNeoModelDefaults(t *testing.T) {
	t.Run("applies the configured model defaults", func(t *testing.T) {
		deps, requests := neoDepsWithModels(t)

		payload := runNeoRequest(t, deps, requests, "--model", "krea").JSONBody(t)

		checks := map[string]any{
			"sampler_name": "Euler",
			"scheduler":    "Karras",
			"steps":        float64(30),
			"cfg_scale":    float64(5),
			"enable_hr":    true,
		}

		for key, want := range checks {
			if payload[key] != want {
				t.Errorf("%s = %v, want %v", key, payload[key], want)
			}
		}

		override, ok := payload["override_settings"].(map[string]any)
		if !ok {
			t.Fatalf("override_settings = %v, want the model override", payload["override_settings"])
		}

		if override["sd_model_checkpoint"] != "hyphoria_v002.safetensors" ||
			override["forge_preset"] != "krea" {
			t.Errorf("override_settings = %v, want the checkpoint and preset", override)
		}

		modules, _ := override["forge_additional_modules"].([]any)
		if !slices.Equal(
			modules,
			[]any{"qwen_3_06b_base.safetensors", "qwen_3_4b_base.safetensors"},
		) {
			t.Errorf("forge_additional_modules = %v, want both encoders", override)
		}
	})

	t.Run("lets flags override the configured defaults", func(t *testing.T) {
		deps, requests := neoDepsWithModels(t)

		payload := runNeoRequest(t, deps, requests, "--model", "krea", "--sampler-name", "DPM++ 2M").
			JSONBody(t)

		if payload["sampler_name"] != "DPM++ 2M" {
			t.Errorf("sampler_name = %v, want the flag value", payload["sampler_name"])
		}

		if payload["scheduler"] != "Karras" {
			t.Errorf("scheduler = %v, want the configured default kept", payload["scheduler"])
		}
	})

	t.Run("matches a checkpoint filename", func(t *testing.T) {
		deps, requests := neoDepsWithModels(t)

		payload := runNeoRequest(t, deps, requests, "--model", "hyphoria_v002.safetensors").
			JSONBody(t)

		if payload["steps"] != float64(30) {
			t.Errorf("steps = %v, want the configured default", payload["steps"])
		}
	})

	t.Run("falls back to the checkpoint name", func(t *testing.T) {
		server, requests := neoForgeServer(t)
		deps := cli.Deps{Forge: forge.New(server.URL, io.Discard)}

		payload := runNeoRequest(
			t,
			deps,
			requests,
			"--model",
			"custom_v1.safetensors",
			"--forge-preset",
			"krea",
		).JSONBody(t)

		override, ok := payload["override_settings"].(map[string]any)
		if !ok {
			t.Fatalf("override_settings = %v, want the checkpoint override", payload["override_settings"])
		}

		if override["sd_model_checkpoint"] != "custom_v1.safetensors" {
			t.Errorf("sd_model_checkpoint = %v, want the given checkpoint", override["sd_model_checkpoint"])
		}
	})
}

func neoDepsWithModels(t *testing.T) (cli.Deps, <-chan testutil.Request) {
	t.Helper()

	modelsFile := filepath.Join(t.TempDir(), "models.toml")
	if err := os.WriteFile(modelsFile, []byte(neoModelsFixture), 0o600); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	server, requests := neoForgeServer(t)

	return cli.Deps{
		Forge:  forge.New(server.URL, io.Discard),
		Config: config.Config{Neo: config.NeoConfig{ModelsFile: modelsFile}},
	}, requests
}

func neoForgeServer(t *testing.T) (*httptest.Server, <-chan testutil.Request) {
	t.Helper()

	png := base64.StdEncoding.EncodeToString(testutil.ValidPNG(t, 4, 4))

	return testutil.CaptureServer(t, http.StatusOK, fmt.Sprintf(`{"images":[%q]}`, png))
}

func runNeoRequest(
	t *testing.T, deps cli.Deps, requests <-chan testutil.Request, args ...string,
) testutil.Request {
	t.Helper()

	full := append([]string{"a cat", "--format", "png", "--out", t.TempDir()}, args...)

	if _, _, err := runNeo(t, deps, full...); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	return testutil.NextRequest(t, requests)
}
