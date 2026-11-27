# Migration

`genwish` replaces the CLI parts of `nano-mcp`, `neo-mcp`, `novel-mcp`, and `booru-mcp`. The MCP
servers, Docker images, and HTTP file hosting are not part of genwish.

## Command mapping

| Old                                   | New                                            |
| ------------------------------------- | ---------------------------------------------- |
| `nano-cli gen image ...`              | `genwish gen image nano ...`                   |
| `nano-cli gen video ...`              | `genwish gen video nano ...`                   |
| `nano-cli models list [image\|video]` | `genwish models list nano [image\|video]`      |
| `nano-cli models balance`             | `genwish balance nano`                         |
| `nano-cli login` / `logout`           | `genwish login nano` / `genwish logout nano`   |
| `neo-cli gen image ...`               | `genwish gen image neo ...`                    |
| `neo-cli edit ...`                    | `genwish edit ...`                             |
| `novel-cli gen image ...`             | `genwish gen image novel ...`                  |
| `novel-cli login` / `logout`          | `genwish login novel` / `genwish logout novel` |
| `booru-cli tags ...`                  | `genwish util booru ...`                       |

Removed commands and flags:

- `serve` on every old CLI.
- `nano-cli web search`, `web scrape`, `web map`, `web crawl`, and `web youtube`.
- MCP tool names and schemas.

## Configuration mapping

Moved to `~/.genwish/config.toml`, or the file named by `--config-file`:

| Old environment variable  | Config key                 |
| ------------------------- | -------------------------- |
| `OUTPUT_FORMAT`           | `output.format`            |
| `SD_URL`                  | `neo.url`                  |
| `NEO_MODELS_FILE`         | `neo.models_file`          |
| `RATE_LIMIT_RPS`          | `booru.rate_limit.rps`     |
| `RATE_LIMIT_BURST`        | `booru.rate_limit.burst`   |
| `REQUEST_TIMEOUT_SECONDS` | `booru.rate_limit.timeout` |

Removed:

- `LOG_LEVEL`: there is no logging; results go to stdout and warnings to stderr.
- `ERROR_DETAIL`: errors always include the useful upstream detail.
- `CREDENTIAL_BACKEND` (and `GENWISH_CREDENTIAL_BACKEND`): credentials always try the OS
  keyring first.
- `MAX_LIMIT`: the tag page limit is fixed at 100.
- `MAX_OFFSET`: the largest offset is fixed at 1000.
- `BLOCKED_TAGS`: there is no tag denylist.
- `IMAGE_URL_MAP`: input addresses are no longer rewritten.
- `NANOGPT_API_KEY`, `NOVELAI_API_KEY`, `DANBOORU_LOGIN`, `DANBOORU_API_KEY`: credentials are
  login-only.
- `USER_AGENT`: the user agent is hardcoded.
- `HOST`, `PORT`, `API_KEY`, `PUBLIC_HOST`, `FILES_DIR`: server-only.

## Credential migration

genwish does not read the old keyring entries, files, or environment variables. Log in again:

```sh
genwish login nano
genwish login novel
genwish login danbooru
```

Old locations, for reference:

- nano-cli: keyring `nano-mcp` / `nanogpt-api-key`, file `~/.config/nano/api_key`,
  backend `NANO_CREDENTIAL_BACKEND`.
- novel-cli: keyring `novel-cli` / `novelai-api-key`, file `~/.config/novel/novelai_api_key`,
  backend `NOVEL_CREDENTIAL_BACKEND`.
- booru-cli: no keyring or file; environment only: `DANBOORU_LOGIN`, `DANBOORU_API_KEY`.

New locations:

| Provider | Keyring service | Keyring user      | File fallback                                |
| -------- | --------------- | ----------------- | -------------------------------------------- |
| nano     | `genwish`       | `nanogpt-api-key` | `os.UserConfigDir()/genwish/nanogpt_api_key` |
| novel    | `genwish`       | `novelai-api-key` | `os.UserConfigDir()/genwish/novelai_api_key` |
| danbooru | `genwish`       | `danbooru`        | `os.UserConfigDir()/genwish/danbooru`        |

Backend: always auto. The OS keyring is tried first, then a user-only file. `genwish login` warns
before prompting when the keyring is unavailable.

The Danbooru entry stores the login and API key together. `genwish login danbooru` prompts for the
login first, then the API key. Danbooru still works anonymously without a stored credential.

## Behavior changes

- `models` takes the provider as its first argument.
- `balance` is top-level and takes the provider.
- `login` and `logout` take the provider.
- API keys are login-only; no key environment variable is read.
- The Danbooru user agent is hardcoded at `genwish/0.1.0`.
- Generated and edited media is written under
  `<output.path>/<images|videos>/<date>/<uuid>.<ext>`; `--out` replaces the base directory.
- `gen image` and `gen video` take the provider as a subcommand.
- `gen video neo` and `gen video novel` fail with an unsupported-provider error.
- `models list neo` reads `neo.models_file` instead of querying Forge.
