# nowctl usage

## Configuration

`nowctl` needs two things for every ServiceNow request: an **instance**
and a **username** (the password is resolved separately, via the OS
credential store). These are resolved in order, highest priority first:

1. `--instance` / `--username` flags
2. `NOWCTL_INSTANCE` / `NOWCTL_USERNAME` environment variables
3. `instance:` / `username:` keys in the config file (default
   `$HOME/.nowctl.yaml`, override with `--config <path>`)
4. unset — commands that need them will error out

`nowctl auth login` writes the config file for you (see below), so in
practice you set instance/username once and never pass them again.

## Authentication

```bash
nowctl auth login --instance dev12345.service-now.com --username admin
```

This prompts for a password (hidden input on a terminal), verifies it by
calling the instance (`GET /api/now/table/sys_user?sysparm_limit=1`), and
only then:

- stores the password in the OS credential store — Keychain on macOS,
  Credential Manager on Windows, Secret Service (D-Bus) on Linux
- writes `instance` and `username` to the config file

Nothing is written if authentication fails.

```bash
nowctl auth logout [--instance ... --username ...]
```

Removes the stored password for the given (or configured) instance and
username. It does not remove the config file's `instance`/`username`
values.

Credentials are stored per `instance|username` pair, so you can be logged
into multiple instances/users at once — just pass `--instance`/`--username`
to switch between them (or log in again to update the default).

## Commands

### `nowctl table list <table>`

Lists records via `GET /api/now/table/{table}`.

| Flag | Description |
|---|---|
| `--query` | Encoded ServiceNow query (`sysparm_query`), e.g. `active=true^priority=1` |
| `--fields` | Comma-separated list of fields to return (`sysparm_fields`) |
| `--limit` | Max records to return (`sysparm_limit`), default `10` |
| `--format` | `json` (default), `csv`, `xml`, `xlsx` |
| `--output` | Write to this file instead of stdout (required for `xlsx`) |

```bash
nowctl table list incident --limit 25 --query "active=true"
nowctl table list incident --fields "number,short_description,state" --format csv
```

### `nowctl record get <table> <sys_id>`

Fetches a single record via `GET /api/now/table/{table}/{sys_id}`.

| Flag | Description |
|---|---|
| `--fields` | Comma-separated list of fields to return (`sysparm_fields`) |
| `--format` | `json` (default), `csv`, `xml`, `xlsx` |
| `--output` | Write to this file instead of stdout (required for `xlsx`) |

```bash
nowctl record get incident 1c741bd70b2322007518478d83673af3
```

JSON output is array-wrapped (`[{...}]`), consistent with `table list` and
the tabular formats.

## Export formats

All export-capable commands share `--format`/`--output`:

- **json** (default) — pretty-printed, to stdout unless `--output` is set
- **csv** — header row + one row per record; columns are the alphabetically
  sorted union of every field seen across all records, so output is
  deterministic across runs
- **xml** — `<result><record><field>value</field>...</record></result>`
- **xlsx** — one sheet named `Records`; **requires `--output`** since it's a
  binary format and can't be written to stdout

For the tabular formats (csv/xml/xlsx), ServiceNow reference fields — which
the API returns as `{"link": "...", "value": "<sys_id>"}` — are flattened to
just the `sys_id`. The `json` format keeps the full nested object.

```bash
nowctl table list incident --format xlsx --output incidents.xlsx
nowctl table list incident --format xml | less
```

## Troubleshooting

**`Error: no stored credentials for <instance> (<username>): run nowctl auth login`**
The instance/username resolved (from flags/env/config) don't have a
password in the credential store yet. Run `auth login` for that exact
instance/username pair.

**`Error: --format xlsx requires --output <file>...`**
XLSX is a binary zip-based format; it can't be printed to a terminal.
Pass `--output <path>.xlsx`.

**`Error: could not authenticate against <instance>: ...`**
`auth login` failed its verification call — check the instance hostname
and password. Network/DNS errors surface here too.

**Piping JSON output to another tool and getting a parse error**
Make sure you're only capturing stdout — informational messages (like
"Using config file: ...") are written to stderr, but redirecting both
streams together (e.g. `2>&1`) will still mix them.
