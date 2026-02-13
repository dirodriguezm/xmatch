import { useQuery } from "@tanstack/react-query";

import { apiFetch } from "@/app/lib/api/client";
import type { components } from "@/types/xwave-api";

type Allwise = components["schemas"]["repository.Allwise"];

export interface MetadataParams {
  id: string;
  catalog: string;
}

async function fetchMetadata(params: MetadataParams): Promise<Allwise | null> {
  const searchParams = new URLSearchParams({
    id: params.id,
    catalog: params.catalog,
  });

  return apiFetch<Allwise>(`/metadata?${searchParams}`);
}

export function useMetadata(params: MetadataParams | null) {
  return useQuery({
    queryKey: ["metadata", params],
    queryFn: () => fetchMetadata(params!),
    enabled: params !== null,
  });
}
