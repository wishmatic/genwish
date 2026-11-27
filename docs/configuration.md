# Configuration

`genwish` reads a TOML config file and three environment variables. Credentials are not configuration; use `genwish login`.

## Config file

The default path is `~/.genwish/config.toml`. Pass `--config-file <path>` to any command to read a different file.

- A missing default file is not an error; the built-in defaults apply.
- A missing `--config-file` path is an error.
- Malformed TOML and unknown keys are errors.

Command flags override the file, which overrides the built-in defaults.

```toml
[output]
format = "webp"        # png, jpeg, jxl, or webp
path = "~/.genwish/out"

[neo]
url = "http://127.0.0.1:7860"
models_file = ""

[booru.rate_limit]
rps = 1.0
burst = 1
timeout = 20           # seconds
```

| Key                        | Type    | Default                 | Purpose                                                                                |
| -------------------------- | ------- | ----------------------- | -------------------------------------------------------------------------------------- |
| `output.format`            | string  | `webp`                  | Default image format for generated and edited images: `png`, `jpeg`, `jxl`, or `webp`. |
| `output.path`              | string  | `~/.genwish/out`        | Base directory for generated and edited media. A leading `~` is expanded.              |
| `neo.url`                  | string  | `http://127.0.0.1:7860` | Forge Neo API base URL.                                                                |
| `neo.models_file`          | string  | empty                   | Newline-delimited checkpoint names for `models list neo`. A leading `~` is expanded.   |
| `booru.rate_limit.rps`     | float   | `1`                     | Danbooru requests per second.                                                          |
| `booru.rate_limit.burst`   | integer | `1`                     | Danbooru burst size.                                                                   |
| `booru.rate_limit.timeout` | integer | `20`                    | Danbooru request timeout in seconds.                                                   |

Every image is converted to `output.format` when the provider returns a different type.

The tag page limit is fixed at 100 and the largest tag offset at 1000. There is no tag denylist.

## Output layout

Every image and video is written under the configured base directory:

```text
<output.path>/images/<YYYY-MM-DD>/<uuid>.<ext>
<output.path>/videos/<YYYY-MM-DD>/<uuid>.<ext>
```

`--out <path>` replaces `output.path` for that command. The date is UTC and the file name is a UUID. Commands print each written path.

## Environment

| Variable             | Default  | Purpose                                                         |
| -------------------- | -------- | --------------------------------------------------------------- |
| `LOG_LEVEL`          | `info`   | Reserved for internal logging.                                  |
| `ERROR_DETAIL`       | `useful` | Set to `verbose` to include upstream response bodies in errors. |
| `CREDENTIAL_BACKEND` | `auto`   | `auto`, `keyring`, or `file`.                                   |

No other environment variable is read. A `.env` file in the working directory is loaded when present.

## Credentials

```sh
genwish login nano
genwish login novel
genwish login danbooru
```

- `nano` and `novel` prompt for one API key.
- `danbooru` prompts for the login first, then the API key.
- `genwish logout <provider>` deletes the stored credential.
- Danbooru also works anonymously when no credential is stored.

`auto` prefers the OS keyring and falls back to a user-only file.

- Keyring service: `genwish`.
- Keyring users: `nanogpt-api-key`, `novelai-api-key`, `danbooru`.
- File fallback, under `os.UserConfigDir()/genwish/`: `nanogpt_api_key`, `novelai_api_key`, `danbooru`.

The Danbooru entry stores the login and API key together.

The Danbooru user agent is fixed at `genwish/0.1.0`.
