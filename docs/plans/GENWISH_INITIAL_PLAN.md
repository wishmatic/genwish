# Genwish Initial Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use /subagent-driven-development (recommended) or /executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build `genwish`, a CLI-only generation tool that folds `nano-cli`, the Neo image tools, the NovelAI image tools, and the Danbooru tag lookup into one binary with mandatory provider selection and no MCP surface.

**Architecture:** Start from the NanoGPT CLI core, strip every MCP, server, and web feature, then add provider-aware command routing on top of a shared config, credential, resolve, and output layer. Provider-specific packages (`nanogpt`, `forge`, `novelai`, `booru`) stay separate behind the `cmd/genwish` command tree, and each provider keeps its existing flags and behavior.

**Tech Stack:** Go 1.27.0, `spf13/cobra`, `caarlos0/env/v11`, `joho/godotenv`, `go.uber.org/zap`, `zalando/go-keyring`, `golang.org/x/term`, `gen2brain/jxl`, `gen2brain/vpx`, `vmihailenco/msgpack/v5`, `golang.org/x/time/rate`.

**Spec:** The user brief is reproduced under "Requirements" below; there is no separate design document.

## Status

Implemented natively. Automated verification passes: `go build ./...`, `go vet ./...`, `go test ./... -race -count=1`, and `gofmt -l .` is clean. The human smoke tests in Unit 10 remain outstanding. The plan stays in `docs/plans` until they pass.

The image URL map was removed after implementation: `internal/sourcemap` and `IMAGE_URL_MAP` no longer exist, and `internal/resolve` fetches URLs unchanged. References to it below record the original build.

## Requirements

- Deprecate `booru-mcp`, `nano-mcp`, `neo-mcp`, and `novel-mcp` in favor of two projects: `genwish` (generation) and `lukwish` (the rest of NanoGPT). This plan covers only `genwish`.
- Start by copying all of `nano-cli` into `genwish` and renaming the binary to `genwish`.
- Remove MCP support entirely. No MCP SDK, no server, no `serve` command. Only the CLI remains.
- `gen image` and `gen video` take a mandatory provider: `nano`, `neo`, or `novel`.
- Video works only with `gen video nano` for now, but the provider must still be stated.
- Strip all web tools (`search`, `scrape`, `map`, `crawl`, `youtube`). Leave the originals in `nano-mcp` untouched.
- Add `genwish util booru`. It does what `booru-cli tags` does today; `tags` is not a separate command.
- Add Neo and NovelAI. Neo edit commands keep their existing names and flags. NovelAI generation keeps its existing flags.
- `models` takes a mandatory provider. `models balance` moves to top-level `balance`, which takes a provider.
- `balance` accepts `nano` (also `nano-gpt`, `nanogpt`, `nano gpt`, case-insensitive) and `novel` (also `novel-ai`, `novel ai`, `novelai`, case-insensitive).
- `balance novel` reports the account's Anlas.
- Do not lose documentation detail. Consolidate it now; the user will clean it up later.

## Global Constraints

- Module path: `github.com/wishmatic/genwish`; Go directive: `go 1.27.0`.
- Binary and root command: `genwish`.
- No MCP dependency, import, command, route, server, Dockerfile, or CI workflow.
- No `internal/mcp`, `internal/server`, `internal/auth`, `internal/present`, `internal/filestore`, or `internal/store` packages.
- Provider is mandatory for `gen image`, `gen video`, `models`, `balance`, `login`, and `logout`.
- Provider aliases are case-insensitive:
  - Nano: `nano`, `nano-gpt`, `nanogpt`, `nano gpt`.
  - Novel: `novel`, `novel-ai`, `novel ai`, `novelai`.
  - Neo: `neo`.
- Video is implemented only for Nano. `gen video neo` and `gen video novel` must fail clearly and must not call a client.
- Neo edit command names and flags stay as they are today: `trim`, `square`, `circle`, `convert`, `bgkill`.
- Nano, Neo, NovelAI, and Booru generation flags stay as they are today unless a task explicitly changes them.
- Do not edit `booru-mcp`, `nano-mcp`, `neo-mcp`, or `novel-mcp`; copy from them.
- API keys are login-only. No API key environment variable is read.
- Danbooru credentials come from `genwish login danbooru`, which prompts for the login first and the API key second. Anonymous access still works when none is stored.
- The Danbooru user agent is hardcoded at `genwish/0.1.0`.
- Preserve every documented env var, flag, JSON shape, and behavior in the new docs. Server-only env vars are removed and recorded in the migration doc.
- Do not run git write commands (no commits, branches, stashes) unless the user explicitly asks.

## Review Focus

1. Provider alias parsing with mixed case, spaces, and hyphens (`NANO GPT`, `novel-ai`, `Novel AI`) must resolve identically for `balance`, `models`, `login`, and `logout`.
2. `gen video neo` and `gen video novel` must fail with a clear unsupported-provider error and must not silently fall through to Nano.
3. NovelAI balance parsing must accept both known shapes of `subscription.trainingStepsLeft`: an object with `fixedTrainingStepsLeft` and `purchasedTrainingSteps`, and a plain integer. A missing field and a 401 must surface as errors.
4. Credential isolation: `login novel` must not overwrite the NanoGPT key, `logout novel` must not delete the NanoGPT key, and environment variables must still take precedence over stored keys.
5. Image inputs: local files (including directories and missing paths), URLs rewritten by `IMAGE_URL_MAP`, data URIs, and raw base64 must behave consistently for `gen image` init images and all `edit` commands.
6. Neo model listing must fail cleanly when `NEO_MODELS_FILE` is set but unreadable or malformed, and must default to an empty list when it is unset.

