package edit

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/format"
)

type convertFlags struct {
	format string
	out    string
}

func NewConvertCommand(deps cli.Deps) *cobra.Command {
	var flags convertFlags

	cmd := &cobra.Command{
		Use:   "convert <image>",
		Short: "Re-encode an image as png, jpeg, jxl, or webp",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return convertRun(cmd, deps, flags, args[0])
		},
	}

	cmd.Flags().
		StringVar(&flags.format, "format", "", "format to convert to: png, jpeg, jxl, or webp")
	cmd.Flags().StringVar(&flags.out, "out", "",
		"base directory for the converted image (default: output.path)")

	return cmd
}

func convertRun(cmd *cobra.Command, deps cli.Deps, flags convertFlags, image string) error {
	target, err := convertTarget(flags.format)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	cfg, err := cli.ResolveConfig(deps)
	if err != nil {
		return err
	}

	resolver := cli.ImageResolver(deps)

	source, err := cli.ReadImage(cmd.Context(), resolver, image)
	if err != nil {
		return fmt.Errorf("convert: read image: %w", err)
	}

	out, err := format.Reformat(source, target)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	path, err := cli.WriteImage(cli.OutputDir(cfg, flags.out), "images", out, target.MediaType())
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	return cli.PrintPaths(deps.Out, []string{path})
}

func convertTarget(value string) (format.Format, error) {
	if value == "" {
		return "", fmt.Errorf("--format is required")
	}

	return format.Parse(value)
}
