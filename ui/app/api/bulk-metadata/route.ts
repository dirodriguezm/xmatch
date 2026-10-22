import { NextRequest, NextResponse } from "next/server";

import type { MetadataRecord } from "@/app/hooks/queries/useBulkMetadata";
import { ApiError, apiFetch } from "@/app/lib/api/client";

/**
 * Proxy for `POST /bulk-metadata`.
 *
 * Note the upstream path. The *deployed* swagger documents this operation as
 * `/metadata/bulk`, which 404s; `/bulk-metadata` is the route that actually
 * answers, and is what `service/docs/swagger.json` in this repo declares.
 */
const MAX_IDS = 500;

export async function POST(request: NextRequest) {
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return NextResponse.json({ error: "Invalid JSON body" }, { status: 400 });
  }

  const { catalog, ids } = (body ?? {}) as { catalog?: unknown; ids?: unknown };

  if (typeof catalog !== "string" || catalog.length === 0) {
    return NextResponse.json(
      { error: "Missing required field: catalog" },
      { status: 400 }
    );
  }

  if (!Array.isArray(ids) || ids.some((id) => typeof id !== "string")) {
    return NextResponse.json(
      { error: "Missing or invalid field: ids (expected string[])" },
      { status: 400 }
    );
  }

  if (ids.length === 0) {
    return NextResponse.json([]);
  }

  if (ids.length > MAX_IDS) {
    return NextResponse.json(
      { error: `Too many ids: ${ids.length} (max ${MAX_IDS})` },
      { status: 400 }
    );
  }

  try {
    const result = await apiFetch<MetadataRecord[]>("/bulk-metadata", {
      method: "POST",
      body: JSON.stringify({ catalog, ids }),
    });
    // A 204 makes apiFetch return null. Normalize to [] so the client never has
    // to tell "no metadata for these ids" apart from "the request failed".
    return NextResponse.json(Array.isArray(result) ? result : []);
  } catch (error) {
    console.error("Bulk metadata proxy error:", error);
    if (error instanceof ApiError) {
      return NextResponse.json(
        { error: error.message },
        { status: error.status }
      );
    }
    return NextResponse.json(
      { error: "Failed to fetch from bulk metadata service" },
      { status: 500 }
    );
  }
}
