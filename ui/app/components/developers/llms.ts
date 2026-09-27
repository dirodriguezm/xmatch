/**
 * Plain-text docs for language models, following https://llmstxt.org:
 * `/llms.txt` is the short index, `/llms-full.txt` the full API reference.
 * Hand-written from service/docs/swagger.json and the Go handlers
 * (service/internal/api) — keep them in sync when the API changes.
 */

import { API_BASE_URL } from "@/app/lib/api/client";
import {
  CATALOG_META,
  type CatalogMeta,
} from "@/app/lib/constants/catalogMeta";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import {
  REPO_URL,
  SITE_NAME,
  SITE_URL,
  SWAGGER_URL,
} from "@/app/lib/constants/site";

/** M31 — the example target used throughout the docs. */
const EX = { ra: 10.6847, dec: 41.269 };

export const LLMS_EXAMPLES = {
  conesearch: `${API_BASE_URL}/conesearch?ra=${EX.ra}&dec=${EX.dec}&radius=10&catalog=all&nneighbor=5`,
  metadata: `${API_BASE_URL}/metadata?id=0098p408_ac51-043708&catalog=allwise`,
  lightcurve: `${API_BASE_URL}/lightcurve?ra=${EX.ra}&dec=${EX.dec}&radius=5`,
};

const BULK_BODY = `{"ra": [10.6847, 83.6331], "dec": [41.269, 22.0145], "radius": 5, "catalog": "all", "nneighbor": 1}`;

function catalogTitle(m: CatalogMeta): string {
  return m.release.startsWith(m.name) ? m.release : `${m.name} ${m.release}`;
}

function catalogLines(): string {
  return Object.values(CATALOG_META)
    .map(
      (m) =>
        `- [${catalogTitle(m)}](${m.homepage}): \`catalog=${m.slug}\` — ${m.wavelength}; ${m.coverage}; ${m.sources.startsWith("~") ? "" : "~"}${m.sources} sources; astrometry ${m.astrometry}.`
    )
    .join("\n");
}

const SUMMARY = `> ${SITE_NAME} is a free, public astronomical cross-match service. Given a sky position (RA/Dec in degrees, J2000/ICRS) it returns nearby sources from Gaia DR3, AllWISE and eROSITA eRASS1, per-catalog metadata, and multi-survey light curves. The REST API needs no key and returns JSON.`;

export function buildLlmsTxt(): string {
  return `# ${SITE_NAME}

${SUMMARY}

Base URL: ${API_BASE_URL}
All radii are in arcseconds. Keep radius ≤ ${MAX_RADIUS_ARCSEC}″ — larger cones may never return. HTTP 204 (empty body) means "no sources found", not an error.

## API endpoints

- [Cone search](${LLMS_EXAMPLES.conesearch}): \`GET /conesearch\` — params \`ra\`, \`dec\` (deg, required), \`radius\` (arcsec, required), \`catalog\` (all|gaia|allwise|erosita, default all), \`nneighbor\` (max sources per catalog, default 1), \`getMetadata\` (bool). Returns \`[{catalog, data: [{id, ra, dec, cat, ipix, distance}]}]\`, distance in arcsec.
- [Metadata](${LLMS_EXAMPLES.metadata}): \`GET /metadata\` — params \`id\` (source id from a cone search), \`catalog\` (gaia|allwise|erosita). Returns the catalog row (photometry, astrometry, flags).
- [Light curve](${LLMS_EXAMPLES.lightcurve}): \`GET /lightcurve\` — params \`ra\`, \`dec\`, \`radius\` (arcsec), \`catalog\` (all|ztf|neowise|allwise), \`nneighbor\`. Returns \`{detections, non_detections, forced_photometry}\`, each entry \`{catalog, id, object_id, mjd, mag, magerr, data}\`.
- [Bulk cone search](${SITE_URL}/developers): \`POST /bulk-conesearch\` — JSON body \`${BULK_BODY}\`. Same response as cone search plus \`index\` (position in the input arrays).
- [Bulk metadata](${SWAGGER_URL}): \`POST /bulk-metadata\` — JSON body \`{"ids": ["…"], "catalog": "allwise"}\`.

## Catalogs

${catalogLines()}

## Docs

- [Full API reference for LLMs](${SITE_URL}/llms-full.txt): every endpoint, parameter, response schema and error, in plain text
- [API playground](${SITE_URL}/developers): build requests, copy curl/Python/JS, run them live
- [Swagger / OpenAPI](${SWAGGER_URL}): machine-readable spec
- [Catalogs](${SITE_URL}/catalogs): releases, coverage, citations and known issues
- [Learn & FAQ](${SITE_URL}/learn): choosing a radius, interpreting matches

## Optional

- [Web app search](${SITE_URL}/search?ra=${EX.ra}&dec=${EX.dec}): interactive results for a position (params \`ra\`, \`dec\`)
- [MCP server](${REPO_URL}/tree/main/mcp): \`xwave-mcp\`, exposes the API as tools for AI agents
- [Source code](${REPO_URL}): Go backend and Next.js frontend
`;
}

