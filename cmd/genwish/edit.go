package main

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/edit"
)

func newEditCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Trim, pad, square, circularise, re-encode, or remove a background",
		Long: "Operate on an image given as a local path or an address the resolver can read. Every " +
			"subcommand writes its result to --out and prints the path.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newEditCropCommand(deps, cropCommand{
			name:  "trim",
			short: "Trim an image to the bounding box of its visible content",
		}),
		newEditCropCommand(deps, cropCommand{
			name:           "square",
			short:          "Trim an image and center it on a transparent square",
			isSquare:       true,
			defaultPadding: edit.DefaultSquarePadding,
		}),
		newEditCropCommand(deps, cropCommand{
			name:     "circle",
			short:    "Trim an image and cut it into the inscribed circle",
			isCircle: true,
		}),
		newEditConvertCommand(deps),
		newEditBgkillCommand(deps),
	)

	return cmd
}

type cropCommand struct {
	name           string
	short          string
	isSquare       bool
	isCircle       bool
	defaultPadding int
}
