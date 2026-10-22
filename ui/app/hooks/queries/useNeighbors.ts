import { useQueries, type UseQueryResult } from "@tanstack/react-query";
import { useCallback } from "react";

import {
  CATALOG_OPTIONS,
  type CatalogOption,
} from "@/app/lib/constants/catalogs";
import { positionAngle } from "@/app/lib/utils/coordinates";

import {
  type CatalogGroup,
  ConeSearchError,
  fetchConeSearch,
} from "./useConeSearch";

/** Radius of the "Nearby sources" search around an object, in arcsec. */
export const NEIGHBOR_RADIUS_ARCSEC = 30;
/** Per-catalog cap; dense fields near the plane can hold far more. */
const NEIGHBOR_MAX_PER_CATALOG = 50;

export interface NeighborsParams {
  ra: number;
  dec: number;
  /** The object itself, left out of the results. */
  self: { id: string; catalog: string };
}

export interface Neighbor {
  catalog: CatalogOption;
  id: string;
  ra: number;
  dec: number;
  separationArcsec: number;
  /** Degrees east of north, as seen from the object. */
  positionAngle: number;
}

export interface NeighborsResult {
  neighbors: Neighbor[];
  isLoading: boolean;
  /** Catalogs whose search failed; the others are still listed. */
  failedCatalogs: CatalogOption[];
  /** Catalogs that hit the per-catalog cap, so may have more sources. */
  truncatedCatalogs: CatalogOption[];
}

/**
 * Every indexed source within NEIGHBOR_RADIUS_ARCSEC of a position, across all
 * search catalogs, nearest first. Each catalog is queried independently so one
 * failing never hides the others.
 */
export function useNeighbors(params: NeighborsParams | null): NeighborsResult {
  const ra = params?.ra;
  const dec = params?.dec;
  const selfId = params?.self.id;
  const selfCatalog = params?.self.catalog.toLowerCase();

  // Combining inside useQueries (with a stable callback) memoises the result,
  // so `neighbors` keeps its identity until the data changes — the sky view
  // re-draws its overlay only then.
  const combine = useCallback(
    (queries: UseQueryResult<CatalogGroup[]>[]): NeighborsResult => {
      const neighbors: Neighbor[] = [];
      const failedCatalogs: CatalogOption[] = [];
      const truncatedCatalogs: CatalogOption[] = [];

      CATALOG_OPTIONS.forEach((catalog, i) => {
        const q = queries[i];
        if (q.isError) failedCatalogs.push(catalog);
        if (!q.data || ra === undefined || dec === undefined) return;
        const rows = q.data.flatMap((g) => g.data ?? []);
        if (rows.length >= NEIGHBOR_MAX_PER_CATALOG)
          truncatedCatalogs.push(catalog);
        for (const r of rows) {
          if (!r.id || typeof r.ra !== "number" || typeof r.dec !== "number")
            continue;
          if (catalog === selfCatalog && r.id === selfId) continue;
          neighbors.push({
            catalog,
            id: r.id,
            ra: r.ra,
            dec: r.dec,
            separationArcsec: r.distance ?? 0,
            positionAngle: positionAngle(ra, dec, r.ra, r.dec),
          });
        }
      });
      neighbors.sort((a, b) => a.separationArcsec - b.separationArcsec);

      return {
        neighbors,
        isLoading: queries.some((q) => q.isPending && q.fetchStatus !== "idle"),
        failedCatalogs,
        truncatedCatalogs,
      };
    },
    [ra, dec, selfId, selfCatalog]
  );

  return useQueries({
    queries: CATALOG_OPTIONS.map((catalog) => ({
      queryKey: [
        "neighbors",
        params?.ra,
        params?.dec,
        catalog,
        NEIGHBOR_RADIUS_ARCSEC,
      ],
      queryFn: () =>
        fetchConeSearch({
          ra: params!.ra,
          dec: params!.dec,
          radius: NEIGHBOR_RADIUS_ARCSEC,
          catalog,
          nneighbor: NEIGHBOR_MAX_PER_CATALOG,
        }),
      enabled: params !== null,
      staleTime: Infinity,
      retry: (failureCount: number, error: unknown) => {
        if (
          error instanceof ConeSearchError &&
          error.status >= 400 &&
          error.status < 500
        )
          return false;
        return failureCount < 2;
      },
    })),
    combine,
  });
}
