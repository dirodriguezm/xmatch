import type { components } from "@/types/xwave-api";

import { csvCell } from "./csv";

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

const SENTINEL = -999;

const CATALOG_LABELS: Record<string, string> = {
  crts: "Catalina (CRTS)",
  gaia: "Gaia DR3",
  neowise: "NEOWISE",
  ps1: "Pan-STARRS DR2",
  swift: "Swift",
  vlass: "VLASS",
  ztf: "ZTF",
};

/**
 * Panel order for light curves: infrared before optical, then by when each
 * survey observed (Catalina 2005–16, Pan-STARRS 2009–14, Gaia 2014–17,
 * ZTF 2018–). Surveys not listed go last, alphabetically.
 */
const LIGHTCURVE_ORDER = ["neowise", "allwise", "crts", "ps1", "gaia", "ztf"];

export function compareLightcurveCatalogs(a: string, b: string): number {
  const rank = (c: string) => {
    const i = LIGHTCURVE_ORDER.indexOf(c.toLowerCase());
    return i === -1 ? LIGHTCURVE_ORDER.length : i;
  };
  return rank(a) - rank(b) || a.localeCompare(b);
}

export function getCatalogLabel(catalog: string): string {
  return CATALOG_LABELS[catalog.toLowerCase()] ?? catalog.toUpperCase();
}

/**
 * Magnitude system of each light-curve survey. Panels are per survey, so AB
 * and Vega never share an axis, but the axis and the CSV say which is which.
 */
const MAG_SYSTEMS: Record<string, "AB" | "Vega"> = {
  // Unfiltered, calibrated to Johnson V.
  crts: "Vega",
  gaia: "Vega",
  neowise: "Vega",
  ps1: "AB",
  ztf: "AB",
};

export function getMagSystem(catalog: string): "AB" | "Vega" | undefined {
  return MAG_SYSTEMS[catalog.toLowerCase()];
}

export interface LightcurveDetection {
  catalog: string;
  id?: string;
  object_id?: string;
  mjd?: number;
  mag?: number;
  magerr?: number;
  data?: Record<string, number | string | null | undefined>;
}

export interface DetectionPoint {
  mjd?: number;
  mag?: number;
  magerr?: number;
  band?: string;
}

function isValidNumber(v: unknown): v is number {
  return typeof v === "number" && Number.isFinite(v) && v !== SENTINEL;
}

const ZTF_FILTER_BAND: Record<number, string> = { 1: "g", 2: "r", 3: "i" };

export function expandDetection(det: LightcurveDetection): DetectionPoint[] {
  const catalog = (det.catalog || "").toLowerCase();
  const data = det.data ?? {};

  if (catalog === "ztf") {
    if (!isValidNumber(det.mjd) || !isValidNumber(det.mag)) return [];
    const fid = data.filterid;
    const band = isValidNumber(fid) ? ZTF_FILTER_BAND[fid] : undefined;
    return [
      {
        mjd: det.mjd,
        mag: det.mag,
        magerr: isValidNumber(det.magerr) ? det.magerr : undefined,
        band: band ?? "ZTF",
      },
    ];
  }

  if (catalog === "neowise") {
    const points: DetectionPoint[] = [];
    const w1 = data.w1mpro;
    if (isValidNumber(det.mjd) && isValidNumber(w1)) {
      const w1err = data.w1sigmpro;
      points.push({
        mjd: det.mjd,
        mag: w1,
        magerr: isValidNumber(w1err) ? w1err : undefined,
        band: "W1",
      });
    }
    const w2 = data.w2mpro;
    if (isValidNumber(det.mjd) && isValidNumber(w2)) {
      const w2err = data.w2sigmpro;
      points.push({
        mjd: det.mjd,
        mag: w2,
        magerr: isValidNumber(w2err) ? w2err : undefined,
        band: "W2",
      });
    }
    return points;
  }

  if (!isValidNumber(det.mjd) || !isValidNumber(det.mag)) return [];
  return [
    {
      mjd: det.mjd,
      mag: det.mag,
      magerr: isValidNumber(det.magerr) ? det.magerr : undefined,
      band: catalog.toUpperCase(),
    },
  ];
}

/**
 * Serialize grouped detection points into a CSV string with columns
 * `survey,band,mag_system,mjd,mag,magerr`. The survey column uses the human-readable
 * catalog label via {@link getCatalogLabel}.
 */
export function detectionPointsToCsv(
  groups: Record<string, DetectionPoint[]>
): string {
  const rows: string[] = ["survey,band,mag_system,mjd,mag,magerr"];
  for (const [catalog, points] of Object.entries(groups)) {
    const survey = getCatalogLabel(catalog);
    for (const p of points) {
      rows.push(
        [
          csvCell(survey),
          csvCell(p.band),
          csvCell(getMagSystem(catalog)),
          csvCell(p.mjd),
          csvCell(p.mag),
          csvCell(p.magerr),
        ].join(",")
      );
    }
  }
  return rows.join("\n");
}

// downloadCsv now lives in ./csv alongside the escaping helpers; re-exported
// here so existing importers keep working.
export { downloadCsv } from "./csv";

export function groupDetectionsByCatalog(
  lc: Lightcurve | null | undefined,
  exclude: string[] = []
): Record<string, DetectionPoint[]> {
  if (!lc?.detections?.length) return {};
  const skip = new Set(exclude.map((c) => c.toLowerCase()));
  const out: Record<string, DetectionPoint[]> = {};
  for (const raw of lc.detections as unknown[]) {
    if (typeof raw !== "object" || raw === null) continue;
    const det = raw as LightcurveDetection;
    const catalog = (det.catalog || "").toLowerCase();
    if (!catalog || skip.has(catalog)) continue;
    const points = expandDetection(det);
    if (points.length === 0) continue;
    (out[catalog] ??= []).push(...points);
  }
  return out;
}
