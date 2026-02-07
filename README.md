# CrossWave Xmatch

API for astronomical cross-matching and catalog queries developed by [ALeRCE](https://alerce.science).

## Overview

CrossWave provides fast cone search and metadata retrieval across multiple astronomical catalogs. The service is optimized for high-throughput queries using HEALPix spatial indexing.

### Key Features

- **Cone Search**: Find objects within a radius of given celestial coordinates
- **Bulk Operations**: Process multiple queries in a single request
- **Metadata Retrieval**: Get detailed catalog information for specific objects
- **Lightcurves**: Retrieve time-series photometry data

## Prerequisites

- [Nix](https://nixos.org/) with flakes enabled (used via [devenv](https://devenv.sh/) for reproducible development)
- [Go](https://go.dev/) 1.24+
- [just](https://github.com/casey/just) command runner
- [swag](https://github.com/swaggo/swag) for Swagger docs generation (`go install github.com/swaggo/swag/cmd/swag@latest`)

## Getting Started

```bash
# Build the service
just build

# Run the HTTP server
just run http-server

# Run the test suite
just test

# Regenerate Swagger documentation
just docs
```

### Error Handling

All validation errors return a JSON object with the following structure:

```json
{
  "field": "ra",
  "reason": "value out of range",
  "value": "400.0"
}
```

### Rate Limits

- Bulk endpoints accept up to **1000 coordinates** per request
- Results are processed in parallel for optimal performance

## API Reference (Swagger)

The full OpenAPI/Swagger specification is available at [`service/docs/swagger.yaml`](service/docs/swagger.yaml). When the HTTP server is running, interactive Swagger UI documentation is served at the `/swagger/index.html` endpoint.

### Generating the docs

From the `service/` directory, run:

```bash
swag init -g cmd/start_http_server.go -o docs --md docs/markdown
```

| Flag | Description |
|------|-------------|
| `-g` | Entry point file with the general API annotations |
| `-o` | Output directory for the generated `docs.go`, `swagger.json`, and `swagger.yaml` |
| `--md` | Directory containing markdown files used for tag descriptions |

Or, using `just` from the repository root:

```bash
just docs
```

## Project Structure

```
xmatch/
├── service/                  # Main Go service
│   ├── cmd/                  # Application entry points (HTTP server, catalog indexer)
│   ├── internal/
│   │   ├── api/              # HTTP API handlers and routing
│   │   ├── search/
│   │   │   ├── conesearch/   # Cone search logic
│   │   │   ├── knn/          # K-nearest neighbors search
│   │   │   ├── lightcurve/   # Lightcurve retrieval
│   │   │   └── metadata/     # Metadata service
│   │   ├── catalog_indexer/  # Catalog ingestion (CSV, FITS, Parquet readers)
│   │   ├── repository/       # Database models and queries
│   │   ├── config/           # Configuration management
│   │   ├── db/               # Database migrations
│   │   ├── di/               # Dependency injection
│   │   └── web/              # Web UI (templates, i18n, middleware)
│   ├── docs/                 # Generated Swagger docs
│   └── ui/                   # Frontend templates and static assets
├── healpix/                  # HEALPix spatial indexing library (Go + C++ via SWIG)
├── deploy/                   # Deployment tooling and systemd service
└── justfile                  # Task runner commands
```

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
