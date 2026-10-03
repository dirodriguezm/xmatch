import { NextRequest, NextResponse } from "next/server";

import { ApiError, apiFetch } from "@/app/lib/api/client";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";
import { type BulkApiMatch, MAX_BULK_ROWS } from "@/app/lib/utils/bulkInput";

/**
 * Proxy for `POST /bulk-conesearch`.
 *
 * `radius` is in ARCSECONDS. The Go request struct documents it as degrees, but
 * the service converts with `arcsecToRadians` and returns `distance` in arcsec
 * (verified: radius 5 against Betelgeuse returns a match at 1.17″).
 */
const CATALOGS = new Set(["all", "gaia", "allwise", "erosita"]);
const MAX_NNEIGHBOR = 100;

function isNumberArray(value: unknown): value is number[] {
  return (
    Array.isArray(value) &&
    value.every((v) => typeof v === "number" && Number.isFinite(v))
  );
}

export async function POST(request: NextRequest) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ error: "Invalid JSON body" }, { status: 400 });
  }

  const { ra, dec, radius, catalog, nneighbor } = (body ?? {}) as Record<
    string,
    unknown
  >;

  if (!isNumberArray(ra) || !isNumberArray(dec)) {
    return NextResponse.json(
      { error: "Missing or invalid fields: ra, dec (expected number[])" },
      { status: 400 }
    );
  }
  if (ra.length !== dec.length) {
    return NextResponse.json(
      { error: `ra and dec lengths differ (${ra.length} vs ${dec.length})` },
      { status: 400 }
    );
  }
  if (ra.length === 0) {
    return NextResponse.json([]);
  }
  if (ra.length > MAX_BULK_ROWS) {
    return NextResponse.json(
      { error: `Too many positions: ${ra.length} (max ${MAX_BULK_ROWS})` },
      { status: 400 }
    );
  }
  if (
    typeof radius !== "number" ||
    !(radius > 0) ||
    radius > MAX_RADIUS_ARCSEC
  ) {
    return NextResponse.json(
      {
        error: `Invalid radius: expected arcsec in (0, ${MAX_RADIUS_ARCSEC}]`,
      },
      { status: 400 }
    );
  }
  const cat = typeof catalog === "string" && catalog ? catalog : "all";
  if (!CATALOGS.has(cat)) {
    return NextResponse.json(
      { error: `Unknown catalog: ${cat}` },
      { status: 400 }
    );
  }
  const n = nneighbor === undefined ? 1 : nneighbor;
  if (
    typeof n !== "number" ||
    !Number.isInteger(n) ||
    n < 1 ||
    n > MAX_NNEIGHBOR
  ) {
    return NextResponse.json(
      {
        error: `Invalid nneighbor: expected an integer in [1, ${MAX_NNEIGHBOR}]`,
      },
      { status: 400 }
    );
  }

  try {
    const result = await apiFetch<BulkApiMatch[]>("/bulk-conesearch", {
      method: "POST",
      body: JSON.stringify({ ra, dec, radius, catalog: cat, nneighbor: n }),
    });
    // 204 (no matches anywhere) → apiFetch returns null.
    return NextResponse.json(Array.isArray(result) ? result : []);
  } catch (error) {
    console.error("Bulk conesearch proxy error:", error);
    if (error instanceof ApiError) {
      return NextResponse.json(
        { error: error.message },
        { status: error.status }
      );
    }
    return NextResponse.json(
      { error: "Failed to fetch from bulk conesearch service" },
      { status: 500 }
    );
  }
}
