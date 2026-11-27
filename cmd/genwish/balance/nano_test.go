package balance

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestRunNano(t *testing.T) {
	t.Run("prints the account credit", func(t *testing.T) {
		server, requests := testutil.CaptureServer(t, http.StatusOK,
			`{"usd_balance":"12.34","nano_balance":"5.6"}`)

		out := &bytes.Buffer{}
		deps := cli.Deps{Out: out, Nano: testutil.NewNanoClient(t, server.URL)}

		if err := RunNano(context.Background(), deps); err != nil {
			t.Fatalf("RunNano() error: %v", err)
		}

		request := testutil.NextRequest(t, requests)

		if request.Method != http.MethodPost || request.Path != "/api/check-balance" {
			t.Errorf("request = %s %s, want POST /api/check-balance", request.Method, request.Path)
		}
		if key := request.Header.Get("x-api-key"); key != "test-key" {
			t.Errorf("x-api-key = %q, want test-key", key)
		}

		for _, want := range []string{"USD: $12.34", "Nano: 5.6"} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}

		if strings.Contains(out.String(), "Deposit address") {
			t.Errorf("stdout contains a deposit address:\n%s", out)
		}
	})

	t.Run("reports an upstream failure", func(t *testing.T) {
		server, _ := testutil.CaptureServer(t, http.StatusUnauthorized, `{"error":"bad key"}`)

		deps := cli.Deps{Out: &bytes.Buffer{}, Nano: testutil.NewNanoClient(t, server.URL)}

		err := RunNano(context.Background(), deps)
		if err == nil || !strings.Contains(err.Error(), "balance nano") {
			t.Errorf("RunNano() error = %v, want the upstream failure", err)
		}
	})
}
