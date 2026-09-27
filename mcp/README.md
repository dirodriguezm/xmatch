# xwave-mcp

A [Model Context Protocol](https://modelcontextprotocol.io) server for the
public **XWave** cross-match API (`https://xwave-astro.udp.cl/v1`). It lets AI
agents (Claude Code, Claude Desktop, Cursor, and other MCP clients) resolve
object names and cross-match sky positions against:

| Catalog | Coverage | Recommended radius |
|---|---|---|
| **Gaia DR3** (`gaia`) | Optical astrometry and photometry | 3″ |
| **AllWISE** (`allwise`) | Mid-infrared W1–W4 plus 2MASS JHK; NEOWISE light curves | 3″ |
| **eROSITA eRASS1** (`erosita`) | Soft X-ray, western Galactic hemisphere | 20″ |

Every result includes a link to the XWave web page for that source
(`https://xwave-rho.vercel.app/object/<id>?catalog=<cat>`).

## Install and build

Requires Node.js 20 or newer.

```bash
cd mcp
npm install
npm run build        # writes dist/index.js
npm test             # unit tests for the pure helpers (vitest)
npm run smoke        # end-to-end: spawns the server over stdio and calls every tool against the live API
```

The server uses stdio. Its entry point is `dist/index.js`.

## Configuration

| Env var | Default | Purpose |
|---|---|---|
| `XWAVE_API_URL` | `https://xwave-astro.udp.cl/v1` | XWave REST API base URL. Point it at `http://localhost:8080/v1` to use a local backend. |
| `XWAVE_WEB_URL` | `https://xwave-rho.vercel.app` | Base URL for the object-page links in results |
| `XWAVE_SESAME_URL` | `https://cds.unistra.fr/cgi-bin/nph-sesame/-oI/A` | CDS Sesame name resolver |
| `XWAVE_TIMEOUT_MS` | `30000` | Per-request timeout |

### Claude Code

```bash
claude mcp add xwave -- node /absolute/path/to/xmatch/mcp/dist/index.js
# with a custom API:
claude mcp add xwave -e XWAVE_API_URL=http://localhost:8080/v1 -- node /absolute/path/to/xmatch/mcp/dist/index.js
```

### Claude Desktop

Add this to `claude_desktop_config.json` (macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "xwave": {
      "command": "node",
      "args": ["/absolute/path/to/xmatch/mcp/dist/index.js"],
      "env": { "XWAVE_API_URL": "https://xwave-astro.udp.cl/v1" }
    }
  }
}
```

### Cursor

Add this to `~/.cursor/mcp.json`, or to `.cursor/mcp.json` in a project:

```json
{
  "mcpServers": {
    "xwave": {
      "command": "node",
      "args": ["/absolute/path/to/xmatch/mcp/dist/index.js"]
    }
  }
}
```

## Units

- Coordinates are **J2000 decimal degrees**: RA 0–360, Dec −90 to +90.
- Radii are **arcseconds**, up to 120″. This applies to every tool, including
  `cross_match_list`. The backend's `BulkConesearchRequest` comment says
  "degrees", but the live service treats the bulk radius as arcseconds too
  (checked with curl: a radius of `5` matches a source 0.22″ away, while
  `0.0014`, which is 5″ written in degrees, returns 204).
- Returned `distance_arcsec` values are angular separations in arcseconds.

## Tools

Each tool returns a short text summary and a `structuredContent` JSON payload.
The same JSON is also included as a second text block, for clients that only
show text.

| Tool | What it does |
|---|---|
| `resolve_name` | Converts an object name to RA/Dec using CDS Sesame (SIMBAD, then NED, then VizieR). |
| `cone_search` | Finds sources within `radius_arcsec` of a position, sorted by distance. |
| `search_by_name` | Resolves a name, then runs a cone search at that position. |
| `cross_match_list` | Runs a bulk cross-match on up to 1000 `{name?, ra, dec}` positions and reports the nearest match for each input. |
| `get_object` | Returns the full catalog record for an `id` and `catalog` (Gaia astrometry and photometry, AllWISE and 2MASS magnitudes, eROSITA fluxes). |
| `get_lightcurve` | Returns time-domain photometry near a position (currently NEOWISE W1), with a per-survey summary. |

`cone_search`, `search_by_name` and `cross_match_list` accept `catalog`
(`all` | `gaia` | `allwise` | `erosita`), `nneighbor` and `per_catalog`. With
`catalog=all`, the API returns the *N nearest sources across all catalogs*, so a
bright AllWISE source can crowd out Gaia and eROSITA matches. `per_catalog=true`
is the default: it queries each catalog separately so each one reports its own
nearest neighbours.

It also provides:

- **Resource** `xwave://catalogs`: JSON describing each catalog, its id format and its recommended radius.
- **Prompt** `multiwavelength_summary(name)`: asks the agent to run search, then object lookup, then light curve, and to summarise the object across wavelengths.

### Examples

```jsonc
// resolve_name
{ "name": "Betelgeuse" }
// → Betelgeuse → RA 88.792939, Dec +7.407064 (J2000, deg) — main id * alf Ori, type s*r, via Simbad

// cone_search
{ "ra": 10.6847, "dec": 41.2690, "radius_arcsec": 5 }
// → 1 source(s): AllWISE 1
//   - [AllWISE] 0098p408_ac51-043708 at 0.22" … https://xwave-rho.vercel.app/object/0098p408_ac51-043708?catalog=allwise

// search_by_name
{ "name": "Betelgeuse", "radius_arcsec": 5, "catalog": "all" }

// cross_match_list
{ "positions": [ { "name": "M31", "ra": 10.6847, "dec": 41.269 },
                 { "name": "Betelgeuse", "ra": 88.7929, "dec": 7.4071 } ],
  "radius_arcsec": 3 }

// get_object
{ "id": "Gaia DR3 381266999950756352", "catalog": "gaia" }

// get_lightcurve
{ "ra": 10.6847, "dec": 41.269, "radius_arcsec": 3, "max_points": 200 }
```

Things you can ask an agent once the server is connected:

- "Is there an X-ray counterpart to Betelgeuse in eROSITA?"
- "Cross-match these 20 coordinates against Gaia and AllWISE within 2 arcsec."
- "Show me the NEOWISE light curve of the AllWISE source closest to M31's nucleus."

## Errors

- If an input is out of range (for example `ra` > 360, or a radius above 120″ or ≤ 0), the tool call fails with a message that names the problem.
- If the API returns HTTP 204 (no sources), the tool returns an empty result, not an error.
- If the API or Sesame returns an error, times out or can't be reached, the tool returns `isError: true` with the HTTP status and the reason the backend gave.

## Layout

```
src/index.ts         stdio entry point
src/server.ts        McpServer: tools, resource, prompt
src/client.ts        HTTP client (fetch with AbortController timeout, 204 handling)
src/helpers.ts       pure helpers: validation, Sesame parsing, formatting
src/helpers.test.ts  vitest unit tests
scripts/smoke.mjs    end-to-end JSON-RPC smoke test
```
