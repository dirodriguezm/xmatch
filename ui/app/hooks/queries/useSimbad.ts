import { useQuery } from "@tanstack/react-query";

import type { SimbadResponse } from "@/app/api/simbad/route";
import type { SimbadRefsResponse } from "@/app/api/simbad-refs/route";

export interface SimbadParams {
  ra: number;
  dec: number;
  radius?: number;
}

async function fetchSimbad(p: SimbadParams): Promise<SimbadResponse> {
  const qs = new URLSearchParams({ ra: String(p.ra), dec: String(p.dec) });
  if (p.radius !== undefined) qs.set("radius", String(p.radius));
  const r = await fetch(`/api/simbad?${qs}`);
  if (!r.ok) {
    const body = await r.json().catch(() => ({}));
    throw new Error(body.error || "SIMBAD lookup failed");
  }
  return r.json();
}

/** Nearest SIMBAD object to a position. */
export function useSimbad(params: SimbadParams | null) {
  return useQuery({
    queryKey: ["simbad", params],
    queryFn: () => fetchSimbad(params!),
    enabled: params !== null,
    staleTime: Infinity,
    retry: 1,
  });
}

/** Papers about a SIMBAD object, newest first. */
export function useSimbadRefs(oid: number | null, limit: number) {
  return useQuery({
    queryKey: ["simbad-refs", oid, limit],
    queryFn: async (): Promise<SimbadRefsResponse> => {
      const qs = new URLSearchParams({
        oid: String(oid),
        limit: String(limit),
      });
      const r = await fetch(`/api/simbad-refs?${qs}`);
      if (!r.ok) {
        const body = await r.json().catch(() => ({}));
        throw new Error(body.error || "SIMBAD references lookup failed");
      }
      return r.json();
    },
    enabled: oid !== null,
    staleTime: Infinity,
    // Keep the shorter list on screen while a longer one loads.
    placeholderData: (previous) => previous,
    retry: 1,
  });
}
