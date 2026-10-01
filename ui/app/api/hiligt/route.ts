import { NextRequest, NextResponse } from "next/server";

import {
  type HiligtMission,
  type HiligtPoint,
  isHiligtMission,
  parseHiligt,
} from "@/app/lib/utils/hiligt";

/**
 * X-ray fluxes and upper limits for one mission from ESA's HILIGT upper limit
 * server. One mission per request: each takes 5–20 s, so the page shows them
 * as they arrive instead of waiting for the slowest.
 */

const ULS_URL = "https://xmmuls.esac.esa.int/ULSservice_passthru";
const TIMEOUT_MS = 55_000;
/** New observations reach HILIGT slowly; a day of caching is plenty. */
const REVALIDATE_S = 86_400;

export const maxDuration = 60;

export interface HiligtResponse {
  mission: HiligtMission;
  points: HiligtPoint[];
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const mission = sp.get("mission") ?? "";
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));
  // HILIGT answers an unknown mission with an empty 200, so check it here.
  if (!isHiligtMission(mission)) {
    return NextResponse.json(
      { error: "Unknown or missing 'mission'" },
      { status: 400 }
    );
  }
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

  const qs = new URLSearchParams({
    mission,
    ra: String(ra),
    dec: String(dec),
    label: "xwave",
    FORMAT: "JSON",
  });

  try {
    const resp = await fetch(`${ULS_URL}?${qs}`, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: "HILIGT request failed", status: resp.status },
        { status: 502 }
      );
    }
    const result: HiligtResponse = {
      mission,
      points: parseHiligt(await resp.text(), mission),
    };
    return NextResponse.json(result, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("HILIGT error:", err);
    return NextResponse.json(
      { error: timedOut ? "HILIGT timed out" : "HILIGT request failed" },
      { status: timedOut ? 504 : 502 }
    );
  }
}
