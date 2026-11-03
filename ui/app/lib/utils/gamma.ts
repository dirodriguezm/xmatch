/**
 * Gamma-ray and hard X-ray light curves for associated sources: the Fermi-LAT
 * Light Curve Repository (0.1–100 GeV) and the Swift-BAT Hard X-ray
 * Transient Monitor (15–50 keV).
 *
 * Gamma-ray positions are good to arcminutes, so an object is never matched
 * to a 4FGL source by position. It is matched to the source's associated
 * counterpart (4FGL's ASSOC1, positioned via NED), which is how the LAT
 * team links a gamma-ray source to a blazar or pulsar. The BAT monitor lists
 * its sources at their counterpart's position, so the same check applies.
 * Both lists are built by scripts/build-gamma-sources.cjs.
 */

import gammaSources from "@/app/lib/constants/gamma-sources.json";

import { sepArcsec } from "./xmm5";

/** First LAT science data, 2008-08-04, and first BAT monitor day, 2005-02-12. */
export const FERMI_START_MJD = 54682;
export const BAT_START_MJD = 53414;

/** An object this close to a counterpart is that counterpart. */
export const GAMMA_ASSOC_RADIUS_ARCSEC = 5;

export interface FermiSource {
  /** 4FGL name, e.g. "4FGL J1229.0+0202". */
  name: string;
  /** Associated counterpart, e.g. "3C 273". */
  assoc: string;
  /** 4FGL class, e.g. "fsrq", "bll", "psr". */
  cls: string;
  ra: number;
  dec: number;
}

export interface BatSource {
  name: string;
  /** Path of the source page under the monitor, e.g. "weak/3C273". */
  path: string;
  type: string;
  ra: number;
  dec: number;
}

const FERMI: FermiSource[] = gammaSources.fermi;
const BAT: BatSource[] = gammaSources.bat;

export const GAMMA_SOURCES_DATE = gammaSources.generated;

export interface GammaAssociations {
  fermi?: FermiSource & { separationArcsec: number };
  bat?: BatSource & { separationArcsec: number };
}

function nearest<T extends { ra: number; dec: number }>(
  list: T[],
  ra: number,
  dec: number
): (T & { separationArcsec: number }) | undefined {
  let best: (T & { separationArcsec: number }) | undefined;
  for (const s of list) {
    // Cheap reject before the trig: 5″ is far below a degree.
    if (Math.abs(s.dec - dec) > 0.01) continue;
    const sep = sepArcsec(ra, dec, s.ra, s.dec);
    if (
      sep <= GAMMA_ASSOC_RADIUS_ARCSEC &&
      (!best || sep < best.separationArcsec)
    ) {
      best = { ...s, separationArcsec: sep };
    }
  }
  return best;
}

export function findGammaAssociations(
  ra: number,
  dec: number
): GammaAssociations {
  return { fermi: nearest(FERMI, ra, dec), bat: nearest(BAT, ra, dec) };
}

export function isFermiSource(name: string): boolean {
  return FERMI.some((s) => s.name === name);
}

export function isBatPath(path: string): boolean {
  return BAT.some((s) => s.path === path);
}

// ── Fermi-LAT Light Curve Repository ──────────────────────────────────────

export const LCR_CADENCES = {
  daily: "3 days",
  weekly: "1 week",
  monthly: "1 month",
} as const;

export type LcrCadence = keyof typeof LCR_CADENCES;

export function isLcrCadence(c: string): c is LcrCadence {
  return Object.prototype.hasOwnProperty.call(LCR_CADENCES, c);
}

/** Fermi mission elapsed time (s since 2001-01-01 UTC) to MJD. */
export function metToMjd(met: number): number {
  return 51910 + met / 86_400;
}

export interface GammaPoint {
  mjd: number;
  /** Photon flux 0.1–100 GeV (ph cm⁻² s⁻¹), or a count rate for BAT. */
  value: number;
  lo?: number;
  hi?: number;
  /** True for a 95% upper limit (value is the limit). */
  upperLimit?: boolean;
}

