package models

import (
	"fmt"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/novelai"
)

func getNovel(deps cli.Deps, id string, args []string) error {
	if err := noVideoKind("novel", args); err != nil {
		return err
	}

	for _, model := range novelai.Models() {
		if model.ID == id {
			_, err := fmt.Fprintf(deps.Out, "id: %s\nkind: %s\n", model.ID, model.Kind)

			return err
		}
	}

	return fmt.Errorf("model %q is not served by novel", id)
}
