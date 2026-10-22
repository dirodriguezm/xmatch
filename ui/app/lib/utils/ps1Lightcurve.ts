/**
 * Pan-STARRS1 DR2 detections (per-epoch grizy photometry) from the MAST
 * catalogs API.
 */

import type { DetectionPoint } from "@/app/lib/utils/lightcurve";

/** The 3π survey covers the sky north of δ = −30°. */
export const PS1_MIN_DEC = -30;

const FILTERS = ["g", "r", "i", "z", "y"];
/** AB zero point for fluxes in Jy: m = −2.5 log10(F) + 8.90. */
const AB_ZP_JY = 8.9;
const MAG_PER_SNR = 2.5 / Math.LN10;

function angularSepDeg(ra1: number, dec1: number, ra2: number, dec2: number) {
  const r = Math.PI / 180;
  const cos =
    Math.sin(dec1 * r) * Math.sin(dec2 * r) +
    Math.cos(dec1 * r) * Math.cos(dec2 * r) * Math.cos((ra1 - ra2) * r);
  return Math.acos(Math.min(1, Math.max(-1, cos))) / r;
}

/**
 * Turn a MAST `detection.csv` cone response into light-curve points (AB
 * magnitudes). CSV rather than JSON because MAST sometimes serialises the
 * 64-bit `objID` as a JSON number, which JSON.parse rounds; CSV keeps the
 * digits as text.
 *
 * When the cone holds detections of several PS1 objects (PS1 has duplicates
 * near stack boundaries), only the object whose mean position is nearest the
 * target is kept. Non-positive fluxes are dropped.
 */
export function parsePs1DetectionsCsv(
  csv: string,
  ra: number,
  dec: number
): { objId: string | null; points: DetectionPoint[] } {
  const lines = csv.split(/\r?\n/).filter((l) => l.trim() !== "");
  const names = (lines[0] ?? "").split(",");
  const col = (name: string) => names.indexOf(name);
  const [iObj, iTime, iFilter, iFlux, iErr, iRa, iDec] = [
    "objID",
    "obsTime",
    "filterID",
    "psfFlux",
    "psfFluxErr",
    "ra",
    "dec",
  ].map(col);
  if ([iObj, iTime, iFilter, iFlux].some((i) => i === -1)) {
    return { objId: null, points: [] };
  }
  const rows = lines.slice(1).map((l) => l.split(","));

  // Mean position per object → nearest object wins.
  const positions = new Map<string, { ra: number; dec: number; n: number }>();
  for (const r of rows) {
    const id = r[iObj];
    const p = positions.get(id) ?? { ra: 0, dec: 0, n: 0 };
    if (iRa !== -1 && iDec !== -1) {
      p.ra += Number(r[iRa]);
      p.dec += Number(r[iDec]);
    }
    p.n += 1;
    positions.set(id, p);
  }
  let objId: string | null = null;
  let best = Infinity;
  for (const [id, p] of positions) {
    const sep =
      iRa !== -1 ? angularSepDeg(ra, dec, p.ra / p.n, p.dec / p.n) : 0;
    if (sep < best) {
      best = sep;
      objId = id;
    }
  }

  const points: DetectionPoint[] = [];
  for (const r of rows) {
    if (r[iObj] !== objId) continue;
    const flux = Number(r[iFlux]);
    const mjd = Number(r[iTime]);
    const band = FILTERS[Number(r[iFilter]) - 1];
    if (!(flux > 0) || !Number.isFinite(mjd) || !band) continue;
    const err = iErr !== -1 ? Number(r[iErr]) : NaN;
    points.push({
      mjd,
      mag: -2.5 * Math.log10(flux) + AB_ZP_JY,
      magerr: err > 0 ? (MAG_PER_SNR * err) / flux : undefined,
      band,
    });
  }
  return { objId, points: points.sort((a, b) => a.mjd! - b.mjd!) };
}
