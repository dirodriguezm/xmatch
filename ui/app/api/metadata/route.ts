import { NextRequest, NextResponse } from "next/server";

import { ApiError, apiFetch } from "@/app/lib/api/client";
import type { components } from "@/types/xwave-api";

type Allwise = components["schemas"]["repository.Allwise"];

export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;

  const id = searchParams.get("id");
  const catalog = searchParams.get("catalog");

  if (!id || !catalog) {
    return NextResponse.json(
      { error: "Missing required parameters: id, catalog" },
      { status: 400 }
    );
  }

  // Forward the parameters to the external API
  const externalParams = new URLSearchParams({
    id,
    catalog,
  });

  try {
    const result = await apiFetch<Allwise>(`/metadata?${externalParams}`);
    return NextResponse.json(result ?? null);
  } catch (error) {
    console.error("Metadata proxy error:", error);
    if (error instanceof ApiError) {
      return NextResponse.json(
        { error: error.message },
        { status: error.status }
      );
    }
    return NextResponse.json(
      { error: "Failed to fetch from metadata service" },
      { status: 500 }
    );
  }
}
