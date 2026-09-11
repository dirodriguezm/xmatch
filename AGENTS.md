# ALeRCE xmatch - Agent Guidelines

## Project Overview

- **Language**: Go (1.24+)
- **Framework**: Gin web framework, actor pattern
- **Database**: SQLite, Parquet
- **Build tool**: Devenv (scripts for tasks, Nix for the environment)

## Environment Setup

The project uses [Devenv](https://devenv.sh) for reproducible development environments (no flakes, no direnv). Ensure you have Nix installed, then enable auto-activation by adding the hook to your shell config (e.g. `~/.bashrc` or `~/.zshrc`):

```bash
eval "$(devenv hook bash)"
```

Then trust the project directory once:

```bash
devenv allow
```

The environment now activates automatically when you enter the project directory. Without the hook, use `devenv shell` to enter it manually. The `xmatch:init-healpix` task runs automatically on shell entry when needed (submodule init + SWIG bindings generation).

## Build Commands

All commands are devenv scripts prefixed with `xwave`, available inside the dev shell (run them from anywhere in the repo; they `cd` into the right directory themselves). Outside the shell, prefix with `devenv shell`, e.g. `devenv shell xwave-build`.

| Command | Description |
|---------|-------------|
| `xwave-build` | Build the Go binary (`service/build/main`) |
| `xwave-run <application> [flags]` | Build and run an application (`server`, `indexer`) |
| `xwave-live-server` | Run with `air` for live reload |
| `xwave-docs` | Generate Swagger documentation |
| `xwave-mock` | Generate mocks with mockery |
| `xwave-migrate <db>` | Run database migrations on `<db>.db` |
| `xwave-clean-build` | Remove `service/build/` |
| `xwave-clean-all` | Clean Go caches and build artifacts |
| `xwave-clean-db <db>` | Remove the specified database file |
| `xwave-init-healpix` | Initialize healpix submodule and generate SWIG bindings manually |

## Testing

| Command | Description |
|---------|-------------|
| `xwave-test` | Run all tests with `grc` colorization and race detector |
| `xwave-test-verbose` | Run all tests verbosely with race detector |
| `devenv shell xwave-test` | Run the test suite from outside the shell (used in CI) |

To run a single test or a specific package, use Go directly:

```bash
cd service
go test ./internal/search/conesearch -v -race
go test -run TestReceive_Mastercat ./internal/catalog_indexer/writer/sqlite
```

Tests use the `testify` framework (assertions, mocks). Integration tests may require a database; see `*_integration_test.go` files.

## Linting and Code Quality

The development shell includes `golangci-lint`. Run it with:

```bash
cd service
golangci-lint run ./...
```

There is no custom `.golangci.yml`; default linting rules apply.

Format Go code with `gofmt`:

```bash
cd service
gofmt -w .
```

## Agent skills

### Issue tracker

Issues tracked as local markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Not used — no formal triage state machine on this repo. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context — `CONTEXT.md` + `docs/adr/` at repo root. See `docs/agents/domain.md`.

---

