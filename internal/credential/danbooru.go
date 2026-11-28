package credential

import (
	"fmt"
	"strings"
)

type DanbooruCredentials struct {
	Login  string
	APIKey string
}

func SaveDanbooru(store Store, creds DanbooruCredentials) error {
	login := strings.TrimSpace(creds.Login)
	apiKey := strings.TrimSpace(creds.APIKey)

	if login == "" || apiKey == "" {
		return fmt.Errorf("Danbooru login and API key are both required")
	}

	if strings.ContainsAny(login, "\r\n") || strings.ContainsAny(apiKey, "\r\n") {
		return fmt.Errorf("Danbooru login and API key must each be a single line")
	}

	return store.Save(fmt.Sprintf("%s\n%s", login, apiKey))
}

func LoadDanbooru(store Store) (DanbooruCredentials, error) {
	raw, err := store.Load()
	if err != nil {
		return DanbooruCredentials{}, err
	}

	login, apiKey, ok := strings.Cut(raw, "\n")
	login = strings.TrimSpace(login)
	apiKey = strings.TrimSpace(apiKey)

	if !ok || login == "" || apiKey == "" {
		return DanbooruCredentials{}, ErrNotFound
	}

	return DanbooruCredentials{Login: login, APIKey: apiKey}, nil
}
