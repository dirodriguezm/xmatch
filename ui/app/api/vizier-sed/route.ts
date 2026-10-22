import { NextRequest, NextResponse } from "next/server";

import {
  aggregateVizierSed,
  parseVizierSedVotable,
  type VizierSedPoint,
} from "@/app/lib/utils/vizierSed";

/**
 * Server-side proxy to the VizieR SED service (CDS, Strasbourg). The service
 * answers in ~10 s with a VOTable of every published photometric point near
 * the position (often hundreds of rows, ~100 KB); it is parsed and collapsed
 * to one point per filter here so the browser receives a small JSON body.
 */

const SED_URL = "https://vizier.cds.unistra.fr/viz-bin/sed";
const DEFAULT_RADIUS_ARCSEC = 2;
const MAX_RADIUS_ARCSEC = 5;
const TIMEOUT_MS = 20_000;
/** Published photometry does not change; cache for a day. */
const REVALIDATE_S = 86_400;

export interface VizierSedResponse {
  found: boolean;
  points: VizierSedPoint[];
  /** Raw rows VizieR returned, before merging per filter. */
  rowCount: number;
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));
  const radius = Number(sp.get("radius") ?? DEFAULT_RADIUS_ARCSEC);

  if (!Number.isFinite(ra) || !Number.isFinite(dec)) {
    return NextResponse.json(
      { error: "Missing or invalid 'ra' / 'dec'" },
      { status: 400 }
    );
  }
  if (!Number.isFinite(radius) || radius <= 0 || radius > MAX_RADIUS_ARCSEC) {
    return NextResponse.json(
      {
        error: `'radius' must be a positive number ≤ ${MAX_RADIUS_ARCSEC} arcsec`,
      },
      { status: 400 }
    );
  }

  const url = `${SED_URL}?-c=${encodeURIComponent(`${ra} ${dec}`)}&-c.rs=${radius}`;

  try {
    const resp = await fetch(url, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    const body = await resp.text();
    if (!resp.ok) {
      return NextResponse.json(
        { error: "VizieR SED request failed", status: resp.status },
        { status: 502 }
      );
    }

    const rows = parseVizierSedVotable(body);
    const points = aggregateVizierSed(rows);
    const result: VizierSedResponse = {
      found: points.length > 0,
      points,
      rowCount: rows.length,
    };
    return NextResponse.json(result, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("VizieR SED error:", err);
    return NextResponse.json(
      {
        error: timedOut
          ? "VizieR SED timed out"
          : "Failed to read VizieR SED response",
      },
      { status: timedOut ? 504 : 502 }
    );
  }
}
