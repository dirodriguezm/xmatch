import { useQuery } from "@tanstack/react-query";

import type { GalacticReddening } from "@/app/lib/utils/irsaDust";

async function fetchReddening(
  ra: number,
  dec: number
): Promise<GalacticReddening> {
  const qs = new URLSearchParams({ ra: String(ra), dec: String(dec) });
  const r = await fetch(`/api/galactic-dust?${qs}`);
  if (!r.ok) {
    const body = await r.json().catch(() => ({}));
    throw new Error(body.error || "Galactic reddening lookup failed");
  }
  return r.json();
}

export function useGalacticReddening(
  params: { ra: number; dec: number } | null
) {
  return useQuery({
    queryKey: ["galactic-reddening", params?.ra, params?.dec],
    queryFn: () => fetchReddening(params!.ra, params!.dec),
    enabled: params !== null,
    // The dust maps are static.
    staleTime: Infinity,
    retry: 1,
  });
}
