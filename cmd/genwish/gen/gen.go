package gen

import (
	"github.com/spf13/cobra"

	genimage "github.com/wishmatic/genwish/cmd/genwish/gen/image"
	genvideo "github.com/wishmatic/genwish/cmd/genwish/gen/video"
	"github.com/wishmatic/genwish/internal/cli"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate"},
		Short:   "Generate images and video",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		genimage.NewCommand(deps),
		genvideo.NewCommand(deps),
	)

	return cmd
}
