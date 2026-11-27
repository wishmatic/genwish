package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/novelai"
)

func runModelsListNovel(deps cli.Deps, args []string, query string, limit int) error {
	if len(args) > 0 {
		kind := strings.ToLower(strings.TrimSpace(args[0]))
		if kind != "image" {
			return fmt.Errorf("novel does not serve %s models", kind)
		}
	}

	all := novelai.Models()

	needle := strings.ToLower(strings.TrimSpace(query))
	matched := make([]novelai.Model, 0, len(all))

	for _, model := range all {
		if needle == "" || strings.Contains(strings.ToLower(model.ID), needle) {
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

	fmt.Fprintln(deps.Out, "NOVEL MODELS")

	writer := tabwriter.NewWriter(deps.Out, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(writer, "MODEL"); err != nil {
		return err
	}

	for _, model := range shown {
		if _, err := fmt.Fprintln(writer, model.ID); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	fmt.Fprintf(deps.Out, "%d of %d models shown (%d matched)\n", len(shown), len(all), len(matched))

	return nil
}
