import { NextRequest, NextResponse } from "next/server";

import {
  parseXmmSourceCsv,
  xmmConeQuery,
  type XmmSource,
} from "@/app/lib/utils/xmm5";

/**
 * Nearest 5XMM-DR15 source (stacked XMM-Newton serendipitous catalogue) from
 * the XMM-Newton Science Archive TAP service.
 */

const TAP_URL = "https://nxsa.esac.esa.int/tap-server/tap/sync";
const TIMEOUT_MS = 20_000;
/** The catalogue release is frozen; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export interface XmmSourceResponse {
  found: boolean;
  reason?: "no_source";
  source?: XmmSource;
}

function respond(result: XmmSourceResponse) {
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
    // ra/dec are validated numbers, so inlining them in ADQL is safe.
    const resp = await fetch(TAP_URL, {
      method: "POST",
      body: new URLSearchParams({
        REQUEST: "doQuery",
        LANG: "ADQL",
        FORMAT: "csv",
        QUERY: xmmConeQuery(ra, dec),
      }),
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "XMM-Newton archive request failed", status: resp.status },
        { status: 502 }
      );
    }
    const source = parseXmmSourceCsv(await resp.text(), ra, dec);
    return respond(
      source ? { found: true, source } : { found: false, reason: "no_source" }
    );
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("5XMM lookup error:", err);
    return NextResponse.json(
      {
        error: timedOut
          ? "XMM-Newton archive timed out"
          : "XMM-Newton archive request failed",
      },
      { status: timedOut ? 504 : 502 }
    );
  }
}
