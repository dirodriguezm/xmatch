import { useQueries } from "@tanstack/react-query";

import {
  CATALOG_OPTIONS,
  type CatalogOption,
} from "@/app/lib/constants/catalogs";

import { ConeSearchError, fetchConeSearch } from "./useConeSearch";

/**
 * Positional uncertainty budget per catalog, in arcsec. A counterpart search
 * uses the larger of the two catalogs involved, so Gaia→AllWISE is limited by
 * the ~6″ WISE PSF and anything involving eROSITA by the X-ray positions.
 */
export const COUNTERPART_RADIUS_ARCSEC: Record<CatalogOption, number> = {
  gaia: 1,
  allwise: 3,
  erosita: 10,
};

export interface CounterpartsParams {
  ra: number;
  dec: number;
  /** Catalog the object itself comes from; it is not searched again. */
  sourceCatalog: string;
}

export interface Counterpart {
  catalog: CatalogOption;
  radiusArcsec: number;
  status: "loading" | "found" | "none" | "error";
  id?: string;
  separationArcsec?: number;
  record?: Record<string, unknown>;
  error?: Error;
}

/** Nearest source within the radius, or null. Separation comes back in arcsec. */
async function fetchNearest(
  ra: number,
  dec: number,
  catalog: CatalogOption,
  radiusArcsec: number
): Promise<Record<string, unknown> | null> {
  const groups = await fetchConeSearch({
    ra,
    dec,
    radius: radiusArcsec,
    catalog,
    nneighbor: 1,
    getMetadata: true,
  });
  const rows = groups.flatMap(
    (g) => (g.data ?? []) as unknown as Record<string, unknown>[]
  );
  if (rows.length === 0) return null;
  const distance = (r: Record<string, unknown>) =>
    typeof r.distance === "number" ? r.distance : Infinity;
  return rows.reduce((best, r) => (distance(r) < distance(best) ? r : best));
}

/**
 * Nearest counterpart of an object in every other search catalog. Each catalog
 * resolves independently, so a failing one never hides the others.
 */
export function useCounterparts(params: CounterpartsParams | null) {
  const source = params?.sourceCatalog.toLowerCase();
  const targets = CATALOG_OPTIONS.filter((c) => c !== source);
  const sourceRadius = COUNTERPART_RADIUS_ARCSEC[source as CatalogOption] ?? 0;

  const queries = useQueries({
    queries: targets.map((catalog) => {
      const radiusArcsec = Math.max(
        sourceRadius,
        COUNTERPART_RADIUS_ARCSEC[catalog]
      );
      return {
        queryKey: [
          "counterpart",
          params?.ra,
          params?.dec,
          catalog,
          radiusArcsec,
        ],
        queryFn: () =>
          fetchNearest(params!.ra, params!.dec, catalog, radiusArcsec),
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
      };
    }),
  });

  return targets.map((catalog, i): Counterpart => {
    const q = queries[i];
    const radiusArcsec = Math.max(
      sourceRadius,
      COUNTERPART_RADIUS_ARCSEC[catalog]
    );
    if (q.isPending) return { catalog, radiusArcsec, status: "loading" };
    if (q.isError)
      return { catalog, radiusArcsec, status: "error", error: q.error };
    const record = q.data;
    if (!record) return { catalog, radiusArcsec, status: "none" };
    return {
      catalog,
      radiusArcsec,
      status: "found",
      id: typeof record.id === "string" ? record.id : undefined,
      separationArcsec:
        typeof record.distance === "number" ? record.distance : undefined,
      record,
    };
  });
}
