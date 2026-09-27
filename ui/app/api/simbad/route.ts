import { NextRequest, NextResponse } from "next/server";

/**
 * Server-side proxy to the SIMBAD TAP service (CDS, Strasbourg): the nearest
 * SIMBAD object to a position, with its main identifier and, when SIMBAD has
 * one, its redshift or radial velocity.
 */

const TAP_URL = "https://simbad.cds.unistra.fr/simbad/sim-tap/sync";
const DEFAULT_RADIUS_ARCSEC = 2;
const MAX_RADIUS_ARCSEC = 10;
const TIMEOUT_MS = 15_000;
/** SIMBAD entries change rarely; cache for a day. */
const REVALIDATE_S = 86_400;

export interface SimbadMatch {
  /** SIMBAD main identifier, whitespace-collapsed ("M  31" → "M 31"). */
  mainId: string;
  separationArcsec: number;
  /** Set when SIMBAD quotes the measurement as a redshift. */
  redshift?: number;
  /** Set when SIMBAD quotes the measurement as a velocity (km/s). */
  radialVelocity?: number;
  /** SIMBAD quality grade of the redshift / velocity, A (best) to E. */
  rvzQuality?: string;
}

export interface SimbadResponse {
  found: boolean;
  match?: SimbadMatch;
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  const ra = Number(sp.get("ra"));
  const dec = Number(sp.get("dec"));
  const radius = Number(sp.get("radius") ?? DEFAULT_RADIUS_ARCSEC);

  if (
    !Number.isFinite(ra) ||
    !Number.isFinite(dec) ||
    ra < 0 ||
    ra >= 360 ||
    dec < -90 ||
    dec > 90
  ) {
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

  // ra/dec/radius are validated finite numbers, so interpolating them is safe.
  const point = `POINT('ICRS', ${ra}, ${dec})`;
  const query =
    `SELECT TOP 1 main_id, rvz_redshift, rvz_radvel, rvz_type, rvz_qual, ` +
    `DISTANCE(POINT('ICRS', ra, dec), ${point}) AS dist FROM basic ` +
    `WHERE CONTAINS(POINT('ICRS', ra, dec), CIRCLE('ICRS', ${ra}, ${dec}, ${radius / 3600})) = 1 ` +
    `ORDER BY dist`;
  const body = new URLSearchParams({
    request: "doQuery",
    lang: "adql",
    format: "json",
    query,
  });

  try {
    const resp = await fetch(`${TAP_URL}?${body}`, {
      signal: AbortSignal.timeout(TIMEOUT_MS),
      next: { revalidate: REVALIDATE_S },
    });
    if (!resp.ok) {
      return NextResponse.json(
        { error: `SIMBAD returned ${resp.status}` },
        { status: 502 }
      );
    }
    const json = (await resp.json()) as { data?: unknown[][] };
    return NextResponse.json(parseSimbadRow(json.data?.[0]));
  } catch (error) {
    console.error("SIMBAD proxy error:", error);
    return NextResponse.json(
      { error: "Failed to contact SIMBAD" },
      { status: 502 }
    );
  }
}

function parseSimbadRow(row: unknown[] | undefined): SimbadResponse {
  if (!row || typeof row[0] !== "string") return { found: false };
  const [mainId, redshift, radvel, type, qual, dist] = row;
  const num = (v: unknown) =>
    typeof v === "number" && Number.isFinite(v) ? v : undefined;
  return {
    found: true,
    match: {
      mainId: mainId.replace(/\s+/g, " ").trim(),
      separationArcsec: (num(dist) ?? 0) * 3600,
      redshift: type === "z" ? num(redshift) : undefined,
      radialVelocity: type === "v" ? num(radvel) : undefined,
      rvzQuality: typeof qual === "string" && qual ? qual : undefined,
    },
  };
}
