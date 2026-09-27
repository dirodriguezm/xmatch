/**
 * Pure helpers: validation, parsing and formatting. No I/O here so that
 * everything in this file is unit-testable.
 */

export const CATALOGS = ["gaia", "allwise", "erosita"] as const;
export type Catalog = (typeof CATALOGS)[number];
export type CatalogOrAll = Catalog | "all";

export const MAX_RADIUS_ARCSEC = 120;
export const MAX_LIST_SIZE = 1000;
export const MAX_NNEIGHBOR = 100;

export const DEFAULT_WEB_URL = "https://xwave-rho.vercel.app";

export interface CatalogInfo {
  id: Catalog;
  name: string;
  band: string;
  recommendedRadiusArcsec: number;
  idExample: string;
  notes: string;
}

export const CATALOG_INFO: Record<Catalog, CatalogInfo> = {
  gaia: {
    id: "gaia",
    name: "Gaia DR3",
    band: "Optical (G, BP, RP), astrometry",
    recommendedRadiusArcsec: 3,
    idExample: "Gaia DR3 381266999950756352",
    notes:
      "~1.8 billion sources with sub-mas astrometry. Positions are epoch J2016.0, so high proper-motion stars may need a larger radius.",
  },
  allwise: {
    id: "allwise",
    name: "AllWISE",
    band: "Mid-infrared (W1 3.4um, W2 4.6um, W3 12um, W4 22um) + 2MASS JHK",
    recommendedRadiusArcsec: 3,
    idExample: "0098p408_ac51-043708",
    notes:
      "~750 million sources. Light curves come from NEOWISE single-exposure photometry.",
  },
  erosita: {
    id: "erosita",
    name: "eROSITA eRASS1",
    band: "Soft X-ray (0.2-2.3 keV)",
    recommendedRadiusArcsec: 20,
    idExample: "1eRASS J055510.3+072427",
    notes:
      "~900k X-ray sources from the first all-sky survey (German half of the sky, Galactic longitude > 180 deg). Positional errors are several arcsec, so use ~15-20 arcsec.",
  },
};

export class ValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ValidationError";
  }
}

export function validateRa(ra: number, label = "ra"): void {
  if (!Number.isFinite(ra) || ra < 0 || ra > 360) {
    throw new ValidationError(
      `${label} must be a J2000 right ascension in decimal degrees between 0 and 360 (got ${ra}). ` +
        `If you have sexagesimal (hh:mm:ss) coordinates, convert them first: RA_deg = 15 * (h + m/60 + s/3600).`,
    );
  }
}

export function validateDec(dec: number, label = "dec"): void {
  if (!Number.isFinite(dec) || dec < -90 || dec > 90) {
    throw new ValidationError(
      `${label} must be a J2000 declination in decimal degrees between -90 and 90 (got ${dec}).`,
    );
  }
}

export function validateRadius(radiusArcsec: number): void {
  if (!Number.isFinite(radiusArcsec) || radiusArcsec <= 0) {
    throw new ValidationError(
      `radius_arcsec must be a positive number of arcseconds (got ${radiusArcsec}).`,
    );
  }
  if (radiusArcsec > MAX_RADIUS_ARCSEC) {
    throw new ValidationError(
      `radius_arcsec must be at most ${MAX_RADIUS_ARCSEC} arcsec (got ${radiusArcsec}). ` +
        `Note the unit is arcseconds, not degrees: 0.01 deg = 36 arcsec. Recommended: Gaia/AllWISE 3", eROSITA 20".`,
    );
  }
}

export function validateNneighbor(n: number): void {
  if (!Number.isInteger(n) || n < 1 || n > MAX_NNEIGHBOR) {
    throw new ValidationError(
      `nneighbor must be an integer between 1 and ${MAX_NNEIGHBOR} (got ${n}).`,
    );
  }
}

export function validatePosition(ra: number, dec: number, radiusArcsec?: number): void {
  validateRa(ra);
  validateDec(dec);
  if (radiusArcsec !== undefined) validateRadius(radiusArcsec);
}

/** Link to the XWave web object page for a catalog source. */
export function objectUrl(id: string, catalog: string, webBase = DEFAULT_WEB_URL): string {
  return `${webBase.replace(/\/+$/, "")}/object/${encodeURIComponent(id)}?catalog=${encodeURIComponent(catalog)}`;
}

export interface SesameResult {
  name: string;
  ra: number;
  dec: number;
  mainId?: string;
  objectType?: string;
  source?: string;
}

/**
 * Parse CDS Sesame plain-text output (`-oI` format). Coordinates come from the
 * first `%J ra dec` line (decimal degrees, J2000).
 */
export function parseSesame(text: string, name: string): SesameResult | null {
  let ra: number | undefined;
  let dec: number | undefined;
  let mainId: string | undefined;
  let objectType: string | undefined;
  let source: string | undefined;
  let lastSource: string | undefined;

  for (const raw of text.split(/\r?\n/)) {
    const line = raw.trim();
    if (ra === undefined && line.startsWith("%J ")) {
      const parts = line.split(/\s+/);
      const r = Number.parseFloat(parts[1] ?? "");
      const d = Number.parseFloat(parts[2] ?? "");
      if (Number.isFinite(r) && Number.isFinite(d)) {
        ra = r;
        dec = d;
        source = lastSource;
      }
    } else if (!mainId && line.startsWith("%I.0 ")) {
      mainId = line.slice(5).trim();
    } else if (!objectType && line.startsWith("%C.0 ")) {
      objectType = line.slice(5).trim();
    } else if (line.startsWith("#=")) {
      // Database header preceding its answer, e.g. "#=Sc=Simbad (CDS, via client/server): 1 134ms"
      const m = /^#=[^=]*=([^(:#]+)/.exec(line);
      lastSource = m ? m[1]!.trim() : undefined;
    }
  }

  if (ra === undefined || dec === undefined) return null;
  return { name, ra, dec, mainId, objectType, source };
}

