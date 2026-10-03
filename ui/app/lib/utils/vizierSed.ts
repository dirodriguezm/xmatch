/**
 * Photometry from the VizieR SED service (https://vizier.cds.unistra.fr/vizier/sed/).
 *
 * The service answers with a VOTable (TABLEDATA) holding one row per published
 * measurement near the position — often hundreds, many of them the same band
 * re-published by several catalogs, and a few from a neighbouring source. The
 * raw rows are parsed here and collapsed to one robust point per filter.
 */

/** One measurement as VizieR publishes it. */
export interface VizierSedRow {
  /** VizieR table the value comes from, e.g. `II/328/allwise`. */
  table: string;
  frequencyGHz: number;
  fluxJy: number;
  fluxErrJy?: number;
  /** Filter identifier, e.g. `SDSS:g`, `WISE:W1`, `GAIA/GAIA3:G`. */
  filter: string;
}

/** One filter after merging every row that measured it. */
export interface VizierSedPoint {
  filter: string;
  wavelengthUm: number;
  /** Median flux density of the rows kept, in Jy. */
  fluxJy: number;
  /** νFν in erg s⁻¹ cm⁻². */
  nuFnu: number;
  nuFnuErr?: number;
  /** Rows kept after rejecting outliers. */
  nMeasurements: number;
  /** Rows rejected as outliers (typically a neighbour or a bad match). */
  nRejected: number;
  tables: string[];
  /**
   * More than {@link INCONSISTENT_DEX} away from the median of neighbouring
   * filters — usually a mis-scaled table or a different source. Kept in the
   * data so the UI can still show it on request.
   */
  inconsistent: boolean;
}

const REQUIRED_FIELDS = [
  "_tabname",
  "_sed_freq",
  "_sed_flux",
  "_sed_eflux",
  "_sed_filter",
] as const;

/** Speed of light in µm·GHz, so λ[µm] = C / ν[GHz]. */
const C_UM_GHZ = 2.99792458e5;
const JY_TO_CGS = 1e-23;
/** Scale MAD to a Gaussian σ. */
const MAD_TO_SIGMA = 1.4826;
const OUTLIER_SIGMA = 3;
/** Neighbouring filters: within this many dex in wavelength. */
const NEIGHBOUR_DEX = 0.2;
const MIN_NEIGHBOURS = 3;
/**
 * A factor of ~30. Real broadband SEDs change far less than this across
 * 0.2 dex in wavelength (a star's UV drop is under 1 dex), but a table off by
 * a unit conversion or matched to another source is typically 2–3 dex off.
 */
export const INCONSISTENT_DEX = 1.5;

function decodeEntities(text: string): string {
  return text
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/&quot;/g, '"')
    .replace(/&apos;/g, "'")
    .replace(/&amp;/g, "&");
}

function toNumber(text: string | undefined): number | undefined {
  if (text == null || text.trim() === "") return undefined;
  const n = Number(text);
  return Number.isFinite(n) ? n : undefined;
}

/**
 * Parse the VOTable the SED service returns. Columns are located by FIELD
 * name, never by position. Throws when the payload is not a VOTable or lacks
 * the SED columns; returns [] for a VOTable with no rows.
 */
export function parseVizierSedVotable(xml: string): VizierSedRow[] {
  if (!xml.includes("<VOTABLE")) {
    throw new Error("VizieR SED: response is not a VOTable");
  }

  const tableStart = xml.indexOf("<TABLE");
  if (tableStart === -1) return [];
  const table = xml.slice(tableStart);

  const fieldNames = [...table.matchAll(/<FIELD\b[^>]*\bname="([^"]+)"/g)].map(
    (m) => m[1]
  );
  const col = Object.fromEntries(
    REQUIRED_FIELDS.map((name) => [name, fieldNames.indexOf(name)])
  ) as Record<(typeof REQUIRED_FIELDS)[number], number>;
  const missing = REQUIRED_FIELDS.filter((name) => col[name] === -1);
  if (missing.length > 0) {
    throw new Error(`VizieR SED: missing columns ${missing.join(", ")}`);
  }

  const rows: VizierSedRow[] = [];
  for (const tr of table.matchAll(/<TR>([\s\S]*?)<\/TR>/g)) {
    const cells = [...tr[1].matchAll(/<TD>([\s\S]*?)<\/TD>|<TD\s*\/>/g)].map(
      (m) => decodeEntities(m[1] ?? "")
    );
    const frequencyGHz = toNumber(cells[col._sed_freq]);
    const fluxJy = toNumber(cells[col._sed_flux]);
    const filter = cells[col._sed_filter]?.trim();
    if (!frequencyGHz || frequencyGHz <= 0 || fluxJy == null || !filter) {
      continue;
    }
    const fluxErrJy = toNumber(cells[col._sed_eflux]);
    rows.push({
      table: cells[col._tabname]?.trim() ?? "",
      frequencyGHz,
      fluxJy,
      fluxErrJy: fluxErrJy != null && fluxErrJy > 0 ? fluxErrJy : undefined,
      filter,
    });
  }
  return rows;
}

