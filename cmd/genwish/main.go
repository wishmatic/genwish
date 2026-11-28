package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/cmd/genwish/balance"
	"github.com/wishmatic/genwish/cmd/genwish/config"
	"github.com/wishmatic/genwish/cmd/genwish/edit"
	"github.com/wishmatic/genwish/cmd/genwish/gen"
	"github.com/wishmatic/genwish/cmd/genwish/models"
	"github.com/wishmatic/genwish/cmd/genwish/util"
	"github.com/wishmatic/genwish/cmd/genwish/version"
	"github.com/wishmatic/genwish/internal/cli"
)

func main() {
	if err := newRoot(cli.DefaultDeps()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRoot(deps cli.Deps) *cobra.Command {
	deps = cli.WithDefaults(deps)

	root := &cobra.Command{
		Use:           "genwish",
		Short:         "Generate images, video, and names from the command line",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.Err)

	config.AddConfigFlag(root, deps)

	root.AddCommand(
		balance.NewCommand(deps),
		config.NewLoginCommand(deps),
		config.NewLogoutCommand(deps),
		config.NewStatusCommand(deps),
		edit.NewCommand(deps),
		gen.NewCommand(deps),
		models.NewCommand(deps),
		util.NewCommand(deps),
		version.NewCommand(deps),
	)

	cli.NormalizeFlagDefaults(root)
	cli.WithConfigDefaults(root, deps)

	return root
}
