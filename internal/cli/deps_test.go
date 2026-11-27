package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/wishmatic/genwish/internal/config"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/nanogpt"
	"github.com/wishmatic/genwish/internal/provider"
)

type memoryStore struct {
	key     string
	loadErr error
}

func (s *memoryStore) Load() (string, error) {
	if s.loadErr != nil {
		return "", s.loadErr
	}
	if s.key == "" {
		return "", credential.ErrNotFound
	}

	return s.key, nil
}

func (s *memoryStore) Save(key string) error {
	s.key = strings.TrimSpace(key)

	return nil
}

func (s *memoryStore) Delete() error {
	s.key = ""

	return nil
}

func TestOpenCredentials(t *testing.T) {
	t.Run("uses the injected store", func(t *testing.T) {
		store := &memoryStore{}

		got, err := OpenCredentials(Deps{Credentials: store}, provider.Nano)
		if err != nil {
			t.Fatalf("OpenCredentials() error: %v", err)
		}

		if got != store {
			t.Error("OpenCredentials() did not return the injected store")
		}
	})
}

func newTestClient(t *testing.T) *nanogpt.Client {
	t.Helper()

	client, err := nanogpt.New(nanogpt.Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("nanogpt.New() error: %v", err)
	}

	return client
}

func TestNanoClient(t *testing.T) {
	t.Run("uses an injected client", func(t *testing.T) {
		client := newTestClient(t)

		got, err := NanoClient(Deps{Nano: client})
		if err != nil {
			t.Fatalf("NanoClient() error: %v", err)
		}

		if got != client {
			t.Error("NanoClient() did not return the injected client")
		}
	})

	t.Run("uses the stored key", func(t *testing.T) {
		client, err := NanoClient(Deps{Credentials: &memoryStore{key: "stored-key"}})
		if err != nil {
			t.Fatalf("NanoClient() error: %v", err)
		}

		if client == nil {
			t.Error("NanoClient() = nil, want a client")
		}
	})

	t.Run("reports a missing key", func(t *testing.T) {
		_, err := NanoClient(Deps{Credentials: &memoryStore{}})
		if err == nil || !strings.Contains(err.Error(), "genwish login nano") {
			t.Errorf("NanoClient() error = %v, want a login hint", err)
		}
	})
}

func TestDanbooruCredentials(t *testing.T) {
	t.Run("reads the stored login and key", func(t *testing.T) {
		creds, err := danbooruCredentials(Deps{Credentials: &memoryStore{key: "user\nkey"}})
		if err != nil {
			t.Fatalf("danbooruCredentials() error: %v", err)
		}

		if creds.Login != "user" || creds.APIKey != "key" {
			t.Errorf("credentials = %+v, want user and key", creds)
		}
	})

	t.Run("falls back to anonymous access", func(t *testing.T) {
		creds, err := danbooruCredentials(Deps{Credentials: &memoryStore{}})
		if err != nil {
			t.Fatalf("danbooruCredentials() error: %v", err)
		}

		if creds.Login != "" || creds.APIKey != "" {
			t.Errorf("credentials = %+v, want none", creds)
		}
	})
}

func TestCatalogService(t *testing.T) {
	t.Run("builds without stored danbooru credentials", func(t *testing.T) {
		cfg := config.Config{
			Booru: config.BooruConfig{
				RateLimit: config.RateLimitConfig{RPS: 1, Burst: 1, Timeout: 20},
			},
		}

		service, err := CatalogService(Deps{Config: cfg, Credentials: &memoryStore{}})
		if err != nil {
			t.Fatalf("CatalogService() error: %v", err)
		}

		if service == nil {
			t.Error("CatalogService() = nil, want a service")
		}
	})
}

func TestNovelClient(t *testing.T) {
	t.Run("reports a missing key", func(t *testing.T) {
		_, err := NovelClient(Deps{Credentials: &memoryStore{}})
		if err == nil || !strings.Contains(err.Error(), "genwish login novel") {
			t.Errorf("NovelClient() error = %v, want a login hint", err)
		}
	})

	t.Run("uses the stored key", func(t *testing.T) {
		client, err := NovelClient(Deps{Credentials: &memoryStore{key: "stored-key"}})
		if err != nil {
			t.Fatalf("NovelClient() error: %v", err)
		}

		if client == nil {
			t.Error("NovelClient() = nil, want a client")
		}
	})
}

func TestOpenCredentialsWithAStoreError(t *testing.T) {
	store := &memoryStore{loadErr: errors.New("keyring unavailable")}

	_, err := NanoClient(Deps{Credentials: store})
	if err == nil || !strings.Contains(err.Error(), "keyring unavailable") {
		t.Errorf("NanoClient() error = %v, want the store failure", err)
	}
}