## Command Tree

```text
genwish
  gen | generate
    image
      nano <prompt...>     # NanoGPT, existing nano-cli flags
      neo <prompt...>      # Forge Neo, existing neo-cli flags
      novel <prompt...>    # NovelAI, existing novel-cli flags
    video
      nano <prompt...>     # implemented
      neo <prompt...>      # explicit unsupported-provider error
      novel <prompt...>    # explicit unsupported-provider error
  models
    list <provider> [image|video]
  balance <provider>       # nano or novel
  login <provider>         # nano or novel
  logout <provider>        # nano or novel
  util
    booru <search...>      # same behavior as booru-cli tags
```

Defaults preserved from the source CLIs:

- `models list nano` defaults to both image and video kinds; `models list neo` and `models list novel` default to image only.
- `gen image novel --model` defaults to `nai-diffusion-5-full`.
- `gen image novel --format` and `gen image neo --format` default to WebP via `OUTPUT_FORMAT`.
- `gen image nano` keeps writing the media type the provider returned.
- `genwish util booru` keeps `--offset`, `--limit`, `--min-count`, `--exact`, repeatable `--category`, and `--json`.

## File Map

```text
genwish/
  go.mod
  go.sum
  .env.example
  README.md
  cmd/genwish/
    main.go
    provider_args.go
    gen.go
    gen_image.go
    gen_image_nano.go
    gen_image_neo.go
    gen_image_novel.go
    gen_video.go
    gen_video_nano.go
    gen_video_neo.go
    gen_video_novel.go
    models.go
    models_list.go
    models_list_nano.go
    models_list_neo.go
    models_list_novel.go
    balance.go
    balance_nano.go
    balance_novel.go
    login.go
    logout.go
    util.go
    util_booru.go
  internal/
    booru/           # from booru-mcp
    catalog/         # from booru-mcp
    cli/             # shared CLI helpers and Deps
    config/          # unified config
    credential/      # provider-aware key storage
    edit/            # from neo-mcp
    fetch/           # from booru-mcp
    forge/           # from neo-mcp
    format/          # from neo-mcp (identical to novel-mcp's copy)
    nanogpt/         # from nano-mcp, web files removed
    novelai/         # from novel-mcp
    provider/        # provider names and aliases
    resolve/         # unified resolver
    sourcemap/       # from the source CLIs
    utils/           # from neo-mcp or novel-mcp (identical)
    view/            # from booru-mcp
  docs/
    plans/GENWISH_INITIAL_PLAN.md
    README.md
    configuration.md
    migration.md
    nano.md
    neo.md
    novel.md
    booru.md
    prompts/neo.md
    prompts/novel.md
    reference/nano-README.md
    reference/neo-README.md
    reference/novel-README.md
    reference/booru-README.md
```

---

## Implementation Unit 1: Scaffold Genwish from the NanoGPT CLI

**Files:**

- Create: `go.mod`, `go.sum`, `.env.example`
- Create: `cmd/genwish/main.go`
- Create: `cmd/genwish/gen.go`
- Create: `cmd/genwish/gen_image.go`
- Create: `cmd/genwish/gen_video.go`
- Create: `cmd/genwish/login.go`
- Create: `cmd/genwish/logout.go`
- Create: `cmd/genwish/models.go`
- Create: `cmd/genwish/models_list.go`
- Create: `cmd/genwish/models_balance.go`
- Copy: `nano-mcp/internal/nanogpt/{client.go,balance.go,images.go,models.go,video.go}` and their tests
- Copy: `nano-mcp/internal/cli/{deps.go,flags.go,media.go,asset.go,output.go}` and their tests
- Copy: `nano-mcp/internal/config/config.go`, `nano-mcp/internal/credential/*`, `nano-mcp/internal/resolve/*`, `nano-mcp/internal/sourcemap/*`, `nano-mcp/internal/utils/*`
- Create: `cmd/genwish/main_test.go` (root help tree)

**Interfaces:**

- Consumes: nothing outside this task.
- Produces: a buildable `genwish` with the NanoGPT generation core and the temporary legacy command shapes.

- [ ] Copy the listed files and rewrite every import to `github.com/wishmatic/genwish/...`.
- [ ] Rename the root command `Use` to `genwish`; change user-facing `nano-cli` strings to `genwish`.
- [ ] Do not copy `internal/auth`, `internal/mcp`, `internal/server`, or `internal/filestore`.
- [ ] Do not copy `cmd/nano-cli/serve.go`, `cmd/nano-cli/web*.go`, or `internal/nanogpt/{search.go,scrape.go,firecrawl.go,youtube.go}`.
- [ ] Remove `newServeCommand` and `newWebCommand` from the root wiring.
- [ ] Copy `internal/config/config.go` unchanged for now; Unit 2 replaces it with the unified schema. Do not delete fields that the copied constructors still reference during this unit.
- [ ] Run `go mod tidy`; the dependency set must not contain `go-chi`, `cors`, `modelcontextprotocol`, `jsonschema-go`, or `google/uuid`.
- [ ] Write `TestNewRoot` in `cmd/genwish/main_test.go` with a subtest `"lists the surviving commands"` asserting the help output contains `gen`, `models`, `login`, `logout` and does not contain `serve`, `web`, `mcp`.

