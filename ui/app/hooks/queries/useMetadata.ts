import { useQuery } from "@tanstack/react-query";

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

  const response = await fetch(`/api/metadata?${searchParams}`);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || "Failed to fetch metadata");
  }

  return response.json();
}

export function useMetadata(params: MetadataParams | null) {
  return useQuery({
    queryKey: ["metadata", params],
    queryFn: () => fetchMetadata(params!),
    enabled: params !== null,
  });
}
