import { NextRequest, NextResponse } from "next/server";

import {
  parseSimbadRefs,
  SIMBAD_REFS_MAX,
  type SimbadReference,
  simbadRefsCountQuery,
  simbadRefsQuery,
} from "@/app/lib/utils/simbadRefs";

/** Papers about a SIMBAD object, newest first, plus the total count. */

const TAP_URL = "https://simbad.cds.unistra.fr/simbad/sim-tap/sync";
const TIMEOUT_MS = 15_000;
/** New papers are attached daily; a day of caching is fine. */
const REVALIDATE_S = 86_400;

export interface SimbadRefsResponse {
  total: number;
  references: SimbadReference[];
}

async function tap(query: string): Promise<unknown[][] | undefined> {
  const qs = new URLSearchParams({
    request: "doQuery",
    lang: "adql",
    format: "json",
    query,
  });
  const resp = await fetch(`${TAP_URL}?${qs}`, {
    signal: AbortSignal.timeout(TIMEOUT_MS),
    next: { revalidate: REVALIDATE_S },
  });
  if (!resp.ok) throw new Error(`SIMBAD returned ${resp.status}`);
  return ((await resp.json()) as { data?: unknown[][] }).data;
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const oid = Number(sp.get("oid"));
  const limit = Number(sp.get("limit") ?? 25);
  if (!Number.isInteger(oid) || oid <= 0) {
    return NextResponse.json({ error: "Invalid 'oid'" }, { status: 400 });
  }
  if (!Number.isInteger(limit) || limit < 1 || limit > SIMBAD_REFS_MAX) {
    return NextResponse.json(
      { error: `'limit' must be an integer from 1 to ${SIMBAD_REFS_MAX}` },
      { status: 400 }
    );
  }

  try {
    // oid and limit are validated integers, so inlining them is safe.
    const [rows, count] = await Promise.all([
      tap(simbadRefsQuery(oid, limit)),
      tap(simbadRefsCountQuery(oid)),
    ]);
    const total = Number(count?.[0]?.[0]);
    const references = parseSimbadRefs(rows);
    const result: SimbadRefsResponse = {
      total: Number.isFinite(total) ? total : references.length,
      references,
    };
    return NextResponse.json(result, {
      headers: {
        "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
      },
    });
  } catch (err) {
    console.error("SIMBAD references error:", err);
    return NextResponse.json(
      { error: "Failed to load references from SIMBAD" },
      { status: 502 }
    );
  }
}
