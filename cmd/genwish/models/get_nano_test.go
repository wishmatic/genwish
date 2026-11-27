package models

import (
	"net/http"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

const nanoDetailBody = `{"data": [{
	"id": "nano-banana",
	"name": "Nano Banana",
	"owned_by": "google",
	"description": "Fast image editing.",
	"tags": ["editing"],
	"capabilities": {"chat": true},
	"supported_parameters": {"size": true},
	"pricing": {"image": 0.02}
}]}`

const videoDetailBody = `{"data": [{
	"id": "veo2-video",
	"name": "Veo 2",
	"owned_by": "google",
	"capabilities": {"audio": true}
}]}`

func nanoDetailServer(t *testing.T) cli.Deps {
	t.Helper()

	server, _ := testutil.CaptureServerFunc(t, func(r *http.Request) (int, string) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/image-models") {
			return http.StatusOK, nanoDetailBody
		}

		return http.StatusOK, videoDetailBody
	})

	return cli.Deps{Nano: testutil.NewNanoClient(t, server.URL)}
}

func TestModelsGetNano(t *testing.T) {
	t.Run("finds a model in either kind", func(t *testing.T) {
		out, err := runGet(t, nanoDetailServer(t), "nano", "nano-banana")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{
			"id: nano-banana",
			"name: Nano Banana",
			"owner: google",
			"description: Fast image editing.",
			"tags: editing",
			`capabilities: {"chat":true}`,
			`supported_parameters: {"size":true}`,
			`pricing: {"image":0.02}`,
		} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("uses the requested kind", func(t *testing.T) {
		out, err := runGet(t, nanoDetailServer(t), "nano", "veo2-video", "video")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out, "id: veo2-video") || !strings.Contains(out, `"audio":true`) {
			t.Errorf("stdout = %q, want the video model", out)
		}
	})

	t.Run("reports an unserved model", func(t *testing.T) {
		_, err := runGet(t, nanoDetailServer(t), "nano", "absent")
		if err == nil || !strings.Contains(err.Error(), "is not served by nano") {
			t.Errorf("Execute() error = %v, want the unserved model", err)
		}
	})

	t.Run("rejects an invalid kind", func(t *testing.T) {
		_, err := runGet(t, nanoDetailServer(t), "nano", "nano-banana", "audio")
		if err == nil || !strings.Contains(err.Error(), "not one of") {
			t.Errorf("Execute() error = %v, want the invalid kind", err)
		}
	})

	t.Run("reports an upstream failure", func(t *testing.T) {
		server, _ := testutil.CaptureServer(t, http.StatusUnauthorized, `{"error":"bad key"}`)

		_, err := runGet(
			t,
			cli.Deps{Nano: testutil.NewNanoClient(t, server.URL)},
			"nano", "nano-banana", "image",
		)
		if err == nil || !strings.Contains(err.Error(), "models get image") {
			t.Errorf("Execute() error = %v, want the upstream failure", err)
		}
	})
}