/** Raw source returned by /conesearch and /bulk-conesearch. */
export interface ApiSource {
  id: string;
  ipix?: number;
  ra: number;
  dec: number;
  cat: string;
  distance: number;
}

export interface ApiCatalogResult {
  catalog: string;
  data: ApiSource[] | null;
  index?: number;
}

export interface Match {
  id: string;
  catalog: string;
  ra: number;
  dec: number;
  distance_arcsec: number;
  url: string;
}

/** Flatten a (bulk-)conesearch response into matches sorted by distance. */
export function flattenMatches(
  results: ApiCatalogResult[] | null | undefined,
  webBase = DEFAULT_WEB_URL,
): Match[] {
  const out: Match[] = [];
  for (const group of results ?? []) {
    for (const s of group.data ?? []) {
      const catalog = s.cat || group.catalog;
      out.push({
        id: s.id,
        catalog,
        ra: s.ra,
        dec: s.dec,
        distance_arcsec: round(s.distance, 3),
        url: objectUrl(s.id, catalog, webBase),
      });
    }
  }
  out.sort((a, b) => a.distance_arcsec - b.distance_arcsec);
  return out;
}

/** Group a bulk-conesearch response by input index. */
export function groupBulkByIndex(
  results: ApiCatalogResult[] | null | undefined,
  count: number,
  webBase = DEFAULT_WEB_URL,
): Match[][] {
  const groups: ApiCatalogResult[][] = Array.from({ length: count }, () => []);
  for (const r of results ?? []) {
    const i = r.index ?? 0;
    if (i >= 0 && i < count) groups[i]!.push(r);
  }
  return groups.map((g) => flattenMatches(g, webBase));
}

export function round(x: number, digits: number): number {
  const f = 10 ** digits;
  return Math.round(x * f) / f;
}

export function fmtCoord(ra: number, dec: number): string {
  return `RA ${ra.toFixed(6)}, Dec ${dec >= 0 ? "+" : ""}${dec.toFixed(6)}`;
}

export function catalogLabel(cat: string): string {
  return (CATALOG_INFO as Record<string, CatalogInfo>)[cat]?.name ?? cat;
}

export function formatMatchLine(m: Match): string {
  return `- [${catalogLabel(m.catalog)}] ${m.id} at ${m.distance_arcsec.toFixed(2)}" (${fmtCoord(m.ra, m.dec)}) ${m.url}`;
}

/** Summarise matches per catalog for text output. */
export function summarizeMatches(matches: Match[], maxLines = 20): string {
  if (matches.length === 0) return "No sources found.";
  const counts = new Map<string, number>();
  for (const m of matches) counts.set(m.catalog, (counts.get(m.catalog) ?? 0) + 1);
  const header = `${matches.length} source(s): ` +
    [...counts.entries()].map(([c, n]) => `${catalogLabel(c)} ${n}`).join(", ");
  const lines = matches.slice(0, maxLines).map(formatMatchLine);
  if (matches.length > maxLines) lines.push(`... and ${matches.length - maxLines} more (see structured content).`);
  return [header, ...lines].join("\n");
}

export interface LightcurvePoint {
  catalog: string;
  mjd: number;
  mag: number;
  magerr: number;
}

export interface LightcurveSummaryBand {
  catalog: string;
  n: number;
  mjd_min: number;
  mjd_max: number;
  mag_min: number;
  mag_max: number;
  mag_median: number;
}

export function summarizeLightcurve(points: LightcurvePoint[]): LightcurveSummaryBand[] {
  const byCat = new Map<string, LightcurvePoint[]>();
  for (const p of points) {
    if (!Number.isFinite(p.mjd) || !Number.isFinite(p.mag)) continue;
    const arr = byCat.get(p.catalog) ?? [];
    arr.push(p);
    byCat.set(p.catalog, arr);
  }
  return [...byCat.entries()].map(([catalog, pts]) => {
    const mags = pts.map((p) => p.mag).sort((a, b) => a - b);
    const mjds = pts.map((p) => p.mjd);
    const mid = Math.floor(mags.length / 2);
    const median = mags.length % 2 ? mags[mid]! : (mags[mid - 1]! + mags[mid]!) / 2;
    return {
      catalog,
      n: pts.length,
      mjd_min: round(Math.min(...mjds), 5),
      mjd_max: round(Math.max(...mjds), 5),
      mag_min: round(mags[0]!, 3),
      mag_max: round(mags[mags.length - 1]!, 3),
      mag_median: round(median, 3),
    };
  });
}

/** Drop null/undefined fields from a metadata record. */
export function compactMetadata(meta: Record<string, unknown>): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(meta)) {
    if (v === null || v === undefined) continue;
    out[k] = v;
  }
  return out;
}