**Acceptance Criteria:**

- `go build ./...` passes.
- `go vet ./...` passes.
- `go test ./...` passes.
- `go run ./cmd/genwish help` lists only `gen`, `models`, `login`, and `logout`.
- `grep -ri "modelcontextprotocol\|go-chi\|jsonschema\|google/uuid" --include="*.go" .` returns nothing.
- No file under `genwish/` imports a source module path (`nano-mcp`, `neo-mcp`, `novel-cli`, `booru-mcp`).

---

## Implementation Unit 2: Provider, Config, and Credential Foundations

**Files:**

- Create: `internal/provider/provider.go`
- Create: `internal/provider/provider_test.go`
- Rewrite: `internal/config/config.go`
- Create: `internal/config/config_test.go`
- Rewrite: `internal/credential/store.go`
- Rewrite: `internal/credential/keyring.go`
- Create: `internal/credential/credential_test.go`
- Rewrite: `internal/cli/deps.go`
- Create: `internal/cli/deps_test.go`

**Interfaces:**

- Consumes: nothing from earlier tasks except the copied files from Unit 1.
- Produces:

```go
package provider

type Name string

const (
	Nano  Name = "nano"
	Neo   Name = "neo"
	Novel Name = "novel"
)

func Parse(value string) (Name, error)         // nano, neo, novel and aliases
func ParseBalance(value string) (Name, error)  // nano or novel only
func ParseCredential(value string) (Name, error) // nano or novel only
func (n Name) String() string
```

```go
package credential

var ErrNotFound = errors.New("no stored API key")
var ErrEmptyKey = errors.New("API key is empty")

func Open(p provider.Name) (Store, error)
```

```go
package config

type Config struct {
	LogLevel string
	ImageURLMap string
	DefaultOutput string
	SDURL string
	ErrorDetail string
	RateLimitRPS float64
	RateLimitBurst int
	RequestTimeoutSeconds int
	MaxLimit int
	MaxOffset int
	NEOModelsFile string
	BlockedTagsRaw string
}
```

- [ ] Write failing tests for `provider.Parse` with the alias table: `nano`, `nano-gpt`, `nanogpt`, `nano gpt`, `NANO GPT`, `neo`, `novel`, `novel-ai`, `novel ai`, `novelai`, `Novel AI`; unknown and empty values must error.
- [ ] Write `ParseBalance` and `ParseCredential` tests: `neo` must error; nano, novel, and danbooru must pass `ParseCredential`, while `Parse` rejects danbooru.
- [ ] Implement `provider` with `strings.ToLower`, `strings.TrimSpace`, and removal of spaces and hyphens before matching.
- [ ] Write credential tests with `NewFileStoreAt(t.TempDir()+"/key")`: save, load, delete, empty key, missing key; separate `Open(provider.Nano)`, `Open(provider.Novel)`, and `Open(provider.Danbooru)` stores must not see each other's key; `GENWISH_CREDENTIAL_BACKEND=bogus` must error.
- [ ] Implement the credential store:
  - Keyring service: `genwish`.
  - Keyring users: `nanogpt-api-key`, `novelai-api-key`, `danbooru`.
  - File fallback: `os.UserConfigDir()/genwish/nanogpt_api_key`, `.../novelai_api_key`, and `.../danbooru`.
  - The Danbooru entry stores a login and API key together.
  - Backend env: `GENWISH_CREDENTIAL_BACKEND` with `auto`, `keyring`, `file`.
  - Keep the auto store behavior from the source: keyring first, file fallback, delete both.
- [ ] Write config tests for defaults and `VerboseErrors`. `BlockedTags` is added and tested in Unit 4 with the copied `booru` package.
- [ ] Implement config env tags with these names: `IMAGE_URL_MAP`, `OUTPUT_FORMAT` (default `webp`), `SD_URL` (default `http://127.0.0.1:7860`), `ERROR_DETAIL`, `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`, `REQUEST_TIMEOUT_SECONDS`, `MAX_LIMIT`, `MAX_OFFSET`, `NEO_MODELS_FILE` (default empty), `BLOCKED_TAGS`, `LOG_LEVEL`. No API key, Danbooru credential, or user agent variable is read. Booru fields are inert until Unit 4 wires them.
- [ ] Rewrite `internal/cli/deps.go` with the fields available at this point:

```go
type Deps struct {
	In io.Reader
	Out io.Writer
	Err io.Writer
	Config config.Config
	Credentials credential.Store
	Nano *nanogpt.Client
	Resolver *resolve.Client
}
```

Later units add exactly one block each: Unit 4 adds `Catalog *catalog.Service`; Unit 5 adds `Forge *forge.Client`; Unit 6 adds `Novel *novelai.Client`. Each addition is accompanied by its constructor in the same unit.

