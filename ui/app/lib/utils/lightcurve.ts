import type { components } from "@/types/xwave-api";

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

const SENTINEL = -999;

const CATALOG_LABELS: Record<string, string> = {
  neowise: "NEOWISE",
  swift: "Swift",
  vlass: "VLASS",
  ztf: "ZTF",
};

export function getCatalogLabel(catalog: string): string {
  return CATALOG_LABELS[catalog.toLowerCase()] ?? catalog.toUpperCase();
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