type Pair = [number, number];
type Triple = [number, number, number];

/**
 * Points from the repository's `lightCurveData` JSON: detections with their
 * 1σ bounds (bins with TS ≥ 4) and 95% upper limits for the rest. Failed
 * fits are dropped: bins flagged as not converged, and detections whose
 * error interval collapsed to the value itself (Mrk 421 has one at 1e-3,
 * four decades above the rest).
 */
export function parseLcrLightcurve(json: unknown): GammaPoint[] {
  const d = (json ?? {}) as {
    flux?: Pair[];
    flux_error?: Triple[];
    flux_upper_limits?: Pair[];
    fit_convergence?: Pair[];
  };
  const bounds = new Map(
    (d.flux_error ?? []).map(([t, lo, hi]) => [t, { lo, hi }])
  );
  const failed = new Set(
    (d.fit_convergence ?? []).filter(([, c]) => c !== 0).map(([t]) => t)
  );
  const points: GammaPoint[] = [];
  for (const [t, value] of d.flux ?? []) {
    const b = bounds.get(t);
    if (failed.has(t) || !b || !(b.hi > b.lo)) continue;
    points.push({ mjd: metToMjd(t), value, lo: b.lo, hi: b.hi });
  }
  for (const [t, value] of d.flux_upper_limits ?? []) {
    if (failed.has(t)) continue;
    points.push({ mjd: metToMjd(t), value, upperLimit: true });
  }
  return points.sort((a, b) => a.mjd - b.mjd);
}

// ── Swift-BAT Hard X-ray Transient Monitor ────────────────────────────────

export const BAT_BASE_URL = "https://swift.gsfc.nasa.gov/results/transients/";

/** Daily-average ASCII light curve of a monitor source. */
export function batLightcurveUrl(path: string): string {
  return `${BAT_BASE_URL}${path}.lc.txt`;
}

/**
 * Daily points (count rate cm⁻² s⁻¹, 15–50 keV) from the monitor's ASCII
 * file. Days flagged as bad (DATA_FLAG ≠ 0) are dropped.
 */
export function parseBatLightcurve(text: string): GammaPoint[] {
  const points: GammaPoint[] = [];
  for (const line of text.split("\n")) {
    if (line.startsWith("#")) continue;
    const [mjd, rate, err, , , , , flag] = line.trim().split(/\s+/).map(Number);
    if (
      mjd === undefined ||
      rate === undefined ||
      err === undefined ||
      !Number.isFinite(mjd) ||
      !Number.isFinite(rate) ||
      !Number.isFinite(err) ||
      flag !== 0
    ) {
      continue;
    }
    points.push({ mjd, value: rate, lo: rate - err, hi: rate + err });
  }
  return points;
}

/**
 * Inverse-variance weighted means over consecutive `days`-long bins. A faint
 * BAT source is noise day by day; a month of days shows its level.
 */
export function binPoints(points: GammaPoint[], days: number): GammaPoint[] {
  const bins = new Map<number, { sw: number; swx: number; swt: number }>();
  for (const p of points) {
    if (p.upperLimit || p.lo === undefined || p.hi === undefined) continue;
    const sigma = (p.hi - p.lo) / 2;
    if (!(sigma > 0)) continue;
    const w = 1 / sigma ** 2;
    const key = Math.floor(p.mjd / days);
    const b = bins.get(key) ?? { sw: 0, swx: 0, swt: 0 };
    b.sw += w;
    b.swx += w * p.value;
    b.swt += w * p.mjd;
    bins.set(key, b);
  }
  return [...bins.entries()]
    .sort(([a], [b]) => a - b)
    .map(([, b]) => {
      const value = b.swx / b.sw;
      const err = 1 / Math.sqrt(b.sw);
      return { mjd: b.swt / b.sw, value, lo: value - err, hi: value + err };
    });
}
