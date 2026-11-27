# NovelAI

`genwish` uses NovelAI for image generation and Anlas balance reporting.

## Authentication

Store a key once:

```sh
genwish login novel
```

`genwish logout novel` deletes the stored key.

Calls spend Anlas on the account the key belongs to.

## Generate an image

```sh
genwish gen image novel "a cat" --out ./media
genwish gen image novel "a cat" --init-image ./photo.png --out ./media
```

Defaults:

- Model: `nai-diffusion-5-full`.
- Steps: `20`.
- CFG scale: `7`.
- Sampler: `k_euler_ancestral`.
- Seed: `-1` (random).
- Denoising strength: `0.75`.
- Format: `output.format` (`webp`).

Flags:

- `--model` (`nai-diffusion-5-full`): NovelAI model id; must not be empty.
- `--negative-prompt` (provider default): things to avoid.
- `--sampler-name` (`k_euler_ancestral`): NovelAI sampler.
- `--steps` (`20`): sampling steps.
- `--width` / `--height` (`0`): output size; 0 uses 512 or the init image's size.
- `--cfg-scale` (`7`): classifier-free guidance scale.
- `--seed` (`-1`): random seed; -1 is random.
- `--init-image` (empty): local image to transform.
- `--init-image-url` (empty): address to transform; use only one init image flag.
- `--denoising-strength` (`0.75`): how much of the init image to change.
- `--noise` (`0`): extra image noise; only used with an init image.
- `--format` (`output.format`, default `webp`): `png`, `jpeg`, `jxl`, or `webp`.
- `--out` (`output.path`): base output directory.

Images are written under `<out>/images/<date>/<uuid>.<ext>`.

## Models

The model list is fixed:

```sh
genwish models list novel
genwish models list novel --query 4-5 --limit 5
genwish models get novel nai-diffusion-5-full
```

The list covers the V5, V4.5, V4, V3, and Furry V3 image families. `--query` filters ids and
`--limit` caps the output. `get` prints the id and kind of one model.

## Balance and Anlas

```sh
genwish balance novel
```

Prints `Anlas: N`, and the tier (`paper`, `tablet`, `scroll`, or `opus`) when NovelAI reports one.

The balance comes from `GET https://image.novelai.net/user/data` with the API key as a Bearer
token. The Anlas total is read from `subscription.trainingStepsLeft`, either as an object with
`fixedTrainingStepsLeft` plus `purchasedTrainingSteps`, or as a plain integer on older responses.
The tier is NovelAI's numeric level: 0 Paper, 1 Tablet, 2 Scroll, 3 Opus.
