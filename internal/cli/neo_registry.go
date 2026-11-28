package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/wishmatic/genwish/internal/config"
	"github.com/wishmatic/genwish/internal/neomodels"
)

// NeoRegistry loads the Forge Neo model registry. A missing default file is an empty registry, so
// generation still works with a checkpoint named on the command line.
func NeoRegistry(deps Deps) (neomodels.Registry, error) {
	cfg, err := ResolveConfig(deps)
	if err != nil {
		return nil, err
	}

	defaultPath := config.DefaultNeoModelsFile()
	path := strings.TrimSpace(cfg.Neo.ModelsFile)
	if path == "" {
		path = defaultPath
	}

	registry, err := neomodels.Load(path)
	if err == nil {
		return registry, nil
	}

	if errors.Is(err, fs.ErrNotExist) && path == defaultPath {
		return neomodels.Registry{}, nil
	}

	return nil, fmt.Errorf("neo.models_file %s: %w", path, err)
}
