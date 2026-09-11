import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import type { CatalogGroup } from "@/app/hooks/queries/useConeSearch";

/**
 * Map API CatalogGroup response to frontend CrossmatchResult format
 *
 * `catalogSlug` is the slug the request was made with. It is threaded through
 * as a join key for photometry enrichment because `catalog` below is an
 * upstream-controlled value — the same endpoint returns "AllWISE" instead of
 * "allwise" when `getMetadata=true`, so it is not a safe key on its own.
 */
export function mapConeSearchResults(
  results: CatalogGroup[],
  catalogSlug?: string
): CrossmatchResult[] {
  return results.flatMap((group) =>
    group.data.map((item, index) => ({
      key: item.id ?? `${group.catalog}-${index}`,
      objectId: item.id ?? `unknown-${group.catalog}-${index}`,
      ra: item.ra ?? 0,
      dec: item.dec ?? 0,
      angularDistance: item.distance ?? 0,
      // Unchanged: row-click navigation and the Tag color both read this.
      catalog: item.cat ?? group.catalog ?? "Unknown",
      ipix: item.ipix,
      catalogSlug: (catalogSlug ?? group.catalog ?? "").toLowerCase(),
    }))
  );
}
