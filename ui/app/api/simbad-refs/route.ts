import { NextRequest, NextResponse } from "next/server";

import {
  parseSimbadRefs,
  SIMBAD_REFS_MAX,
  type SimbadReference,
  simbadRefsCountQuery,
  type SimbadRefsOptions,
  simbadRefsQuery,
} from "@/app/lib/utils/simbadRefs";

/** Papers about a SIMBAD object, plus how many match the year filter. */

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
  const sort = sp.get("sort") ?? "newest";
  if (sort !== "newest" && sort !== "focused") {
    return NextResponse.json(
      { error: "'sort' must be 'newest' or 'focused'" },
      { status: 400 }
    );
  }
  const year = (name: string) => {
    const v = sp.get(name);
    if (v === null) return undefined;
    const n = Number(v);
    return Number.isInteger(n) && n >= 1800 && n <= 2200 ? n : NaN;
  };
  const fromYear = year("from");
  const toYear = year("to");
  if (Number.isNaN(fromYear) || Number.isNaN(toYear)) {
    return NextResponse.json(
      { error: "'from' and 'to' must be years" },
      { status: 400 }
    );
  }
  const options: SimbadRefsOptions = { sort, fromYear, toYear };

  try {
    // Every inlined value is a validated integer or one of two sort names.
    const [rows, count] = await Promise.all([
      tap(simbadRefsQuery(oid, limit, options)),
      tap(simbadRefsCountQuery(oid, options)),
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