- [ ] Provide constructors:
  - `NanoClient(deps)`: reads `credential.Open(provider.Nano)`; error message says `run `genwish login nano``.
  - `NovelClient(deps)` (Unit 6): reads `credential.Open(provider.Novel)`; error message says `run `genwish login novel``.
  - `ForgeClient(deps)` (Unit 5): from `Config.SDURL`.
  - `CatalogService(deps)` (Unit 4): loads the Danbooru login and API key from `provider.Danbooru`, falls back to anonymous when absent, and uses the hardcoded user agent.
  - `ImageResolver(deps)`: parses `IMAGE_URL_MAP` and returns a `resolve.Client`.
  - `ReadImage(ctx, resolver, address)`: local file if it exists, otherwise resolve.

**Acceptance Criteria:**

- Alias tests cover all listed spellings and both rejected providers.
- Credential isolation test proves nano and novel keys are independent.
- `go build ./...`, `go vet ./...`, and `go test ./...` pass.
- `docs/configuration.md` is not part of this unit yet; the env names above are the source of truth for Unit 9.

---

## Implementation Unit 3: Provider-Aware Command Tree and NanoGPT Commands

**Files:**

- Create: `cmd/genwish/gen.go`
- Create: `cmd/genwish/gen_image.go`
- Create: `cmd/genwish/gen_image_nano.go`
- Create: `cmd/genwish/gen_video.go`
- Create: `cmd/genwish/gen_video_nano.go`
- Create: `cmd/genwish/models.go`
- Create: `cmd/genwish/models_list.go`
- Create: `cmd/genwish/models_list_nano.go`
- Create: `cmd/genwish/balance.go`
- Create: `cmd/genwish/balance_nano.go`
- Rewrite: `cmd/genwish/login.go`
- Rewrite: `cmd/genwish/logout.go`
- Delete: `cmd/genwish/models_balance.go`
- Move Nano request and output helpers into `gen_image_nano.go` and `gen_video_nano.go`
- Create: `cmd/genwish/gen_image_nano_test.go`, `gen_video_nano_test.go`, `models_list_nano_test.go`, `balance_nano_test.go`, `login_test.go`, `logout_test.go`

**Interfaces:**

- Consumes: `provider`, `config`, `credential`, `cli.Deps`.
- Produces:

```go
func newGenCommand(deps cli.Deps) *cobra.Command
func newGenImageCommand(deps cli.Deps) *cobra.Command
func newGenImageNanoCommand(deps cli.Deps) *cobra.Command
func newGenVideoCommand(deps cli.Deps) *cobra.Command
func newGenVideoNanoCommand(deps cli.Deps) *cobra.Command
func newModelsCommand(deps cli.Deps) *cobra.Command
func newModelsListCommand(deps cli.Deps) *cobra.Command
func newModelsListNanoCommand(deps cli.Deps) *cobra.Command
func newBalanceCommand(deps cli.Deps) *cobra.Command
func newBalanceNanoCommand(deps cli.Deps) *cobra.Command
func newLoginCommand(deps cli.Deps) *cobra.Command
func newLogoutCommand(deps cli.Deps) *cobra.Command
```

- [ ] Write root/command tests:
  - `TestNewGenCommand` subtest `"lists the provider subcommands"` asserts `image` and `video`.
  - `TestNewGenImageCommand` subtest `"requires a provider"` runs `gen image` and asserts help output is printed and the command succeeds.
  - `TestNewGenVideoCommand` subtest `"requires a provider"` does the same for `gen video`.
- [ ] Implement `gen image` as a parent that registers `nano` (and, in later units, `neo` and `novel`). `gen image nano <prompt...>` keeps the current nano flags and behavior from `nano-mcp/cmd/nano-cli/gen_image.go`.
- [ ] Intermediate state: Unit 3 registers only `nano` for `gen image`, `gen video`, and `models list`, and only `nano` for `balance`. Providers added in later units must fail as unknown subcommands or a clear provider error, never panic. Neo generation and `models list neo` arrive in Unit 5; NovelAI generation, `models list novel`, and `login/logout novel` arrive in Unit 6; `balance novel` arrives in Unit 7; video stubs arrive in Unit 8.
- [ ] Implement `gen video` as a parent that registers `nano` (and unsupported `neo`/`novel` stubs in Unit 8). `gen video nano <prompt...>` keeps the current nano flags and behavior from `nano-mcp/cmd/nano-cli/gen_video.go`.
- [ ] Restructure `models`:
  - `genwish models list <provider> [image|video]`.
  - `--query` and `--limit` stay.
  - `models list nano` keeps the existing kind defaulting (`image` and `video`).
  - Provider lookup uses `provider.Parse`; unknown providers error before any client call.
  - `balance` is no longer a subcommand of `models`.
- [ ] Move the balance implementation to `balance`:
  - `genwish balance <provider>`.
  - Provider uses `provider.ParseBalance`.
  - Nano prints `USD: ...` and `Nano: ...`.
  - `balance neo` errors: `balance is not available for provider neo`.
- [ ] Make `login` and `logout` provider-aware:
  - `genwish login <provider>` and `genwish logout <provider>`.
  - Provider uses `provider.ParseCredential`.
  - The prompt says `NanoGPT API key:` or `NovelAI API key:`.
  - `login neo` and `logout neo` error: `credentials are not used by provider neo`.
  - Success output stays `logged in` and `logged out`.
- [ ] Port the existing Nano CLI tests for image, video, models, and balance, adapting the command paths and injected deps.

**Acceptance Criteria:**

