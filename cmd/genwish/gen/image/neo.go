package genimage

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/forge"
	"github.com/wishmatic/genwish/internal/resolve"
)

const (
	forgeDefaultSampler   = "DPM++ 2M"
	forgeDefaultScheduler = "Automatic"
)

type forgeFlags struct {
	model             string
	preset            string
	vaeAndTextModels  []string
	negativePrompt    string
	sampler           string
	scheduler         string
	steps             int
	width             int
	height            int
	cfgScale          float64
	seed              int
	initImage         string
	denoisingStrength float64
	enableHR          bool
	hrScale           float64
	hrUpscaler        string
	hrSecondPassSteps int
	hrCFG             float64
	outputFormat      string
	out               string
}

func NewNeoCommand(deps cli.Deps) *cobra.Command {
	var flags forgeFlags

	cmd := &cobra.Command{
		Use:   "neo <prompt...>",
		Short: "Generate an image with the local Forge Neo instance",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return forgeRun(cmd, deps, flags, args)
		},
	}

	cmd.Flags().StringVar(
		&flags.model,
		"model",
		"",
		"model id from models.toml, or a checkpoint filename to load",
	)
	cmd.Flags().StringVar(
		&flags.preset,
		"forge-preset",
		"",
		"Forge UI preset; defaults to the model's configured preset",
	)
	cmd.Flags().StringArrayVar(
		&flags.vaeAndTextModels,
		"vae-and-text-model",
		nil,
		"VAE or text encoder filename to load; repeat for more; defaults to the model's encoders",
	)
	cmd.Flags().StringVar(
		&flags.negativePrompt,
		"negative-prompt",
		"",
		"things to avoid in the generated image",
	)
	cmd.Flags().StringVar(
		&flags.sampler,
		"sampler-name",
		forgeDefaultSampler,
		"sampler to use; defaults to the model's configured sampler",
	)
	cmd.Flags().StringVar(
		&flags.scheduler,
		"scheduler",
		forgeDefaultScheduler,
		"scheduler to use; defaults to the model's configured scheduler",
	)
	cmd.Flags().IntVar(
		&flags.steps,
		"steps",
		20,
		"number of sampling steps; defaults to the model's configured steps",
	)
	cmd.Flags().IntVar(&flags.width, "width", 0, "output width in pixels")
	cmd.Flags().IntVar(&flags.height, "height", 0, "output height in pixels")
	cmd.Flags().Float64Var(
		&flags.cfgScale,
		"cfg-scale",
		7,
		"classifier-free guidance scale; defaults to the model's configured value",
	)
	cmd.Flags().IntVar(&flags.seed, "seed", -1, "random seed")
	cmd.Flags().StringVar(
		&flags.initImage,
		"init-image",
		"",
		"image to transform instead of generating from scratch: a local path, http(s) URL, or data URI",
	)
	cmd.Flags().Float64Var(
		&flags.denoisingStrength,
		"denoising-strength",
		0,
		"how much of the init image to change",
	)
	cmd.Flags().BoolVar(
		&flags.enableHR,
		"enable-hr",
		false,
		"enable hi-res second-pass upscaling; defaults to the model's configured setting",
	)
	cmd.Flags().Float64Var(&flags.hrScale, "hr-scale", 0, "hi-res upscaling factor")
	cmd.Flags().StringVar(&flags.hrUpscaler, "hr-upscaler", "", "hi-res upscaler to use")
	cmd.Flags().IntVar(
		&flags.hrSecondPassSteps,
		"hr-second-pass-steps",
		0,
		"steps for the hi-res second pass",
	)
	cmd.Flags().Float64Var(&flags.hrCFG, "hr-cfg", 0, "CFG scale for the hi-res second pass")
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
		"base directory for the generated image",
	)
	cli.SetConfigDefault(cmd, "out", "output.path")

	return cmd
}

