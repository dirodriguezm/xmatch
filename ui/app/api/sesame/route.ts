import { NextRequest, NextResponse } from "next/server";

import {
  fetchFromMirrors,
  SESAME_MIRRORS,
  UPSTREAM_HOST_HEADER,
} from "@/app/lib/utils/mirrors";

/** Per mirror; Sesame answers in 1–3 s when it is up. */
const TIMEOUT_MS = 8_000;

/**
 * API route to proxy Sesame name resolution requests, trying the CDS mirrors
 * in turn. This avoids CORS issues when calling Sesame from the browser.
 */
export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;
  const name = searchParams.get("name");

  if (!name) {
    return NextResponse.json(
      { error: "Missing 'name' parameter" },
      { status: 400 }
    );
  }

  const encodedName = encodeURIComponent(name.trim());
  try {
    const { response, host } = await fetchFromMirrors(
      SESAME_MIRRORS,
      `/-oI/SNV?${encodedName}`,
      {},
      TIMEOUT_MS
    );

    if (!response.ok) {
      return NextResponse.json(
        { error: "Sesame service unavailable" },
        { status: 502 }
      );
    }

    const text = await response.text();
    const result = parseSesameResponse(text, name);

    const headers = { [UPSTREAM_HOST_HEADER]: host };
    if (!result) {
      return NextResponse.json(
        { error: `Could not resolve "${name}"` },
        { status: 404, headers }
      );
    }

    return NextResponse.json(result, { headers });
  } catch (error) {
    console.error("Sesame proxy error:", error);
    return NextResponse.json(
      { error: "Failed to contact Sesame service" },
      { status: 500 }
    );
  }
}

/**
 * Parse Sesame response text to extract coordinates
 * Format: %J ra dec (in degrees)
 */
function parseSesameResponse(
  text: string,
  originalName: string
): { ra: number; dec: number; objectName: string } | null {
  const lines = text.split("\n");

  for (const line of lines) {
    if (line.startsWith("%J")) {
      const parts = line.trim().split(/\s+/);
      if (parts.length >= 3) {
        const ra = parseFloat(parts[1] ?? "");
        const dec = parseFloat(parts[2] ?? "");

        if (!isNaN(ra) && !isNaN(dec)) {
          return {
            ra,
            dec,
            objectName: originalName,
          };
        }
      }
    }
  }

  return null;
}
