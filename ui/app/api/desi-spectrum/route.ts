import { NextRequest, NextResponse } from "next/server";

/**
 * Server-side proxy to NOIRLab's SPARCL service. Resolves a DESI DR1 spectrum
 * for a given TARGETID and returns its full-resolution wavelength / flux arrays.
 *
 * Two-step REST flow (no auth required):
 *   1. POST /api/find/      — TARGETID → sparcl_id (uuid) + redshift / spectype.
 *   2. POST /api/spectras/  — sparcl_id → wavelength / flux arrays.
 *
 * Gotchas baked in below:
 * - TARGETID is a 64-bit BIGINT that overflows JS numeric precision, so we keep
 *   it as a digit string and splice it into the find body as a raw numeric
 *   literal (never via Number()/JSON.stringify of the number).
 * - `/api/spectras/` wants a BARE JSON array body (`["uuid"]`), not an object,
 *   and `format=json` must be set or the server returns Python pickle.
 * - Both find and spectras responses are lists whose element [0] is a
 *   META/HEADER object; the actual records follow at [1..].
 */

const SPARCL_BASE = "https://astrosparcl.datalab.noirlab.edu";
const SPARCL_HEADERS = {
  "Content-Type": "application/json",
  "User-Agent": "xmatch-ui/1.0 (DESI spectrum lookup)",
};
const SPARCL_TIMEOUT_MS = 15000;

// POST to SPARCL with a hard timeout and a single retry on a transient failure
// (network error / abort / 5xx). Keeps normal request volume to ≤2 per object
// while riding out brief upstream blips; a thrown error bubbles to the 500 path.
async function sparclPost(url: string, body: string): Promise<Response> {
  let lastErr: unknown;
  for (let attempt = 0; attempt <= 1; attempt++) {
    if (attempt > 0) await new Promise((r) => setTimeout(r, 600));
    try {
      const resp = await fetch(url, {
        method: "POST",
        headers: SPARCL_HEADERS,
        body,
        signal: AbortSignal.timeout(SPARCL_TIMEOUT_MS),
      });
      if (resp.status >= 500 && attempt === 0) continue; // retry once on 5xx
      return resp;
    } catch (err) {
      lastErr = err; // network error / timeout — retry once, then rethrow
    }
  }
  throw lastErr;
}

interface SparclRecord {
  sparcl_id?: string;
  redshift?: number;
  spectype?: string;
  wavelength?: number[];
  flux?: number[];
  model?: number[];
  ivar?: number[];
  [key: string]: unknown;
}

// Drop the leading META object and return the data records.
function dataRecords(body: unknown): SparclRecord[] {
  if (!Array.isArray(body)) return [];
  return body.slice(1) as SparclRecord[];
}

export async function GET(request: NextRequest) {
  const targetid = request.nextUrl.searchParams.get("targetid") ?? "";

  // Keep TARGETID as a string; only accept plain digits so we can safely splice
  // it into the find body as a numeric literal without precision loss.
  if (!/^\d+$/.test(targetid)) {
    return NextResponse.json(
      { error: "Missing or invalid 'targetid' (must be a positive integer)" },
      { status: 400 }
    );
  }

  try {
    // Step 1: TARGETID → sparcl_id, constrained to DESI DR1.
    const findBody =
      `{"outfields":["sparcl_id","targetid","data_release","redshift","spectype"],` +
      `"search":[["targetid",${targetid}],["data_release","DESI-DR1"]]}`;
    const findResp = await sparclPost(
      `${SPARCL_BASE}/api/find/?limit=5`,
      findBody
    );
    if (!findResp.ok) {
      const detail = await findResp.text();
      return NextResponse.json(
        {
          error: "SPARCL find failed",
          status: findResp.status,
          detail: detail.slice(0, 500),
        },
        { status: 502 }
      );
    }
    const found = dataRecords(await findResp.json());
    const match = found.find((r) => typeof r.sparcl_id === "string");
    if (!match?.sparcl_id) {
      return NextResponse.json(
        { found: false },
        {
          headers: {
            "Cache-Control":
              "public, s-maxage=86400, stale-while-revalidate=86400",
          },
        }
      );
    }
    const sparclId = match.sparcl_id;

    // Step 2: sparcl_id → wavelength / flux arrays (full resolution).
    // `model` is the redrock best-fit, `ivar` the inverse variance (→ noise).
    const include = "wavelength,flux,ivar,model,redshift,spectype";
    const specResp = await sparclPost(
      `${SPARCL_BASE}/api/spectras/?format=json&include=${include}`,
      JSON.stringify([sparclId])
    );
    if (!specResp.ok) {
      const detail = await specResp.text();
      return NextResponse.json(
        {
          error: "SPARCL spectras failed",
          status: specResp.status,
          detail: detail.slice(0, 500),
        },
        { status: 502 }
      );
    }
    const spec = dataRecords(await specResp.json())[0];
    if (!spec || !Array.isArray(spec.wavelength) || !Array.isArray(spec.flux)) {
      return NextResponse.json(
        { error: "Unexpected SPARCL spectras response shape" },
        { status: 502 }
      );
    }

    return NextResponse.json(
      {
        found: true,
        sparclId,
        targetid,
        spectype: spec.spectype ?? match.spectype ?? null,
        redshift: spec.redshift ?? match.redshift ?? null,
        wavelength: spec.wavelength,
        flux: spec.flux,
        // Optional overlays — may be absent on edge records; chart degrades gracefully.
        model: Array.isArray(spec.model) ? spec.model : null,
        ivar: Array.isArray(spec.ivar) ? spec.ivar : null,
      },
      {
        headers: {
          "Cache-Control":
            "public, s-maxage=86400, stale-while-revalidate=86400",
        },
      }
    );
  } catch (err) {
    console.error("DESI spectrum lookup error:", err);
    return NextResponse.json(
      { error: "Failed to contact NOIRLab SPARCL" },
      { status: 500 }
    );
  }
}
