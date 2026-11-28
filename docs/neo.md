# Forge Neo

`genwish` uses a local Forge Neo instance for image generation and background removal, and runs the
crop and convert edits itself.

Only [Haoming02's Neo fork](https://github.com/Haoming02/sd-webui-forge-classic) (the `neo`
branch) is supported. Launch Forge with `--api`, then set `neo.url` in the config file to its API.
The default is `http://127.0.0.1:7860`.

## Generate an image

```sh
genwish gen image neo "a red fox in snow" --model krea --out ./media
```

`--model` names a model in the registry. Its checkpoint, encoders, preset, sampler, scheduler,
steps, cfg, and hi-res setting are used unless a flag overrides them. A name that is not configured
is used as the checkpoint filename as given.

Flags:

- `--model` (required): model id from the registry, or a checkpoint filename.
- `--forge-preset` (empty unless the model sets one): Forge UI preset to apply.
- `--vae-and-text-model` (the model's encoders): VAE or text encoder filename, repeatable.
- `--negative-prompt` (empty): things to avoid.
- `--sampler-name` (the model's sampler, else `DPM++ 2M`): sampler.
- `--scheduler` (the model's scheduler, else `Automatic`): scheduler.
- `--steps` (the model's steps, else `20`): sampling steps.
- `--width` / `--height` (`0`): output size; 0 uses the init image or the model default.
- `--cfg-scale` (the model's cfg, else `7`): classifier-free guidance scale.
- `--seed` (`-1`): random seed; -1 is random.
- `--init-image` (empty): image to transform; a local path, an http(s) URL, or a data URI.
- `--denoising-strength` (`0`): how much of the init image to change.
- `--enable-hr` (the model's hi-res setting, else false): enable hi-res second-pass upscaling.
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

Model metadata lives in a TOML registry, `~/.genwish/models.toml` by default. Set
`neo.models_file` to read another path. Each model is a table under `[neo]`:

```toml
[neo.krea]
checkpoint = "hyphoria_v002.safetensors"
encoders = ["qwen_3_06b_base.safetensors", "qwen_3_4b_base.safetensors"]
preset = "krea"
sizes = ["1024x1024", "832x1216"]

[neo.krea.samplers]
default = "Euler"
options = ["Euler", "Euler a"]

[neo.krea.schedulers]
default = "Automatic"
options = ["Automatic", "Karras"]

[neo.krea.steps]
default = 20
min = 10
max = 40

[neo.krea.cfg]
default = 5.0
min = 1.0
max = 15.0

[neo.krea.hires]
enabled = true

[neo.krea.instructions]
usage = """
Load the krea preset, then pick a size and copy one style prefix.
Keep the negative prompt short.
"""
style = """
Anime: `masterpiece, best quality, anime screencap, 1girl, ...`
Painterly: `illustration, painterly, thick brush strokes, ...`
"""
```

Every field is optional. `checkpoint` and each encoder are verbatim file names with their
extensions. `samplers` and `schedulers` take a `default` and `options`, and the default must
appear in the options. `steps` and `cfg` take `default`, `min`, and `max`: a bound requires both
bounds and a default between them. `sizes` are `WxH` strings. `hires.enabled` defaults to false.

`instructions.usage` says how to use the model. `instructions.style` lists the styles the model
supports, with prompt examples for each. Write both as TOML multi-line strings so each line
survives; `models get` prints them line by line. [example/models.toml](../example/models.toml)
has two complete models.

```sh
genwish models list neo
genwish models get neo krea
```

`list` prints the id, checkpoint, and preset; `--query` filters those and `--limit` caps the
output. `get` prints the whole entry, including the instructions.

A missing default registry means no models are configured; a missing `neo.models_file` path is an
error. Generation does not require a registry entry, so any checkpoint can be used.
