# nowctl usage

## Configuration

`nowctl` needs two things for every ServiceNow request: an **instance**
and a **username** (the password is resolved separately, from the config
file). These are resolved in order, highest priority first:

1. `--instance` / `--username` flags
2. `NOWCTL_INSTANCE` / `NOWCTL_USERNAME` environment variables
3. `instance:` / `username:` keys in the config file (default
   `$HOME/.nowctl.yaml`, override with `--config <path>`)
4. unset — commands that need them will error out

`nowctl auth login` writes the config file for you (see below), so in
practice you set instance/username once and never pass them again.

Pass `--verbose` on any command to print diagnostic info, such as which
config file was loaded — this is suppressed by default to keep output
quiet.

## Authentication

```bash
nowctl auth login --instance dev12345.service-now.com --username admin
```

This prompts for a password (hidden input on a terminal), verifies it by
calling the instance (`GET /api/now/table/sys_user?sysparm_limit=1`), and
only then writes `instance`, `username`, and the password to the config
file (default `$HOME/.nowctl.yaml`), locking its permissions down to
owner-only read/write (`0600`) since it now holds a secret.

Nothing is written if authentication fails.

The password is stored in plaintext, protected only by that file
permission (and whatever disk encryption your OS provides) — there's no
OS keychain involved, by design: it needs to work the same on a headless
Linux VM, a container, or CI, none of which reliably have a D-Bus session
or Secret Service provider available. Treat `~/.nowctl.yaml` like any
other credentials file (e.g. `~/.aws/credentials`, `~/.netrc`) — don't
commit it, don't copy it around loosely.

```bash
nowctl auth logout [--instance ... --username ...]
```

Removes the stored password for the given (or configured) instance and
username. It does not remove the config file's `instance`/`username`
values.

Credentials are stored per `instance|username` pair under `credentials:`
in the config file, so you can be logged into multiple instances/users at
once — just pass `--instance`/`--username` to switch between them (or log
in again to update the default).

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

### `nowctl table describe <table>`

Lists a table's fields via `sys_dictionary`, including fields inherited
from parent tables in the class hierarchy (e.g. `incident` inherits
`number`, `short_description`, `state`, etc. from `task`) -- resolved by
walking `sys_db_object.super_class`.

Each field's `internal_type` and `reference` are shown as display values
(e.g. `Reference`, `Configuration Item`) rather than raw sys_ids.

| Flag | Description |
|---|---|
| `--format` | `json` (default), `csv`, `xml`, `xlsx` |
| `--output` | Write to this file instead of stdout (required for `xlsx`) |

```bash
nowctl table describe incident
nowctl table describe incident --format csv
```

Requires read access to `sys_dictionary` and `sys_db_object`; some
ServiceNow roles restrict these, which surfaces as a `403` error.

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

### `nowctl record create <table>`

Creates a record via `POST /api/now/table/{table}`.

| Flag | Description |
|---|---|
| `--field key=value` | A single field (repeatable) |
| `--data '<json>'` | Inline JSON object of fields |
| `--file <path.json>` | JSON file of fields |
| `--format` | `json` (default), `csv`, `xml`, `xlsx` |
| `--output` | Write to this file instead of stdout (required for `xlsx`) |

`--field`, `--data`, and `--file` can be combined; on a key conflict,
later wins in this order: `--file` → `--data` → `--field`. At least one
must be given.

```bash
nowctl record create incident --field short_description="VPN down" --field priority=1
nowctl record create incident --data '{"short_description":"VPN down","priority":"1"}'
nowctl record create incident --file new-incident.json
```

The created record's `sys_id` is printed to stderr as confirmation, and
the record itself is written via `--format`/`--output` like `record get`.

### `nowctl record update <table> <sys_id>`

Patches a record via `PATCH /api/now/table/{table}/{sys_id}`, updating
only the fields supplied — everything else on the record is left alone.
Takes the same `--field`/`--data`/`--file` and `--format`/`--output`
flags as `record create`.

```bash
nowctl record update incident <sys_id> --field state=2 --field work_notes="picked up"
```

Note: journal fields like `work_notes` and `comments` are write-only
through the Table API — you can set them here, but reading them back via
`record get`/`table list` returns empty (ServiceNow stores the actual
entries in a separate journal table).

### `nowctl record delete <table> <sys_id>`

Deletes a record via `DELETE /api/now/table/{table}/{sys_id}`. This is
destructive and irreversible.

| Flag | Description |
|---|---|
| `--yes` / `-y` | Skip the confirmation prompt |

```bash
nowctl record delete incident <sys_id>          # prompts "Delete incident record <sys_id>? [y/N]"
nowctl record delete incident <sys_id> --yes     # no prompt
```

When stdin isn't a terminal (piped/scripted) and `--yes` isn't passed,
the command errors immediately instead of hanging on a prompt no one can
answer.

### `nowctl doctor`

Checks, in order: the config file (exists? permissions `0600`?), the
resolved instance and username, whether a password is stored for that
pair, and finally live authentication via `GET /api/now/table/sys_user`
(the same call `auth login` uses to verify credentials). Prints an
`OK`/`FAIL` line per check and exits non-zero if anything failed --
handy for scripting or just confirming a fresh setup works before
running real commands.

```bash
nowctl doctor
```

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

## Table aliases

Long CMDB table names like `cmdb_ci_linux_server` get tedious to type.
Aliases give them a short name, stored under `aliases` in the config file
(default `$HOME/.nowctl.yaml`) alongside `instance`/`username`.

```bash
nowctl alias set computer cmdb_ci_computer
nowctl alias set linux_server cmdb_ci_linux_server
nowctl alias list
nowctl alias remove linux_server
```

Once set, an alias can be used wherever a `<table>` argument is expected —
`table list`, `record get/create/update/delete`:

```bash
nowctl table list computer --limit 5
nowctl record get computer <sys_id>
```

The command prints `Resolved alias "computer" -> "cmdb_ci_computer"` to
stderr so it's clear which real table a command actually ran against.

You can also edit the config file directly:

```yaml
instance: dev12345.service-now.com
username: admin
aliases:
  computer: cmdb_ci_computer
  linux_server: cmdb_ci_linux_server
```

## Troubleshooting

**`Error: no stored credentials for <instance> (<username>): run nowctl auth login`**
The instance/username resolved (from flags/env/config) don't have a
password in the config file yet. Run `auth login` for that exact
instance/username pair.

**`Error: --format xlsx requires --output <file>...`**
XLSX is a binary zip-based format; it can't be printed to a terminal.
Pass `--output <path>.xlsx`.

**`Error: could not authenticate against <instance>: ...`**
`auth login` failed its verification call — check the instance hostname
and password. Network/DNS errors surface here too.

**Piping JSON output to another tool and getting a parse error**
Make sure you're only capturing stdout — informational messages (e.g.
`--verbose`'s "Using config file: ...", or the "Resolved alias ..." note)
are written to stderr, but redirecting both streams together (e.g.
`2>&1`) will still mix them.
