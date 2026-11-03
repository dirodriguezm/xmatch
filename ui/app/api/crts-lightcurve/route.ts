import { NextRequest, NextResponse } from "next/server";

import {
  type CrtsSource,
  parseCrtsCsv,
  parseCrtsQueryPage,
} from "@/app/lib/utils/crtsLightcurve";
import type { DetectionPoint } from "@/app/lib/utils/lightcurve";

/**
 * Catalina (CRTS) per-epoch photometry from Caltech's cone search. The
 * service is HTTP-only and answers with an HTML page linking a temporary CSV,
 * so it has to go through this proxy. Asking for the page without its plot
 * and image cutouts brings the query from ~5 s down to ~0.5 s.
 */

const SEARCH_URL =
  "http://nunuku.caltech.edu/cgi-bin/getcssconedb_release_img.cgi";
/** Cone radius in arcmin (6″); sources are then matched within 3″. */
const RADIUS_ARCMIN = 0.1;
const TIMEOUT_MS = 20_000;
/** The data release is frozen; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export interface CrtsLightcurveResponse {
  found: boolean;
  /** Why nothing was found: area outside the survey, or no source there. */
  reason?: "not_covered" | "no_source";
  points: DetectionPoint[];
  sources: CrtsSource[];
  outliers: number;
  saturated: boolean;
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
  const nothing = (
    reason: CrtsLightcurveResponse["reason"]
  ): CrtsLightcurveResponse => ({
    found: false,
    reason,
    points: [],
    sources: [],
    outliers: 0,
    saturated: false,
  });

  const form = new FormData();
  form.set("RADec", `${ra} ${dec}`);
  form.set("Rad", String(RADIUS_ARCMIN));
  form.set("DB", "photcat");
  form.set("OUT", "csv");
  form.set("SHORT", "short");
  form.set(".submit", "Submit");

  try {
    const page = await fetch(SEARCH_URL, {
      method: "POST",
      body: form,
      signal: AbortSignal.timeout(TIMEOUT_MS),
    });
    if (!page.ok) {
      return NextResponse.json(
        { error: "CRTS request failed", status: page.status },
        { status: 502 }
      );
    }
    const { status, csvUrl } = parseCrtsQueryPage(await page.text());
    if (status === "not_covered") {
      return NextResponse.json(nothing("not_covered"), { headers });
    }
    if (status === "empty") {
      return NextResponse.json(nothing("no_source"), { headers });
    }
    if (!csvUrl) {
      return NextResponse.json(
        { error: "Unexpected CRTS response" },
        { status: 502 }
      );
    }

    const csv = await fetch(csvUrl, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
    });
    if (!csv.ok) {
      return NextResponse.json(
        { error: "CRTS CSV download failed", status: csv.status },
        { status: 502 }
      );
    }
    const lc = parseCrtsCsv(await csv.text(), ra, dec);
    const result: CrtsLightcurveResponse =
      lc.points.length > 0
        ? { found: true, ...lc }
        : { ...nothing("no_source"), outliers: lc.outliers };
    return NextResponse.json(result, { headers });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("CRTS light curve error:", err);
    return NextResponse.json(
      { error: timedOut ? "CRTS timed out" : "CRTS request failed" },
      { status: timedOut ? 504 : 502 }
    );
  }
}
