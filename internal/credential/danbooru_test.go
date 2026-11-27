package credential

import (
	"errors"
	"path/filepath"
	"testing"
)

type danbooruMemoryStore struct {
	value string
	err   error
}

func (s *danbooruMemoryStore) Load() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	if s.value == "" {
		return "", ErrNotFound
	}

	return s.value, nil
}

func (s *danbooruMemoryStore) Save(value string) error {
	s.value = value

	return nil
}

func (s *danbooruMemoryStore) Delete() error {
	s.value = ""

	return nil
}

func TestDanbooruCredentials(t *testing.T) {
	t.Run("round trips", func(t *testing.T) {
		store := &danbooruMemoryStore{}

		if err := SaveDanbooru(store, DanbooruCredentials{Login: "user", APIKey: "key"}); err != nil {
			t.Fatalf("SaveDanbooru() error: %v", err)
		}

		got, err := LoadDanbooru(store)
		if err != nil {
			t.Fatalf("LoadDanbooru() error: %v", err)
		}

		if got.Login != "user" || got.APIKey != "key" {
			t.Errorf("LoadDanbooru() = %+v, want user and key", got)
		}
	})

	t.Run("requires both values", func(t *testing.T) {
		store := &danbooruMemoryStore{}

		if err := SaveDanbooru(store, DanbooruCredentials{Login: "user"}); err == nil {
			t.Error("SaveDanbooru() error = nil, want a missing-key error")
		}
	})

	t.Run("reports missing credentials", func(t *testing.T) {
		_, err := LoadDanbooru(&danbooruMemoryStore{})
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("LoadDanbooru() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("reports a malformed stored value", func(t *testing.T) {
		_, err := LoadDanbooru(&danbooruMemoryStore{value: "only-a-login"})
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("LoadDanbooru() error = %v, want ErrNotFound", err)
		}
	})
}

func TestDanbooruFileStore(t *testing.T) {
	store := NewFileStoreAt(filepath.Join(t.TempDir(), "danbooru"))

	if err := SaveDanbooru(store, DanbooruCredentials{Login: "user", APIKey: "key"}); err != nil {
		t.Fatalf("SaveDanbooru() error: %v", err)
	}

	got, err := LoadDanbooru(store)
	if err != nil {
		t.Fatalf("LoadDanbooru() error: %v", err)
	}

	if got.Login != "user" || got.APIKey != "key" {
		t.Errorf("LoadDanbooru() = %+v, want user and key", got)
	}

	if err := store.Delete(); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	if _, err := LoadDanbooru(store); !errors.Is(err, ErrNotFound) {
		t.Errorf("LoadDanbooru() after Delete error = %v, want ErrNotFound", err)
	}
}
