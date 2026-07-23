# nowctl

A cross-platform command-line interface for ServiceNow: list table records,
fetch individual records, and export results to JSON, CSV, XML, or Excel —
all from the terminal, on Windows, Linux, and macOS.

## Features

- **Table queries** — list records from any table with ServiceNow query syntax
- **Record lookup** — fetch a single record by `sys_id`
- **Export** — JSON, CSV, XML, or XLSX, to stdout or a file
- **Secure credentials** — passwords are verified against the instance, then
  stored in the OS credential store (Keychain / Credential Manager / Secret
  Service) — never written to disk in plaintext
- **Zero-flag workflow** — log in once and the instance/username are
  remembered; only the command itself is needed afterwards

## Installation

Requires [Go](https://go.dev) 1.26+.

```bash
git clone https://github.com/karastoyanov/nowcontrol.git
cd nowcontrol
go build -o nowctl .
```

Move the resulting `nowctl` binary onto your `PATH` (e.g. `/usr/local/bin`
on macOS/Linux, or any directory in `%PATH%` on Windows).

## Quick start

```bash
# 1. Authenticate once — verifies credentials, stores the password in the
#    OS credential store, and remembers the instance/username in ~/.nowctl.yaml
nowctl auth login --instance dev12345.service-now.com --username admin

# 2. Everything else needs no flags
nowctl table list incident --limit 10
nowctl table list incident --query "active=true^priority=1"
nowctl record get incident <sys_id>

# 3. Export instead of printing JSON
nowctl table list incident --format csv --output incidents.csv
nowctl table list incident --format xlsx --output incidents.xlsx

# 4. Remove stored credentials
nowctl auth logout
```

See [docs/USAGE.md](docs/USAGE.md) for the full command reference,
configuration precedence, and export details.

## Development

```bash
go build -o nowctl .   # build
go vet ./...           # static checks
gofmt -l .             # formatting check
```

## Project layout

```
cmd/              Cobra command definitions (CLI surface)
internal/auth/    OS credential store integration (go-keyring)
internal/client/  ServiceNow Table API client (GET/POST/PATCH)
internal/export/  JSON/CSV/XML/XLSX renderers
```

## License

GPL-3.0. See [LICENSE](LICENSE).
