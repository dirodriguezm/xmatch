import { NextRequest, NextResponse } from "next/server";

import {
  type ChandraSource,
  cscConeUrl,
  parseCscCone,
} from "@/app/lib/utils/csc";

/**
 * Nearest Chandra Source Catalog 2.2 master source, from the CXC cone search.
 */

const TIMEOUT_MS = 20_000;
/** The catalogue release is frozen; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export interface ChandraSourceResponse {
  found: boolean;
  reason?: "no_source";
  source?: ChandraSource;
}

function respond(result: ChandraSourceResponse) {
  return NextResponse.json(result, {
    headers: {
      "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
    },
  });
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));
  if (
    !sp.has("ra") ||
    !sp.has("dec") ||
    !Number.isFinite(ra) ||
    !Number.isFinite(dec) ||
    ra < 0 ||
    ra > 360 ||
    dec < -90 ||
    dec > 90
  ) {
    return NextResponse.json(
      { error: "Missing or invalid 'ra' / 'dec'" },
      { status: 400 }
    );
  }

  try {
    const resp = await fetch(cscConeUrl(ra, dec), {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "Chandra Source Catalog request failed", status: resp.status },
        { status: 502 }
      );
    }
    const source = parseCscCone(await resp.text(), ra, dec);
    return respond(
      source ? { found: true, source } : { found: false, reason: "no_source" }
    );
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("CSC lookup error:", err);
    return NextResponse.json(
      {
        error: timedOut
          ? "Chandra Source Catalog timed out"
          : "Chandra Source Catalog request failed",
      },
      { status: timedOut ? 504 : 502 }
    );
  }
}
