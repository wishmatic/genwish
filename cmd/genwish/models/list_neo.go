package models

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/neomodels"
)

func runNeo(deps cli.Deps, args []string, query string, limit int) error {
	if err := noVideoKind("neo", args); err != nil {
		return err
	}

	registry, err := cli.NeoRegistry(deps)
	if err != nil {
		return err
	}

	ids := registry.IDs()
	needle := strings.ToLower(strings.TrimSpace(query))
	matched := make([]string, 0, len(ids))

	for _, id := range ids {
		model, _ := registry.Get(id)
		if needle == "" || neoModelMatches(id, model, needle) {
			matched = append(matched, id)
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
	if _, err := fmt.Fprintln(writer, "ID\tCHECKPOINT\tPRESET"); err != nil {
		return err
	}

	for _, id := range shown {
		model, _ := registry.Get(id)
		if _, err := fmt.Fprintf(
			writer,
			"%s\t%s\t%s\n",
			tableField(id),
			tableField(model.Checkpoint),
			tableField(model.Preset),
		); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Fprintf(
		deps.Out,
		"%d of %d models shown (%d matched)\n",
		len(shown),
		len(ids),
		len(matched),
	)

	return nil
}

func neoModelMatches(id string, model neomodels.Model, needle string) bool {
	for _, field := range []string{id, model.Checkpoint, model.Preset} {
		if strings.Contains(strings.ToLower(field), needle) {
			return true
		}
	}

	return false
}
