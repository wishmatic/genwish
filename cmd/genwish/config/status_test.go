package config

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/provider"
	"github.com/wishmatic/genwish/internal/testutil"
)

func runStatusWith(
	t *testing.T, deps cli.Deps, stores map[provider.Name]credential.Store,
) (string, error) {
	t.Helper()

	out := &bytes.Buffer{}
	deps.Out = out
	if deps.Err == nil {
		deps.Err = &bytes.Buffer{}
	}

	open := func(name provider.Name) (credential.Store, error) {
		if store, ok := stores[name]; ok {
			return store, nil
		}

		return &memoryStore{}, nil
	}

	err := runStatus(deps, open)

	return out.String(), err
}

func TestStatusShowsEachProvider(t *testing.T) {
	stores := map[provider.Name]credential.Store{
		provider.Nano:     &memoryStore{key: "nano-key"},
		provider.Novel:    &memoryStore{},
		provider.Danbooru: &memoryStore{key: "someone\nkey"},
	}

	out, err := runStatusWith(t, cli.Deps{}, stores)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 4 {
		t.Fatalf("stdout = %q, want a header and three providers", out)
	}

	if got := strings.Join(strings.Fields(lines[0]), " "); got != "PROVIDER STATUS" {
		t.Errorf("header = %q, want PROVIDER STATUS", got)
	}

	want := map[string]string{
		"nano":     "logged in",
		"novel":    "not logged in",
		"danbooru": "logged in as someone",
	}

	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			t.Errorf("line = %q, want a provider and a status", line)

			continue
		}

		status := strings.Join(fields[1:], " ")
		if want[fields[0]] != status {
			t.Errorf("status for %s = %q, want %q", fields[0], status, want[fields[0]])
		}
	}
}

func TestStatusReportsAStoreFailure(t *testing.T) {
	stores := map[provider.Name]credential.Store{
		provider.Nano: &memoryStore{loadErr: errors.New("keyring exploded")},
	}

	_, err := runStatusWith(t, cli.Deps{}, stores)
	if err == nil || !strings.Contains(err.Error(), "keyring exploded") {
		t.Errorf("Execute() error = %v, want the store failure", err)
	}
}

func TestStatusUsesTheInjectedStore(t *testing.T) {
	out := &bytes.Buffer{}

	err := testutil.Run(
		t,
		NewStatusCommand(cli.Deps{
			Out:         out,
			Err:         &bytes.Buffer{},
			Credentials: &memoryStore{key: "stored-key"},
		}),
	)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(out.String(), "nano") ||
		!strings.Contains(out.String(), "logged in") {
		t.Errorf("stdout = %q, want the stored login listed", out)
	}

	if !strings.Contains(out.String(), "not logged in") {
		t.Errorf("stdout = %q, want the danbooru row without a login", out)
	}
}