function median(values: number[]): number {
  const sorted = [...values].sort((a, b) => a - b);
  const mid = Math.floor(sorted.length / 2);
  return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
}

/**
 * Collapse rows to one point per filter: median flux after rejecting values
 * more than 3σ (from the MAD) away from the median. Non-positive fluxes are
 * dropped because they cannot sit on a log axis. Sorted by wavelength.
 */
export function aggregateVizierSed(rows: VizierSedRow[]): VizierSedPoint[] {
  const byFilter = new Map<string, VizierSedRow[]>();
  for (const row of rows) {
    if (row.fluxJy <= 0) continue;
    const list = byFilter.get(row.filter) ?? [];
    list.push(row);
    byFilter.set(row.filter, list);
  }

  const points: VizierSedPoint[] = [];
  for (const [filter, group] of byFilter) {
    const fluxes = group.map((r) => r.fluxJy);
    const center = median(fluxes);
    const sigma =
      MAD_TO_SIGMA * median(fluxes.map((f) => Math.abs(f - center)));
    // With MAD = 0 (most rows identical) there is no scale to reject against.
    const kept =
      sigma > 0
        ? group.filter(
            (r) => Math.abs(r.fluxJy - center) <= OUTLIER_SIGMA * sigma
          )
        : group;
    const used = kept.length > 0 ? kept : group;

    const fluxJy = median(used.map((r) => r.fluxJy));
    const frequencyGHz = median(used.map((r) => r.frequencyGHz));
    const errors = used
      .map((r) => r.fluxErrJy)
      .filter((e): e is number => e != null);
    // Reported errors when there are any; otherwise the scatter between rows.
    const fluxErrJy =
      errors.length > 0
        ? median(errors)
        : used.length > 2 && sigma > 0
          ? sigma
          : undefined;

    const toNuFnu = (f: number) => f * JY_TO_CGS * frequencyGHz * 1e9;
    points.push({
      filter,
      wavelengthUm: C_UM_GHZ / frequencyGHz,
      fluxJy,
      nuFnu: toNuFnu(fluxJy),
      nuFnuErr: fluxErrJy != null ? toNuFnu(fluxErrJy) : undefined,
      nMeasurements: used.length,
      nRejected: group.length - used.length,
      tables: [...new Set(used.map((r) => r.table))].sort(),
      inconsistent: false,
    });
  }

  return flagInconsistent(
    points.sort((a, b) => a.wavelengthUm - b.wavelengthUm)
  );
}

/**
 * Mark filters whose νFν disagrees with their spectral neighbours by more
 * than {@link INCONSISTENT_DEX}. Filters with too few neighbours are never
 * flagged: there is nothing to compare them against.
 */
function flagInconsistent(points: VizierSedPoint[]): VizierSedPoint[] {
  const logLambda = points.map((p) => Math.log10(p.wavelengthUm));
  const logFlux = points.map((p) => Math.log10(p.nuFnu));
  return points.map((p, i) => {
    const neighbours = logFlux.filter(
      (_, j) => j !== i && Math.abs(logLambda[j] - logLambda[i]) < NEIGHBOUR_DEX
    );
    if (neighbours.length < MIN_NEIGHBOURS) return p;
    const deviation = Math.abs(logFlux[i] - median(neighbours));
    return deviation > INCONSISTENT_DEX ? { ...p, inconsistent: true } : p;
  });
}
