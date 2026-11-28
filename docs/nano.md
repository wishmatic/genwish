# NanoGPT

`genwish` uses NanoGPT for image and video generation, model listing, and balance.

## Authentication

Store a key once:

```sh
genwish login nano
```

`genwish logout nano` deletes the stored key.

## Generate an image

```sh
genwish gen image nano "a red fox in snow" --out ./media
```

Flags:

- `--model` (`bytedance/seedream-v5.0-pro`): image model id.
- `--n` (provider default): number of images to generate.
- `--size` (provider default): output size or model-specific resolution.
- `--image` (empty): image to transform; a local path, an http(s) URL, or a data URI.
- `--images` (empty): additional input image, same forms, repeatable.
- `--mask` (empty): inpainting mask; a local path, an http(s) URL, or a data URI.
- `--format` (`output.format`): `png`, `jpeg`, `jxl`, or `webp`.
- `--out` (`output.path`): base output directory.
- `--seed` (unset): seed for models that take one.
- `--strength` (unset): how far the output may move from the input image.
- `--guidance-scale` (unset): how closely the model follows the prompt.
- `--num-inference-steps` (unset): denoising steps.

Input images accept a local path, an http(s) URL, a data URI, or raw base64.

Generated images are written under `<out>/images/<date>/<uuid>.<ext>` in the configured format. The
command prints each path and `Cost: $...`.

## Generate a video

```sh
genwish gen video nano "waves at sunset" --out ./media
genwish gen video nano "waves at sunset" --model veo2-video --out ./media
```

`--model` defaults to `alibaba/wan-3.0/text-to-video`.

Flags:

- `--model` (`alibaba/wan-3.0/text-to-video`): video model id.
- `--image` (empty): image to animate; a local path, an http(s) URL, or a data URI. Sets
  `--mode image-to-video` when no mode is given.
- `--negative-prompt` (empty): content to suppress.
- `--duration` (provider default): length in seconds.
- `--aspect-ratio` (provider default): `16:9`, `9:16`, `1:1`, `4:3`, or `3:4`.
- `--resolution` (provider default): `480p`, `720p`, `1080p`, `2k`, or `4k`.
- `--mode` (empty): `text-to-video`, `image-to-video`, `reference-to-video`, or `video-edit`.
- `--seed` (unset): seed for providers that honour one.
- `--generate-audio` (false): add AI audio on models that support it.
- `--video-url` (empty): source video to extend, edit, or upscale.
- `--out` (`output.path`): base output directory.

The result is written under `<out>/videos/<date>/<uuid>.<ext>`, and the command prints the path and
cost.

## Models

```sh
genwish models list nano
genwish models list nano image --query anime --limit 10
genwish models list nano video
genwish models get nano <id> [image|video]
```

Flags: `--query` filters by id, name, description, owner, or tags; `--limit` defaults to 25 and is
capped at 100. With no kind, both image and video models are listed. `get` prints the full entry,
including capabilities, supported parameters, and pricing.

## Balance

```sh
genwish balance nano
```

Prints `USD` and `Nano`.
