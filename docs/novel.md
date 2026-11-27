# NovelAI

`genwish` uses NovelAI for image generation and Anlas balance reporting.

## Authentication

Store a key once:

```sh
genwish login novel
```

`genwish logout novel` deletes the stored key. API keys are login-only; the environment is not read.

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

| Flag                   | Default                  | Purpose                                             |
| ---------------------- | ------------------------ | --------------------------------------------------- |
| `--model`              | `nai-diffusion-5-full`   | NovelAI model id. Must not be empty.                |
| `--negative-prompt`    | provider default         | Things to avoid.                                    |
| `--sampler-name`       | `k_euler_ancestral`      | NovelAI sampler.                                    |
| `--steps`              | `20`                     | Sampling steps.                                     |
| `--width` / `--height` | `0`                      | Output size; 0 uses 512 or the init image's size.   |
| `--cfg-scale`          | `7`                      | Classifier-free guidance scale.                     |
| `--seed`               | `-1`                     | Random seed; -1 is random.                          |
| `--init-image`         | empty                    | Local image to transform.                           |
| `--init-image-url`     | empty                    | Address to transform. Use only one init image flag. |
| `--denoising-strength` | `0.75`                   | How much of the init image to change.               |
| `--noise`              | `0`                      | Extra image noise; only used with an init image.    |
| `--format`             | `output.format` (`webp`) | `png`, `jpeg`, `jxl`, or `webp`.                    |
| `--out`                | `output.path`            | Base output directory.                              |

Images are written under `<out>/images/<date>/<uuid>.<ext>`.

## Models

NovelAI publishes no model-list endpoint, so genwish ships a fixed list:

```sh
genwish models list novel
genwish models list novel --query 4-5 --limit 5
```

The list covers the V5, V4.5, V4, V3, and Furry V3 image families. `--query` filters ids and `--limit` caps the output. `genwish models list novel video` is rejected because NovelAI serves only image models.

## Balance and Anlas

```sh
genwish balance novel
```

Prints `Anlas: N`, and `Tier: ...` when NovelAI reports one.

The balance comes from `GET https://image.novelai.net/user/data` with the API key as a Bearer token. The Anlas total is read from `subscription.trainingStepsLeft`, either as an object with `fixedTrainingStepsLeft` plus `purchasedTrainingSteps`, or as a plain integer on older responses.

## Agent model knowledge

[prompts/novel.md](prompts/novel.md) is the source prompt template for model and sampler knowledge.
