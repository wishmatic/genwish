# Booru

`genwish util booru` searches Danbooru tag names and prints the same text or JSON that the old MCP `tags` tool returned.

```sh
genwish util booru blue_hair
genwish util booru blue hair --limit 5 --category general
genwish util booru blue_hair --json
```

Flags:

| Flag          | Default | Purpose                                                                                     |
| ------------- | ------- | ------------------------------------------------------------------------------------------- |
| `--offset`    | `0`     | How many matching tags to skip.                                                             |
| `--limit`     | `25`    | How many matching tags to return.                                                           |
| `--min-count` | `1`     | Lowest work count a substring match may have. `0` includes the zero-work canonical tail.    |
| `--exact`     | `false` | The search is the whole tag name and aliases resolve to their target.                       |
| `--category`  | all     | Keep only this category, repeatable: `general`, `artist`, `copyright`, `character`, `meta`. |
| `--json`      | `false` | Print JSON instead of text.                                                                 |

The search is a literal substring of the tag name, so `blue hair` and `blue_hair` are the same search.

## Results

Results are ordered by work count, name ascending on ties, and form one page of the match list. Walk the list with `--offset` and stop when `more` is false.

Substring results exclude tags below `--min-count` works by default because Danbooru's tag table also holds concatenated and punctuation-mangled zero-work names. `withheld` and `withheld_best_count` report what the floor removed. `--min-count 0` returns them.

An empty page reports a status:

- `ok`: matches were returned.
- `no_substring_match`: the substring genuinely matched nothing.
- `exact_not_found`: the exact name is absent.
- `offset_past_end`: the offset ran past the last match.
- `unknown`: an unanswerable case was reached.

An offset past the end is never reported as "no such tag". Use `--exact` to ask whether a tag exists; a substring miss proves nothing.

## Synonyms and implications

The MCP server built an alias and implication graph in the background. The CLI does not build that index, so `synonyms` and per-tag `implications` are empty. Counts are read live, and `snapshot_date` records the date they were read.

## Authentication

Danbooru answers anonymously, so `genwish util booru` works with nothing configured. To raise the anonymous rate limit, log in once:

```sh
genwish login danbooru
```

The command prompts for the Danbooru login first, then the API key. `genwish logout danbooru` deletes the stored pair.

## Configuration

See [configuration.md](configuration.md). The relevant settings are `ERROR_DETAIL` and the `booru.rate_limit` table (`rps`, `burst`, `timeout`). The tag page size is fixed at 100 and the largest offset at 1000. The user agent is hardcoded.
