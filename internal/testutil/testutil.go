package testutil

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/nanogpt"
)

var PNG = []byte("\x89PNG\r\n\x1a\n test image")

type Request struct {
	Method string
	Path   string
	Query  string
	Header http.Header
	Body   []byte
}

func CaptureServer(t *testing.T, status int, response string) (*httptest.Server, <-chan Request) {
	t.Helper()

	return CaptureServerFunc(t, func(*http.Request) (int, string) {
		return status, response
	})
}

func CaptureServerFunc(
	t *testing.T, respond func(*http.Request) (int, string),
) (*httptest.Server, <-chan Request) {
	t.Helper()

	requests := make(chan Request, 16)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			http.Error(w, "read body", http.StatusInternalServerError)

			return
		}

		requests <- Request{
			Method: r.Method,
			Path:   r.URL.Path,
			Query:  r.URL.RawQuery,
			Header: r.Header.Clone(),
			Body:   body,
		}

		status, response := respond(r)

		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	return server, requests
}

func NextRequest(t *testing.T, requests <-chan Request) Request {
	t.Helper()

	select {
	case request := <-requests:
		return request
	default:
		t.Fatal("upstream received no request")

		return Request{}
	}
}

func NextRequestPath(t *testing.T, requests <-chan Request, path string) Request {
	t.Helper()

	for range 16 {
		request := NextRequest(t, requests)
		if request.Path == path {
			return request
		}
	}

	t.Fatalf("upstream received no request for %s", path)

	return Request{}
}

func (r Request) JSONBody(t *testing.T) map[string]any {
	t.Helper()

	var body map[string]any
	if err := json.Unmarshal(r.Body, &body); err != nil {
		t.Fatalf("decode request body %q: %v", r.Body, err)
	}

	return body
}

func ValidPNG(t *testing.T, width, height int) []byte {
	t.Helper()

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, width, height))); err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	return buffer.Bytes()
}

func NewNanoClient(t *testing.T, baseURL string) *nanogpt.Client {
	t.Helper()

	client, err := nanogpt.New(nanogpt.Config{APIKey: "test-key", BaseURL: baseURL})
	if err != nil {
		t.Fatalf("nanogpt.New() error: %v", err)
	}

	return client
}

func SetHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	return home
}

func Run(t *testing.T, cmd *cobra.Command, args ...string) error {
	t.Helper()

	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd.Execute()
}
