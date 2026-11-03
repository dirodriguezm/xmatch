import { useQuery } from "@tanstack/react-query";

import type { VizierSedResponse } from "@/app/api/vizier-sed/route";

export interface VizierSedParams {
  ra: number;
  dec: number;
  radius?: number;
}

async function fetchVizierSed(p: VizierSedParams): Promise<VizierSedResponse> {
  const qs = new URLSearchParams({ ra: String(p.ra), dec: String(p.dec) });
  if (p.radius !== undefined) qs.set("radius", String(p.radius));
  const r = await fetch(`/api/vizier-sed?${qs}`);
  if (!r.ok) {
    const body = await r.json().catch(() => ({}));
    throw new Error(body.error || "VizieR SED lookup failed");
  }
  return r.json();
}

export function useVizierSed(params: VizierSedParams | null) {
  return useQuery({
    queryKey: ["vizier-sed", params],
    queryFn: () => fetchVizierSed(params!),
    enabled: params !== null,
    // Published photometry does not change within a session.
    staleTime: Infinity,
    retry: 1,
  });
}
