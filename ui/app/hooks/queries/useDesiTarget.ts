import { useQuery } from "@tanstack/react-query";

export interface DesiTargetResult {
  targetid: string | null;
  separationArcsec?: number;
}

export interface DesiTargetParams {
  ra: number;
  dec: number;
  radius?: number;
}

async function fetchDesiTarget(p: DesiTargetParams): Promise<DesiTargetResult> {
  const qs = new URLSearchParams({ ra: String(p.ra), dec: String(p.dec) });
  if (p.radius !== undefined) qs.set("radius", String(p.radius));
  const r = await fetch(`/api/desi-target?${qs}`);
  if (!r.ok) throw new Error("DESI target lookup failed");
  return r.json();
}

export function useDesiTarget(params: DesiTargetParams | null) {
  return useQuery({
    queryKey: ["desi-target", params],
    queryFn: () => fetchDesiTarget(params!),
    enabled: params !== null,
    staleTime: 24 * 60 * 60 * 1000,
  });
}
