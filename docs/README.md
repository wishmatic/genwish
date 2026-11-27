# Documentation

`genwish` is a CLI-only generation tool with three providers and one utility.

## Command tree

```text
genwish
  gen | generate
    image
      nano <prompt...>
      neo <prompt...>
      novel <prompt...>
    video
      nano <prompt...>
      neo <prompt...>      unsupported for now
      novel <prompt...>    unsupported for now
  edit
    trim <image>
    square <image>
    circle <image>
    convert <image>
    bgkill <image>
  models
    list <provider> [image|video]
  balance <provider>
  login <provider>
  logout <provider>
  util
    booru <search...>
```

Provider names are case-insensitive. `nano` also accepts `nano-gpt`, `nanogpt`, and `nano gpt`. `novel` also accepts `novel-ai`, `novel ai`, and `novelai`.

## Where to look

- [configuration.md](configuration.md): the TOML config file and environment variables.
- [migration.md](migration.md): old commands and settings mapped to genwish.
- [nano.md](nano.md): NanoGPT image and video generation.
- [neo.md](neo.md): Forge Neo generation, editing, and model listing.
- [novel.md](novel.md): NovelAI generation, balance, and Anlas.
- [booru.md](booru.md): tag search semantics.
- [prompts/](prompts): agent prompts for Forge Neo and NovelAI model knowledge.
- [reference/](reference): the original READMEs and `.env.example` files, copied verbatim.
