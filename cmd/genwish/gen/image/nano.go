package genimage

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/wishmatic/genwish/internal/cli"
	"github.com/wishmatic/genwish/internal/format"
	"github.com/wishmatic/genwish/internal/nanogpt"
)

const defaultImageModel = "bytedance/seedream-v5.0-pro"

func NewNanoCommand(deps cli.Deps) *cobra.Command {
	var (
		model        string
		size         string
		image        string
		images       []string
		mask         string
		outputFormat string
		out          string
	)

	cmd := &cobra.Command{
		Use:   "nano <prompt...>",
		Short: "Generate images through NanoGPT",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			prompt, err := cli.PromptFromArgs(args)
			if err != nil {
				return err
			}

			if strings.TrimSpace(model) == "" {
				model = defaultImageModel
			}

			cfg, err := cli.ResolveConfig(deps)
			if err != nil {
				return err
			}

			target, err := cli.OutputFormat(cfg, outputFormat)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			resolver := cli.ImageResolver(deps)

			imageDataURL, err := cli.ResolveDataURL(cmd.Context(), resolver, image)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			maskDataURL, err := cli.ResolveDataURL(cmd.Context(), resolver, mask)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			imageDataURLs, err := cli.ResolveDataURLs(cmd.Context(), resolver, images)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			client, err := cli.NanoClient(deps)
			if err != nil {
				return err
			}

			result, err := client.GenerateImages(
				cmd.Context(),
				nanogpt.ImageRequest{
					Model:  model,
					Prompt: prompt,
					N:      cli.OptionalInt(cmd, "n"),
					Size:   size,

					ResponseFormat: nanogpt.ResponseFormatBase64,

					ImageDataURL:  imageDataURL,
					ImageDataURLs: imageDataURLs,
					MaskDataURL:   maskDataURL,

					Strength:          cli.OptionalFloat(cmd, "strength"),
					GuidanceScale:     cli.OptionalFloat(cmd, "guidance-scale"),
					NumInferenceSteps: cli.OptionalInt(cmd, "num-inference-steps"),
					Seed:              cli.OptionalInt(cmd, "seed"),
				},
			)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			assets, err := readyImages(cmd, result.Images)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			paths, err := writeImages(cli.OutputDir(cfg, out), assets, target)
			if err != nil {
				return fmt.Errorf("gen image: %w", err)
			}

			for _, path := range paths {
				fmt.Fprintln(deps.Out, path)
			}

			fmt.Fprintf(deps.Out, "Cost: %s\n", cli.FormatUSD(result.CostUSD))

			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", defaultImageModel, "image model id")
	cmd.Flags().Int("n", 0, "number of images to generate")
	cmd.Flags().StringVar(&size, "size", "", "output size or model-specific resolution")
	cmd.Flags().StringVar(&image, "image", "", "image to transform")
	cmd.Flags().StringArrayVar(&images, "images", nil, "additional input image")
	cmd.Flags().StringVar(&mask, "mask", "", "inpainting mask")
	cmd.Flags().StringVar(
		&outputFormat,
		"format",
		"",
		"output format: png, jpeg, jxl, or webp",
	)
	cli.SetConfigDefault(cmd, "format", "output.format")
	cmd.Flags().StringVar(
		&out,
		"out",
		"",
		"base directory for generated images",
	)
	cli.SetConfigDefault(cmd, "out", "output.path")
	cmd.Flags().Int("seed", 0, "seed for the models that take one")
	cmd.Flags().Float64("strength", 0, "how far the output may move from the input image")
	cmd.Flags().Float64("guidance-scale", 0, "how closely the model follows the prompt")
	cmd.Flags().Int("num-inference-steps", 0, "denoising steps")

	return cmd
}

type readyImage struct {
	data      []byte
	mediaType string
}

func readyImages(cmd *cobra.Command, images []nanogpt.GeneratedImage) ([]readyImage, error) {
	ready := make([]readyImage, 0, len(images))

	for index, image := range images {
		if image.URL != "" {
			data, mediaType, err := cli.DownloadAsset(cmd.Context(), image.URL, "")
			if err != nil {
				return nil, fmt.Errorf("image %d: %w", index+1, err)
			}

			ready = append(ready, readyImage{data: data, mediaType: mediaType})

			continue
		}

		mediaType, ok := cli.DetectAssetMediaType(image.Data)
		if !ok {
			return nil, fmt.Errorf("image %d is not a supported image type", index+1)
		}

		ready = append(ready, readyImage{data: image.Data, mediaType: mediaType})
	}

	return ready, nil
}

func writeImages(base string, images []readyImage, target format.Format) ([]string, error) {
	paths := make([]string, 0, len(images))

	for _, image := range images {
		data := image.data

		if image.mediaType != target.MediaType() {
			converted, err := format.Convert(data, target)
			if err != nil {
				return nil, err
			}

			data = converted
		}

		path, err := cli.WriteImage(base, "images", data, target.MediaType())
		if err != nil {
			return nil, err
		}

		paths = append(paths, path)
	}

	return paths, nil
}
