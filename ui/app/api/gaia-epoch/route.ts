import { NextRequest, NextResponse } from "next/server";

import { parseGaiaEpochCsv, parseGaiaTapCone } from "@/app/lib/utils/gaiaEpoch";
import type { DetectionPoint } from "@/app/lib/utils/lightcurve";

/**
 * Gaia DR3 epoch photometry from the ESA Gaia archive.
 *
 * - `?source_id=` fetches it directly (the caller already knows the source
 *   has epoch photometry, from the XWave counterpart record).
 * - `?ra=&dec=` first cone-searches gaiadr3.gaia_source (1″) for the nearest
 *   source; used when XWave has no Gaia counterpart for the position.
 */

const DATALINK_URL = "https://gea.esac.esa.int/data-server/data";
const TAP_URL = "https://gea.esac.esa.int/tap-server/tap/sync";
const CONE_RADIUS_DEG = 1 / 3600;
const TIMEOUT_MS = 25_000;
/** DR3 is frozen; cache for a week. */
const REVALIDATE_S = 7 * 86_400;

export interface GaiaEpochResponse {
  found: boolean;
  sourceId: string | null;
  points: DetectionPoint[];
  reason?: "no_gaia_source" | "no_epoch_photometry";
}

async function coneSearch(ra: number, dec: number) {
  const adql =
    `SELECT TOP 1 source_id, has_epoch_photometry, ` +
    `DISTANCE(POINT(ra, dec), POINT(${ra}, ${dec})) AS dist ` +
    `FROM gaiadr3.gaia_source ` +
    `WHERE 1 = CONTAINS(POINT(ra, dec), CIRCLE(${ra}, ${dec}, ${CONE_RADIUS_DEG})) ` +
    `ORDER BY dist`;
  const body = new URLSearchParams({
    REQUEST: "doQuery",
    LANG: "ADQL",
    FORMAT: "csv",
    QUERY: adql,
  });
  const resp = await fetch(TAP_URL, {
    method: "POST",
    body,
    signal: AbortSignal.timeout(TIMEOUT_MS),
    next: { revalidate: REVALIDATE_S },
  });
  if (!resp.ok) throw new Error(`Gaia TAP failed: ${resp.status}`);
  return parseGaiaTapCone(await resp.text());
}

async function epochPhotometry(sourceId: string) {
  const qs = new URLSearchParams({
    RETRIEVAL_TYPE: "EPOCH_PHOTOMETRY",
    ID: sourceId,
    DATA_STRUCTURE: "INDIVIDUAL",
    FORMAT: "CSV",
    VALID_DATA: "true",
  });
  const resp = await fetch(`${DATALINK_URL}?${qs}`, {
    signal: AbortSignal.timeout(TIMEOUT_MS),
    next: { revalidate: REVALIDATE_S },
  });
  if (!resp.ok) throw new Error(`Gaia DataLink failed: ${resp.status}`);
  // Sources without epoch photometry come back as an empty body.
  return parseGaiaEpochCsv(await resp.text());
}

function respond(result: GaiaEpochResponse) {
  return NextResponse.json(result, {
    headers: {
      "Cache-Control": `public, s-maxage=${REVALIDATE_S}, stale-while-revalidate=${REVALIDATE_S}`,
    },
  });
}

export async function GET(request: NextRequest) {
  const sp = request.nextUrl.searchParams;
  let sourceId = sp.get("source_id");

  try {
    if (sourceId == null) {
      const ra = Number(sp.get("ra"));
      const dec = Number(sp.get("dec"));
      if (
        !sp.has("ra") ||
        !sp.has("dec") ||
        !Number.isFinite(ra) ||
        !Number.isFinite(dec)
      ) {
        return NextResponse.json(
          { error: "Pass 'source_id', or 'ra' and 'dec'" },
          { status: 400 }
        );
      }
      const nearest = await coneSearch(ra, dec);
      if (!nearest) {
        return respond({
          found: false,
          sourceId: null,
          points: [],
          reason: "no_gaia_source",
        });
      }
      if (!nearest.hasEpochPhotometry) {
        return respond({
          found: false,
          sourceId: nearest.sourceId,
          points: [],
          reason: "no_epoch_photometry",
        });
      }
      sourceId = nearest.sourceId;
    }

    if (!/^\d{1,20}$/.test(sourceId)) {
      return NextResponse.json(
        { error: "Invalid 'source_id'" },
        { status: 400 }
      );
    }

    const points = await epochPhotometry(sourceId);
    return respond({
      found: points.length > 0,
      sourceId,
      points,
      ...(points.length === 0
        ? { reason: "no_epoch_photometry" as const }
        : {}),
    });
  } catch (err) {
    const timedOut = err instanceof DOMException && err.name === "TimeoutError";
    console.error("Gaia epoch photometry error:", err);
    return NextResponse.json(
      {
        error: timedOut
          ? "Gaia archive timed out"
          : "Gaia archive request failed",
      },
      { status: timedOut ? 504 : 502 }
    );
  }
}
