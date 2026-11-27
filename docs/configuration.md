# Configuration

`genwish` reads configuration from the environment, optionally loaded from `.env` in the working directory. Credentials are not read from the environment; use `genwish login`.

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

| Variable                     | Default | Purpose                       |
| ---------------------------- | ------- | ----------------------------- |
| `GENWISH_CREDENTIAL_BACKEND` | `auto`  | `auto`, `keyring`, or `file`. |

`auto` prefers the OS keyring and falls back to a user-only file.

- Keyring service: `genwish`.
- Keyring users: `nanogpt-api-key`, `novelai-api-key`, `danbooru`.
- File fallback, under `os.UserConfigDir()/genwish/`: `nanogpt_api_key`, `novelai_api_key`, `danbooru`.

The Danbooru entry stores the login and API key together.

## Generation and input

| Variable          | Default                 | Purpose                                                                                                                                   |
| ----------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `IMAGE_URL_MAP`   | empty                   | Comma-separated `public=private` address rewrites for input images. A mapped directory is readable by anyone who can call the tool.       |
| `OUTPUT_FORMAT`   | `webp`                  | Default image format for Forge Neo and NovelAI: `png`, `jpeg`, `jxl`, or `webp`.                                                          |
| `SD_URL`          | `http://127.0.0.1:7860` | Forge Neo API base URL.                                                                                                                   |
| `NEO_MODELS_FILE` | empty                   | Newline-delimited checkpoint names for `models list neo`. Unset means an empty list. Blank lines and lines starting with `#` are ignored. |

## Booru

| Variable                  | Default  | Purpose                                                         |
| ------------------------- | -------- | --------------------------------------------------------------- |
| `ERROR_DETAIL`            | `useful` | Set to `verbose` to include upstream response bodies in errors. |
| `RATE_LIMIT_RPS`          | `1`      | Danbooru requests per second.                                   |
| `RATE_LIMIT_BURST`        | `1`      | Danbooru burst size.                                            |
| `REQUEST_TIMEOUT_SECONDS` | `20`     | Danbooru request timeout.                                       |
| `MAX_LIMIT`               | `100`    | Largest tag page, at most 100.                                  |
| `MAX_OFFSET`              | `1000`   | Largest tag offset.                                             |
| `BLOCKED_TAGS`            | empty    | Comma-separated tag denylist.                                   |

The Danbooru user agent is fixed at `genwish/0.1.0`.

## Logging

| Variable    | Default | Purpose                        |
| ----------- | ------- | ------------------------------ |
| `LOG_LEVEL` | `info`  | Reserved for internal logging. |

## Removed settings

The MCP servers used `HOST`, `PORT`, `API_KEY`, `PUBLIC_HOST`, and `FILES_DIR`. `genwish` does not run a server and does not read them.

`genwish` also does not read `NANOGPT_API_KEY`, `NOVELAI_API_KEY`, `DANBOORU_LOGIN`, `DANBOORU_API_KEY`, or `USER_AGENT`. API keys are login-only, and the user agent is hardcoded.
