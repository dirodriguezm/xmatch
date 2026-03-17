import { NextRequest, NextResponse } from "next/server";

import { ApiError, apiFetch } from "@/app/lib/api/client";
import type { components } from "@/types/xwave-api";

type Mastercat = components["schemas"]["repository.Mastercat"];

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;

  const ra = searchParams.get("ra");
  const dec = searchParams.get("dec");
  const radius = searchParams.get("radius");

  if (!ra || !dec || !radius) {
    return NextResponse.json(
      { error: "Missing required parameters: ra, dec, radius" },
      { status: 400 }
    );
  }

  // Forward the parameters to the external API
  const externalParams = new URLSearchParams({
    ra,
    dec,
    radius,
  });

  const catalog = searchParams.get("catalog");
  if (catalog) externalParams.set("catalog", catalog);

  const nneighbor = searchParams.get("nneighbor");
  if (nneighbor) externalParams.set("nneighbor", nneighbor);

  const getMetadata = searchParams.get("getMetadata");
  if (getMetadata) externalParams.set("getMetadata", getMetadata);

  try {
    const result = await apiFetch<Mastercat[]>(`/conesearch?${externalParams}`);
    return NextResponse.json(result ?? []);
  } catch (error) {
    console.error("Conesearch proxy error:", error);
    if (error instanceof ApiError) {
      return NextResponse.json(
        { error: error.message },
        { status: error.status }
      );
    }
    return NextResponse.json(
      { error: "Failed to fetch from conesearch service" },
      { status: 500 }
    );
  }
}