- `genwish gen image` and `genwish gen video` print help and do not generate.
- `genwish gen image nano ...` produces the same `nanogpt.ImageRequest` and output files as `nano-cli gen image ...`.
- `genwish gen video nano ...` produces the same `nanogpt.VideoRequest` and output file as `nano-cli gen video ...`.
- `genwish models list nano image --query x --limit 5` prints the same table as before.
- `genwish balance nano` prints the same three lines as before.
- `genwish login nano` and `genwish logout nano` use the provider-aware credential store.
- `go test ./...` passes.

**Tests to add:**

- `TestGenImageNano` subtests: `"sends the prompt and flags"`, `"writes the returned image"`, `"reports the cost"`, `"requires a prompt"`.
- `TestGenVideoNano` subtests: `"requires a model"`, `"animates an init image"`, `"writes the returned video"`.
- `TestModelsList` subtests: `"nano defaults to both kinds"`, `"rejects an unknown provider"`, `"rejects an unknown kind"`.
- `TestBalance` subtests: `"nano prints the balances"`, `"rejects neo"`.
- `TestLogin` and `TestLogout` subtests: `"stores the nano key"`, `"deletes the nano key"`, `"rejects neo"`.

---

## Implementation Unit 4: `util booru`

**Files:**

- Copy: `booru-mcp/internal/booru/*` and tests
- Copy: `booru-mcp/internal/catalog/*` and tests
- Copy: `booru-mcp/internal/fetch/*` and tests
- Copy: `booru-mcp/internal/view/*` and tests
- Create: `cmd/genwish/util.go`
- Create: `cmd/genwish/util_booru.go`
- Create: `cmd/genwish/util_booru_test.go`

**Interfaces:**

- Consumes: `cli.CatalogService`, `booru`, `catalog`, `view`.
- Produces:

```go
func newUtilCommand(deps cli.Deps) *cobra.Command
func newUtilBooruCommand(deps cli.Deps) *cobra.Command
func runUtilBooru(ctx context.Context, deps cli.Deps, opts booruOptions, args []string) error
```

- [ ] Copy the Booru packages, rewrite imports to `github.com/wishmatic/genwish/...`, and keep their tests.
- [ ] Add `Catalog *catalog.Service` to `cli.Deps` and implement `cli.CatalogService(deps)`.
- [ ] Add `Config.BlockedTags()` and its test now that `booru` exists.
- [ ] Implement `util` as a parent command with `booru` as its child.
- [ ] Implement `util booru <search...>` with the same flags and behavior as `booru-cli tags`: `--offset`, `--limit`, `--min-count`, `--exact`, repeatable `--category`, `--json`.
- [ ] Keep the error prefix `tags:` for service failures, matching the source CLI.
- [ ] Port `booru-mcp/cmd/booru-cli/tags_test.go` to `util_booru_test.go`, injecting a fake `catalog.Service` through `cli.Deps`.
- [ ] Add `TestNewUtilCommand` subtest `"lists booru"` and `TestNewUtilBooruCommand` subtest `"lists its flags"`.

**Acceptance Criteria:**

- `genwish util booru blue_hair` prints the same text as `booru-cli tags blue_hair`.
- `genwish util booru blue_hair --json` prints the same JSON as `booru-cli tags blue_hair --json`.
- Category parsing, min-count, offset, and limit behave as before.
- `go test ./...` passes.

---

## Implementation Unit 5: Neo (Forge, Edit, Models)

**Files:**

- Copy: `neo-mcp/internal/format/*` and tests
- Copy: `neo-mcp/internal/edit/*` and tests
- Copy: `neo-mcp/internal/forge/*` and tests
- Create: `cmd/genwish/gen_image_neo.go`
- Create: `cmd/genwish/edit.go`
- Create: `cmd/genwish/edit_crop.go`
- Create: `cmd/genwish/edit_convert.go`
- Create: `cmd/genwish/edit_bgkill.go`
- Create: `cmd/genwish/models_list_neo.go`
- Create: corresponding `_test.go` files
- Modify: `cmd/genwish/gen_image.go` and `cmd/genwish/models_list.go` registration

**Interfaces:**

- Consumes: `format`, `edit`, `forge`, `cli`.
- Produces:

```go
func newGenImageNeoCommand(deps cli.Deps) *cobra.Command
func newEditCommand(deps cli.Deps) *cobra.Command
func newEditCropCommand(deps cli.Deps, spec cropCommand) *cobra.Command
func newEditConvertCommand(deps cli.Deps) *cobra.Command
func newEditBgkillCommand(deps cli.Deps) *cobra.Command
func newModelsListNeoCommand(deps cli.Deps) *cobra.Command
```

- [ ] Copy the Neo packages and rewrite imports. The `format` package is shared with NovelAI; copy it once here.
- [ ] Add `Forge *forge.Client` to `cli.Deps` and implement `cli.ForgeClient(deps)`.
- [ ] Implement `gen image neo <prompt...>` with the exact flags and behavior from `neo-mcp/cmd/neo-cli/gen_image.go`.
- [ ] Implement the `edit` tree exactly as in `neo-mcp/cmd/neo-cli/edit.go`: `trim`, `square`, `circle`, `convert`, `bgkill`, with their existing flags and output.
- [ ] Implement `models list neo [image]` from configuration, not from Forge:
  - New env var `NEO_MODELS_FILE`, default empty.
  - The file is newline-delimited checkpoint names; blank lines and lines starting with `#` are ignored.
  - Unset `NEO_MODELS_FILE` means an empty list and exits successfully.
  - A set but unreadable or missing file is an error naming the path and the variable.
  - `video` errors: `neo does not serve video models`.
  - `--query` filters on the checkpoint name; `--limit` applies.
