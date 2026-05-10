# ALeRCE xmatch - Agent Guidelines

## Project Overview

- **Language**: Go (1.24+)
- **Framework**: Gin web framework, actor pattern
- **Database**: SQLite, Parquet
- **Build tool**: Just (command runner), Nix/Devenv for environment

## Environment Setup

The project uses Nix/Devenv for reproducible development environments. Ensure you have Nix installed and then run:

```bash
nix develop
```

This will drop you into a shell with all dependencies (Go, golangci-lint, healpix libraries, etc.) and run the `init-healpix` script automatically.

## Build Commands

All commands are defined in the `justfile`. Most commands assume you are in the `service/` directory (they set `working-directory`). Use `just` to run them from the repository root.

| Command | Description | Working Directory |
|---------|-------------|-------------------|
| `just build` | Build the Go binary (`build/main`) | `service/` |
| `just run application flags=''` | Build and run a specific application (e.g., `server`, `indexer`) | `service/` |
| `just live-server` | Run with `air` for live reload | `service/` |
| `just build-css` | Compile Tailwind CSS | `service/` |
| `just build-css-watch` | Watch and compile CSS | `service/` |
| `just docs` | Generate Swagger documentation | `service/` |
| `just mock` | Generate mocks with mockery | `service/` |
| `just migrate db` | Run database migrations on `db`.db | root |
| `just clean-build` | Remove `service/build/` | root |
| `just clean-all` | Clean Go caches and build artifacts | `service/` |
| `just clean-db db` | Remove the specified database file | root |

## Testing

| Command | Description |
|---------|-------------|
| `just test` | Run all tests with `grc` colorization and race detector |
| `just test-verbose` | Run all tests verbosely with race detector |

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

