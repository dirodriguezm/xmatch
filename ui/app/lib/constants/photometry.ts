/**
 * Representative magnitude per search catalog.
 *
 * The service publishes a different metadata schema for every catalog — AllWISE
 * returns `w1mpro…w4mpro` plus 2MASS associations, Gaia returns
 * `phot_*_mean_mag` — so a fixed `W1 | G | RP | …` column set would be mostly
 * empty for any given row. Instead each catalog names the magnitudes it can
 * offer, in preference order, and the results table shows the first one present.
 */

/** One candidate magnitude for a catalog. */
export interface MagSpec {
  /** Display label, e.g. "W1", "G". */
  band: string;
  /** Metadata field holding the magnitude. */
  field: string;
  /** 1-sigma uncertainty, surfaced in the tooltip when present. */
  errField?: string;
  description: string;
}

/** Same "no measurement" convention the light-curve payloads use. */
const SENTINEL = -999;

export const CATALOG_MAG_SPECS: Record<string, MagSpec[]> = {
  allwise: [
    {
      band: "W1",
      field: "w1mpro",
      errField: "w1sigmpro",
      description: "WISE W1 (3.4 µm) profile-fit magnitude",
    },
    {
      band: "W2",
      field: "w2mpro",
      errField: "w2sigmpro",
      description: "WISE W2 (4.6 µm) profile-fit magnitude",
    },
    {
      band: "J",
      field: "j_m_2mass",
      errField: "j_msig_2mass",
      description: "2MASS J (1.24 µm) via the AllWISE 2MASS association",
    },
  ],
  gaia: [
    {
      band: "G",
      field: "phot_g_mean_mag",
      description: "Gaia DR3 G-band mean magnitude",
    },
    {
      band: "BP",
      field: "phot_bp_mean_mag",
      description: "Gaia DR3 BP mean magnitude",
    },
    {
      band: "RP",
      field: "phot_rp_mean_mag",
      description: "Gaia DR3 RP mean magnitude",
    },
  ],
  // Intentionally empty: eROSITA is an X-ray catalog and publishes fluxes, not
  // magnitudes. Its metadata schema is also unverified — every position probed
  // returned 204. An empty spec makes the Mag cell render "—" and, more
  // usefully, makes buildBulkMetadataRequests skip the round trip entirely.
  erosita: [],
};

export interface MagValue {
  band: string;
  mag: number;
  magErr?: number;
  description: string;
}

/**
 * A magnitude of exactly 0 is legal, so this tests for finiteness rather than
 * truthiness — `!0` would silently drop a real measurement.
 */
function isMeasured(value: unknown): value is number {
  return (
    typeof value === "number" && Number.isFinite(value) && value !== SENTINEL
  );
}

/**
 * First available magnitude for `catalogSlug`, or null when the catalog has no
 * spec or every candidate field is missing.
 */
export function pickMag(
  catalogSlug: string,
  meta: Record<string, unknown> | undefined
): MagValue | null {
  const specs = CATALOG_MAG_SPECS[catalogSlug.toLowerCase()];
  if (!meta || !specs?.length) return null;

  for (const spec of specs) {
    const value = meta[spec.field];
    if (!isMeasured(value)) continue;
    const err = spec.errField ? meta[spec.errField] : undefined;
    return {
      band: spec.band,
      mag: value,
      magErr: isMeasured(err) ? err : undefined,
      description: spec.description,
    };
  }
  return null;
}
