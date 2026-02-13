import { useQuery } from "@tanstack/react-query";

import { apiFetch } from "@/app/lib/api/client";
import type { components } from "@/types/xwave-api";

type Mastercat = components["schemas"]["repository.Mastercat"];

export interface ConeSearchParams {
  ra: number;
  dec: number;
  radius: number;
  catalog?: string;
  nneighbor?: number;
  getMetadata?: boolean;
}

async function fetchConeSearch(params: ConeSearchParams): Promise<Mastercat[]> {
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

  const result = await apiFetch<Mastercat[]>(`/conesearch?${searchParams}`);
  return result ?? [];
}

export function useConeSearch(params: ConeSearchParams | null) {
  return useQuery({
    queryKey: ["conesearch", params],
    queryFn: () => fetchConeSearch(params!),
    enabled: params !== null,
  });
}
