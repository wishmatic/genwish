# Documentation

`genwish` is a CLI for image, video, and name generation, with three providers and a Danbooru
tag utility.

## Command tree

```text
genwish
  gen | generate
    image
      nano <prompt...>
      neo <prompt...>
      novel <prompt...>
    name
    video
      nano <prompt...>
  edit
    trim <image>
    square <image>
    circle <image>
    convert <image>
    bgkill <image>
  models
    list <provider> [image|video]
    get <provider> <id> [image|video]
  balance <provider>
  login <provider>
  logout <provider>
  status
  util
    booru <search...>
```

Provider names are case-insensitive. `nano` also accepts `nano-gpt`, `nanogpt`, and `nano gpt`.
`novel` also accepts `novel-ai`, `novel ai`, and `novelai`.

## Where to look

- [configuration.md](configuration.md): the TOML config file.
- [nano.md](nano.md): NanoGPT image and video generation.
- [neo.md](neo.md): Forge Neo generation, editing, and model listing.
- [novel.md](novel.md): NovelAI generation, balance, and Anlas.
- [booru.md](booru.md): tag search semantics.
- [name.md](name.md): local name generation.
