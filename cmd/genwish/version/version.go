package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	buildversion "github.com/wishmatic/genwish/internal/version"
)

func NewCommand(deps cli.Deps) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the genwish version",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			_, err := fmt.Fprintln(deps.Out, buildversion.String())

			return err
		},
	}
}
