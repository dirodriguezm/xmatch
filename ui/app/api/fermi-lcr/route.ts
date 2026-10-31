import { NextRequest, NextResponse } from "next/server";

import {
  type GammaPoint,
  isFermiSource,
  isLcrCadence,
  parseLcrLightcurve,
} from "@/app/lib/utils/gamma";

/**
 * Light curve of a 4FGL source from the Fermi-LAT Light Curve Repository.
 * Only sources in app/lib/constants/gamma-sources.json are served.
 */

const LCR_URL =
  "https://fermi.gsfc.nasa.gov/ssc/data/access/lat/LightCurveRepository/queryDB.php";
const TIMEOUT_MS = 30_000;
/** The repository adds a bin every few days at most. */
const REVALIDATE_S = 86_400;

export interface FermiLcrResponse {
  source: string;
  cadence: string;
  points: GammaPoint[];
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const source = sp.get("source") ?? "";
  const cadence = sp.get("cadence") ?? "weekly";
  if (!isFermiSource(source) || !isLcrCadence(cadence)) {
    return NextResponse.json(
      { error: "Unknown 'source' or 'cadence'" },
      { status: 400 }
    );
  }

  try {
    const params = new URLSearchParams({
      typeOfRequest: "lightCurveData",
      source_name: source,
      cadence,
      flux_type: "photon",
      index_type: "fixed",
      ts_min: "4",
    });
    const resp = await fetch(`${LCR_URL}?${params}`, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "Fermi LCR request failed", status: resp.status },
        { status: 502 }
      );
    }
    const result: FermiLcrResponse = {
      source,
      cadence,
      points: parseLcrLightcurve(await resp.json()),
    };
    return NextResponse.json(result, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("Fermi LCR error:", err);
    return NextResponse.json(
      { error: timedOut ? "Fermi LCR timed out" : "Fermi LCR request failed" },
      { status: timedOut ? 504 : 502 }
    );
  }
}
