package balance

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/novelai"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestNewCommand(t *testing.T) {
	t.Run("shows help without a provider", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewCommand(cli.Deps{Out: out})
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := testutil.Run(t, cmd); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(out.String(), "Providers: nano, novel") {
			t.Errorf("stdout = %q, want the provider options", out)
		}
	})

	t.Run("accepts the nano aliases", func(t *testing.T) {
		for _, alias := range []string{"nano-gpt", "nanogpt", "nano gpt"} {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte(`{"usd_balance":"1","nano_balance":"2"}`))
				}),
			)

			deps := cli.Deps{Out: &bytes.Buffer{}, Nano: testutil.NewNanoClient(t, server.URL)}
			err := testutil.Run(t, NewCommand(deps), alias)
			server.Close()

			if err != nil {
				t.Errorf("balance %s error: %v", alias, err)
			}
		}
	})

	t.Run("accepts the novel aliases", func(t *testing.T) {
		for _, alias := range []string{"novel-ai", "novel ai", "novelai"} {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, _ = w.Write([]byte(`{"subscription":{"trainingStepsLeft":7}}`))
				}),
			)

			out := &bytes.Buffer{}
			deps := cli.Deps{Out: out, Novel: novelai.New(server.URL, "key")}
			err := testutil.Run(t, NewCommand(deps), alias)
			server.Close()

			if err != nil {
				t.Errorf("balance %s error: %v", alias, err)
				continue
			}

			if !strings.Contains(out.String(), "Anlas: 7") {
				t.Errorf("balance %s stdout = %q, want the Anlas balance", alias, out)
			}
		}
	})

	t.Run("rejects neo", func(t *testing.T) {
		err := testutil.Run(t, NewCommand(cli.Deps{Out: &bytes.Buffer{}}), "neo")
		if err == nil || !strings.Contains(err.Error(), "no balance") {
			t.Errorf("Execute() error = %v, want a provider rejection", err)
		}
	})
}