- [ ] Add `Config.NEOModelsFile` for `NEO_MODELS_FILE` and its default-empty test.
- [ ] Port every Neo CLI test, adapting command paths and deps. Keep the existing error prefixes (`forge:`, `edit:`, `bgkill:`, `convert:`).
- [ ] Add `TestNeoModels` subtests: `"reads the configured file"`, `"ignores comments and blank lines"`, `"defaults to empty when unset"`, `"reports a missing configured file"`.

**Acceptance Criteria:**

- `genwish gen image neo "..." --model x --forge-preset y` sends the same Forge payload as `neo-cli`.
- `genwish edit trim|square|circle|convert|bgkill` behaves exactly as `neo-cli edit ...`.
- `genwish models list neo` prints the checkpoints in `NEO_MODELS_FILE`; an unset variable prints nothing, and a missing configured file produces a useful error.
- `genwish models list neo video` errors.
- `go test ./...` passes.

---

## Implementation Unit 6: NovelAI (Generation, Login, Models)

**Files:**

- Copy: `novel-mcp/internal/novelai/*` and tests, except `_AUTHORSHIP`
- Create: `internal/novelai/models.go`
- Create: `internal/novelai/models_test.go`
- Create: `cmd/genwish/gen_image_novel.go`
- Create: `cmd/genwish/models_list_novel.go`
- Create: corresponding `_test.go` files
- Modify: `internal/cli/deps.go` to add `NovelClient`
- Modify: `cmd/genwish/gen_image.go`, `cmd/genwish/models_list.go`, `cmd/genwish/login.go`, `cmd/genwish/logout.go` registration

**Interfaces:**

- Consumes: `novelai`, `format`, `cli`.
- Produces:

```go
func newGenImageNovelCommand(deps cli.Deps) *cobra.Command
func newModelsListNovelCommand(deps cli.Deps) *cobra.Command

func Models() []Model

type Model struct {
	ID string
	Kind string // "image"
}
```

- [ ] Copy the NovelAI package and rewrite imports. Keep the streaming, msgpack, and SSE tests.
- [ ] Add `Novel *novelai.Client` to `cli.Deps` and implement `cli.NovelClient(deps)`.
- [ ] Implement `gen image novel <prompt...>` with the exact flags and behavior from `novel-cli/cmd/novel-cli/gen_image.go`, including `--init-image` local path support, `--init-image-url`, `--denoising-strength` default `0.75`, and `--format` default from `OUTPUT_FORMAT`.
- [ ] Add a static NovelAI model list in `internal/novelai/models.go`:

```text
nai-diffusion-5-full
nai-diffusion-5-full-medium
nai-diffusion-5-full-inpainting
nai-diffusion-5-full-medium-inpainting
nai-diffusion-5-curated
nai-diffusion-4-5-full
nai-diffusion-4-5-full-inpainting
nai-diffusion-4-5-curated
nai-diffusion-4-5-curated-inpainting
nai-diffusion-4-full
nai-diffusion-4-full-inpainting
nai-diffusion-4-curated-preview
nai-diffusion-4-curated-inpainting
nai-diffusion-3
nai-diffusion-3-inpainting
nai-diffusion-furry-3
nai-diffusion-furry-3-inpainting
```

- [ ] Implement `models list novel [image]`:
  - Default kind is image.
  - `video` errors: `novel does not serve video models`.
  - `--query` filters model IDs.
  - `--limit` applies.
- [ ] Add `NovelClient` to `cli` using `provider.Novel`; the error message names `genwish login novel`.
- [ ] Port the NovelAI CLI tests, adapting command paths. Keep `TestGenImageNovel` subtests for: default model, init image from file, init image from URL, denoising default, format conversion, out directory.

**Acceptance Criteria:**

- `genwish gen image novel "a cat"` uses `nai-diffusion-5-full` and writes a WebP by default.
- `genwish gen image novel "a cat" --init-image ./photo.png` reads the local file.
- `genwish login novel` stores a NovelAI key separately from the Nano key; `genwish logout novel` removes only the NovelAI key.
- `genwish models list novel` prints the static list; `genwish models list novel video` errors.
- `go test ./...` passes.

---

## Implementation Unit 7: NovelAI Balance and Anlas

**Files:**

- Create: `internal/novelai/balance.go`
- Create: `internal/novelai/balance_test.go`
- Create: `cmd/genwish/balance_novel.go`
- Modify: `cmd/genwish/balance.go`
- Create: `cmd/genwish/balance_novel_test.go`

**Interfaces:**

- Consumes: `novelai.Client`.
- Produces:

```go
type Balance struct {
	Anlas int64
	Tier  string
}

func (c *Client) CheckBalance(ctx context.Context) (*Balance, error)
```

- [ ] Write failing tests with `httptest`:
  - Object shape: `{"subscription":{"trainingStepsLeft":{"fixedTrainingStepsLeft":9898,"purchasedTrainingSteps":32}}}` → `Anlas == 9930`.
  - Integer shape: `{"subscription":{"trainingStepsLeft":1200}}` → `Anlas == 1200`.
  - Missing `subscription` or `trainingStepsLeft` → error mentioning `trainingStepsLeft`.
  - HTTP 401 → error containing the status.
