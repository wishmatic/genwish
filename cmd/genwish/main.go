package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		Short:         "Generate images and video from the command line",
		Long:          "Generate images and video through NanoGPT, Forge Neo, and NovelAI, and look up booru tags.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.SetIn(deps.In)
	root.SetOut(deps.Out)
	root.SetErr(deps.Err)

	root.PersistentFlags().StringVar(&deps.Flags.ConfigPath, "config-file", "",
		"path to config.toml (default ~/.genwish/config.toml)")

	root.AddCommand(
		newLoginCommand(deps),
		newLogoutCommand(deps),
		newModelsCommand(deps),
		newBalanceCommand(deps),
		newGenCommand(deps),
		newEditCommand(deps),
		newUtilCommand(deps),
	)

	return root
}
