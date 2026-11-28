package edit

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/forge"
)

type bgkillFlags struct {
	model        string
	fullMode     bool
	outputFormat string
	out          string
}

func NewBgkillCommand(deps cli.Deps) *cobra.Command {
	var flags bgkillFlags

	cmd := &cobra.Command{
		Use:   "bgkill <image>",
		Short: "Remove the background from an image with AI (Forge Neo + Birefnet)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return runBgkill(cmd, deps, flags, args[0])
		},
	}

	cmd.Flags().StringVar(&flags.model, "model", "General", "BiRefNet model to load")
	cmd.Flags().BoolVar(
		&flags.fullMode,
		"full-mode",
		false,
		"run the model in fp32 instead of fp16",
	)
	cmd.Flags().StringVar(
		&flags.outputFormat,
		"format",
		"",
		"output image format: png, jpeg, jxl, or webp",
	)
	cli.SetConfigDefault(cmd, "format", "output.format")
	cmd.Flags().StringVar(
		&flags.out,
		"out",
		"",
		"base directory for the foreground image",
	)
	cli.SetConfigDefault(cmd, "out", "output.path")

	return cmd
}

func runBgkill(cmd *cobra.Command, deps cli.Deps, flags bgkillFlags, image string) error {
	cfg, err := cli.ResolveConfig(deps)
	if err != nil {
		return err
	}

	target, err := cli.OutputFormat(cfg, flags.outputFormat)
	if err != nil {
		return fmt.Errorf("bgkill: %w", err)
	}

	resolver := cli.ImageResolver(deps)

	source, err := cli.ReadImage(cmd.Context(), resolver, image)
	if err != nil {
		return fmt.Errorf("bgkill: read image: %w", err)
	}

	out, err := cli.ForgeClient(deps, cfg).Bgkill(
		cmd.Context(),
		forge.BgkillRequest{
			ModelName:  flags.model,
			ImageData:  source,
			IsFullMode: flags.fullMode,
		},
	)
	if err != nil {
		return fmt.Errorf("bgkill: %w", err)
	}

	paths, err := cli.WriteImages(cli.OutputDir(cfg, flags.out), "images", [][]byte{out}, target)
	if err != nil {
		return fmt.Errorf("bgkill: %w", err)
	}

	return cli.PrintPaths(deps.Out, paths)
}