- [ ] Implement `CheckBalance`:
  - `GET {baseURL}/user/data`.
  - Header: `Authorization: Bearer <key>`.
  - Accept 200 only.
  - Decode `subscription.trainingStepsLeft` as `json.RawMessage`; try an object with `fixedTrainingStepsLeft` and `purchasedTrainingSteps`, then a plain integer.
  - Read `subscription.tier` when present into `Tier`.
- [ ] Implement `balance novel` to print:

```text
Anlas: 9930
```

Print `Tier: <tier>` only when the API reported one.

- [ ] `balance neo` remains an error from Unit 3.

**Acceptance Criteria:**

- `genwish balance novel` prints the Anlas total for a real key.
- Both documented response shapes are covered by tests.
- A missing field and a 401 are reported as errors, not zero.
- Aliases `novel-ai`, `novel ai`, and `novelai` all reach the same command.
- `go test ./...` passes.
- Human verification: run `genwish balance novel` with a real NovelAI key and compare with the NovelAI web UI.

---

## Implementation Unit 8: Video Provider Consistency

**Files:**

- Create: `cmd/genwish/gen_video_neo.go`
- Create: `cmd/genwish/gen_video_novel.go`
- Create: `cmd/genwish/gen_video_provider_test.go`
- Modify: `cmd/genwish/gen_video.go` registration

**Interfaces:**

- Consumes: `provider`.
- Produces:

```go
func newGenVideoNeoCommand(deps cli.Deps) *cobra.Command
func newGenVideoNovelCommand(deps cli.Deps) *cobra.Command
```

- [ ] Write tests asserting `gen video neo "x"` and `gen video novel "x"` return an error containing `video generation is not supported for provider neo` (and `novel`), and that no client is called.
- [ ] Implement both stubs with `cobra.MinimumNArgs(1)` and a `RunE` that returns the unsupported-provider error.
- [ ] Register both under `gen video` alongside `nano`.
- [ ] Add a subtest asserting `gen video` with no subcommand prints help listing `nano`, `neo`, and `novel`.

**Acceptance Criteria:**

- `gen video nano "waves"` still works.
- `gen video neo "waves"` and `gen video novel "waves"` fail clearly and do not generate.
- `gen video` with no provider prints help.
- `go test ./...` passes.

---

## Implementation Unit 9: Documentation Consolidation

**Files:**

- Rewrite: `README.md`
- Create: `docs/README.md`
- Create: `docs/configuration.md`
- Create: `docs/migration.md`
- Create: `docs/nano.md`
- Create: `docs/neo.md`
- Create: `docs/novel.md`
- Create: `docs/booru.md`
- Create: `docs/prompts/neo.md`
- Create: `docs/prompts/novel.md`
- Create: `docs/reference/nano-README.md`
- Create: `docs/reference/neo-README.md`
- Create: `docs/reference/novel-README.md`
- Create: `docs/reference/booru-README.md`
- Create: `.env.example`

**Interfaces:**

- Consumes: the implemented command tree and config.
- Produces: documentation that preserves every source detail.

- [ ] Copy the four source READMEs verbatim into `docs/reference/` so no detail is lost.
- [ ] Copy `neo-mcp/docs/PROMPT.md` and `novel-mcp/docs/PROMPT.md` into `docs/prompts/` and update command names only.
- [ ] Write `README.md` as a short, high-level overview: install, quick start, provider list, command tree, and links to the detailed docs.
- [ ] Write `docs/nano.md` covering every nano generation flag, output behavior, media type handling, `IMAGE_URL_MAP`, and the NanoGPT balance fields.
- [ ] Write `docs/neo.md` covering Forge Neo setup (`--api`), `SD_URL`, `OUTPUT_FORMAT`, `IMAGE_URL_MAP`, every generation flag, every edit command and flag, `bgkill`, `models list neo` from `NEO_MODELS_FILE`, and the Haoming02 Neo fork requirement.
- [ ] Write `docs/novel.md` covering NovelAI login, every generation flag, local `--init-image` versus `--init-image-url`, `models list novel`, `balance novel`, and Anlas.
- [ ] Write `docs/booru.md` covering the tag search arguments, the JSON shape, status values (`ok`, `no_substring_match`, `exact_not_found`, `offset_past_end`, `unknown`), `min_count`, synonyms, implications, and the empty `synonyms`/`implications` behavior of the CLI.
- [ ] Write `docs/configuration.md` with a table of every env var read by `internal/config` and `internal/credential`, its default, and its purpose:

```text
IMAGE_URL_MAP
OUTPUT_FORMAT            default webp
SD_URL                   default http://127.0.0.1:7860
ERROR_DETAIL             default useful
RATE_LIMIT_RPS           default 1
RATE_LIMIT_BURST         default 1
REQUEST_TIMEOUT_SECONDS  default 20
MAX_LIMIT                default 100
MAX_OFFSET               default 1000
NEO_MODELS_FILE          default empty
BLOCKED_TAGS
LOG_LEVEL                default info
GENWISH_CREDENTIAL_BACKEND  default auto
```

- [ ] Document the login-only credentials: `genwish login nano`, `genwish login novel`, and `genwish login danbooru` (login first, then API key). Record that no key environment variable and no user agent variable is read.

