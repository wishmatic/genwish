package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTxt2ImgOverrideSettings(t *testing.T) {
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"images": ["cG5n"]}`))
	}))
	defer server.Close()

	c := New(server.URL, io.Discard)

	_, err := c.txt2img(
		context.Background(),
		txt2imgRequest{
			Checkpoint:             "model",
			ForgePreset:            "flux",
			ForgeAdditionalModules: []string{"vae", "t5xxl.ckpt"},
			Prompt:                 "test",
		},
	)
	if err != nil {
		t.Fatalf("txt2img() error: %v", err)
	}

	var payload struct {
		OverrideSettings struct {
			Checkpoint             string   `json:"sd_model_checkpoint"`
			ForgePreset            string   `json:"forge_preset"`
			ForgeAdditionalModules []string `json:"forge_additional_modules"`
		} `json:"override_settings"`
	}
	if err := json.Unmarshal(gotBody, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}

	if payload.OverrideSettings.Checkpoint != "model.safetensors" {
		t.Errorf(
			"sd_model_checkpoint = %q, want model.safetensors",
			payload.OverrideSettings.Checkpoint,
		)
	}
	if payload.OverrideSettings.ForgePreset != "flux" {
		t.Errorf("forge_preset = %q, want flux", payload.OverrideSettings.ForgePreset)
	}
	if len(payload.OverrideSettings.ForgeAdditionalModules) != 2 {
		t.Fatalf(
			"forge_additional_modules = %v, want 2 entries",
			payload.OverrideSettings.ForgeAdditionalModules,
		)
	}
	if payload.OverrideSettings.ForgeAdditionalModules[0] != "vae.safetensors" {
		t.Errorf(
			"forge_additional_modules[0] = %q, want vae.safetensors",
			payload.OverrideSettings.ForgeAdditionalModules[0],
		)
	}
	if payload.OverrideSettings.ForgeAdditionalModules[1] != "t5xxl.ckpt" {
		t.Errorf(
			"forge_additional_modules[1] = %q, want t5xxl.ckpt (already has extension)",
			payload.OverrideSettings.ForgeAdditionalModules[1],
		)
	}
}

func TestTxt2ImgNoOverrideSettings(t *testing.T) {
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"images": ["cG5n"]}`))
	}))
	defer server.Close()

	c := New(server.URL, io.Discard)

	_, err := c.txt2img(context.Background(), txt2imgRequest{Prompt: "test"})
	if err != nil {
		t.Fatalf("txt2img() error: %v", err)
	}

	if strings.Contains(string(gotBody), "override_settings") {
		t.Errorf("payload should not contain override_settings when empty, got: %s", gotBody)
	}
}

func TestTxt2ImgRejectsEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"images": []}`))
	}))
	defer server.Close()

	c := New(server.URL, io.Discard)

	_, err := c.txt2img(context.Background(), txt2imgRequest{Prompt: "test"})
	if err == nil || !strings.Contains(err.Error(), "no images") {
		t.Fatalf("txt2img() error = %v, want an error about the missing images", err)
	}
}

func TestTxt2ImgWarnsAboutExtraImages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"images": ["cG5n", "Z2lm"]}`))
	}))
	defer server.Close()

	var warnings bytes.Buffer

	images, err := New(
		server.URL,
		&warnings,
	).txt2img(context.Background(), txt2imgRequest{Prompt: "test"})
	if err != nil {
		t.Fatalf("txt2img() error: %v", err)
	}

	if len(images) != 1 || string(images[0]) != "png" {
		t.Errorf("images = %q, want only the first image", images)
	}

	if want := "warning: txt2img returned 2 images"; !strings.Contains(warnings.String(), want) {
		t.Errorf("warnings = %q, want %q", warnings.String(), want)
	}
}
