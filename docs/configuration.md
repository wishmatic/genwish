# Configuration

`genwish` reads a TOML config file. Credentials come from `genwish login`.

## Config file

The default path is `~/.genwish/config.toml`. Any command creates an empty file there if it is
missing. Pass `--config-file <path>` to read a different file; that file must exist. A complete
example lives in [example/config.toml](../example/config.toml).

- The empty file is valid; omitted keys keep their built-in defaults.
- A missing `--config-file` path is an error.
- Malformed TOML and unknown keys are errors.

Command flags override the file, which overrides the built-in defaults.

```toml
[output]
format = "webp"        # png, jpeg, jxl, or webp
path = "~/.genwish/out"

[neo]
url = "http://127.0.0.1:7860"
models_file = "~/.genwish/models.toml"

[booru.rate_limit]
rps = 1.0
burst = 1
timeout = 20           # seconds
```

- `output.format` (string, default `webp`): `png`, `jpeg`, `jxl`, or `webp`.
- `output.path` (string, default `~/.genwish/out`): base directory for generated and edited media.
- `neo.url` (string, default `http://127.0.0.1:7860`): Forge Neo API base URL.
- `neo.models_file` (string, default `~/.genwish/models.toml`): TOML model registry for
  `models list neo` and `models get neo`. The format is described in [neo.md](neo.md).
- `booru.rate_limit.rps` (float, default `1`): Danbooru requests per second.
- `booru.rate_limit.burst` (integer, default `1`): Danbooru burst size.
- `booru.rate_limit.timeout` (integer, default `20`): Danbooru request timeout in seconds.

A leading `~` in `output.path` or `neo.models_file` is expanded.

Every image is converted to `output.format` when the provider returns a different type.

The tag page limit is fixed at 100 and the largest tag offset at 1000.

## Output layout

Every image and video is written under the configured base directory:

```text
<output.path>/images/<YYYY-MM-DD>/<uuid>.<ext>
<output.path>/videos/<YYYY-MM-DD>/<uuid>.<ext>
```

`--out <path>` replaces `output.path` for that command. The date is UTC and the file name is a
UUID. Commands print each written path.

## Credentials

```sh
genwish login nano
genwish login novel
genwish login danbooru
```

- `nano` and `novel` prompt for one API key.
- `danbooru` prompts for the login first, then the API key.
- `genwish logout <provider>` deletes the stored credential.
- Danbooru also answers anonymous requests.

Credentials use the OS keyring, falling back to a user-only file when the keyring is unavailable.
`genwish login` warns before prompting when that fallback will be used.

- Keyring service: `genwish`.
- Keyring users: `nanogpt-api-key`, `novelai-api-key`, `danbooru`.
- File fallback, under `os.UserConfigDir()/genwish/`:
  `nanogpt_api_key`, `novelai_api_key`, `danbooru`.

The Danbooru entry stores the login and API key together.

The Danbooru user agent is fixed at `genwish/0.1.0`.
