package models

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/nanogpt"
)

const (
	defaultModelLimit = 25
	maxModelLimit     = 100
)

func runNano(ctx context.Context, deps cli.Deps, args []string, query string, limit int) error {
	kinds, err := modelKinds(args)
	if err != nil {
		return err
	}

	shownLimit, err := modelLimit(limit)
	if err != nil {
		return err
	}

	client, err := cli.NanoClient(deps)
	if err != nil {
		return err
	}

	for _, kind := range kinds {
		all, err := client.Models(ctx, kind)
		if err != nil {
			return fmt.Errorf("models list %s: %w", kind, err)
		}

		matched := filterModels(all, query)
		shown := matched
		if len(shown) > shownLimit {
			shown = shown[:shownLimit]
		}

		if err := printModels(deps, kind, all, matched, shown); err != nil {
			return err
		}
	}

	return nil
}

func modelKinds(args []string) ([]string, error) {
	if len(args) == 0 {
		return nanogpt.ModelKinds, nil
	}

	kind := strings.ToLower(strings.TrimSpace(args[0]))
	if !slices.Contains(nanogpt.ModelKinds, kind) {
		return nil, fmt.Errorf("model kind %q is not one of %v", kind, nanogpt.ModelKinds)
	}

	return []string{kind}, nil
}

func modelLimit(limit int) (int, error) {
	switch {
	case limit == 0:
		return defaultModelLimit, nil
	case limit < 0 || limit > maxModelLimit:
		return 0, fmt.Errorf("limit must be 1 to %d, got %d", maxModelLimit, limit)
	default:
		return limit, nil
	}
}

func filterModels(models []nanogpt.Model, query string) []nanogpt.Model {
	needle := strings.ToLower(strings.TrimSpace(query))
	if needle == "" {
		return models
	}

	matched := make([]nanogpt.Model, 0, len(models))

	for _, model := range models {
		if modelMatches(model, needle) {
			matched = append(matched, model)
		}
	}

	return matched
}

func modelMatches(model nanogpt.Model, needle string) bool {
	fields := append([]string{model.ID, model.Name, model.Description, model.OwnedBy}, model.Tags...)

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), needle) {
			return true
		}
	}

	return false
}

func printModels(
	deps cli.Deps, kind string, all, matched, shown []nanogpt.Model,
) error {
	fmt.Fprintf(deps.Out, "%s MODELS\n", strings.ToUpper(kind))

	writer := tabwriter.NewWriter(deps.Out, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(writer, "ID\tNAME\tOWNER\tDESCRIPTION"); err != nil {
		return err
	}

	for _, model := range shown {
		if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n",
			tableField(model.ID), tableField(model.Name), tableField(model.OwnedBy),
			tableField(model.Description),
		); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Fprintf(deps.Out, "%d of %d models shown (%d matched)\n\n", len(shown), len(all), len(matched))

	return nil
}

func tableField(value string) string {
	value = strings.Join(strings.Fields(value), " ")
	if value == "" {
		return "-"
	}

	return value
}
