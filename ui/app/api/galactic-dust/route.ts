import { NextRequest, NextResponse } from "next/server";

import { parseIrsaDust } from "@/app/lib/utils/irsaDust";

/**
 * Server-side proxy to the IRSA Galactic Dust Reddening service. Returns the
 * E(B−V) at the position from Schlafly & Finkbeiner (2011) and SFD (1998).
 */

const DUST_URL = "https://irsa.ipac.caltech.edu/cgi-bin/DUST/nph-dust";
const TIMEOUT_MS = 15_000;
/** The dust maps are static; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));

  if (!Number.isFinite(ra) || !Number.isFinite(dec)) {
    return NextResponse.json(
      { error: "Missing or invalid 'ra' / 'dec'" },
      { status: 400 }
    );
  }

  // regSize is the size of the statistics box; 2° is the service minimum.
  // Only the reference-pixel value at the position is used.
  const url = `${DUST_URL}?locstr=${encodeURIComponent(`${ra} ${dec}`)}&regSize=2.0`;

  try {
    const resp = await fetch(url, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "IRSA dust request failed", status: resp.status },
        { status: 502 }
      );
    }
    const reddening = parseIrsaDust(await resp.text());
    return NextResponse.json(reddening, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("IRSA dust error:", err);
    return NextResponse.json(
      { error: timedOut ? "IRSA dust timed out" : "Failed to read IRSA dust" },
      { status: timedOut ? 504 : 502 }
    );
  }
}
