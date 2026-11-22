import { NextRequest, NextResponse } from "next/server";

import {
  fetchFromMirrors,
  UPSTREAM_HOST_HEADER,
  VIZIER_SED_MIRRORS,
} from "@/app/lib/utils/mirrors";
import {
  aggregateVizierSed,
  parseVizierSedVotable,
  type VizierSedPoint,
  type VizierSedRow,
} from "@/app/lib/utils/vizierSed";

/**
 * Server-side proxy to the VizieR SED service (CDS, Strasbourg, or the CfA
 * mirror). CDS answers in ~10 s with a VOTable of every published
 * photometric point near the position (often hundreds of rows, ~100 KB);
 * CfA can take 1–2 min for a position it hasn't cached. The rows are parsed
 * and collapsed to one point per filter here so the browser receives a
 * small JSON body.
 */

/**
 * Room for CfA's slow first answer plus retries: the last retry starts by
 * RETRY_DEADLINE_MS and takes at most ~66 s (CDS 404 + CfA's 65 s timeout).
 */
export const maxDuration = 200;

const DEFAULT_RADIUS_ARCSEC = 2;
const MAX_RADIUS_ARCSEC = 5;
/** Default per mirror; VIZIER_SED_MIRRORS sets each one's own. */
const TIMEOUT_MS = 20_000;
/**
 * Published photometry barely changes; the CDN keeps successful answers for
 * a week, so a slow mirror lookup is paid once per position per week, not
 * per visitor.
 */
const REVALIDATE_S = 7 * 86_400;
/**
 * CfA's gateway gives up after 60 s (504) on a position it hasn't cached,
 * and right after may serve a VOTable without the SED columns, but it keeps
 * computing and has the answer moments later. So such answers are retried.
 */
const RETRY_DELAY_MS = 8_000;
/** No retry starts after this, so the answer fits within `maxDuration`. */
const RETRY_DEADLINE_MS = 120_000;

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

type SedFetch =
  | { ok: true; rows: VizierSedRow[]; host: string }
  | { ok: false; status: number };

/** VizieR SED rows from the first mirror that gives a complete answer. */
async function fetchSedRows(query: string): Promise<SedFetch> {
  const started = Date.now();
  for (;;) {
    // Not stored in Next's data cache: a 504 or a half-built VOTable must
    // not be replayed for a week. The CDN caches the parsed success instead.
    const { response, host } = await fetchFromMirrors(
      VIZIER_SED_MIRRORS,
      query,
      { cache: "no-store" },
      TIMEOUT_MS
    );
    const body = await response.text();
    let failure: number | null = response.ok ? null : response.status;
    if (failure === null) {
      try {
        return { ok: true, rows: parseVizierSedVotable(body), host };
      } catch (err) {
        console.warn(`VizieR SED from ${host} unreadable, retrying:`, err);
        failure = 502;
      }
    }
    const retryable = failure >= 500;
    if (!retryable || Date.now() - started + RETRY_DELAY_MS > RETRY_DEADLINE_MS)
      return { ok: false, status: failure };
    await sleep(RETRY_DELAY_MS);
  }
}

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

  const query = `?-c=${encodeURIComponent(`${ra} ${dec}`)}&-c.rs=${radius}`;

  try {
    const sed = await fetchSedRows(query);
    if (!sed.ok) {
      return NextResponse.json(
        { error: "VizieR SED request failed", status: sed.status },
        { status: 502 }
      );
    }

    const { rows, host } = sed;
    const points = aggregateVizierSed(rows);
    const result: VizierSedResponse = {
      found: points.length > 0,
      points,
      rowCount: rows.length,
    };
    return NextResponse.json(result, {
      headers: {
        [UPSTREAM_HOST_HEADER]: host,
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
