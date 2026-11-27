package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestBalanceNano(t *testing.T) {
	t.Run("prints the account credit", func(t *testing.T) {
		server, requests := captureServer(t, http.StatusOK,
			`{"usd_balance":"12.34","nano_balance":"5.6"}`)

		out, _, err := runCommand(t, cli.Deps{Nano: newClientAt(t, server.URL)}, "balance", "nano")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		request := nextRequest(t, requests)

		if request.Method != http.MethodPost || request.Path != "/api/check-balance" {
			t.Errorf("request = %s %s, want POST /api/check-balance", request.Method, request.Path)
		}
		if key := request.Header.Get("x-api-key"); key != "test-key" {
			t.Errorf("x-api-key = %q, want test-key", key)
		}

		for _, want := range []string{"USD: $12.34", "Nano: 5.6"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}

		if strings.Contains(out, "Deposit address") {
			t.Errorf("stdout contains a deposit address:\n%s", out)
		}
	})

	t.Run("reports an upstream failure", func(t *testing.T) {
		server, _ := captureServer(t, http.StatusUnauthorized, `{"error":"bad key"}`)

		_, _, err := runCommand(t, cli.Deps{Nano: newClientAt(t, server.URL)}, "balance", "nano")
		if err == nil || !strings.Contains(err.Error(), "balance nano") {
			t.Errorf("Execute() error = %v, want the upstream failure", err)
		}
	})

	t.Run("accepts the provider aliases", func(t *testing.T) {
		for _, alias := range []string{"nano-gpt", "nanogpt", "nano gpt"} {
			server, _ := captureServer(t, http.StatusOK, `{"usd_balance":"1","nano_balance":"2"}`)

			if _, _, err := runCommand(t, cli.Deps{Nano: newClientAt(t, server.URL)}, "balance", alias); err != nil {
				t.Errorf("balance %s error: %v", alias, err)
			}
		}
	})

	t.Run("rejects neo", func(t *testing.T) {
		_, _, err := runCommand(t, cli.Deps{}, "balance", "neo")
		if err == nil || !strings.Contains(err.Error(), "no balance") {
			t.Errorf("Execute() error = %v, want a provider rejection", err)
		}
	})
}
