import { NextRequest, NextResponse } from "next/server";

import {
  batLightcurveUrl,
  type GammaPoint,
  isBatPath,
  parseBatLightcurve,
} from "@/app/lib/utils/gamma";

/**
 * Daily 15–50 keV light curve from the Swift-BAT Hard X-ray Transient
 * Monitor. Only sources in app/lib/constants/gamma-sources.json are served.
 */

const TIMEOUT_MS = 30_000;
/** The monitor updates within hours; a day of caching is plenty here. */
const REVALIDATE_S = 86_400;

export interface BatLightcurveResponse {
  path: string;
  points: GammaPoint[];
}

export async function GET(request: NextRequest) {
  const path = request.nextUrl.searchParams.get("path") ?? "";
  if (!isBatPath(path)) {
    return NextResponse.json({ error: "Unknown 'path'" }, { status: 400 });
  }

  try {
    const resp = await fetch(batLightcurveUrl(path), {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "Swift-BAT monitor request failed", status: resp.status },
        { status: 502 }
      );
    }
    const result: BatLightcurveResponse = {
      path,
      points: parseBatLightcurve(await resp.text()),
    };
    return NextResponse.json(result, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("Swift-BAT monitor error:", err);
    return NextResponse.json(
      {
        error: timedOut
          ? "Swift-BAT monitor timed out"
          : "Swift-BAT monitor request failed",
      },
      { status: timedOut ? 504 : 502 }
    );
  }
}
