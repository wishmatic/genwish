package config

import (
	"errors"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/cli"
)

func TestLogoutRejectsNeo(t *testing.T) {
	t.Run("reports the provider does not use a key", func(t *testing.T) {
		_, _, err := runLogout(t, cli.Deps{}, "neo")
		if err == nil || !strings.Contains(err.Error(), "does not use an API key") {
			t.Errorf("Execute() error = %v, want neo rejected", err)
		}
	})
}

func TestLogoutDeletesTheStoredKey(t *testing.T) {
	store := &memoryStore{key: "secret-key"}

	out, _, err := runLogout(t, cli.Deps{Credentials: store}, "nano")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !store.deleted || store.key != "" {
		t.Errorf("store after logout = %+v, want the key deleted", store)
	}

	if !strings.Contains(out, "logged out") {
		t.Errorf("stdout = %q, want a logout confirmation", out)
	}
}

func TestLogoutSucceedsWithoutAStoredKey(t *testing.T) {
	store := &memoryStore{}

	if _, _, err := runLogout(t, cli.Deps{Credentials: store}, "nano"); err != nil {
		t.Errorf("Execute() error: %v", err)
	}
}

func TestLogoutReportsADeleteFailure(t *testing.T) {
	store := &memoryStore{deleteErr: errors.New("delete failed")}

	_, _, err := runLogout(t, cli.Deps{Credentials: store}, "nano")
	if err == nil || !strings.Contains(err.Error(), "delete API key") {
		t.Errorf("Execute() error = %v, want a delete failure", err)
	}
}

func TestLogoutDeletesTheDanbooruCredential(t *testing.T) {
	store := &memoryStore{key: "user\nkey"}

	out, _, err := runLogout(t, cli.Deps{Credentials: store}, "danbooru")
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if !store.deleted || store.key != "" {
		t.Errorf("store after logout = %+v, want the credential deleted", store)
	}

	if !strings.Contains(out, "logged out") {
		t.Errorf("stdout = %q, want a logout confirmation", out)
	}
}
