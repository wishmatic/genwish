package config

import (
	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
)

func AddConfigFlag(root *cobra.Command, deps cli.Deps) {
	root.PersistentFlags().StringVar(
		&deps.Flags.ConfigPath,
		"config-file",
		"",
		"path to config.toml (default: ~/.genwish/config.toml)",
	)

	root.PersistentPreRunE = func(*cobra.Command, []string) error {
		return cli.EnsureDefaultConfig(deps)
	}
}
