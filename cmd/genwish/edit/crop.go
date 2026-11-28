package edit

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	imageedit "github.com/wishmatic/genwish/internal/edit"
	"github.com/wishmatic/genwish/internal/format"
)

type cropCommand struct {
	name           string
	short          string
	isSquare       bool
	isCircle       bool
	defaultPadding int
}

func NewTrimCommand(deps cli.Deps) *cobra.Command {
	return newCropCommand(
		deps,
		cropCommand{
			name:  "trim",
			short: "Trim an image to the bounding box of its visible content",
		},
	)
}

func NewSquareCommand(deps cli.Deps) *cobra.Command {
	return newCropCommand(
		deps,
		cropCommand{
			name:           "square",
			short:          "Trim an image and center it on a transparent square",
			isSquare:       true,
			defaultPadding: imageedit.DefaultSquarePadding,
		},
	)
}

func NewCircleCommand(deps cli.Deps) *cobra.Command {
	return newCropCommand(
		deps,
		cropCommand{
			name:     "circle",
			short:    "Trim an image and cut it into the inscribed circle",
			isCircle: true,
		},
	)
}

func newCropCommand(deps cli.Deps, spec cropCommand) *cobra.Command {
	var (
		padding      int
		outputFormat string
		out          string
	)

	cmd := &cobra.Command{
		Use:   spec.name + " <image>",
		Short: spec.short,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			opts, _ := imageedit.OptionsFor(imageedit.Request{
				IsCrop:   true,
				IsSquare: spec.isSquare,
				IsCircle: spec.isCircle,
				Padding:  &padding,
			})

			return runCrop(cmd, deps, opts, outputFormat, out, args[0])
		},
	}

	cmd.Flags().IntVar(
		&padding,
		"padding",
		spec.defaultPadding,
		"transparent padding in pixels around the trimmed content",
	)
	cmd.Flags().StringVar(
		&outputFormat,
		"format",
		"",
		"output image format: png, jpeg, jxl, or webp",
	)
	cmd.Flags().StringVar(
		&out,
		"out",
		"",
		"base directory for the edited image (default: output.path)",
	)

	return cmd
}

func runCrop(
	cmd *cobra.Command, deps cli.Deps, opts imageedit.Options, outputFormat, out, image string,
) error {
	cfg, err := cli.ResolveConfig(deps)
	if err != nil {
		return err
	}

	target, err := cli.OutputFormat(cfg, outputFormat)
	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}

	resolver := cli.ImageResolver(deps)

	source, err := cli.ReadImage(cmd.Context(), resolver, image)
	if err != nil {
		return fmt.Errorf("edit: read image: %w", err)
	}

	edited, err := editImage(source, opts, target)
	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}

	path, err := cli.WriteImage(cli.OutputDir(cfg, out), "images", edited, target.MediaType())
	if err != nil {
		return fmt.Errorf("edit: %w", err)
	}

	return cli.PrintPaths(deps.Out, []string{path})
}

func editImage(data []byte, opts imageedit.Options, target format.Format) ([]byte, error) {
	src, err := format.Decode(data)
	if err != nil {
		return nil, err
	}

	out, err := imageedit.Apply(src, opts)
	if err != nil {
		return nil, err
	}

	return format.Encode(out, target)
}
