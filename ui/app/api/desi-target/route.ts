import { NextRequest, NextResponse } from "next/server";

/**
 * Server-side proxy to NOIRLab Astro Data Lab's TAP service. Cone-searches
 * `desi_dr1.zpix` for the closest DESI DR1 spectrum to (RA, Dec) and returns
 * its TARGETID (or null if none in the cone).
 *
 * - TAP `POINT/CIRCLE` rejected by the underlying SQL translator, so we use
 *   a small RA/Dec bounding box and refine by angular distance in JS.
 * - Position columns are `mean_fiber_ra` / `mean_fiber_dec` (not `ra`/`dec`).
 * - TARGETID is kept as a string — it's a 64-bit BIGINT that overflows JS
 *   numeric precision.
 */

const TAP_SYNC_URL = "https://datalab.noirlab.edu/tap/sync";
const DEFAULT_RADIUS_ARCSEC = 1.5;

function arcsecBetween(
  ra1: number,
  dec1: number,
  ra2: number,
  dec2: number
): number {
  const toRad = (d: number) => (d * Math.PI) / 180;
  const sinDec = Math.sin(toRad(dec1)) * Math.sin(toRad(dec2));
  const cosDec =
    Math.cos(toRad(dec1)) * Math.cos(toRad(dec2)) * Math.cos(toRad(ra1 - ra2));
  const cosTheta = Math.min(1, Math.max(-1, sinDec + cosDec));
  return (Math.acos(cosTheta) * 180 * 3600) / Math.PI;
}

// Parse a TAP CSV body: returns { header: string[], rows: string[][] }.
// The body may also be a VOTable XML error envelope; in that case we return
// the raw text in `error` so callers can surface it.
function parseTapCsv(
  text: string
): { rows: string[][]; header: string[] } | { error: string } {
  if (text.startsWith("<?xml") || text.startsWith("<VOTABLE")) {
    const errMatch = text.match(/QUERY_STATUS"[^>]*>([\s\S]*?)<\/INFO>/);
    return { error: errMatch ? errMatch[1].trim() : text.slice(0, 500) };
  }
  const lines = text.split("\n").filter((l) => l.length > 0);
  if (lines.length === 0) return { header: [], rows: [] };
  const header = lines[0].split(",");
  const rows = lines.slice(1).map((l) => l.split(","));
  return { header, rows };
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
  if (!Number.isFinite(radius) || radius <= 0 || radius > 60) {
    return NextResponse.json(
      { error: "'radius' must be a positive number ≤ 60 arcsec" },
      { status: 400 }
    );
  }

  const ddec = radius / 3600;
  const cosDec = Math.max(Math.cos((dec * Math.PI) / 180), 1e-6);
  const dra = ddec / cosDec;

  const adql =
    `SELECT TOP 20 targetid, mean_fiber_ra, mean_fiber_dec ` +
    `FROM desi_dr1.zpix ` +
    `WHERE mean_fiber_ra BETWEEN ${ra - dra} AND ${ra + dra} ` +
    `AND mean_fiber_dec BETWEEN ${dec - ddec} AND ${dec + ddec}`;
  const url =
    `${TAP_SYNC_URL}?REQUEST=doQuery&LANG=ADQL&FORMAT=csv` +
    `&QUERY=${encodeURIComponent(adql)}`;

  try {
    const tapResp = await fetch(url);
    const body = await tapResp.text();

    if (!tapResp.ok) {
      return NextResponse.json(
        {
          error: "TAP request failed",
          status: tapResp.status,
          detail: body.slice(0, 500),
        },
        { status: 502 }
      );
    }

    const parsed = parseTapCsv(body);
    if ("error" in parsed) {
      return NextResponse.json(
        { error: "TAP query error", detail: parsed.error },
        { status: 502 }
      );
    }

    const iTid = parsed.header.indexOf("targetid");
    const iRa = parsed.header.indexOf("mean_fiber_ra");
    const iDec = parsed.header.indexOf("mean_fiber_dec");
    if (iTid < 0 || iRa < 0 || iDec < 0) {
      return NextResponse.json(
        { error: "Unexpected TAP response shape", header: parsed.header },
        { status: 502 }
      );
    }

    let best: { targetid: string; separationArcsec: number } | null = null;
    for (const row of parsed.rows) {
      const tid = row[iTid];
      const recRa = Number(row[iRa]);
      const recDec = Number(row[iDec]);
      if (!tid || !Number.isFinite(recRa) || !Number.isFinite(recDec)) continue;
      const sep = arcsecBetween(ra, dec, recRa, recDec);
      if (sep <= radius && (!best || sep < best.separationArcsec)) {
        best = { targetid: tid, separationArcsec: sep };
      }
    }

    return NextResponse.json(best ?? { targetid: null }, {
      headers: {
        "Cache-Control": "public, s-maxage=86400, stale-while-revalidate=86400",
      },
    });
  } catch (err) {
    console.error("DESI target lookup error:", err);
    return NextResponse.json(
      { error: "Failed to contact NOIRLab TAP" },
      { status: 500 }
    );
  }
}
