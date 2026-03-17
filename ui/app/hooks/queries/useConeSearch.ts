import { useQueries, useQuery } from "@tanstack/react-query";

import {
  type CatalogRadiusConfig,
  convertRadiusToDegrees,
} from "@/app/lib/constants/search";
import type { components } from "@/types/xwave-api";

type Mastercat = components["schemas"]["repository.Mastercat"];

export interface CatalogGroupItem {
  id?: string;
  ipix?: number;
  ra?: number;
  dec?: number;
  cat?: string;
  distance?: number;
}

export interface CatalogGroup {
  catalog: string;
  data: CatalogGroupItem[];
}

export interface ConeSearchParams {
  ra: number;
  dec: number;
  radius: number;
  catalog?: string;
  nneighbor?: number;
  getMetadata?: boolean;
}

class ConeSearchError extends Error {
  constructor(
    message: string,
    public readonly status: number
  ) {
    super(message);
    this.name = "ConeSearchError";
  }
}

async function fetchConeSearch(
  params: ConeSearchParams
): Promise<CatalogGroup[]> {
  const searchParams = new URLSearchParams({
    ra: params.ra.toString(),
    dec: params.dec.toString(),
    radius: params.radius.toString(),
  });

  if (params.catalog) {
    searchParams.set("catalog", params.catalog);
  }
  if (params.nneighbor !== undefined) {
    searchParams.set("nneighbor", params.nneighbor.toString());
  }
  if (params.getMetadata !== undefined) {
    searchParams.set("getMetadata", params.getMetadata.toString());
  }

  const response = await fetch(`/api/conesearch?${searchParams}`);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new ConeSearchError(
      errorData.error || "Failed to fetch cone search results",
      response.status
    );
  }

  const result = await response.json();
  return result ?? [];
}

export function useConeSearch(params: ConeSearchParams | null) {
  return useQuery({
    queryKey: ["conesearch", params],
    queryFn: () => fetchConeSearch(params!),
    enabled: params !== null,
  });
}

export function useParallelConeSearch(
  base: { ra: number; dec: number } | null,
  configs: CatalogRadiusConfig[]
) {
  return useQueries({
    queries: configs
      .filter((c) => c.enabled)
      .map((c) => ({
        queryKey: [
          "conesearch",
          base?.ra,
          base?.dec,
          c.catalog,
          c.radius,
          c.unit,
        ],
        queryFn: () =>
          fetchConeSearch({
            ra: base!.ra,
            dec: base!.dec,
            radius: convertRadiusToDegrees(c.radius, c.unit),
            catalog: c.catalog,
          }),
        enabled: base !== null,
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
  });
}
