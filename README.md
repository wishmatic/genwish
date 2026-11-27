# genwish

Bespoke personal CLI for image and video generation, plus Danbooru tag lookup.

## Install

```sh
go install github.com/wishmatic/genwish/cmd/genwish@latest
```

## Providers

| Provider | Generation                   | Balance      | API key                |
| -------- | ---------------------------- | ------------ | ---------------------- |
| `nano`   | image, video                 | USD and Nano | `genwish login nano`   |
| `neo`    | image, edit, convert, bgkill | none         | none, uses local Forge |
| `novel`  | image                        | Anlas        | `genwish login novel`  |

## Quick start

```sh
genwish login nano
genwish gen image nano "a red fox in snow" --out ./media
genwish gen video nano "waves at sunset" --model veo2-video --out ./media

genwish gen image neo "a red fox" --model sd.safetensors --forge-preset default --out ./media
genwish edit square ./logo.png --out ./media

genwish login novel
genwish gen image novel "a cat" --out ./media
genwish balance novel

genwish util booru blue_hair --limit 5
```

## Documentation

- [Command and provider reference](docs/README.md)
- [Configuration](docs/configuration.md)
- [NanoGPT](docs/nano.md)
- [Forge Neo](docs/neo.md)
- [NovelAI](docs/novel.md)
- [Booru](docs/booru.md)

## License

Apache-2.0. See [LICENSE](LICENSE).

This project is not affiliated with or endorsed by NanoGPT, NovelAI, Danbooru, AUTOMATIC1111,
Illyasviel, Haoming02, or any other third party.
