package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
	"github.com/wishmatic/genwish/internal/sourcemap"
	"go.uber.org/zap"
)

type Deps struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	Config      config.Config
	Credentials credential.Store

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

	return deps
}

func ResolveConfig(deps Deps) (config.Config, error) {
	if deps.Config != (config.Config{}) {
		return deps.Config, nil
	}

	return config.Load()
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

	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	key := strings.TrimSpace(cfg.NanoGPTAPIKey)
	if key == "" {
		store, err := OpenCredentials(deps, provider.Nano)
		if err != nil {
			return nil, err
		}

		key, err = store.Load()
		if errors.Is(err, credential.ErrNotFound) {
			return nil, fmt.Errorf("not logged in: run `genwish login nano` or set NANOGPT_API_KEY")
		}
		if err != nil {
			return nil, fmt.Errorf("load API key: %w", err)
		}
	}

	return nanogpt.New(nanogpt.Config{APIKey: key})
}

func NovelClient(deps Deps) (*novelai.Client, error) {
	if deps.Novel != nil {
		return deps.Novel, nil
	}

	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	key := strings.TrimSpace(cfg.NovelAIAPIKey)
	if key == "" {
		store, err := OpenCredentials(deps, provider.Novel)
		if err != nil {
			return nil, err
		}

		key, err = store.Load()
		if errors.Is(err, credential.ErrNotFound) {
			return nil, fmt.Errorf("not logged in: run `genwish login novel` or set NOVELAI_API_KEY")
		}
		if err != nil {
			return nil, fmt.Errorf("load API key: %w", err)
		}
	}

	return novelai.New(novelai.DefaultBaseURL, key), nil
}

func ForgeClient(deps Deps, cfg config.Config) *forge.Client {
	if deps.Forge != nil {
		return deps.Forge
	}

	return forge.New(cfg.SDURL, zap.NewNop())
}

func CatalogService(deps Deps) (*catalog.Service, error) {
	if deps.Catalog != nil {
		return deps.Catalog, nil
	}

	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	provider, err := booruProvider(cfg)
	if err != nil {
		return nil, err
	}

	return catalog.New(provider, catalog.Options{
		BlockedTags: cfg.BlockedTags(),
		MaxLimit:    cfg.MaxLimit,
		MaxOffset:   cfg.MaxOffset,
		Relations:   nil,
	}), nil
}

func booruProvider(cfg config.Config) (*booru.Client, error) {
	transport, err := fetch.New(fetch.Config{
		BaseURL:       booru.DefaultBaseURL,
		UserAgent:     cfg.UserAgent,
		Timeout:       time.Duration(cfg.RequestTimeoutSeconds) * time.Second,
		Limiter:       fetch.NewLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst),
		VerboseErrors: cfg.VerboseErrors(),
	})
	if err != nil {
		return nil, fmt.Errorf("danbooru transport: %w", err)
	}

	return booru.New(booru.Config{
		BaseURL:  booru.DefaultBaseURL,
		Login:    cfg.DanbooruLogin,
		APIKey:   cfg.DanbooruAPIKey,
		MaxLimit: cfg.MaxLimit,
		HTTP:     transport,
	}), nil
}

func ImageResolver(deps Deps) (*resolve.Client, error) {
	if deps.Resolver != nil {
		return deps.Resolver, nil
	}

	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	sources, err := sourcemap.Parse(cfg.ImageURLMap)
	if err != nil {
		return nil, fmt.Errorf("IMAGE_URL_MAP: %w", err)
	}

	return resolve.New(sources), nil
}
