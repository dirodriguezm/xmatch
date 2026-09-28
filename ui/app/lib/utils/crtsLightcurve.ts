/**
 * Catalina Sky Survey photometry (CRTS data release) from Caltech's cone
 * search service. A query returns an HTML page that either says the area is
 * not covered, says it is covered but empty, or links a temporary CSV of
 * per-epoch photometry: MasterID, Mag, Magerr, RA, Dec, MJD, Blend.
 */

import type { DetectionPoint } from "@/app/lib/utils/lightcurve";

export type CrtsQueryStatus = "not_covered" | "empty" | "rows" | "unknown";

/** Classify the service's HTML reply and pull out the CSV link, if any. */
export function parseCrtsQueryPage(html: string): {
  status: CrtsQueryStatus;
  csvUrl: string | null;
} {
  const text = html.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ");
  const csvUrl = html.match(/https?:\/\/[^"'\s<>]+\.csv/)?.[0] ?? null;
  if (/not covered/i.test(text)) return { status: "not_covered", csvUrl: null };
  if (/no objects were found/i.test(text)) {
    return { status: "empty", csvUrl: null };
  }
  if (csvUrl && /There were \d+ lines/i.test(text)) {
    return { status: "rows", csvUrl };
  }
  return { status: "unknown", csvUrl };
}

/**
 * The first digit of a MasterID names the telescope: 1 = Catalina Sky Survey
 * (0.7 m), 2 = Mount Lemmon Survey (1.5 m), 3 = Siding Spring Survey (0.5 m).
 */
const TELESCOPES: Record<string, string> = {
  "1": "CSS",
  "2": "MLS",
  "3": "SSS",
};

export function crtsTelescope(masterId: string): string {
  return TELESCOPES[masterId[0]] ?? "CRTS";
}

/** Brighter than this median the detector saturates; treat with care. */
export const CRTS_SATURATION_MAG = 12.5;
/** Sources whose mean position lies within this of the target are merged. */
export const CRTS_MATCH_RADIUS_ARCSEC = 3;

export interface CrtsSource {
  id: string;
  telescope: string;
  epochs: number;
  sepArcsec: number;
}

export interface CrtsLightcurve {
  points: DetectionPoint[];
  sources: CrtsSource[];
  /** Isolated spikes removed from `points`. */
  outliers: number;
  saturated: boolean;
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

function median(values: number[]): number {
  const s = [...values].sort((a, b) => a - b);
  const mid = s.length >> 1;
  return s.length % 2 ? s[mid] : (s[mid - 1] + s[mid]) / 2;
}

/**
 * Turn the CRTS CSV into light-curve points, one band per telescope.
 *
 * The same star often has several MasterIDs (overlapping fields and different
 * telescopes), so every ID whose mean position is within
 * {@link CRTS_MATCH_RADIUS_ARCSEC} of the target is kept, not just the
 * nearest.
 *
 * CSS takes ~4 exposures per field in ~30 minutes, so a real eclipse or flare
 * shows in all of that night's exposures while a spurious measurement stands
 * alone. A point is dropped only when it is > 1 mag from its telescope's
 * median and > 0.75 mag from the median of its night-mates.
 */
export function parseCrtsCsv(
  csv: string,
  ra: number,
  dec: number
): CrtsLightcurve {
  const lines = csv.split(/\r?\n/).filter((l) => l.trim() !== "");
  const names = (lines[0] ?? "").split(",").map((s) => s.trim());
  const col = (n: string) => names.indexOf(n);
  const [iId, iMag, iErr, iRa, iDec, iMjd] = [
    "MasterID",
    "Mag",
    "Magerr",
    "RA",
    "Dec",
    "MJD",
  ].map(col);
  const empty: CrtsLightcurve = {
    points: [],
    sources: [],
    outliers: 0,
    saturated: false,
  };
  if ([iId, iMag, iRa, iDec, iMjd].some((i) => i === -1)) return empty;

  interface Row {
    id: string;
    mag: number;
    err: number;
    mjd: number;
  }
  const byId = new Map<string, { rows: Row[]; ra: number; dec: number }>();
  for (const line of lines.slice(1)) {
    const c = line.split(",").map((s) => s.trim());
    const row: Row = {
      id: c[iId],
      mag: Number(c[iMag]),
      err: iErr === -1 ? NaN : Number(c[iErr]),
      mjd: Number(c[iMjd]),
    };
    if (!row.id || !Number.isFinite(row.mag) || !Number.isFinite(row.mjd)) {
      continue;
    }
    const g = byId.get(row.id) ?? { rows: [], ra: 0, dec: 0 };
    g.rows.push(row);
    g.ra += Number(c[iRa]);
    g.dec += Number(c[iDec]);
    byId.set(row.id, g);
  }

  const sources: CrtsSource[] = [];
  const kept: Row[] = [];
  for (const [id, g] of byId) {
    const sep = sepArcsec(ra, dec, g.ra / g.rows.length, g.dec / g.rows.length);
    if (sep > CRTS_MATCH_RADIUS_ARCSEC) continue;
    sources.push({
      id,
      telescope: crtsTelescope(id),
      epochs: g.rows.length,
      sepArcsec: Math.round(sep * 100) / 100,
    });
    kept.push(...g.rows);
  }
  if (kept.length === 0) return empty;

  // Outlier test per telescope, against the telescope median and the other
  // exposures of the same night (floor(MJD) never splits a night at these
  // longitudes).
  const byTelescope = new Map<string, Row[]>();
  for (const r of kept) {
    const t = crtsTelescope(r.id);
    byTelescope.set(t, [...(byTelescope.get(t) ?? []), r]);
  }
  const points: DetectionPoint[] = [];
  let outliers = 0;
  for (const [telescope, rows] of byTelescope) {
    const med = median(rows.map((r) => r.mag));
    const nights = new Map<number, Row[]>();
    for (const r of rows) {
      const n = Math.floor(r.mjd);
      nights.set(n, [...(nights.get(n) ?? []), r]);
    }
    for (const r of rows) {
      const mates = nights.get(Math.floor(r.mjd))!.filter((m) => m !== r);
      const spike =
        mates.length > 0 &&
        Math.abs(r.mag - med) > 1 &&
        Math.abs(r.mag - median(mates.map((m) => m.mag))) > 0.75;
      if (spike) {
        outliers += 1;
        continue;
      }
      points.push({
        mjd: r.mjd,
        mag: r.mag,
        magerr: r.err > 0 ? r.err : undefined,
        band: telescope,
      });
    }
  }

  return {
    points: points.sort((a, b) => a.mjd! - b.mjd!),
    sources: sources.sort((a, b) => a.sepArcsec - b.sepArcsec),
    outliers,
    saturated: median(kept.map((r) => r.mag)) < CRTS_SATURATION_MAG,
  };
}
