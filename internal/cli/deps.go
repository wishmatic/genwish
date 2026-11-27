package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/wishmatic/genwish/internal/booru"
	"github.com/wishmatic/genwish/internal/catalog"
	"github.com/wishmatic/genwish/internal/config"
	"github.com/wishmatic/genwish/internal/credential"
	"github.com/wishmatic/genwish/internal/fetch"
	"github.com/wishmatic/genwish/internal/forge"
	"github.com/wishmatic/genwish/internal/nanogpt"
	"github.com/wishmatic/genwish/internal/novelai"
	"github.com/wishmatic/genwish/internal/provider"
	"github.com/wishmatic/genwish/internal/resolve"
)

const userAgent = "genwish/0.1.0"

type Flags struct {
	ConfigPath string
}

type Deps struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	Config      config.Config
	Credentials credential.Store
	Flags       *Flags

	Nano     *nanogpt.Client
	Novel    *novelai.Client
	Forge    *forge.Client
	Resolver *resolve.Client
	Catalog  *catalog.Service
}

func DefaultDeps() Deps {
	return Deps{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

func WithDefaults(deps Deps) Deps {
	if deps.In == nil {
		deps.In = os.Stdin
	}
	if deps.Out == nil {
		deps.Out = os.Stdout
	}
	if deps.Err == nil {
		deps.Err = os.Stderr
	}
	if deps.Flags == nil {
		deps.Flags = &Flags{}
	}

	return deps
}

func ResolveConfig(deps Deps) (config.Config, error) {
	if deps.Config != (config.Config{}) {
		return deps.Config, nil
	}

	return config.Load(deps.configPath())
}

func (d Deps) configPath() string {
	if d.Flags == nil {
		return ""
	}

	return d.Flags.ConfigPath
}

func OpenCredentials(deps Deps, p provider.Name) (credential.Store, error) {
	if deps.Credentials != nil {
		return deps.Credentials, nil
	}

	return credential.Open(p)
}

func NanoClient(deps Deps) (*nanogpt.Client, error) {
	if deps.Nano != nil {
		return deps.Nano, nil
	}

	store, err := OpenCredentials(deps, provider.Nano)
	if err != nil {
		return nil, err
	}

	key, err := store.Load()
	if errors.Is(err, credential.ErrNotFound) {
		return nil, fmt.Errorf("not logged in: run `genwish login nano`")
	}
	if err != nil {
		return nil, fmt.Errorf("load API key: %w", err)
	}

	return nanogpt.New(nanogpt.Config{APIKey: key})
}

func NovelClient(deps Deps) (*novelai.Client, error) {
	if deps.Novel != nil {
		return deps.Novel, nil
	}

	store, err := OpenCredentials(deps, provider.Novel)
	if err != nil {
		return nil, err
	}

	key, err := store.Load()
	if errors.Is(err, credential.ErrNotFound) {
		return nil, fmt.Errorf("not logged in: run `genwish login novel`")
	}
	if err != nil {
		return nil, fmt.Errorf("load API key: %w", err)
	}

	return novelai.New(novelai.DefaultBaseURL, key), nil
}

func ForgeClient(deps Deps, cfg config.Config) *forge.Client {
	if deps.Forge != nil {
		return deps.Forge
	}

	return forge.New(cfg.Neo.URL, deps.Err)
}

func CatalogService(deps Deps) (*catalog.Service, error) {
	if deps.Catalog != nil {
		return deps.Catalog, nil
	}

	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	creds, err := danbooruCredentials(deps)
	if err != nil {
		return nil, err
	}

	danbooruClient, err := booruProvider(cfg, creds)
	if err != nil {
		return nil, err
	}

	return catalog.New(danbooruClient, catalog.Options{}), nil
}

func danbooruCredentials(deps Deps) (credential.DanbooruCredentials, error) {
	store, err := OpenCredentials(deps, provider.Danbooru)
	if err != nil {
		return credential.DanbooruCredentials{}, err
	}

	creds, err := credential.LoadDanbooru(store)
	if errors.Is(err, credential.ErrNotFound) {
		return credential.DanbooruCredentials{}, nil
	}
	if err != nil {
		return credential.DanbooruCredentials{}, err
	}

	return creds, nil
}

func booruProvider(cfg config.Config, creds credential.DanbooruCredentials) (*booru.Client, error) {
	transport, err := fetch.New(fetch.Config{
		BaseURL:   booru.DefaultBaseURL,
		UserAgent: userAgent,
		Timeout:   time.Duration(cfg.Booru.RateLimit.Timeout) * time.Second,
		Limiter:   fetch.NewLimiter(cfg.Booru.RateLimit.RPS, cfg.Booru.RateLimit.Burst),
	})
	if err != nil {
		return nil, fmt.Errorf("danbooru transport: %w", err)
	}

	return booru.New(booru.Config{
		BaseURL: booru.DefaultBaseURL,
		Login:   creds.Login,
		APIKey:  creds.APIKey,
		HTTP:    transport,
	}), nil
}

func ImageResolver(deps Deps) *resolve.Client {
	if deps.Resolver != nil {
		return deps.Resolver
	}

	return resolve.New()
}
