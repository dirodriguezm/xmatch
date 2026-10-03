import { NextRequest, NextResponse } from "next/server";

import type { DetectionPoint } from "@/app/lib/utils/lightcurve";
import {
  parsePs1DetectionsCsv,
  PS1_MIN_DEC,
} from "@/app/lib/utils/ps1Lightcurve";

/**
 * Pan-STARRS1 DR2 per-epoch detections from the MAST catalogs API, cone
 * search of 1″ around the position. Detections with psfQfPerfect < 0.9
 * (significantly masked PSF) are excluded server-side, as MAST's own
 * examples recommend.
 */

const DETECTION_URL =
  "https://catalogs.mast.stsci.edu/api/v0.1/panstarrs/dr2/detection.csv";
const RADIUS_DEG = 1 / 3600;
const COLUMNS = "[objID,obsTime,filterID,psfFlux,psfFluxErr,ra,dec]";
const TIMEOUT_MS = 20_000;
/** DR2 is frozen; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export interface Ps1LightcurveResponse {
  found: boolean;
  objId: string | null;
  points: DetectionPoint[];
  reason?: "outside_footprint";
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));
  if (
    !sp.has("ra") ||
    !sp.has("dec") ||
    !Number.isFinite(ra) ||
    !Number.isFinite(dec)
  ) {
    return NextResponse.json(
      { error: "Missing or invalid 'ra' / 'dec'" },
      { status: 400 }
    );
  }

  const headers = {
    "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
  };

  if (dec < PS1_MIN_DEC - 1) {
    const result: Ps1LightcurveResponse = {
      found: false,
      objId: null,
      points: [],
      reason: "outside_footprint",
    };
    return NextResponse.json(result, { headers });
  }

  const qs = new URLSearchParams({
    ra: String(ra),
    dec: String(dec),
    radius: String(RADIUS_DEG),
    columns: COLUMNS,
    "psfQfPerfect.min": "0.9",
    pagesize: "50000",
  });

  try {
    const resp = await fetch(`${DETECTION_URL}?${qs}`, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "MAST request failed", status: resp.status },
        { status: 502 }
      );
    }
    const { objId, points } = parsePs1DetectionsCsv(await resp.text(), ra, dec);
    const result: Ps1LightcurveResponse = {
      found: points.length > 0,
      objId,
      points,
    };
    return NextResponse.json(result, { headers });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("PS1 light curve error:", err);
    return NextResponse.json(
      { error: timedOut ? "MAST timed out" : "MAST request failed" },
      { status: timedOut ? 504 : 502 }
    );
  }
}
