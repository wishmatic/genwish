package models

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/neomodels"
)

func getNeo(deps cli.Deps, id string, args []string) error {
	if err := noVideoKind("neo", args); err != nil {
		return err
	}

	registry, err := cli.NeoRegistry(deps)
	if err != nil {
		return err
	}

	model, ok := registry.Get(id)
	if !ok {
		return fmt.Errorf("model %q is not configured for neo", id)
	}

	return writeNEOModel(deps, id, model)
}

func writeNEOModel(deps cli.Deps, id string, model neomodels.Model) error {
	var b strings.Builder

	fmt.Fprintf(&b, "id: %s\n", id)

	if model.Checkpoint != "" {
		fmt.Fprintf(&b, "checkpoint: %s\n", model.Checkpoint)
	}

	if len(model.Encoders) > 0 {
		fmt.Fprintf(&b, "encoders: %s\n", strings.Join(model.Encoders, ", "))
	}

	if model.Preset != "" {
		fmt.Fprintf(&b, "preset: %s\n", model.Preset)
	}

	writeChoice(&b, "samplers", model.Samplers)
	writeChoice(&b, "schedulers", model.Schedulers)
	writeIntRange(&b, "steps", model.Steps)
	writeFloatRange(&b, "cfg", model.CFG)

	if model.Hires.Enabled {
		b.WriteString("hires:\n  enabled: true\n")
	}

	if len(model.Sizes) > 0 {
		fmt.Fprintf(&b, "sizes: %s\n", strings.Join(model.Sizes, ", "))
	}

	writeInstructions(&b, model.Instructions)

	_, err := io.WriteString(deps.Out, b.String())

	return err
}

func writeChoice(b *strings.Builder, name string, choice neomodels.Choice) {
	if choice.Default == "" && len(choice.Options) == 0 {
		return
	}

	fmt.Fprintf(b, "%s:\n", name)

	if choice.Default != "" {
		fmt.Fprintf(b, "  default: %s\n", choice.Default)
	}

	if len(choice.Options) > 0 {
		fmt.Fprintf(b, "  options: %s\n", strings.Join(choice.Options, ", "))
	}
}

func writeIntRange(b *strings.Builder, name string, value neomodels.Range[int]) {
	if value.Default == nil && value.Min == nil && value.Max == nil {
		return
	}

	fmt.Fprintf(b, "%s:\n", name)

	if value.Default != nil {
		fmt.Fprintf(b, "  default: %d\n", *value.Default)
	}

	if value.Min != nil {
		fmt.Fprintf(b, "  min: %d\n", *value.Min)
	}

	if value.Max != nil {
		fmt.Fprintf(b, "  max: %d\n", *value.Max)
	}
}

func writeFloatRange(b *strings.Builder, name string, value neomodels.Range[float64]) {
	if value.Default == nil && value.Min == nil && value.Max == nil {
		return
	}

	fmt.Fprintf(b, "%s:\n", name)

	if value.Default != nil {
		fmt.Fprintf(b, "  default: %s\n", formatFloat(*value.Default))
	}

	if value.Min != nil {
		fmt.Fprintf(b, "  min: %s\n", formatFloat(*value.Min))
	}

	if value.Max != nil {
		fmt.Fprintf(b, "  max: %s\n", formatFloat(*value.Max))
	}
}

func writeInstructions(b *strings.Builder, instructions neomodels.Instructions) {
	if instructions.Usage == "" && instructions.Style == "" {
		return
	}

	b.WriteString("instructions:\n")

	writeInstruction(b, "usage", instructions.Usage)
	writeInstruction(b, "style", instructions.Style)
}

func writeInstruction(b *strings.Builder, name, value string) {
	value = strings.TrimRight(value, "\n")
	if value == "" {
		return
	}

	lines := strings.Split(value, "\n")
	if len(lines) == 1 {
		fmt.Fprintf(b, "  %s: %s\n", name, value)

		return
	}

	fmt.Fprintf(b, "  %s: |\n", name)

	for _, line := range lines {
		if line == "" {
			b.WriteString("\n")

			continue
		}

		fmt.Fprintf(b, "    %s\n", line)
	}
}

func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}
