package models

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/wishmatic/genwish/internal/cli"
)

func runNeo(deps cli.Deps, args []string, query string, limit int) error {
	if len(args) > 0 {
		kind := strings.ToLower(strings.TrimSpace(args[0]))
		if kind != "image" {
			return fmt.Errorf("neo does not serve %s models", kind)
		}
	}

	cfg, err := cli.ResolveConfig(deps)
	if err != nil {
		return err
	}

	models, err := readNEOModels(cfg.Neo.ModelsFile)
	if err != nil {
		return err
	}

	needle := strings.ToLower(strings.TrimSpace(query))
	matched := make([]string, 0, len(models))

	for _, model := range models {
		if needle == "" || strings.Contains(strings.ToLower(model), needle) {
			matched = append(matched, model)
		}
	}

	shownLimit, err := modelLimit(limit)
	if err != nil {
		return err
	}

	shown := matched
	if len(shown) > shownLimit {
		shown = shown[:shownLimit]
	}

	fmt.Fprintln(deps.Out, "NEO MODELS")

	writer := tabwriter.NewWriter(deps.Out, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(writer, "MODEL"); err != nil {
		return err
	}

	for _, model := range shown {
		if _, err := fmt.Fprintln(writer, model); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Fprintf(deps.Out, "%d of %d models shown (%d matched)\n", len(shown), len(models), len(matched))

	return nil
}

func readNEOModels(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("neo.models_file %s: %w", path, err)
	}

	models := make([]string, 0)

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		models = append(models, line)
	}

	return models, nil
}
