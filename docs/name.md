# Names

`genwish gen name` generates names with the
[namegen](https://github.com/wishmatic/namegen) library.

```sh
genwish gen name
genwish gen name --gender female --culture Japanese --count 5
genwish gen name --json
```

Flags:

- `--gender` (empty): `male`, `female`, or `non-binary`; unset picks one at random.
- `--culture` (empty): a supported culture; unset picks one at random. An unknown value lists them.
- `--count` (`1`): how many names to generate, 1 to 100.
- `--json` (false): print a JSON array of `{given_name, surname, gender}`.

Cultures cover the real-world regions and languages namegen ships, plus the `phonemes`, `elvish`,
`khuzdul`, `orkind`, and `fae` constructed sets. Names print one per line as `Given Surname`
unless `--json` is given.
