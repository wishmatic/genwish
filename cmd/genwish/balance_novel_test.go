package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/novelai"
)

func TestBalanceNovel(t *testing.T) {
	t.Run("prints the Anlas balance", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(
				`{"subscription":{"tier":"opus","trainingStepsLeft":{"fixedTrainingStepsLeft":9898,"purchasedTrainingSteps":32}}}`,
			))
		}))
		t.Cleanup(server.Close)

		out, _, err := runCommand(t, cli.Deps{Novel: novelai.New(server.URL, "key")}, "balance", "novel")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{"Anlas: 9930", "Tier: opus"} {
			if !strings.Contains(out, want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("accepts the provider aliases", func(t *testing.T) {
		for _, alias := range []string{"novel-ai", "novel ai", "novelai"} {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte(`{"subscription":{"trainingStepsLeft":7}}`))
			}))

			out, _, err := runCommand(t, cli.Deps{Novel: novelai.New(server.URL, "key")}, "balance", alias)
			server.Close()

			if err != nil {
				t.Errorf("balance %s error: %v", alias, err)
				continue
			}

			if !strings.Contains(out, "Anlas: 7") {
				t.Errorf("balance %s stdout = %q, want the Anlas balance", alias, out)
			}
		}
	})

	t.Run("reports a missing key", func(t *testing.T) {
		t.Setenv("NOVELAI_API_KEY", "")

		_, _, err := runCommand(t, cli.Deps{Credentials: &memoryStore{}}, "balance", "novel")
		if err == nil || !strings.Contains(err.Error(), "genwish login novel") {
			t.Errorf("Execute() error = %v, want a login hint", err)
		}
	})
}
