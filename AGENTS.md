# ALeRCE xmatch - Agent Guidelines

This document provides guidelines for AI agents working on the ALeRCE xmatch project. It covers build commands, testing, linting, and code style conventions.

## Project Overview

- **Language**: Go (1.24+)
- **Framework**: Gin web framework, actor pattern
- **Database**: SQLite, Parquet
- **Build tool**: Just (command runner), Nix/Devenv for environment
- **Frontend**: Tailwind CSS (compiled via Tailwind CLI)

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

## Code Style Guidelines

### Imports

Group imports in the following order, separated by a blank line:

1. Standard library
2. Internal packages (project modules)
3. External dependencies

Example:
```go
import (
    "context"
    "fmt"
    "log/slog"

    "github.com/dirodriguezm/xmatch/service/internal/actor"
    "github.com/dirodriguezm/xmatch/service/internal/search/conesearch"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)
```

### Naming Conventions

- **Exported identifiers**: PascalCase
- **Unexported identifiers**: camelCase
- **Acronyms**: Keep all caps (e.g., `SqliteWriter`, `ID`, `Ra`, `Dec`)
- **Interfaces**: Use `er` suffix when appropriate (e.g., `Repository`, `Writer`)
- **Constructors**: `New` for the primary constructor; `NewX` for alternative constructors
- **Test functions**: `TestXxx` where `Xxx` describes the scenario

### Error Handling

- Functions that can fail should return `error` as the last return value.
- Use `if err != nil { return err }` pattern; avoid panics in production code.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Log errors with `slog.Error` providing structured fields.

### Logging

Use `log/slog` for structured logging. Follow these conventions:

- Use `slog.Debug` for debugging information, `slog.Info` for important events, `slog.Error` for errors.
- Provide key‑value pairs as separate arguments after the message:
  ```go
  slog.Info("catalog indexed", "catalog", cat, "rows", n)
  ```

### Context Usage

- Pass `context.Context` as the first parameter to functions that perform I/O, network calls, or long‑running operations.
- Respect cancellation and deadlines from the context.

### Struct Design

- Use pointer receivers for methods that modify the struct or when the struct is large.
- Prefer returning pointers from constructors (`&X{}`).
- Embedding is used for composition (e.g., `gin.Engine`).

### Testing

- Use `testify/assert` for assertions.
- Use `testify/mock` for mocking interfaces (see `just mock`).
- Table‑driven tests are encouraged for multiple test cases.
- Integration tests should be named `*_integration_test.go` and may require external resources.

## Mock Generation

Mocks are generated with `mockery` using the configuration in `service/.mockery.yaml`. Run:

```bash
just mock
```

This creates `mocks.go` files in the same directory as the interface.

## Database Migrations

Migrations are managed with `go-migrate`. The `just migrate` command expects a database file name (without extension) as argument:

```bash
just migrate zeus      # applies migrations to zeus.db
```

Migrations are located in `service/internal/db/migrations/`.

## Environment Variables

- `CONFIG_PATH`: path to the service config file (default `service/config.yaml`)
- `LOG_LEVEL`: log level (`debug`, `info`, `warn`, `error`)
- `ENVIRONMENT`: runtime environment (`local`, `production`)
- `USE_LOGGER`: enable structured logging (`true`/`false`)

## Frontend Assets

Tailwind CSS is used for styling. The source file is `service/ui/static/css/tailwind.css`. It is compiled to `service/ui/static/css/output.css` with the `just build-css` command. The compiled CSS is minified and optimized.

## Submodules

The `healpix` directory is a Git submodule containing a C++ HEALPix library with Go bindings. The Nix environment automatically builds it and sets the necessary `CGO` flags (`CGO_CFLAGS`, `CGO_LDFLAGS`). Do not modify the submodule directly.

## CI/CD

GitHub Actions workflows are in `.github/workflows/`. The `go.yml` workflow runs `just build` and `just test` inside a Nix environment. Ensure any changes pass these checks.

## Commit Messages

Follow conventional commits style: use a type prefix (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`). Start with a lowercase verb and keep the summary line under 72 characters.

## Additional Notes

- The `healpix` Go module is a local replace dependency (points to `../healpix`). Changes to the healpix bindings must be made in that separate module.
- The `dev.db` and `zeus.db` files are SQLite databases used for development; they are ignored by Git.
- The `service/configs/` directory contains YAML configuration files for different catalogs.

---
*Last updated: 2025‑01‑28*