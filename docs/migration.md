# Migration

`genwish` replaces the CLI parts of `nano-mcp`, `neo-mcp`, `novel-mcp`, and `booru-mcp`. The MCP servers, Docker images, and HTTP file hosting are not part of genwish.

## Command mapping

| Old | New |
| --- | --- |
| `nano-cli gen image ...` | `genwish gen image nano ...` |
| `nano-cli gen video ...` | `genwish gen video nano ...` |
| `nano-cli models list [image\|video]` | `genwish models list nano [image\|video]` |
| `nano-cli models balance` | `genwish balance nano` |
| `nano-cli login` / `logout` | `genwish login nano` / `genwish logout nano` |
| `neo-cli gen image ...` | `genwish gen image neo ...` |
| `neo-cli edit ...` | `genwish edit ...` |
| `novel-cli gen image ...` | `genwish gen image novel ...` |
| `novel-cli login` / `logout` | `genwish login novel` / `genwish logout novel` |
| `booru-cli tags ...` | `genwish util booru ...` |

Removed commands and flags:

- `serve` on every old CLI.
- `nano-cli web search`, `web scrape`, `web map`, `web crawl`, and `web youtube`.
- MCP tool names and schemas.

## Environment mapping

Unchanged:

- `NANOGPT_API_KEY`
- `NOVELAI_API_KEY`
- `IMAGE_URL_MAP`
- `OUTPUT_FORMAT`
- `SD_URL`
- `DANBOORU_LOGIN`
- `DANBOORU_API_KEY`
- `USER_AGENT`
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

Removed, server-only:

- `HOST`
- `PORT`
- `API_KEY`
- `PUBLIC_HOST`
- `FILES_DIR`

## Credential migration

genwish does not read the old keyring entries or files. Log in again:

```sh
genwish login nano
genwish login novel
```

Old locations, for reference:

| Tool | Keyring service | Keyring user | File fallback | Backend variable |
| --- | --- | --- | --- | --- |
| nano-cli | `nano-mcp` | `nanogpt-api-key` | `~/.config/nano/api_key` | `NANO_CREDENTIAL_BACKEND` |
| novel-cli | `novel-cli` | `novelai-api-key` | `~/.config/novel/novelai_api_key` | `NOVEL_CREDENTIAL_BACKEND` |

New locations:

| Provider | Keyring service | Keyring user | File fallback |
| --- | --- | --- | --- |
| nano | `genwish` | `nanogpt-api-key` | `os.UserConfigDir()/genwish/nanogpt_api_key` |
| novel | `genwish` | `novelai-api-key` | `os.UserConfigDir()/genwish/novelai_api_key` |

Backend variable: `GENWISH_CREDENTIAL_BACKEND`.

## Behavior changes

- `models` takes the provider as its first argument.
- `balance` is top-level and takes the provider.
- `login` and `logout` take the provider.
- `gen image` and `gen video` take the provider as a subcommand.
- `gen video neo` and `gen video novel` fail with an unsupported-provider error.
- `models list neo` reads `NEO_MODELS_FILE` instead of querying Forge.