func forgeRun(cmd *cobra.Command, deps cli.Deps, flags forgeFlags, args []string) error {
	prompt, err := cli.PromptFromArgs(args)
	if err != nil {
		return err
	}

	if flags.model == "" {
		return fmt.Errorf("forge: --model is required")
	}

	flags, err = resolveNeoFlags(cmd, deps, flags)
	if err != nil {
		return fmt.Errorf("forge: %w", err)
	}

	if flags.preset == "" {
		return fmt.Errorf("forge: --forge-preset is required")
	}

	cfg, err := cli.ResolveConfig(deps)
	if err != nil {
		return err
	}

	target, err := cli.OutputFormat(cfg, flags.outputFormat)
	if err != nil {
		return fmt.Errorf("forge: %w", err)
	}

	resolver := cli.ImageResolver(deps)

	initImage, err := readInitImage(cmd.Context(), resolver, flags.initImage)
	if err != nil {
		return fmt.Errorf("forge: %w", err)
	}

	client := cli.ForgeClient(deps, cfg)

	images, err := client.Generate(cmd.Context(), flags.request(prompt, initImage))
	if err != nil {
		return fmt.Errorf("forge: %w", err)
	}

	paths, err := cli.WriteImages(cli.OutputDir(cfg, flags.out), "images", images, target)
	if err != nil {
		return fmt.Errorf("forge: %w", err)
	}

	return cli.PrintPaths(deps.Out, paths)
}

// resolveNeoFlags fills in the model's configured defaults. An unconfigured id is treated as the
// checkpoint filename itself, so a bare checkpoint still generates without a registry entry.
func resolveNeoFlags(cmd *cobra.Command, deps cli.Deps, flags forgeFlags) (forgeFlags, error) {
	registry, err := cli.NeoRegistry(deps)
	if err != nil {
		return forgeFlags{}, err
	}

	model, ok := registry.Match(flags.model)
	if !ok {
		return flags, nil
	}

	if model.Checkpoint != "" {
		flags.model = model.Checkpoint
	}

	if len(flags.vaeAndTextModels) == 0 {
		flags.vaeAndTextModels = model.Encoders
	}

	if flags.preset == "" {
		flags.preset = model.Preset
	}

	if !cmd.Flags().Changed("sampler-name") && model.Samplers.Default != "" {
		flags.sampler = model.Samplers.Default
	}

	if !cmd.Flags().Changed("scheduler") && model.Schedulers.Default != "" {
		flags.scheduler = model.Schedulers.Default
	}

	if !cmd.Flags().Changed("steps") && model.Steps.Default != nil {
		flags.steps = *model.Steps.Default
	}

	if !cmd.Flags().Changed("cfg-scale") && model.CFG.Default != nil {
		flags.cfgScale = *model.CFG.Default
	}

	if !cmd.Flags().Changed("enable-hr") {
		flags.enableHR = model.Hires.Enabled
	}

	return flags, nil
}

func (f forgeFlags) request(prompt string, initImage []byte) forge.GenerateRequest {
	return forge.GenerateRequest{
		InitImage: initImage,

		Checkpoint:             f.model,
		ForgePreset:            f.preset,
		ForgeAdditionalModules: f.vaeAndTextModels,

		Prompt:         prompt,
		NegativePrompt: f.negativePrompt,

		Sampler:   f.sampler,
		Scheduler: f.scheduler,
		Steps:     f.steps,

		Width:    f.width,
		Height:   f.height,
		CFGScale: f.cfgScale,

		DenoisingStrength: f.denoisingStrength,

		Seed: f.seed,

		EnableHR:          f.enableHR,
		HRScale:           f.hrScale,
		HRUpscaler:        f.hrUpscaler,
		HRSecondPassSteps: f.hrSecondPassSteps,
		HRCFGScale:        f.hrCFG,
	}
}

func readInitImage(ctx context.Context, resolver *resolve.Client, address string) ([]byte, error) {
	if address == "" {
		return nil, nil
	}

	return cli.ReadImage(ctx, resolver, address)
}
