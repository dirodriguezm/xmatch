import type {
  CrossmatchResult,
  EnrichedResult,
} from "@/app/components/results/ResultsTable";
import type {
  BulkMetadataGroup,
  BulkMetadataRequest,
} from "@/app/hooks/queries/useBulkMetadata";
import {
  CATALOG_MAG_SPECS,
  type MagValue,
  pickMag,
} from "@/app/lib/constants/photometry";

/**
 * Gaia designations carry internal whitespace ("Gaia DR3 2425866213898295424"),
 * so collapse whitespace runs and case for the fallback match.
 */
function normalizeId(id: string): string {
  return id.trim().replace(/\s+/g, " ").toLowerCase();
}

/** NUL separates the parts: it appears in neither a slug nor an id. */
function indexKey(catalogSlug: string, id: string): string {
  return `${catalogSlug.toLowerCase()}\u0000${id}`;
}

/**
 * Build a `${catalog}\0${id}` → MagValue lookup from the per-catalog responses.
 *
 * The service does NOT return records in the order the ids were sent, so this
 * is a keyed join and must never become a positional zip. Each record is
 * indexed twice — verbatim and normalized — so an exact match always wins while
 * near-misses still resolve.
 */
export function buildMagIndex(
  groups: BulkMetadataGroup[]
): Map<string, MagValue> {
  const index = new Map<string, MagValue>();

  for (const group of groups) {
    for (const record of group.records) {
      const id = record.id;
      if (typeof id !== "string" || id.length === 0) continue;

      const mag = pickMag(group.catalog, record);
      if (!mag) continue;

      index.set(indexKey(group.catalog, id), mag);
      const normalized = indexKey(group.catalog, normalizeId(id));
      if (!index.has(normalized)) index.set(normalized, mag);
    }
  }

  return index;
}

/**
 * Attach photometry to search results.
 *
 * Non-destructive by construction: objectId, catalog, ra, dec, ipix and key are
 * passed through untouched, so the existing columns and the row-click
 * navigation contract are unaffected. Rows with no match keep
 * `photometry: undefined` and render "—".
 */
export function enrichWithPhotometry(
  rows: CrossmatchResult[],
  index: Map<string, MagValue>
): EnrichedResult[] {
  if (index.size === 0) return rows;

  return rows.map((row) => {
    const slug = row.catalogSlug ?? row.catalog;
    const mag =
      index.get(indexKey(slug, row.objectId)) ??
      index.get(indexKey(slug, normalizeId(row.objectId)));
    return mag ? { ...row, photometry: mag } : row;
  });
}

/**
 * One request per catalog, ids de-duplicated.
 *
 * Catalogs with no magnitude spec (eROSITA) are skipped, so we never spend a
 * round trip fetching data the table has no way to render.
 */
export function buildBulkMetadataRequests(
  rows: CrossmatchResult[]
): BulkMetadataRequest[] {
  const byCatalog = new Map<string, Set<string>>();

  for (const row of rows) {
    const slug = (row.catalogSlug ?? row.catalog).toLowerCase();
    if (!CATALOG_MAG_SPECS[slug]?.length) continue;
    // Synthetic placeholder from mapConeSearchResults when the API omitted an id.
    if (!row.objectId || row.objectId.startsWith("unknown-")) continue;

    let ids = byCatalog.get(slug);
    if (!ids) {
      ids = new Set();
      byCatalog.set(slug, ids);
    }
    ids.add(row.objectId);
  }

  return (
    [...byCatalog.entries()]
      .map(([catalog, ids]) => ({ catalog, ids: [...ids] }))
      // Deterministic order keeps the useQueries array stable across renders.
      .sort((a, b) => a.catalog.localeCompare(b.catalog))
  );
}
