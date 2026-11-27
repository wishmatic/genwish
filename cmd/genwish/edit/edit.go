package edit

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
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
		NewTrimCommand(deps),
		NewSquareCommand(deps),
		NewCircleCommand(deps),
		NewConvertCommand(deps),
		NewBgkillCommand(deps),
	)

	return cmd
}
