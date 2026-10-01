/**
 * 5XMM-DR15, the XMM-Newton serendipitous source catalogue (stacked over all
 * observations), from the XMM-Newton Science Archive TAP service.
 *
 * Only the columns shown on the object page are requested. The catalogue's
 * source-class columns (classx_*, classopt_*) are deliberately left out.
 */

/** EPIC bands 1–5 of the catalogue, in keV. */
export const XMM_BANDS = [
  { id: 1, label: "0.2–0.5 keV" },
  { id: 2, label: "0.5–1 keV" },
  { id: 3, label: "1–2 keV" },
  { id: 4, label: "2–4.5 keV" },
  { id: 5, label: "4.5–12 keV" },
] as const;

/** Sources whose stacked position lies within this of the target are kept. */
export const XMM_MATCH_RADIUS_ARCSEC = 5;

export const XMM_COLUMNS = [
  "srcid",
  "iauname",
  "ra",
  "dec",
  "ep_flux",
  "ep_flux_err",
  ...XMM_BANDS.flatMap((b) => [`ep_${b.id}_flux`, `ep_${b.id}_flux_err`]),
  ...[1, 2, 3, 4].flatMap((i) => [`ep_hr${i}`, `ep_hr${i}_err`]),
  "n_obs",
  "mjd_first",
  "mjd_last",
  "var_flag",
  "var_prob",
  "fvar",
  "fvar_err",
  "stack_gamma",
  "stack_nh",
  "extent",
  "sum_flag",
  "gaiadr3_source_id",
  "gaia_match_prob",
  "wise_name",
  "wise_match_prob",
  "info_counterparts",
] as const;

export interface XmmValue {
  value: number;
  err?: number;
}

export interface XmmSource {
  srcid: string;
  name: string;
  ra: number;
  dec: number;
  separationArcsec: number;
  /** 0.2–12 keV flux (erg s⁻¹ cm⁻²). */
  flux?: XmmValue;
  bandFluxes: { label: string; flux?: XmmValue }[];
  hardnessRatios: { label: string; hr?: XmmValue }[];
  observations?: number;
  mjdFirst?: number;
  mjdLast?: number;
  variable?: boolean;
  /** Probability that the source is constant; small means variable. */
  constantProbability?: number;
  fvar?: XmmValue;
  /** Photon index of the simple power-law fit. */
  gamma?: number;
  /** Absorbing column of that fit (cm⁻²). */
  nh?: number;
  /** Extent (arcsec); 0 for a point source. */
  extentArcsec?: number;
  /** Catalogue quality summary flag, 0 (best) to 4. */
  quality?: number;
  gaia?: { sourceId: string; probability?: number };
  wise?: { name: string; probability?: number };
  /** Source page at the XMM-Newton SSC (IRAP) with all counterparts. */
  infoUrl?: string;
}

export function splitCsv(text: string): { header: string[]; rows: string[][] } {
  const lines = text.split(/\r?\n/).filter((l) => l.trim() !== "");
  if (lines.length === 0) return { header: [], rows: [] };
  return {
    header: lines[0].split(",").map((s) => s.trim()),
    rows: lines.slice(1).map((l) => l.split(",")),
  };
}

function num(v: string | undefined): number | undefined {
  if (v == null || v.trim() === "") return undefined;
  const n = Number(v);
  return Number.isFinite(n) ? n : undefined;
}

function sepArcsec(ra1: number, dec1: number, ra2: number, dec2: number) {
  const r = Math.PI / 180;
  const h =
    Math.sin(((dec2 - dec1) * r) / 2) ** 2 +
    Math.cos(dec1 * r) *
      Math.cos(dec2 * r) *
      Math.sin(((ra2 - ra1) * r) / 2) ** 2;
  return ((2 * Math.asin(Math.sqrt(h))) / r) * 3600;
}

/** ADQL cone over the stacked catalogue, returning only {@link XMM_COLUMNS}. */
export function xmmConeQuery(ra: number, dec: number): string {
  const radiusDeg = (2 * XMM_MATCH_RADIUS_ARCSEC) / 3600;
  return (
    `SELECT ${XMM_COLUMNS.join(", ")} FROM xsa.v_epic_source_cat ` +
    `WHERE 1=CONTAINS(POINT('ICRS', ra, dec), ` +
    `CIRCLE('ICRS', ${ra}, ${dec}, ${radiusDeg})) ` +
    // The table also holds one row per observation of each source; the
    // stacked summary row is the one with an observation count.
    `AND n_obs IS NOT NULL`
  );
}

/**
 * Nearest 5XMM source to (ra, dec) within {@link XMM_MATCH_RADIUS_ARCSEC},
 * or null. Expects the CSV produced by {@link xmmConeQuery}.
 */
export function parseXmmSourceCsv(
  csv: string,
  ra: number,
  dec: number
): XmmSource | null {
  const { header, rows } = splitCsv(csv);
  const col = (name: string) => header.indexOf(name);
  if (col("ra") === -1 || col("dec") === -1 || col("iauname") === -1) {
    return null;
  }

  let best: { row: string[]; sep: number } | null = null;
  const iObs = col("n_obs");
  for (const row of rows) {
    // Skip per-observation rows if the query didn't already exclude them.
    if (iObs !== -1 && num(row[iObs]) === undefined) continue;
    const sRa = num(row[col("ra")]);
    const sDec = num(row[col("dec")]);
    if (sRa === undefined || sDec === undefined) continue;
    const sep = sepArcsec(ra, dec, sRa, sDec);
    if (sep <= XMM_MATCH_RADIUS_ARCSEC && (!best || sep < best.sep)) {
      best = { row, sep };
    }
  }
  if (!best) return null;

  const row = best.row;
  const get = (name: string) => {
    const i = col(name);
    return i === -1 ? undefined : row[i]?.trim();
  };
  const n = (name: string) => num(get(name));
  const pair = (name: string): XmmValue | undefined => {
    const value = n(name);
    return value === undefined ? undefined : { value, err: n(`${name}_err`) };
  };
  const str = (name: string) => get(name) || undefined;

  const gaiaId = str("gaiadr3_source_id");
  const wiseName = str("wise_name");
  const varFlag = str("var_flag")?.toLowerCase();

  return {
    srcid: str("srcid") ?? "",
    name: str("iauname") ?? "",
    ra: n("ra")!,
    dec: n("dec")!,
    separationArcsec: Math.round(best.sep * 100) / 100,
    flux: pair("ep_flux"),
    bandFluxes: XMM_BANDS.map((b) => ({
      label: b.label,
      flux: pair(`ep_${b.id}_flux`),
    })),
    hardnessRatios: [1, 2, 3, 4].map((i) => ({
      label: `HR${i}`,
      hr: pair(`ep_hr${i}`),
    })),
    observations: n("n_obs"),
    mjdFirst: n("mjd_first"),
    mjdLast: n("mjd_last"),
    variable:
      varFlag === "true" ? true : varFlag === "false" ? false : undefined,
    constantProbability: n("var_prob"),
    fvar: pair("fvar"),
    gamma: n("stack_gamma"),
    nh: n("stack_nh"),
    extentArcsec: n("extent"),
    quality: n("sum_flag"),
    gaia: gaiaId
      ? { sourceId: gaiaId, probability: n("gaia_match_prob") }
      : undefined,
    wise: wiseName
      ? { name: wiseName, probability: n("wise_match_prob") }
      : undefined,
    infoUrl: str("info_counterparts"),
  };
}
