import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import type { CallToolResult } from "@modelcontextprotocol/sdk/types.js";
import * as z from "zod";
import { ApiError, XWaveClient, type ClientConfig } from "./client.js";
import {
  CATALOGS,
  CATALOG_INFO,
  DEFAULT_WEB_URL,
  MAX_LIST_SIZE,
  MAX_NNEIGHBOR,
  MAX_RADIUS_ARCSEC,
  ValidationError,
  catalogLabel,
  compactMetadata,
  flattenMatches,
  fmtCoord,
  groupBulkByIndex,
  objectUrl,
  summarizeLightcurve,
  summarizeMatches,
  validateNneighbor,
  validatePosition,
  validateRadius,
  type ApiCatalogResult,
  type Catalog,
  type CatalogOrAll,
  type Match,
} from "./helpers.js";

export const SERVER_VERSION = "0.1.0";

// ---------- shared zod fields ----------
const raField = z.number().min(0).max(360).describe("Right ascension, J2000, decimal degrees (0-360)");
const decField = z.number().min(-90).max(90).describe("Declination, J2000, decimal degrees (-90 to 90)");
const radiusField = (def: number) =>
  z
    .number()
    .positive()
    .max(MAX_RADIUS_ARCSEC)
    .default(def)
    .describe(`Search radius in ARCSECONDS (max ${MAX_RADIUS_ARCSEC}). Recommended: Gaia/AllWISE 3, eROSITA 20.`);
const catalogAllField = z
  .enum(["all", ...CATALOGS])
  .default("all")
  .describe("Catalog to search: all | gaia (Gaia DR3) | allwise (AllWISE) | erosita (eROSITA eRASS1)");
const nneighborField = (def: number) =>
  z
    .number()
    .int()
    .min(1)
    .max(MAX_NNEIGHBOR)
    .default(def)
    .describe("Maximum number of neighbours to return per position (per catalog when per_catalog=true)");
const perCatalogField = z
  .boolean()
  .default(true)
  .describe(
    "When catalog=all, query each catalog separately so every catalog gets its own nearest neighbours. " +
      "If false, the API returns the N nearest sources overall, which can hide a catalog entirely.",
  );

// ---------- result helpers ----------
function ok(text: string, structured: Record<string, unknown>): CallToolResult {
  return {
    content: [
      { type: "text", text },
      // Some clients only surface `content`; mirror the structured payload as JSON text too.
      { type: "text", text: "```json\n" + JSON.stringify(structured, null, 2) + "\n```" },
    ],
    structuredContent: structured,
  };
}

function fail(err: unknown): CallToolResult {
  const msg =
    err instanceof ValidationError
      ? `Invalid input: ${err.message}`
      : err instanceof ApiError
        ? `XWave request failed: ${err.message}`
        : `Unexpected error: ${(err as Error)?.message ?? String(err)}`;
  return { content: [{ type: "text", text: msg }], isError: true };
}

async function guarded(fn: () => Promise<CallToolResult>): Promise<CallToolResult> {
  try {
    return await fn();
  } catch (err) {
    return fail(err);
  }
}

// ---------- core operations (shared by several tools) ----------
async function doConeSearch(
  client: XWaveClient,
  webBase: string,
  ra: number,
  dec: number,
  radius: number,
  catalog: CatalogOrAll,
  nneighbor: number,
  perCatalog: boolean,
): Promise<Match[]> {
  validatePosition(ra, dec, radius);
  validateNneighbor(nneighbor);
  let raw: ApiCatalogResult[];
  if (catalog === "all" && perCatalog) {
    const parts = await Promise.all(CATALOGS.map((c) => client.coneSearch(ra, dec, radius, c, nneighbor)));
    raw = parts.flat();
  } else {
    raw = await client.coneSearch(ra, dec, radius, catalog, nneighbor);
  }
  return flattenMatches(raw, webBase);
}

function coneStructured(ra: number, dec: number, radius: number, catalog: string, matches: Match[]) {
  return {
    query: { ra, dec, radius_arcsec: radius, catalog },
    count: matches.length,
    matches,
  };
}