export function buildLlmsFullTxt(): string {
  return `# ${SITE_NAME} — full API reference

${SUMMARY}

This file documents the public REST API in full. The short index is at ${SITE_URL}/llms.txt; the OpenAPI spec is at ${SWAGGER_URL}.

## Conventions

- Base URL: \`${API_BASE_URL}\` (all paths below are relative to it).
- Authentication: none. No API key, no cookies.
- Coordinates: RA and Dec in decimal degrees, J2000/ICRS. RA 0–360, Dec −90–90. Convert sexagesimal (e.g. "00:42:44.3 +41:16:09") to degrees before calling.
- Radius: arcseconds, for every endpoint including bulk. Stay at or below ${MAX_RADIUS_ARCSEC}″; bigger cones can hang.
- Responses: JSON. \`200\` with a body on success; \`204 No Content\` with an empty body when nothing matched — treat as an empty list; \`400\` with a validation error for bad parameters; \`500\` with a string message on server errors.
- CORS: browser calls are only allowed from the ${SITE_NAME} web origins. Call from a server, script or notebook instead.
- Etiquette: this is a shared academic service with no formal rate limit. Prefer one bulk request over many single ones, keep concurrency modest (≤ 4 parallel requests), and back off on 5xx.

## Catalogs

${catalogLines()}

Source ids look like \`Gaia DR3 381266999950756352\` (gaia) and \`0098p408_ac51-043708\` (allwise). Pass them verbatim (URL-encoded) to /metadata.

## GET /conesearch

Search for sources within a radius of a position. Returns matches grouped by catalog.

Query parameters:
- \`ra\` (number, required): Right Ascension in degrees, 0–360.
- \`dec\` (number, required): Declination in degrees, −90–90.
- \`radius\` (number, required): search radius in arcseconds, > 0.
- \`catalog\` (string, optional, default \`all\`): one of \`all\`, \`gaia\`, \`allwise\`, \`erosita\`.
- \`nneighbor\` (integer, optional, default 1, min 1): maximum number of sources returned per catalog, nearest first. Latency does not depend on it, so 100 is fine.
- \`getMetadata\` (boolean, optional, default false): replace each match with its full catalog row (as returned by /metadata) plus \`distance\`. Note the group's \`catalog\` label is then the display name (e.g. "AllWISE").

Example:
\`\`\`
GET ${LLMS_EXAMPLES.conesearch}
\`\`\`

Response 200:
\`\`\`json
[
  {
    "catalog": "allwise",
    "data": [
      {"id": "0098p408_ac51-043708", "ipix": 45450658820, "ra": 10.6846947, "dec": 41.2689392, "cat": "allwise", "distance": 0.219}
    ]
  },
  {
    "catalog": "gaia",
    "data": [
      {"id": "Gaia DR3 381266999950756352", "ipix": 45450570725, "ra": 10.684112, "dec": 41.266848, "cat": "gaia", "distance": 7.909}
    ]
  }
]
\`\`\`

Fields: \`id\` source identifier; \`ra\`/\`dec\` catalog position in degrees (Gaia at epoch J2016.0); \`cat\` catalog slug; \`ipix\` HEALPix index (nested) used internally; \`distance\` angular separation from the query position in arcseconds.

Response 400 (conesearch.ValidationError): \`{"field": "radius", "reason": "...", "errValue": "..."}\`.

## POST /bulk-conesearch

Cone search many positions in one request with a shared radius. Positions are searched in parallel.

JSON body (BulkConesearchRequest):
- \`ra\` (number[], required): RA values in degrees.
- \`dec\` (number[], required): Dec values in degrees; same length as \`ra\`.
- \`radius\` (number, required): radius in arcseconds.
- \`catalog\` (string, optional, default \`all\`).
- \`nneighbor\` (integer, optional, default 1): max sources per catalog per position.

Example:
\`\`\`
curl -X POST "${API_BASE_URL}/bulk-conesearch" -H "Content-Type: application/json" -d '${BULK_BODY}'
\`\`\`

Response 200: same shape as /conesearch, one group per (position, catalog) with an extra \`index\` field giving the 0-based position in the input arrays. Positions with no match are omitted. 204 if nothing matched at all.

## GET /metadata

Full catalog row for one source.

Query parameters:
- \`id\` (string, required): source identifier as returned by /conesearch.
- \`catalog\` (string, required): \`gaia\`, \`allwise\` or \`erosita\`.

Example:
\`\`\`
GET ${LLMS_EXAMPLES.metadata}
\`\`\`

Response 200 (AllWISE): \`{"id", "cntr", "ra", "dec", "w1mpro", "w1sigmpro", "w2mpro", "w2sigmpro", "w3mpro", "w3sigmpro", "w4mpro", "w4sigmpro", "j_m_2mass", "j_msig_2mass", "h_m_2mass", "h_msig_2mass", "k_m_2mass", "k_msig_2mass"}\` — WISE profile-fit Vega magnitudes W1–W4 (3.4, 4.6, 12, 22 µm) and associated 2MASS J/H/Ks magnitudes with 1σ errors; nulls where not measured.

Response 200 (Gaia): Gaia DR3 gaia_source columns such as \`source_id\`, \`ra_error\`, \`dec_error\` (mas), \`parallax\`, \`parallax_error\` (mas), \`pmra\`, \`pmdec\` (mas/yr), \`ruwe\`, \`phot_g_mean_mag\`, \`phot_bp_mean_mag\`, \`phot_rp_mean_mag\` and flux columns.

Response 204: id not found in that catalog. Response 400 (metadata.ValidationError): \`{"field", "reason", "value"}\`.

## POST /bulk-metadata

JSON body (BulkMetadataRequest): \`{"ids": ["id1", "id2"], "catalog": "allwise"}\`. Returns an array of rows as in /metadata; ids queried in parallel. 204 if none found.

## GET /lightcurve

Time-series photometry near a position, merged across surveys.

Query parameters:
- \`ra\` (number, required): degrees.
- \`dec\` (number, required): degrees.
- \`radius\` (number, required): arcseconds (≈5″ recommended).
- \`catalog\` (string, optional, default \`all\`): \`all\`, \`ztf\`, \`neowise\` or \`allwise\`.
- \`nneighbor\` (integer, optional, default 1, min 1): number of objects to include.

Example:
\`\`\`
GET ${LLMS_EXAMPLES.lightcurve}
\`\`\`

Response 200 (LightcurveResponse):
\`\`\`json
{
  "detections": [
    {"catalog": "neowise", "id": "23895r192-000059", "object_id": "98140801351043708", "mjd": 59211.1665, "mag": 3.866, "magerr": 0.04, "data": {"w1mpro": 3.866, "w2mpro": 6.852, "...": "..."}}
  ],
  "non_detections": [],
  "forced_photometry": []
}
\`\`\`

Entry fields (LightcurveEntry): \`catalog\` survey, \`id\` measurement id, \`object_id\` survey object id, \`mjd\` Modified Julian Date, \`mag\`/\`magerr\` magnitude and error in the survey's primary band, \`data\` the raw survey row (e.g. NEOWISE W1/W2).

## Python quickstart

\`\`\`python
import requests
import pandas as pd

BASE = "${API_BASE_URL}"
r = requests.get(f"{BASE}/conesearch", params={"ra": ${EX.ra}, "dec": ${EX.dec}, "radius": 10, "catalog": "all", "nneighbor": 5}, timeout=60)
r.raise_for_status()
groups = r.json() if r.status_code == 200 else []
df = pd.DataFrame([row for g in groups for row in g["data"]])
print(df.sort_values("distance"))
\`\`\`

## Tips for agents

- Resolve object names (e.g. "M31", "Betelgeuse") to coordinates first, e.g. with CDS Sesame (https://cds.unistra.fr/cgi-bin/nph-sesame/-oJ/SNV?M31), then call /conesearch.
- Pick the radius from catalog astrometry: 1–3″ for Gaia/AllWISE, 15–30″ for eROSITA.
- Link users to the web view: ${SITE_URL}/search?ra=<ra>&dec=<dec>
- An MCP server, \`xwave-mcp\`, wraps these endpoints as tools: ${REPO_URL}/tree/main/mcp
`;
}
