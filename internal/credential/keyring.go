package credential

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"

	"github.com/wishmatic/genwish/internal/provider"
)

const keyringService = "genwish"

type KeyringStore struct {
	provider provider.Name
}

func (s KeyringStore) Load() (string, error) {
	user, err := keyringUser(s.provider)
	if err != nil {
		return "", err
	}

	key, err := keyring.Get(keyringService, user)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read API key from keyring: %w", err)
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrNotFound
	}

	return key, nil
}

func (s KeyringStore) Save(key string) error {
	user, err := keyringUser(s.provider)
	if err != nil {
		return err
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return ErrEmptyKey
	}

	if err := keyring.Set(keyringService, user, key); err != nil {
		return fmt.Errorf("store API key in keyring: %w", err)
	}

	return nil
}

func (s KeyringStore) Delete() error {
	user, err := keyringUser(s.provider)
	if err != nil {
		return err
	}

	err = keyring.Delete(keyringService, user)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete API key from keyring: %w", err)
	}

	return nil
}
