package config

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/testutil"
)

func TestLoginArgs(t *testing.T) {
	t.Run("shows help without a provider", func(t *testing.T) {
		out := &bytes.Buffer{}
		cmd := NewLoginCommand(cli.Deps{})
		cmd.SetOut(out)
		cmd.SetErr(&bytes.Buffer{})

		if err := testutil.Run(t, cmd); err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		for _, want := range []string{
			"login <provider>", "Reads credentials from STDIN",
		} {
			if !strings.Contains(out.String(), want) {
				t.Errorf("stdout does not contain %q:\n%s", want, out)
			}
		}
	})

	t.Run("rejects more than one provider", func(t *testing.T) {
		_, _, err := runLogin(t, cli.Deps{}, "nano", "novel")
		if err == nil || !strings.Contains(err.Error(), "accepts at most 1") {
			t.Errorf("Execute() error = %v, want the extra argument rejected", err)
		}
	})
}

func TestLoginRejectsNeo(t *testing.T) {
	t.Run("reports the provider does not use a key", func(t *testing.T) {
		_, _, err := runLogin(t, cli.Deps{}, "neo")
		if err == nil || !strings.Contains(err.Error(), "does not use an API key") {
			t.Errorf("Execute() error = %v, want neo rejected", err)
		}
	})
}

func TestLoginStoresTrimmedStdin(t *testing.T) {
	store := &memoryStore{}

	out, errOut, err := runLogin(t, cli.Deps{
		In:          strings.NewReader("  secret-key\n"),
		Credentials: store,
	}, "nano")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if store.key != "secret-key" {
		t.Errorf("stored key = %q, want secret-key", store.key)
	}

	if !strings.Contains(out, "logged in") {
		t.Errorf("stdout = %q, want a login confirmation", out)
	}

	if strings.Contains(out+errOut, "secret-key") {
		t.Errorf("output leaked the API key: stdout %q, stderr %q", out, errOut)
	}

	if strings.Contains(errOut, "already logged in") {
		t.Errorf("stderr = %q, want no replacement warning", errOut)
	}
}

func TestLoginWarnsWhenAlreadyLoggedIn(t *testing.T) {
	store := &memoryStore{key: "old-key"}

	out, errOut, err := runLogin(t, cli.Deps{
		In:          strings.NewReader("new-key\n"),
		Credentials: store,
	}, "nano")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !strings.Contains(errOut, "nano is already logged in") {
		t.Errorf("stderr = %q, want the replacement warning", errOut)
	}

	if store.key != "new-key" {
		t.Errorf("stored key = %q, want the replacement", store.key)
	}

	if !strings.Contains(out, "logged in") {
		t.Errorf("stdout = %q, want a login confirmation", out)
	}
}

func TestLoginRejectsMultipleLinesOfPipedInput(t *testing.T) {
	store := &memoryStore{}

	_, _, err := runLogin(t, cli.Deps{
		In:          strings.NewReader("key-from-pipe\nignored\n"),
		Credentials: store,
	}, "nano")
	if err == nil || !strings.Contains(err.Error(), "single line") {
		t.Fatalf("Execute() error = %v, want a single-line rejection", err)
	}

	if store.key != "" {
		t.Errorf("stored key = %q, want nothing", store.key)
	}

	if _, _, err := runLogin(t, cli.Deps{
		In:          strings.NewReader("  one-line-key  \n"),
		Credentials: store,
	}, "nano"); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if store.key != "one-line-key" {
		t.Errorf("stored key = %q, want one-line-key", store.key)
	}
}

func TestLoginRejectsAnEmptyKey(t *testing.T) {
	store := &memoryStore{}

	_, _, err := runLogin(t, cli.Deps{
		In:          strings.NewReader(" \n\t "),
		Credentials: store,
	}, "nano")
	if !errors.Is(err, credential.ErrEmptyKey) {
		t.Fatalf("Execute() error = %v, want ErrEmptyKey", err)
	}

	if store.key != "" {
		t.Errorf("stored key = %q, want nothing", store.key)
	}
}

func TestLoginReportsAStoreFailure(t *testing.T) {
	store := &memoryStore{saveErr: errors.New("keyring exploded")}

	_, _, err := runLogin(t, cli.Deps{
		In:          strings.NewReader("secret-key\n"),
		Credentials: store,
	}, "nano")
	if err == nil || !strings.Contains(err.Error(), "store API key") {
		t.Errorf("Execute() error = %v, want a store failure", err)
	}
}

func TestLoginDanbooru(t *testing.T) {
	t.Run("stores the login and key", func(t *testing.T) {
		store := &memoryStore{}

		out, errOut, err := runLogin(t, cli.Deps{
			In:          strings.NewReader("user\nkey\n"),
			Credentials: store,
		}, "danbooru")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if store.key != "user\nkey" {
			t.Errorf("stored credentials = %q, want the login then the key", store.key)
		}

		if !strings.Contains(out, "logged in") {
			t.Errorf("stdout = %q, want a login confirmation", out)
		}

		if strings.Contains(out+errOut, "key\n") {
			t.Errorf("output leaked the API key: stdout %q, stderr %q", out, errOut)
		}
	})

	t.Run("requires both lines", func(t *testing.T) {
		store := &memoryStore{}

		_, _, err := runLogin(t, cli.Deps{
			In:          strings.NewReader("user\n"),
			Credentials: store,
		}, "danbooru")
		if err == nil || !strings.Contains(err.Error(), "login line then an API key line") {
			t.Errorf("Execute() error = %v, want a two-line rejection", err)
		}

		if store.key != "" {
			t.Errorf("stored credentials = %q, want nothing", store.key)
		}
	})

	t.Run("warns when already logged in", func(t *testing.T) {
		store := &memoryStore{key: "old\nkey"}

		_, errOut, err := runLogin(t, cli.Deps{
			In:          strings.NewReader("user\nkey\n"),
			Credentials: store,
		}, "danbooru")
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}

		if !strings.Contains(errOut, "danbooru is already logged in") {
			t.Errorf("stderr = %q, want the replacement warning", errOut)
		}

		if store.key != "user\nkey" {
			t.Errorf("stored credentials = %q, want the replacement", store.key)
		}
	})
}
