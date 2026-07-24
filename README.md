# nowctl

A cross-platform command-line interface for ServiceNow: list table records,
fetch individual records, and export results to JSON, CSV, XML, or Excel —
all from the terminal, on Windows, Linux, and macOS.

## Features

- **Table queries** — list records from any table with ServiceNow query syntax
- **Record CRUD** — get, create, update, and delete individual records by `sys_id`
- **Export** — JSON, CSV, XML, or XLSX, to stdout or a file
- **Credentials** — passwords are verified against the instance, then stored
  in the config file (`$HOME/.nowctl.yaml`, permissions locked to owner-only)
  — no OS keychain dependency, so it works the same on a headless Linux VM,
  a container, or a desktop
- **Zero-flag workflow** — log in once and the instance/username are
  remembered; only the command itself is needed afterwards
- **Table aliases** — define short names for long table names (e.g. `computer`
  for `cmdb_ci_computer`); usable anywhere a `<table>` argument is expected

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap karastoyanov/nowctl https://github.com/karastoyanov/NowControl
brew install nowctl

# to update later:
brew upgrade nowctl
```

### Debian/Ubuntu (apt repository)

A self-hosted apt repository is published on every release, so `apt`
handles install/upgrade/removal like any other package:

```bash
curl -fsSL https://karastoyanov.github.io/NowControl/nowctl-archive-keyring.asc \
  | sudo gpg --dearmor -o /usr/share/keyrings/nowctl-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/nowctl-archive-keyring.gpg] https://karastoyanov.github.io/NowControl stable main" \
  | sudo tee /etc/apt/sources.list.d/nowctl.list
sudo apt update
sudo apt install nowctl

# to update later:
sudo apt update && sudo apt upgrade nowctl

# to remove:
sudo apt remove nowctl
```

Prefer a one-off install without adding a repo? Download the `.deb`
directly from the [latest release](https://github.com/karastoyanov/NowControl/releases/latest)
and `sudo apt install ./nowctl_<version>_linux_<arch>.deb` — but then
you're back to manually re-downloading for updates.

### Manual (any platform)

Download the archive for your platform from the
[latest release](https://github.com/karastoyanov/NowControl/releases/latest),
extract it, and move the `nowctl` binary onto your `PATH`. To update,
repeat these steps with the new version's archive.

```bash
# macOS (Apple Silicon) example — swap in your OS/arch as needed
curl -LO https://github.com/karastoyanov/NowControl/releases/latest/download/nowctl_<version>_darwin_arm64.tar.gz
tar -xzf nowctl_<version>_darwin_arm64.tar.gz
sudo mv nowctl /usr/local/bin/
nowctl --version
```

Windows/Linux users: download the matching `.zip`/`.tar.gz` from the
releases page and put `nowctl`(`.exe`) anywhere on `%PATH%`/`$PATH`.

### Build from source

Requires [Go](https://go.dev) 1.26+.

```bash
git clone https://github.com/karastoyanov/NowControl.git
cd NowControl
go build -o nowctl .
```

## Quick start

```bash
# 1. Authenticate once — verifies credentials, then stores instance,
#    username, and password in ~/.nowctl.yaml (permissions locked to owner-only)
nowctl auth login --instance dev12345.service-now.com --username admin

# 2. Everything else needs no flags
nowctl table list incident --limit 10
nowctl table list incident --query "active=true^priority=1"
nowctl record get incident <sys_id>

# 3. Create, update, delete
nowctl record create incident --field short_description="VPN down" --field priority=1
nowctl record update incident <sys_id> --field state=2
nowctl record delete incident <sys_id>

# 4. Export instead of printing JSON
nowctl table list incident --format csv --output incidents.csv
nowctl table list incident --format xlsx --output incidents.xlsx

# 5. Alias long table names for less typing
nowctl alias set computer cmdb_ci_computer
nowctl table list computer --limit 5

# 6. Remove stored credentials
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

### Release process

Pushing a `vX.Y.Z` tag triggers `.github/workflows/release.yml`:

1. **goreleaser** builds archives for all platforms, `.deb` packages, and
   the GitHub Release (with a custom header/footer -- edit `release.header`/
   `release.footer` in `.goreleaser.yaml` for per-version notes), and pushes
   a Homebrew formula to this repo's `homebrew` branch.
2. **apt-repo** downloads the `.deb` assets just published, rebuilds the
   apt repository with `aptly` (re-importing every previously published
   `.deb` from `gh-pages/pool` first, so old versions stay installable),
   signs it with the `APT_GPG_PRIVATE_KEY` secret, and pushes the result to
   the `gh-pages` branch, served at https://karastoyanov.github.io/NowControl/.

The apt-signing key is a dedicated, passphrase-less keypair used only for
this repo -- not anyone's personal GPG identity. To rotate it: generate a
new key, update the `APT_GPG_PRIVATE_KEY` secret, and republish
`nowctl-archive-keyring.asc` on `gh-pages` (existing users need to re-run
the `curl | gpg --dearmor` step from the install instructions).

## Project layout

```
cmd/              Cobra command definitions (CLI surface)
internal/client/  ServiceNow Table API client (GET/POST/PATCH)
internal/export/  JSON/CSV/XML/XLSX renderers
```

## License

GPL-3.0. See [LICENSE](LICENSE).
