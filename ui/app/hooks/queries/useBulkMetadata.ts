import { useQueries } from "@tanstack/react-query";

/**
 * Metadata as an open bag.
 *
 * The generated `repository.Allwise` is the only metadata schema the service
 * publishes, yet every catalog answers with a different field set (gaia →
 * `phot_*_mean_mag`, erosita → unverified). Rather than cast that type onto
 * payloads it does not describe, model what is actually guaranteed — an id and
 * a position — and funnel every other field access through `pickMag`.
 */
export interface MetadataRecord extends Record<string, unknown> {
  id?: string;
  ra?: number;
  dec?: number;
}

export interface BulkMetadataRequest {
  catalog: string;
  ids: string[];
}

export interface BulkMetadataGroup {
  catalog: string;
  records: MetadataRecord[];
}

class BulkMetadataError extends Error {
  constructor(
    message: string,
    public readonly status: number
  ) {
    super(message);
    this.name = "BulkMetadataError";
  }
}

async function fetchBulkMetadata(
  req: BulkMetadataRequest
): Promise<MetadataRecord[]> {
  const response = await fetch("/api/bulk-metadata", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new BulkMetadataError(
      errorData.error || "Failed to fetch bulk metadata",
      response.status
    );
  }

  const result = await response.json();
  return Array.isArray(result) ? result : [];
}

/**
 * One POST per catalog, in parallel, to enrich already-rendered search results.
 *
 * ENRICHMENT ONLY: callers must not fold this into the page's loading or error
 * state. A failure here leaves the Mag column empty; it must never blank the
 * results table.
 */
export function useBulkMetadata(requests: BulkMetadataRequest[]) {
  const results = useQueries({
    queries: requests.map(({ catalog, ids }) => {
      // Sorted so the cache key does not depend on cone-search ordering.
      const sortedIds = [...ids].sort();
      return {
        queryKey: ["bulk-metadata", catalog, sortedIds.join("|")],
        queryFn: () => fetchBulkMetadata({ catalog, ids: sortedIds }),
        enabled: sortedIds.length > 0,
        // Catalog photometry is immutable; the app-wide 60s staleTime would
        // refetch it for no reason.
        staleTime: Infinity,
        gcTime: 30 * 60 * 1000,
        retry: (failureCount: number, error: unknown) => {
          if (
            error instanceof BulkMetadataError &&
            error.status >= 400 &&
            error.status < 500
          )
            return false;
          return failureCount < 1;
        },
      };
    }),
  });

  // results[i] pairs with requests[i].
  return {
    groups: results.map((r, i) => ({
      catalog: requests[i].catalog,
      records: r.data ?? [],
    })) as BulkMetadataGroup[],
    isFetching: results.some((r) => r.isFetching),
    isError: results.some((r) => r.isError),
  };
}
