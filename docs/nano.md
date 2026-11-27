# NanoGPT

`genwish` uses NanoGPT for image and video generation, model listing, and balance.

## Authentication

Store a key once:

```sh
genwish login nano
```

`genwish logout nano` deletes the stored key. API keys are login-only; the environment is not read.

## Generate an image

```sh
genwish gen image nano "a red fox in snow" --out ./media
```

Flags:

| Flag                    | Default          | Purpose                                           |
| ----------------------- | ---------------- | ------------------------------------------------- |
| `--model`               | NanoGPT default  | Image model id.                                   |
| `--n`                   | provider default | Number of images to generate.                     |
| `--size`                | provider default | Output size or model-specific resolution.         |
| `--image`               | empty            | Image to transform.                               |
| `--images`              | empty            | Additional input image, repeatable.               |
| `--mask`                | empty            | Inpainting mask.                                  |
| `--format`              | `output.format`  | `png`, `jpeg`, `jxl`, or `webp`.                  |
| `--out`                 | `output.path`    | Base output directory.                            |
| `--seed`                | unset            | Seed for models that take one.                    |
| `--strength`            | unset            | How far the output may move from the input image. |
| `--guidance-scale`      | unset            | How closely the model follows the prompt.         |
| `--num-inference-steps` | unset            | Denoising steps.                                  |

Input images accept an http(s) URL, a data URI, or raw base64.

Generated images are written under `<out>/images/<date>/<uuid>.<ext>` in the configured format. The command prints each path and `Cost: $...`.

## Generate a video

```sh
genwish gen video nano "waves at sunset" --model veo2-video --out ./media
```

`--model` is required.

Flags:

| Flag                | Default          | Purpose                                                                   |
| ------------------- | ---------------- | ------------------------------------------------------------------------- |
| `--model`           | required         | Video model id.                                                           |
| `--image`           | empty            | Image to animate; sets `--mode image-to-video` when no mode is given.     |
| `--negative-prompt` | empty            | Content to suppress.                                                      |
| `--duration`        | provider default | Length in seconds.                                                        |
| `--aspect-ratio`    | provider default | `16:9`, `9:16`, `1:1`, `4:3`, or `3:4`.                                   |
| `--resolution`      | provider default | `480p`, `720p`, `1080p`, `2k`, or `4k`.                                   |
| `--mode`            | empty            | `text-to-video`, `image-to-video`, `reference-to-video`, or `video-edit`. |
| `--seed`            | unset            | Seed for providers that honour one.                                       |
| `--generate-audio`  | false            | Add AI audio on models that support it.                                   |
| `--video-url`       | empty            | Source video to extend, edit, or upscale.                                 |
| `--out`             | `output.path`    | Base output directory.                                                    |

The result is written under `<out>/videos/<date>/<uuid>.<ext>`, and the command prints the path and cost.

## Models

```sh
genwish models list nano
genwish models list nano image --query anime --limit 10
genwish models list nano video
```

Flags: `--query` filters by id, name, description, owner, or tags; `--limit` defaults to 25 and is capped at 100. With no kind, both image and video models are listed.

## Balance

```sh
genwish balance nano
```

Prints `USD` and `Nano`.
