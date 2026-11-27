# Forge Neo

`genwish` uses a local Forge Neo instance for image generation and background removal, and runs the crop and convert edits itself.

Only [Haoming02's Neo fork](https://github.com/Haoming02/sd-webui-forge-classic/tree/neo#stable-diffusion-webui-forge---neo) is supported. Launch Forge with `--api`, then set `neo.url` in the config file to its API. The default is `http://127.0.0.1:7860`.

No Forge installation or API support is provided here.

## Generate an image

```sh
genwish gen image neo "a red fox in snow" --model sd.safetensors --forge-preset default --out ./media
```

`--model` and `--forge-preset` are required.

Flags:

| Flag                     | Default                  | Purpose                                                  |
| ------------------------ | ------------------------ | -------------------------------------------------------- |
| `--model`                | required                 | Checkpoint filename to load.                             |
| `--forge-preset`         | required                 | Forge UI preset to apply.                                |
| `--vae-and-text-model`   | empty                    | VAE or text encoder filename, repeatable.                |
| `--negative-prompt`      | empty                    | Things to avoid.                                         |
| `--sampler-name`         | `DPM++ 2M`               | Sampler.                                                 |
| `--scheduler`            | `Automatic`              | Scheduler.                                               |
| `--steps`                | `20`                     | Sampling steps.                                          |
| `--width` / `--height`   | `0`                      | Output size; 0 uses the init image or the model default. |
| `--cfg-scale`            | `7`                      | Classifier-free guidance scale.                          |
| `--seed`                 | `-1`                     | Random seed; -1 is random.                               |
| `--init-image`           | empty                    | Local image to transform.                                |
| `--init-image-url`       | empty                    | Address to transform. Use only one init image flag.      |
| `--denoising-strength`   | `0`                      | How much of the init image to change.                    |
| `--enable-hr`            | false                    | Enable hi-res second-pass upscaling.                     |
| `--hr-scale`             | `0`                      | Hi-res upscaling factor.                                 |
| `--hr-upscaler`          | empty                    | Hi-res upscaler.                                         |
| `--hr-second-pass-steps` | `0`                      | Steps for the hi-res pass.                               |
| `--hr-cfg`               | `0`                      | CFG scale for the hi-res pass.                           |
| `--format`               | `output.format` (`webp`) | `png`, `jpeg`, `jxl`, or `webp`.                         |
| `--out`                  | `output.path`            | Base output directory.                                   |

Images are written under `<out>/images/<date>/<uuid>.<ext>`.

## Edit

Every edit command takes an image as a local path or an address the resolver can read: an http(s) URL, a data URI, or raw base64.

```sh
genwish edit trim ./photo.png
genwish edit square ./logo.png --padding 8 --format webp
genwish edit circle ./logo.png
genwish edit convert ./photo.png --format jxl
genwish edit bgkill ./photo.jpg --model General
```

### trim, square, circle

Flags: `--padding` (transparent padding in pixels; square has a nonzero default), `--format`, `--out`.

### convert

`--format` is required. Flags: `--format`, `--out`.

### bgkill

Removes the background through the [`sd-webui-birefnet`](https://github.com/dimitribarbot/sd-webui-birefnet) extension.

Flags: `--model` (default `General`), `--full-mode` (fp32 instead of fp16), `--format`, `--out`.

## Models

`genwish` does not query the Forge API for models. It reads a local file:

Set `neo.models_file` in the config file, then run:

```sh
genwish models list neo
```

The file is newline-delimited. Blank lines and lines starting with `#` are ignored. When `neo.models_file` is empty, the list is empty. `--query` filters names and `--limit` caps the output.

`genwish models list neo video` is rejected because Neo serves only image models.

## Agent model knowledge

The agent needs to know checkpoint, VAE, text encoder, and upscaler names. No list ships here. [prompts/neo.md](prompts/neo.md) is the source prompt template.
