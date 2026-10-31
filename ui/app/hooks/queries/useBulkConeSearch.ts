import { useMutation } from "@tanstack/react-query";

import {
  BULK_CHUNK_SIZE,
  type BulkApiMatch,
  type BulkResultRow,
  type BulkTarget,
  groupBulkResults,
} from "@/app/lib/utils/bulkInput";

export interface BulkConeSearchParams {
  targets: BulkTarget[];
  /** Arcsec. */
  radius: number;
  catalog: string;
  nneighbor: number;
}

export class BulkConeSearchError extends Error {
  constructor(
    message: string,
    public readonly status: number
  ) {
    super(message);
    this.name = "BulkConeSearchError";
  }
}

async function postChunk(
  targets: BulkTarget[],
  params: BulkConeSearchParams
): Promise<BulkApiMatch[]> {
  const response = await fetch("/api/bulk-conesearch", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      ra: targets.map((t) => t.ra),
      dec: targets.map((t) => t.dec),
      radius: params.radius,
      catalog: params.catalog,
      nneighbor: params.nneighbor,
    }),
  });
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new BulkConeSearchError(
      typeof errorData.error === "string"
        ? errorData.error
        : "Bulk cone search failed",
      response.status
    );
  }
  const result = await response.json();
  return Array.isArray(result) ? result : [];
}

/**
 * Run a bulk cross-match, splitting the list into BULK_CHUNK_SIZE requests.
 * Each chunk's `index` is relative to the chunk, so it is shifted back to the
 * target's position in the full list before grouping.
 */
export async function fetchBulkConeSearch(
  params: BulkConeSearchParams
): Promise<BulkResultRow[]> {
  const chunks: BulkTarget[][] = [];
  for (let i = 0; i < params.targets.length; i += BULK_CHUNK_SIZE) {
    chunks.push(params.targets.slice(i, i + BULK_CHUNK_SIZE));
  }
  const responses = await Promise.all(
    chunks.map((chunk) => postChunk(chunk, params))
  );
  const merged = responses.flatMap((matches, c) =>
    matches.map((m) => ({ ...m, index: m.index + c * BULK_CHUNK_SIZE }))
  );
  return groupBulkResults(params.targets, merged);
}

/** User-triggered bulk cross-match (a mutation: runs on "Cross-match", not on render). */
export function useBulkConeSearch() {
  return useMutation({
    mutationKey: ["bulk-conesearch"],
    mutationFn: fetchBulkConeSearch,
  });
}