// ---------- server ----------
export function createServer(cfg: ClientConfig, webBase = DEFAULT_WEB_URL): McpServer {
  const client = new XWaveClient(cfg);
  const server = new McpServer(
    { name: "xwave", version: SERVER_VERSION },
    {
      instructions:
        "XWave cross-matches sky positions against Gaia DR3 (optical), AllWISE (mid-IR) and eROSITA eRASS1 (X-ray). " +
        "Coordinates are J2000 decimal degrees; radii are ARCSECONDS (max 120). " +
        "Use resolve_name or search_by_name for object names, cone_search for one position, cross_match_list for many. " +
        "Read xwave://catalogs for catalog descriptions and recommended radii.",
    },
  );

  // resolve_name
  server.registerTool(
    "resolve_name",
    {
      title: "Resolve object name",
      description:
        "Resolve an astronomical object name (e.g. 'M31', 'Betelgeuse', 'NGC 1068', 'SN 1987A') to J2000 RA/Dec in decimal degrees using CDS Sesame (SIMBAD/NED/VizieR).",
      inputSchema: { name: z.string().min(1).describe("Object name or identifier") },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ name }) =>
      guarded(async () => {
        const r = await client.resolveName(name);
        if (!r) {
          return {
            content: [{ type: "text", text: `Sesame could not resolve "${name}". Check the spelling or supply coordinates directly.` }],
            isError: true,
          };
        }
        const extra = [r.mainId && `main id ${r.mainId}`, r.objectType && `type ${r.objectType}`, r.source && `via ${r.source}`]
          .filter(Boolean)
          .join(", ");
        return ok(`${name} → ${fmtCoord(r.ra, r.dec)} (J2000, deg)${extra ? ` — ${extra}` : ""}`, { ...r });
      }),
  );

  // cone_search
  server.registerTool(
    "cone_search",
    {
      title: "Cone search",
      description:
        "Find catalog sources within a radius (arcsec) of a sky position. Returns sources sorted by angular distance, each with an XWave web link.",
      inputSchema: {
        ra: raField,
        dec: decField,
        radius_arcsec: radiusField(5),
        catalog: catalogAllField,
        nneighbor: nneighborField(5),
        per_catalog: perCatalogField,
      },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ ra, dec, radius_arcsec, catalog, nneighbor, per_catalog }) =>
      guarded(async () => {
        const matches = await doConeSearch(client, webBase, ra, dec, radius_arcsec, catalog, nneighbor, per_catalog);
        const text = `Cone search at ${fmtCoord(ra, dec)}, r=${radius_arcsec}", catalog=${catalog}\n${summarizeMatches(matches)}`;
        return ok(text, coneStructured(ra, dec, radius_arcsec, catalog, matches));
      }),
  );

  // search_by_name
  server.registerTool(
    "search_by_name",
    {
      title: "Search by object name",
      description: "Resolve an object name with CDS Sesame, then run a cone search around the resolved position.",
      inputSchema: {
        name: z.string().min(1).describe("Object name, e.g. 'Betelgeuse' or 'M31'"),
        radius_arcsec: radiusField(5),
        catalog: catalogAllField,
        nneighbor: nneighborField(5),
        per_catalog: perCatalogField,
      },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ name, radius_arcsec, catalog, nneighbor, per_catalog }) =>
      guarded(async () => {
        const r = await client.resolveName(name);
        if (!r) {
          return {
            content: [{ type: "text", text: `Sesame could not resolve "${name}". Try cone_search with explicit coordinates.` }],
            isError: true,
          };
        }
        const matches = await doConeSearch(client, webBase, r.ra, r.dec, radius_arcsec, catalog, nneighbor, per_catalog);
        const text =
          `${name} resolved to ${fmtCoord(r.ra, r.dec)}${r.mainId ? ` (${r.mainId})` : ""}.\n` +
          `Cone search r=${radius_arcsec}", catalog=${catalog}\n${summarizeMatches(matches)}`;
        return ok(text, { resolved: r, ...coneStructured(r.ra, r.dec, radius_arcsec, catalog, matches) });
      }),
  );

  // cross_match_list
  server.registerTool(
    "cross_match_list",
    {
      title: "Cross-match a list of positions",
      description:
        `Cross-match up to ${MAX_LIST_SIZE} positions in one request (bulk cone search). ` +
        "Returns, for each input, its matches sorted by distance plus the best (nearest) match.",
      inputSchema: {
        positions: z
          .array(
            z.object({
              name: z.string().optional().describe("Optional label echoed back in the output"),
              ra: raField,
              dec: decField,
            }),
          )
          .min(1)
          .max(MAX_LIST_SIZE)
          .describe("Positions to cross-match"),
        radius_arcsec: radiusField(3),
        catalog: catalogAllField,
        nneighbor: nneighborField(1),
        per_catalog: perCatalogField,
      },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ positions, radius_arcsec, catalog, nneighbor, per_catalog }) =>
      guarded(async () => {
        validateRadius(radius_arcsec);
        validateNneighbor(nneighbor);
        positions.forEach((p) => validatePosition(p.ra, p.dec));
        const ras = positions.map((p) => p.ra);
        const decs = positions.map((p) => p.dec);
        let raw: ApiCatalogResult[];
        if (catalog === "all" && per_catalog) {
          const parts = await Promise.all(
            CATALOGS.map((c) => client.bulkConeSearch(ras, decs, radius_arcsec, c, nneighbor)),
          );
          raw = parts.flat();
        } else {
          raw = await client.bulkConeSearch(ras, decs, radius_arcsec, catalog, nneighbor);
        }
        const grouped = groupBulkByIndex(raw, positions.length, webBase);
        const results = positions.map((p, i) => {
          const matches = grouped[i]!;
          return {
            index: i,
            name: p.name ?? null,
            ra: p.ra,
            dec: p.dec,
            n_matches: matches.length,
            best: matches[0] ?? null,
            matches,
          };
        });
        const matched = results.filter((r) => r.n_matches > 0).length;
        const lines = results.slice(0, 50).map((r) => {
          const label = r.name ?? `#${r.index}`;
          if (!r.best) return `- ${label} (${fmtCoord(r.ra, r.dec)}): no match`;
          const b = r.best;
          return `- ${label}: ${r.n_matches} match(es); nearest [${catalogLabel(b.catalog)}] ${b.id} at ${b.distance_arcsec.toFixed(2)}" ${b.url}`;
        });
        if (results.length > 50) lines.push(`... ${results.length - 50} more inputs in structured content.`);
        const text =
          `Cross-matched ${positions.length} position(s), r=${radius_arcsec}", catalog=${catalog}: ` +
          `${matched} with matches, ${positions.length - matched} without.\n` +
          lines.join("\n");
        return ok(text, {
          query: { radius_arcsec, catalog, nneighbor, per_catalog, count: positions.length },
          matched,
          unmatched: positions.length - matched,
          results,
        });
      }),
  );

  // get_object
  server.registerTool(
    "get_object",
    {
      title: "Get object metadata",
      description:
        "Fetch the full catalog record for a source id returned by a search (e.g. Gaia DR3 photometry/astrometry, AllWISE W1-W4 + 2MASS magnitudes, eROSITA fluxes).",
      inputSchema: {
        id: z.string().min(1).describe("Source id exactly as returned by cone_search, e.g. 'Gaia DR3 381266999950756352'"),
        catalog: z.enum(CATALOGS).describe("Catalog the id belongs to: gaia | allwise | erosita"),
      },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ id, catalog }) =>
      guarded(async () => {
        const meta = await client.metadata(id, catalog);
        if (!meta) {
          return {
            content: [{ type: "text", text: `No ${catalogLabel(catalog)} source with id "${id}". Ids must match the search output exactly.` }],
            isError: true,
          };
        }
        const record = compactMetadata(meta);
        const url = objectUrl(id, catalog, webBase);
        const keys = Object.keys(record);
        const text =
          `${catalogLabel(catalog)} ${id}` +
          (typeof record.ra === "number" && typeof record.dec === "number" ? ` at ${fmtCoord(record.ra, record.dec)}` : "") +
          `\n${keys.length} fields. Web page: ${url}\n` +
          highlight(catalog, record);
        return ok(text, { id, catalog, url, metadata: record });
      }),
  );

  // get_lightcurve
  server.registerTool(
    "get_lightcurve",
    {
      title: "Get light curve",
      description:
        "Fetch time-domain photometry around a position (currently NEOWISE single-epoch W1 detections matched to AllWISE sources). Returns a per-survey summary plus the points (MJD, mag, magerr).",
      inputSchema: {
        ra: raField,
        dec: decField,
        radius_arcsec: radiusField(3),
        nneighbor: nneighborField(1),
        max_points: z.number().int().min(0).max(5000).default(500).describe("Maximum number of detections to include in the output (0 = summary only)"),
      },
      annotations: { readOnlyHint: true, openWorldHint: true },
    },
    async ({ ra, dec, radius_arcsec, nneighbor, max_points }) =>
      guarded(async () => {
        validatePosition(ra, dec, radius_arcsec);
        validateNneighbor(nneighbor);
        const lc = await client.lightcurve(ra, dec, radius_arcsec, nneighbor);
        const det = (lc.detections ?? []).map((d) => ({
          catalog: d.catalog,
          id: d.id,
          object_id: d.object_id,
          mjd: d.mjd,
          mag: d.mag,
          magerr: d.magerr,
        }));
        det.sort((a, b) => a.mjd - b.mjd);
        const summary = summarizeLightcurve(det);
        const nNon = (lc.non_detections ?? []).length;
        const nForced = (lc.forced_photometry ?? []).length;
        const text =
          det.length === 0
            ? `No light-curve detections within ${radius_arcsec}" of ${fmtCoord(ra, dec)}.`
            : `Light curve at ${fmtCoord(ra, dec)}, r=${radius_arcsec}": ${det.length} detection(s), ${nNon} non-detection(s), ${nForced} forced-photometry point(s).\n` +
              summary
                .map(
                  (s) =>
                    `- ${s.catalog}: ${s.n} points, MJD ${s.mjd_min.toFixed(1)}–${s.mjd_max.toFixed(1)}, mag ${s.mag_min}–${s.mag_max} (median ${s.mag_median})`,
                )
                .join("\n");
        return ok(text, {
          query: { ra, dec, radius_arcsec, nneighbor },
          summary,
          counts: { detections: det.length, non_detections: nNon, forced_photometry: nForced },
          truncated: det.length > max_points,
          detections: det.slice(0, max_points),
        });
      }),
  );

  // resource: catalogs
  server.registerResource(
    "catalogs",
    "xwave://catalogs",
    {
      title: "XWave catalogs",
      description: "Catalogs available in XWave, with wavelength coverage, id format and recommended match radii",
      mimeType: "application/json",
    },
    async (uri) => ({
      contents: [
        {
          uri: uri.href,
          mimeType: "application/json",
          text: JSON.stringify(
            {
              units: { coordinates: "J2000 decimal degrees", radius: "arcseconds", max_radius_arcsec: MAX_RADIUS_ARCSEC },
              catalogs: Object.values(CATALOG_INFO),
            },
            null,
            2,
          ),
        },
      ],
    }),
  );

  // prompt: multiwavelength summary
  server.registerPrompt(
    "multiwavelength_summary",
    {
      title: "Multi-wavelength summary",
      description: "Look up an object across Gaia, AllWISE and eROSITA and summarise what XWave knows about it",
      argsSchema: { name: z.string().describe("Object name, e.g. 'Betelgeuse'") },
    },
    ({ name }) => ({
      messages: [
        {
          role: "user",
          content: {
            type: "text",
            text:
              `Use the XWave tools to build a multi-wavelength summary of ${name}. ` +
              `1) Call search_by_name with radius_arcsec=5 (use 20 for eROSITA if nothing is found). ` +
              `2) For the nearest match in each catalog, call get_object. ` +
              `3) If there is an AllWISE match, call get_lightcurve at its position. ` +
              `Report positions, separations, key magnitudes/fluxes and variability, and include the XWave links.`,
          },
        },
      ],
    }),
  );

  return server;
}

/** A few human-meaningful fields per catalog for the text summary. */
function highlight(catalog: Catalog, r: Record<string, unknown>): string {
  const pick: Record<Catalog, string[]> = {
    gaia: ["phot_g_mean_mag", "bp_rp", "parallax", "parallax_error", "pmra", "pmdec", "ruwe", "radial_velocity", "teff_gspphot", "phot_variable_flag"],
    allwise: ["w1mpro", "w2mpro", "w3mpro", "w4mpro", "j_m_2mass", "h_m_2mass", "k_m_2mass"],
    erosita: [],
  };
  let keys = pick[catalog].filter((k) => k in r);
  if (keys.length === 0) keys = Object.keys(r).filter((k) => k !== "id" && k !== "ra" && k !== "dec").slice(0, 10);
  return keys.map((k) => `${k}=${fmtVal(r[k])}`).join(", ");
}

function fmtVal(v: unknown): string {
  if (typeof v === "number") return Number.isInteger(v) ? String(v) : String(Math.round(v * 1e4) / 1e4);
  return String(v);
}

