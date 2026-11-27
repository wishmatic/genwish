package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/nanogpt"
)

func getNano(ctx context.Context, deps cli.Deps, id string, args []string) error {
	kinds, err := modelKinds(args)
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
			return fmt.Errorf("models get %s: %w", kind, err)
		}

		for _, model := range all {
			if model.ID == id {
				return writeNanoModel(deps, model)
			}
		}
	}

	return fmt.Errorf("model %q is not served by nano", id)
}

func writeNanoModel(deps cli.Deps, model nanogpt.Model) error {
	var b strings.Builder

	fmt.Fprintf(&b, "id: %s\n", model.ID)

	if model.Name != "" {
		fmt.Fprintf(&b, "name: %s\n", model.Name)
	}

	if model.OwnedBy != "" {
		fmt.Fprintf(&b, "owner: %s\n", model.OwnedBy)
	}

	if model.Description != "" {
		fmt.Fprintf(&b, "description: %s\n", model.Description)
	}

	if len(model.Tags) > 0 {
		fmt.Fprintf(&b, "tags: %s\n", strings.Join(model.Tags, ", "))
	}

	fields := []struct {
		name  string
		value map[string]any
	}{
		{name: "capabilities", value: model.Capabilities},
		{name: "supported_parameters", value: model.SupportedParameters},
		{name: "pricing", value: model.Pricing},
	}

	for _, field := range fields {
		if len(field.value) == 0 {
			continue
		}

		encoded, err := json.Marshal(field.value)
		if err != nil {
			return err
		}

		fmt.Fprintf(&b, "%s: %s\n", field.name, encoded)
	}

	_, err := io.WriteString(deps.Out, b.String())

	return err
}
