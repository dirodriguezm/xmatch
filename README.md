# XWave (xmatch)

Astronomical cross-match service: cone search, bulk cone search, catalog
metadata and light curves over Gaia DR3, AllWISE and eROSITA eRASS1, indexed
with HEALPix.

The repository holds two applications:

| Directory  | What it is                         | Stack                                 |
| ---------- | ---------------------------------- | ------------------------------------- |
| `service/` | HTTP API and catalog indexer       | Go, Gin, SQLite, Parquet, HEALPix C++ |
| `ui/`      | Web interface, deployed on Vercel  | Next.js, React, Ant Design, Tailwind  |

## Service (`service/`)

### Environment

Development uses [Devenv](https://devenv.sh) (Nix, no flakes). With Nix
installed, enter the shell from the repository root:

```bash
devenv shell
```

Shell entry initialises the `healpix` submodule and generates its SWIG
bindings when needed. See [`AGENTS.md`](AGENTS.md) for the auto-activation
hook.

### Commands

All commands are devenv scripts, available inside the shell from anywhere in
the repository:

| Command                           | Description                                    |
| --------------------------------- | ---------------------------------------------- |
| `xwave-build`                     | Build the binary (`service/build/main`)        |
| `xwave-release`                   | Build the fully static release binary          |
| `xwave-run <application> [flags]` | Build and run `server` or `indexer`            |
| `xwave-live-server`               | Run the server with live reload (`air`)        |
| `xwave-test`                      | Run all tests with the race detector           |
| `xwave-docs`                      | Regenerate the Swagger docs in `service/docs/` |
| `xwave-mock`                      | Regenerate mocks with mockery                  |
| `xwave-migrate <db>`              | Run database migrations on `<db>.db`           |

The full list, including the clean-up scripts, is in [`AGENTS.md`](AGENTS.md).

### API

| Method | Path                  | Description                            |
| ------ | --------------------- | -------------------------------------- |
| GET    | `/v1/conesearch`      | Sources within a radius of a position  |
| POST   | `/v1/bulk-conesearch` | Cone search for many positions at once |
| GET    | `/v1/metadata`        | Catalog row for a source id            |
| POST   | `/v1/bulk-metadata`   | Catalog rows for many source ids       |
| GET    | `/v1/lightcurve`      | Light curve around a position          |

The OpenAPI spec lives in [`service/docs/`](service/docs/); a running server
serves Swagger UI at `/swagger/index.html`.

To add a catalog, see [`ADD_NEW_CATALOG.md`](ADD_NEW_CATALOG.md).

## Web interface (`ui/`)

```bash
cd ui
npm install
npm run dev        # http://localhost:3000
```

| Command                  | Description                                           |
| ------------------------ | ----------------------------------------------------- |
| `npm run build`          | Production build                                      |
| `npm test`               | Unit tests (Vitest)                                   |
| `npm run lint`           | ESLint                                                |
| `npm run knip`           | Unused files, exports and dependencies                |
| `npm run generate:types` | Regenerate API types from `service/docs/swagger.json` |

The interface calls the public API at `https://xwave-astro.udp.cl/v1`. A
pre-commit hook runs lint-staged and knip.

## License

Apache License 2.0. See [LICENSE](LICENSE).
