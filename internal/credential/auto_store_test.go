package credential

import (
	"errors"
	"strings"
	"testing"
)

type stubStore struct {
	key       string
	loadErr   error
	saveErr   error
	deleteErr error
	deleted   bool
}

func (s *stubStore) Load() (string, error) {
	if s.loadErr != nil {
		return "", s.loadErr
	}

	if s.key == "" {
		return "", ErrNotFound
	}

	return s.key, nil
}

func (s *stubStore) Save(key string) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.key = key

	return nil
}

func (s *stubStore) Delete() error {
	s.deleted = true

	if s.deleteErr != nil {
		return s.deleteErr
	}

	s.key = ""

	return nil
}

func TestAutoStoreLoad(t *testing.T) {
	t.Run("prefers the first store with a key", func(t *testing.T) {
		store := newAutoStore(
			&stubStore{key: "first"},
			&stubStore{key: "second"},
		)

		key, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if key != "first" {
			t.Errorf("Load() = %q, want first", key)
		}
	})

	t.Run("falls through not found and errors", func(t *testing.T) {
		store := newAutoStore(
			&stubStore{loadErr: errors.New("unavailable")},
			&stubStore{},
			&stubStore{key: "file-key"},
		)

		key, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
		}

		if key != "file-key" {
			t.Errorf("Load() = %q, want file-key", key)
		}
	})

	t.Run("reports not found when every store is empty", func(t *testing.T) {
		store := newAutoStore(&stubStore{}, &stubStore{})

		if _, err := store.Load(); !errors.Is(err, ErrNotFound) {
			t.Errorf("Load() error = %v, want ErrNotFound", err)
		}
	})

	t.Run("keeps the underlying error", func(t *testing.T) {
		store := newAutoStore(&stubStore{loadErr: errors.New("unavailable")})

		_, err := store.Load()
		if !errors.Is(err, ErrNotFound) || !strings.Contains(err.Error(), "unavailable") {
			t.Errorf("Load() error = %v, want ErrNotFound with the cause", err)
		}
	})
}

func TestAutoStoreSave(t *testing.T) {
	t.Run("uses the first writable store and clears later copies", func(t *testing.T) {
		fallback := &stubStore{key: "old"}
		store := newAutoStore(
			&stubStore{saveErr: errors.New("unavailable")},
			&stubStore{key: "old-pass"},
			fallback,
		)

		if err := store.Save("new-key"); err != nil {
			t.Fatalf("Save() error: %v", err)
		}

		if store.stores[1].(*stubStore).key != "new-key" {
			t.Error("Save() did not write to the first writable store")
		}

		if !fallback.deleted {
			t.Error("Save() left the lower-priority copy behind")
		}
	})

	t.Run("reports the last failure", func(t *testing.T) {
		store := newAutoStore(
			&stubStore{saveErr: errors.New("keyring down")},
			&stubStore{saveErr: errors.New("pass down")},
		)

		err := store.Save("key")
		if err == nil || !strings.Contains(err.Error(), "pass down") {
			t.Errorf("Save() error = %v, want the last failure", err)
		}
	})
}

func TestAutoStoreDelete(t *testing.T) {
	t.Run("clears every reachable store", func(t *testing.T) {
		first := &stubStore{key: "key"}
		second := &stubStore{key: "key"}
		store := newAutoStore(first, second)

		if err := store.Delete(); err != nil {
			t.Fatalf("Delete() error: %v", err)
		}

		if !first.deleted || !second.deleted {
			t.Error("Delete() left a store untouched")
		}
	})

	t.Run("tolerates an unreachable store", func(t *testing.T) {
		unreachable := &stubStore{
			loadErr:   errors.New("unavailable"),
			deleteErr: errors.New("unavailable"),
		}
		reachable := &stubStore{key: "key"}
		store := newAutoStore(unreachable, reachable)

		if err := store.Delete(); err != nil {
			t.Errorf("Delete() error = %v, want the unreachable store tolerated", err)
		}

		if !reachable.deleted {
			t.Error("Delete() left the reachable store untouched")
		}
	})

	t.Run("reports a surviving credential", func(t *testing.T) {
		store := newAutoStore(&stubStore{key: "key", deleteErr: errors.New("locked")})

		err := store.Delete()
		if err == nil || !strings.Contains(err.Error(), "locked") {
			t.Errorf("Delete() error = %v, want the surviving credential reported", err)
		}
	})
}
