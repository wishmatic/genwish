package cli

import (
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/wishmatic/genwish/internal/config"
)

const configDefaultKey = "genwish-config-default"

var trailingDefault = regexp.MustCompile(`\s*\(default:[^)]*\)$`)

func SetConfigDefault(cmd *cobra.Command, flagName, key string) {
	flag := cmd.Flags().Lookup(flagName)
	if flag == nil {
		return
	}

	if flag.Annotations == nil {
		flag.Annotations = map[string][]string{}
	}

	flag.Annotations[configDefaultKey] = []string{key}
}

func WithConfigDefaults(root *cobra.Command, deps Deps) {
	help := root.HelpFunc()

	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cfg := config.Defaults()
		if resolved, err := ResolveConfig(deps); err == nil {
			cfg = resolved
		}

		ApplyConfigDefaults(cmd, cfg)
		help(cmd, args)
	})
}

func ApplyConfigDefaults(cmd *cobra.Command, cfg config.Config) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		keys := flag.Annotations[configDefaultKey]
		if len(keys) == 0 {
			return
		}

		value := configDefault(cfg, keys[0])
		if value == "" {
			return
		}

		flag.Usage = trailingDefault.ReplaceAllString(flag.Usage, "")
		flag.Usage += " (default: " + value + ")"
	})
}

func configDefault(cfg config.Config, key string) string {
	switch key {
	case "output.format":
		return cfg.Output.Format
	case "output.path":
		return cfg.Output.Path
	default:
		return ""
	}
}
