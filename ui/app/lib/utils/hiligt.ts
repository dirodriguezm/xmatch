/**
 * HILIGT, ESA's High-Energy Light-curve Generator (the XMM upper limit
 * server): per-observation X-ray fluxes, or upper limits computed from the
 * images where no source was detected.
 *
 * Fluxes use HILIGT's default conversion: an absorbed power law with
 * NH = 3×10²⁰ cm⁻² and photon index 2. Upper limits are at `sigma` σ.
 */

export const HILIGT_MISSIONS = {
  XMMpnt: "XMM-Newton pointed",
  XMMslew: "XMM-Newton slew",
  RosatSurvey: "ROSAT All-Sky Survey",
} as const;

export type HiligtMission = keyof typeof HILIGT_MISSIONS;

export function isHiligtMission(m: string): m is HiligtMission {
  return Object.prototype.hasOwnProperty.call(HILIGT_MISSIONS, m);
}

export const HILIGT_MODEL_NOTE = "Absorbed power law, NH = 3×10²⁰ cm⁻², Γ = 2";

export interface HiligtPoint {
  mission: HiligtMission;
  instrument?: string;
  obsid?: string;
  /** Midpoint of the observation. */
  mjd: number;
  /** Band in keV, e.g. "0.2-2". */
  band: string;
  /** Detected flux (erg s⁻¹ cm⁻²). */
  flux?: number;
  fluxErr?: number;
  /** Upper limit on the flux (erg s⁻¹ cm⁻²) when not detected. */
  upperLimit?: number;
  sigma?: number;
  exposureS?: number;
}

const MJD_UNIX_EPOCH = 40587;

/** "2012-07-16T22:52:16" or "1990-12-20 22:04:27" (UTC) → MJD. */
export function isoToMjd(value: unknown): number | undefined {
  if (typeof value !== "string" || value.trim() === "") return undefined;
  const ms = Date.parse(`${value.trim().replace(" ", "T")}Z`);
  return Number.isFinite(ms) ? ms / 86_400_000 + MJD_UNIX_EPOCH : undefined;
}

/** HILIGT mixes numbers, numeric strings, " " and null. */
function num(v: unknown): number | undefined {
  if (typeof v === "number") return Number.isFinite(v) ? v : undefined;
  if (typeof v !== "string" || v.trim() === "") return undefined;
  const n = Number(v);
  return Number.isFinite(n) ? n : undefined;
}

function bandLabel(elow: unknown, ehigh: unknown): string | undefined {
  const lo = num(elow);
  const hi = num(ehigh);
  return lo === undefined || hi === undefined ? undefined : `${lo}-${hi}`;
}

/**
 * Parse the JSON array returned by `ULSservice_passthru`. An empty body (what
 * HILIGT sends for no coverage) and records that failed (`_status` ≠ "Ok")
 * yield no points.
 */
export function parseHiligt(
  text: string,
  mission: HiligtMission
): HiligtPoint[] {
  if (text.trim() === "") return [];
  let records: unknown;
  try {
    records = JSON.parse(text);
  } catch {
    return [];
  }
  if (!Array.isArray(records)) return [];

  const points: HiligtPoint[] = [];
  for (const r of records as Record<string, unknown>[]) {
    if (!r || (typeof r._status === "string" && r._status !== "Ok")) continue;
    const start = isoToMjd(r._start_date);
    const end = isoToMjd(r._end_date);
    const mjd =
      start !== undefined && end !== undefined ? (start + end) / 2 : start;
    const band = bandLabel(r._elow, r._ehigh);
    if (mjd === undefined || band === undefined) continue;

    const flux = num(r._crate_flux);
    const upperLimit = num(r._ul_flux);
    if (flux === undefined && upperLimit === undefined) continue;

    points.push({
      mission,
      instrument: typeof r._instrum === "string" ? r._instrum : undefined,
      obsid: r._obsid == null ? undefined : String(r._obsid),
      mjd,
      band,
      ...(flux !== undefined
        ? { flux, fluxErr: num(r._crate_flux_err) }
        : { upperLimit, sigma: num(r._ulsig) }),
      exposureS: num(r._exptime),
    });
  }
  return points.sort((a, b) => a.mjd - b.mjd);
}
