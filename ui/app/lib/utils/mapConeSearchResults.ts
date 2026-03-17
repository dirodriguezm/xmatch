import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import type { CatalogGroup } from "@/app/hooks/queries/useConeSearch";

/**
 * Map API CatalogGroup response to frontend CrossmatchResult format
 */
export function mapConeSearchResults(
  results: CatalogGroup[]
): CrossmatchResult[] {
  return results.flatMap((group) =>
    group.data.map((item, index) => ({
      key: item.id ?? `${group.catalog}-${index}`,
      objectId: item.id ?? `unknown-${group.catalog}-${index}`,
      ra: item.ra ?? 0,
      dec: item.dec ?? 0,
      angularDistance: item.distance ?? 0,
      catalog: item.cat ?? group.catalog ?? "Unknown",
      ipix: item.ipix,
    }))
  );
}
