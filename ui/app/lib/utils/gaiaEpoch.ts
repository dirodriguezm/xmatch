/**
 * Gaia DR3 epoch photometry (G, BP, RP per transit) from the ESA Gaia
 * archive DataLink service. Only sources flagged `has_epoch_photometry`
 * (~11.7 million, mostly variables) have it.
 */

import type { DetectionPoint } from "@/app/lib/utils/lightcurve";

/**
 * Gaia times are TCB days since BJD 2455197.5, i.e. MJD 55197.0. The
 * barycentric/TCB vs UTC difference (≤ ~8 min) is negligible on a light-curve
 * plot spanning years.
 */
const GAIA_TIME_TO_MJD = 55197.0;

/** 2.5 / ln 10: converts flux S/N to a magnitude error. */
const MAG_PER_SNR = 2.5 / Math.LN10;

/**
 * The numeric `source_id` XWave returns loses precision (a 64-bit integer
 * through a double), but the designation string keeps every digit.
 */
export function gaiaSourceIdFromDesignation(
  designation: string
): string | null {
  const m = designation.trim().match(/^Gaia DR3 (\d{1,20})$/);
  return m ? m[1] : null;
}

function splitCsv(text: string): { header: string[]; rows: string[][] } {
  const lines = text.split(/\r?\n/).filter((l) => l.trim() !== "");
  if (lines.length === 0) return { header: [], rows: [] };
  return {
    header: lines[0].split(","),
    rows: lines.slice(1).map((l) => l.split(",")),
  };
}

function num(v: string | undefined): number | undefined {
  if (v == null || v.trim() === "") return undefined;
  const n = Number(v);
  return Number.isFinite(n) ? n : undefined;
}

const BANDS = [
  {
    band: "G",
    time: "g_transit_time",
    mag: "g_transit_mag",
    snr: "g_transit_flux_over_error",
    reject: "variability_flag_g_reject",
  },
  {
    band: "BP",
    time: "bp_obs_time",
    mag: "bp_mag",
    snr: "bp_flux_over_error",
    reject: "variability_flag_bp_reject",
  },
  {
    band: "RP",
    time: "rp_obs_time",
    mag: "rp_mag",
    snr: "rp_flux_over_error",
    reject: "variability_flag_rp_reject",
  },
] as const;

/**
 * Parse the DataLink CSV (DATA_STRUCTURE=INDIVIDUAL): one row per transit.
 * Measurements the Gaia pipeline rejected (photometry or variability flags)
 * are dropped. Magnitudes are Vega (Gaia's photometric system).
 */
export function parseGaiaEpochCsv(csv: string): DetectionPoint[] {
  const { header, rows } = splitCsv(csv);
  const col = (name: string) => header.indexOf(name);
  if (col("g_transit_time") === -1) return [];

  const rejectedAll = col("rejected_by_photometry");
  const points: DetectionPoint[] = [];
  for (const row of rows) {
    if (rejectedAll !== -1 && row[rejectedAll] === "true") continue;
    for (const b of BANDS) {
      if (row[col(b.reject)] === "true") continue;
      const t = num(row[col(b.time)]);
      const mag = num(row[col(b.mag)]);
      if (t == null || mag == null) continue;
      const snr = num(row[col(b.snr)]);
      points.push({
        mjd: t + GAIA_TIME_TO_MJD,
        mag,
        magerr: snr != null && snr > 0 ? MAG_PER_SNR / snr : undefined,
        band: b.band,
      });
    }
  }
  return points.sort((a, b) => a.mjd! - b.mjd!);
}

/**
 * Parse a Gaia TAP cone search (CSV with `source_id,has_epoch_photometry`):
 * the nearest source, or null when the cone was empty.
 */
export function parseGaiaTapCone(
  csv: string
): { sourceId: string; hasEpochPhotometry: boolean } | null {
  const { header, rows } = splitCsv(csv);
  const iId = header.indexOf("source_id");
  const iFlag = header.indexOf("has_epoch_photometry");
  if (iId === -1 || rows.length === 0 || !/^\d+$/.test(rows[0][iId] ?? "")) {
    return null;
  }
  return {
    sourceId: rows[0][iId],
    hasEpochPhotometry: rows[0][iFlag] === "true" || rows[0][iFlag] === "1",
  };
}
