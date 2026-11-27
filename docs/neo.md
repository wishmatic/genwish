# Forge Neo

`genwish` uses a local Forge Neo instance for image generation and background removal, and runs the
crop and convert edits itself.

Only [Haoming02's Neo fork](https://github.com/Haoming02/sd-webui-forge-classic) (the `neo`
branch) is supported. Launch Forge with `--api`, then set `neo.url` in the config file to its API.
The default is `http://127.0.0.1:7860`.

## Generate an image

```sh
genwish gen image neo "a red fox in snow" --model sd.safetensors \
  --forge-preset default --out ./media
```

`--model` and `--forge-preset` are required.

Flags:

- `--model` (required): checkpoint filename to load.
- `--forge-preset` (required): Forge UI preset to apply.
- `--vae-and-text-model` (empty): VAE or text encoder filename, repeatable.
- `--negative-prompt` (empty): things to avoid.
- `--sampler-name` (`DPM++ 2M`): sampler.
- `--scheduler` (`Automatic`): scheduler.
- `--steps` (`20`): sampling steps.
- `--width` / `--height` (`0`): output size; 0 uses the init image or the model default.
- `--cfg-scale` (`7`): classifier-free guidance scale.
- `--seed` (`-1`): random seed; -1 is random.
- `--init-image` (empty): local image to transform.
- `--init-image-url` (empty): address to transform; use only one init image flag.
- `--denoising-strength` (`0`): how much of the init image to change.
- `--enable-hr` (false): enable hi-res second-pass upscaling.
- `--hr-scale` (`0`): hi-res upscaling factor.
- `--hr-upscaler` (empty): hi-res upscaler.
- `--hr-second-pass-steps` (`0`): steps for the hi-res pass.
- `--hr-cfg` (`0`): CFG scale for the hi-res pass.
- `--format` (`output.format`, default `webp`): `png`, `jpeg`, `jxl`, or `webp`.
- `--out` (`output.path`): base output directory.

Images are written under `<out>/images/<date>/<uuid>.<ext>`.

## Edit

Every edit command takes an image as a local path or an address the resolver can read: an http(s)
URL, a data URI, or raw base64.

```sh
genwish edit trim ./photo.png
genwish edit square ./logo.png --padding 8 --format webp
genwish edit circle ./logo.png
genwish edit convert ./photo.png --format jxl
genwish edit bgkill ./photo.jpg --model General
```

### trim, square, circle

Flags: `--padding` (transparent padding in pixels; square has a nonzero default), `--format`,
`--out`.

### convert

`--format` is required. Flags: `--format`, `--out`.

### bgkill

Removes the background through the
[`sd-webui-birefnet`](https://github.com/dimitribarbot/sd-webui-birefnet) extension.

Flags: `--model` (default `General`), `--full-mode` (fp32 instead of fp16), `--format`, `--out`.

## Models

Model names come from a local file. Set `neo.models_file` in the config file, then run:

```sh
genwish models list neo
```

The file is newline-delimited. Blank lines and lines starting with `#` are ignored. `--query`
filters names and `--limit` caps the output.
