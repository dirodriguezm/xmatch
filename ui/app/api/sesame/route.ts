import { NextRequest, NextResponse } from "next/server";

/**
 * API route to proxy Sesame name resolution requests
 * This avoids CORS issues when calling the CDS Sesame service from the browser
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
  const sesameUrl = `https://cds.unistra.fr/cgi-bin/nph-sesame/-oI/SNV?${encodedName}`;

  try {
    const response = await fetch(sesameUrl);

    if (!response.ok) {
      return NextResponse.json(
        { error: "Sesame service unavailable" },
        { status: 502 }
      );
    }

    const text = await response.text();
    const result = parseSesameResponse(text, name);

    if (!result) {
      return NextResponse.json(
        { error: `Could not resolve "${name}"` },
        { status: 404 }
      );
    }

    return NextResponse.json(result);
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
        const ra = parseFloat(parts[1]);
        const dec = parseFloat(parts[2]);

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
