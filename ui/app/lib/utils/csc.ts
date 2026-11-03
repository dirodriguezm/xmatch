/**
 * Chandra Source Catalog 2.2 (master sources), from the CXC cone search.
 *
 * The cone search answers in under a second with a fixed set of summary
 * columns. The CSC TAP service doesn't implement ADQL's CONTAINS/CIRCLE,
 * so it can't run a positional query.
 */

import { sepArcsec } from "./xmm5";

export const CSC_CONE_URL = "https://cda.cfa.harvard.edu/csc22scs/coneSearch";

/** Sources within this of the target are kept; CSC positions are sub-arcsecond. */
export const CSC_MATCH_RADIUS_ARCSEC = 3;

/** Integer columns use this for "no value". */
const INT_NULL = -2147483648;

export interface CscRange {
  value: number;
  lo?: number;
  hi?: number;
}

export interface ChandraSource {
  name: string;
  ra: number;
  dec: number;
  separationArcsec: number;
  /** 95% error ellipse semi-major axis (arcsec). */
  posErrArcsec?: number;
  /** Broad band 0.5–7 keV aperture flux (erg s⁻¹ cm⁻²). */
  flux?: CscRange;
  /** Hard–medium and medium–soft hardness ratios. */
  hardHm?: CscRange;
  hardMs?: CscRange;
  significance?: number;
  /** Gregory–Loredo variability index within observations, 0–10. */
  varIntra?: number;
  /** Variability index between observations, 0–8. */
  varInter?: number;
  /** Extended beyond the PSF. */
  extended: boolean;
  /** Piled up in every observation, so the flux is unreliable. */
  saturated: boolean;
  /** Confused with a nearby source in at least one observation. */
  confused: boolean;
}

export function cscConeUrl(ra: number, dec: number): string {
  const params = new URLSearchParams({
    RA: String(ra),
    DEC: String(dec),
    SR: String((2 * CSC_MATCH_RADIUS_ARCSEC) / 3600),
  });
  return `${CSC_CONE_URL}?${params}`;
}

function unescapeXml(s: string): string {
  return s
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&apos;/g, "'")
    .replace(/&amp;/g, "&");
}

/** Rows of a single-table VOTable as column-name → cell-text records. */
export function parseVoTable(xml: string): Record<string, string>[] {
  const names = [...xml.matchAll(/<FIELD\b[^>]*\bname="([^"]+)"/g)].map(
    (m) => m[1] ?? ""
  );
  const rows: Record<string, string>[] = [];
  for (const tr of xml.matchAll(/<TR>([\s\S]*?)<\/TR>/g)) {
    const cells = [
      ...(tr[1] ?? "").matchAll(/<TD(?:\s*\/>|>([\s\S]*?)<\/TD>)/g),
    ];
    const row: Record<string, string> = {};
    names.forEach((name, i) => {
      row[name] = unescapeXml(cells[i]?.[1]?.trim() ?? "");
    });
    rows.push(row);
  }
  return rows;
}

function num(v: string | undefined): number | undefined {
  if (v == null || v === "") return undefined;
  const n = Number(v);
  return Number.isFinite(n) && n !== INT_NULL ? n : undefined;
}

function range(
  row: Record<string, string>,
  col: string,
  lo: string,
  hi: string
): CscRange | undefined {
  const value = num(row[col]);
  return value === undefined
    ? undefined
    : { value, lo: num(row[lo]), hi: num(row[hi]) };
}

const flag = (v: string | undefined) => v === "T" || v === "true" || v === "1";

/**
 * Nearest CSC master source to (ra, dec) within
 * {@link CSC_MATCH_RADIUS_ARCSEC}, or null.
 */
export function parseCscCone(
  xml: string,
  ra: number,
  dec: number
): ChandraSource | null {
  let best: { row: Record<string, string>; sep: number } | null = null;
  for (const row of parseVoTable(xml)) {
    const sRa = num(row.ra);
    const sDec = num(row.dec);
    if (sRa === undefined || sDec === undefined) continue;
    const sep = sepArcsec(ra, dec, sRa, sDec);
    if (sep <= CSC_MATCH_RADIUS_ARCSEC && (!best || sep < best.sep)) {
      best = { row, sep };
    }
  }
  if (!best) return null;
  const { row, sep } = best;

  return {
    name: row.name ?? "",
    ra: num(row.ra) ?? ra,
    dec: num(row.dec) ?? dec,
    separationArcsec: Math.round(sep * 100) / 100,
    posErrArcsec: num(row.err_ellipse_r0),
    flux: range(row, "flux_aper_b", "flux_aper_lolim_b", "flux_aper_hilim_b"),
    hardHm: range(row, "hard_hm", "hard_hm_lolim", "hard_hm_hilim"),
    hardMs: range(row, "hard_ms", "hard_ms_lolim", "hard_ms_hilim"),
    significance: num(row.significance),
    varIntra: num(row.var_intra_index_b),
    varInter: num(row.var_inter_index_b),
    extended: flag(row.extent_flag),
    saturated: flag(row.sat_src_flag),
    confused: flag(row.conf_flag),
  };
}
