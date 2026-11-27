package credential

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zalando/go-keyring"

	"github.com/wishmatic/genwish/internal/provider"
)

var (
	ErrNotFound = errors.New("no stored API key")
	ErrEmptyKey = errors.New("API key is empty")
)

type Store interface {
	Load() (string, error)
	Save(key string) error
	Delete() error
}

func Open(p provider.Name) (Store, error) {
	if _, err := keyringUser(p); err != nil {
		return nil, err
	}

	file, err := NewFileStore(p)
	if err != nil {
		return nil, err
	}

	stores := []Store{KeyringStore{provider: p}}

	if pass, ok := newPassStore(p); ok {
		stores = append(stores, pass)
	}

	return newAutoStore(append(stores, file)...), nil
}

func StorageWarning(p provider.Name) string {
	user, err := keyringUser(p)
	if err != nil {
		return ""
	}

	_, keyringErr := keyring.Get(keyringService, user)
	if keyringErr == nil || errors.Is(keyringErr, keyring.ErrNotFound) {
		return ""
	}

	if _, ok := newPassStore(p); ok {
		return ""
	}

	path, err := configFilePath(p)
	if err != nil {
		path = ""
	}

	return storageWarning(keyringErr, path)
}

func storageWarning(keyringErr error, path string) string {
	if keyringErr == nil || errors.Is(keyringErr, keyring.ErrNotFound) {
		return ""
	}

	if path == "" {
		return "warning: the OS keyring is unavailable; the credential will be stored in a user-only file"
	}

	return fmt.Sprintf(
		"warning: the OS keyring is unavailable; the credential will be stored in %s",
		path,
	)
}

type FileStore struct {
	path string
}

func NewFileStore(p provider.Name) (*FileStore, error) {
	path, err := configFilePath(p)
	if err != nil {
		return nil, err
	}

	return NewFileStoreAt(path), nil
}

func NewFileStoreAt(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) Load() (string, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read API key: %w", err)
	}

	key := strings.TrimSpace(string(data))
	if key == "" {
		return "", ErrNotFound
	}

	return key, nil
}

func (s *FileStore) Save(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrEmptyKey
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	tmp, err := os.CreateTemp(dir, "."+filepath.Base(s.path)+"-*")
	if err != nil {
		return fmt.Errorf("create temporary API key file: %w", err)
	}

	name := tmp.Name()
	defer func() {
		_ = os.Remove(name)
	}()

	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()

		return fmt.Errorf("set API key permissions: %w", err)
	}

	if _, err := tmp.WriteString(key); err != nil {
		_ = tmp.Close()

		return fmt.Errorf("write API key: %w", err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close API key file: %w", err)
	}

	if err := os.Rename(name, s.path); err != nil {
		return fmt.Errorf("store API key: %w", err)
	}

	name = ""

	return nil
}

func (s *FileStore) Delete() error {
	err := os.Remove(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete API key: %w", err)
	}

	return nil
}

func configFilePath(p provider.Name) (string, error) {
	file, err := credentialFile(p)
	if err != nil {
		return "", err
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find config directory: %w", err)
	}

	return filepath.Join(dir, "genwish", file), nil
}

func credentialFile(p provider.Name) (string, error) {
	switch p {
	case provider.Nano:
		return "nanogpt_api_key", nil
	case provider.Novel:
		return "novelai_api_key", nil
	case provider.Danbooru:
		return "danbooru", nil
	default:
		return "", fmt.Errorf("provider %q does not use an API key", p)
	}
}

func keyringUser(p provider.Name) (string, error) {
	switch p {
	case provider.Nano:
		return "nanogpt-api-key", nil
	case provider.Novel:
		return "novelai-api-key", nil
	case provider.Danbooru:
		return "danbooru", nil
	default:
		return "", fmt.Errorf("provider %q does not use an API key", p)
	}
}

type autoStore struct {
	stores []Store
}

func newAutoStore(stores ...Store) *autoStore {
	return &autoStore{stores: stores}
}

func (s *autoStore) Load() (string, error) {
	var firstErr error

	for _, store := range s.stores {
		key, err := store.Load()
		if err == nil {
			return key, nil
		}

		if !errors.Is(err, ErrNotFound) && firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		return "", fmt.Errorf("%w (%v)", ErrNotFound, firstErr)
	}

	return "", ErrNotFound
}

func (s *autoStore) Save(key string) error {
	var lastErr error

	for index, store := range s.stores {
		if err := store.Save(key); err != nil {
			lastErr = err

			continue
		}

		for _, stale := range s.stores[index+1:] {
			if err := stale.Delete(); err != nil {
				return fmt.Errorf("remove stale credential: %w", err)
			}
		}

		return nil
	}

	return lastErr
}

func (s *autoStore) Delete() error {
	var firstErr error

	for _, store := range s.stores {
		err := store.Delete()
		if err == nil || errors.Is(err, ErrNotFound) {
			continue
		}

		if key, loadErr := store.Load(); loadErr == nil && key != "" {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}
