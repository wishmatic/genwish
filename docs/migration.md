# Migration

`genwish` replaces the CLI parts of `nano-mcp`, `neo-mcp`, `novel-mcp`, and `booru-mcp`. The MCP servers, Docker images, and HTTP file hosting are not part of genwish.

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

## Environment mapping

Unchanged:

- `IMAGE_URL_MAP`
- `OUTPUT_FORMAT`
- `SD_URL`
- `ERROR_DETAIL`
- `RATE_LIMIT_RPS`
- `RATE_LIMIT_BURST`
- `REQUEST_TIMEOUT_SECONDS`
- `MAX_LIMIT`
- `MAX_OFFSET`
- `BLOCKED_TAGS`
- `LOG_LEVEL`

Added:

- `GENWISH_CREDENTIAL_BACKEND`
- `NEO_MODELS_FILE`

Removed, because credentials are login-only:

- `NANOGPT_API_KEY`
- `NOVELAI_API_KEY`
- `DANBOORU_LOGIN`
- `DANBOORU_API_KEY`

Removed, because the user agent is hardcoded:

- `USER_AGENT`

Removed, server-only:

- `HOST`
- `PORT`
- `API_KEY`
- `PUBLIC_HOST`
- `FILES_DIR`

## Credential migration

genwish does not read the old keyring entries, files, or environment variables. Log in again:

```sh
genwish login nano
genwish login novel
genwish login danbooru
```

Old locations, for reference:

| Tool      | Keyring service | Keyring user      | File fallback                     | Backend variable                                       |
| --------- | --------------- | ----------------- | --------------------------------- | ------------------------------------------------------ |
| nano-cli  | `nano-mcp`      | `nanogpt-api-key` | `~/.config/nano/api_key`          | `NANO_CREDENTIAL_BACKEND`                              |
| novel-cli | `novel-cli`     | `novelai-api-key` | `~/.config/novel/novelai_api_key` | `NOVEL_CREDENTIAL_BACKEND`                             |
| booru-cli | none            | none              | none                              | environment only: `DANBOORU_LOGIN`, `DANBOORU_API_KEY` |

New locations:

| Provider | Keyring service | Keyring user      | File fallback                                |
| -------- | --------------- | ----------------- | -------------------------------------------- |
| nano     | `genwish`       | `nanogpt-api-key` | `os.UserConfigDir()/genwish/nanogpt_api_key` |
| novel    | `genwish`       | `novelai-api-key` | `os.UserConfigDir()/genwish/novelai_api_key` |
| danbooru | `genwish`       | `danbooru`        | `os.UserConfigDir()/genwish/danbooru`        |

Backend variable: `GENWISH_CREDENTIAL_BACKEND`.

The Danbooru entry stores the login and API key together. `genwish login danbooru` prompts for the login first, then the API key. Danbooru still works anonymously without a stored credential.

## Behavior changes

- `models` takes the provider as its first argument.
- `balance` is top-level and takes the provider.
- `login` and `logout` take the provider.
- API keys are login-only; no key environment variable is read.
- The Danbooru user agent is hardcoded at `genwish/0.1.0`.
- `gen image` and `gen video` take the provider as a subcommand.
- `gen video neo` and `gen video novel` fail with an unsupported-provider error.
- `models list neo` reads `NEO_MODELS_FILE` instead of querying Forge.