- [ ] Write `docs/migration.md` mapping the old commands, env vars, and credential locations to the new ones:
  - `nano-cli gen image` → `genwish gen image nano`.
  - `nano-cli gen video` → `genwish gen video nano`.
  - `nano-cli models list` → `genwish models list nano`.
  - `nano-cli models balance` → `genwish balance nano`.
  - `novel-cli gen image` → `genwish gen image novel`.
  - `novel-cli login` / `logout` → `genwish login novel` / `genwish logout novel`.
  - `neo-cli gen image` → `genwish gen image neo`.
  - `neo-cli edit ...` → `genwish edit ...`.
  - `booru-cli tags` → `genwish util booru`.
  - Removed commands: `serve`, `web *`.
  - Removed env vars: `HOST`, `PORT`, `API_KEY`, `PUBLIC_HOST`, `FILES_DIR`, `NANOGPT_API_KEY`, `NOVELAI_API_KEY`, `DANBOORU_LOGIN`, `DANBOORU_API_KEY`, `USER_AGENT`.
  - Old credential locations: keyring `nano-mcp`/`nanogpt-api-key`, `novel-cli`/`novelai-api-key`; files `~/.config/nano/api_key`, `~/.config/novel/novelai_api_key`; Danbooru was environment-only; backend envs `NANO_CREDENTIAL_BACKEND`, `NOVEL_CREDENTIAL_BACKEND`. New: keyring service `genwish`, users `nanogpt-api-key`, `novelai-api-key`, and `danbooru`; files `~/.config/genwish/{nanogpt_api_key,novelai_api_key,danbooru}`; backend env `GENWISH_CREDENTIAL_BACKEND`.
- [ ] Add `.env.example` listing the new env vars with comments.

**Acceptance Criteria:**

- Every env var read by the code appears in `docs/configuration.md`.
- Every flag in `genwish help` and each subcommand `--help` appears in the provider docs.
- Every source README is preserved under `docs/reference/`.
- No doc tells the user to run `serve`, `web`, Docker, or an MCP client.
- `docs/migration.md` accounts for every removed command and server-only env var.

---

## Implementation Unit 10: Final Verification and Cleanup

**Files:**

- No new production files.
- Move: `docs/plans/GENWISH_INITIAL_PLAN.md` to `docs/plans/done/` only after every acceptance criterion is checked.

**Interfaces:**

- Consumes: every earlier unit.
- Produces: a verified tree and a handoff note.

- [ ] Run `gofmt -l .` and fix every reported file.
- [ ] Run `go build ./...`.
- [ ] Run `go vet ./...`.
- [ ] Run `go test ./... -race -count=1`.
- [ ] Run `go run ./cmd/genwish help` and inspect the tree.
- [ ] Run `go run ./cmd/genwish gen image --help`, `gen image nano --help`, `gen image neo --help`, `gen image novel --help`, `gen video --help`, `models list --help`, `balance --help`, `login --help`, `logout --help`, and `util booru --help`.
- [ ] Grep the tree for `modelcontextprotocol`, `go-chi`, `jsonschema`, `google/uuid`, `nano-mcp`, `neo-mcp`, `novel-cli`, and `booru-mcp`; all must return nothing.
- [ ] Grep for server-only env names (`PUBLIC_HOST`, `FILES_DIR`, `HOST`, `PORT`, `API_KEY`) in Go sources; none may remain except in `docs/migration.md`.
- [ ] Human smoke tests, marked here because they need real services:
  - `genwish login nano`, then `genwish gen image nano "a red fox" --out ./tmp`.
  - `genwish gen image neo "a red fox" --model <checkpoint> --forge-preset <preset> --out ./tmp` with Forge running.
  - `genwish edit square ./some.png --out ./tmp`.
  - `genwish login novel`, then `genwish gen image novel "a red fox" --out ./tmp`.
  - `genwish balance nano` and `genwish balance novel`.
  - `genwish login danbooru` (login first, then key), then `genwish util booru blue_hair --limit 3`.

**Acceptance Criteria:**

- The five automated commands pass.
- No removed-feature identifiers remain.
- The human smoke tests pass.
- The plan is moved to `docs/plans/done/` by the user or agent only after the acceptance criteria are checked.

---

## Self-Review Notes

- Spec coverage: copy, rename, MCP removal, web removal, provider-mandatory generation, provider-mandatory models, top-level balance with aliases, Anlas, util booru, Neo edit preservation, NovelAI generation, and documentation preservation each map to a unit.
- Deliberate decisions not stated verbatim in the brief:
  - `gen image` and `gen video` use provider subcommands, which makes the provider mandatory while preserving each provider's existing flag set.
  - `models` keeps its `list` subcommand: `genwish models list <provider> [image|video]`.
  - `login` and `logout` took a mandatory provider so both credentials can coexist.
  - `genwish models list neo` reads `NEO_MODELS_FILE` (newline-delimited names) instead of calling Forge, because the user wants Neo models configured locally and not checked in.
  - `models list novel` is a curated static list because NovelAI publishes no model-list endpoint.
  - Server-only config is removed; the migration doc records it.
- Type consistency: `provider.Name` is the single provider type; `cli.Deps` is the single injection point; `format.Format` is the single image format type.
- No task implements a feature it does not test, except the human smoke tests, which are explicitly marked as human verification.
