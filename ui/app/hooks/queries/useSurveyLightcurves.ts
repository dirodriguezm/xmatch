import { useQuery } from "@tanstack/react-query";

import type { GaiaEpochResponse } from "@/app/api/gaia-epoch/route";
import type { Ps1LightcurveResponse } from "@/app/api/ps1-lightcurve/route";

async function getJson<T>(url: string, what: string): Promise<T> {
  const r = await fetch(url);
  if (!r.ok) {
    const body = await r.json().catch(() => ({}));
    throw new Error(body.error || `${what} lookup failed`);
  }
  return r.json();
}

/** Either a known Gaia source_id or a position to cone-search. */
export type GaiaEpochParams =
  | { sourceId: string }
  | { ra: number; dec: number };

export function useGaiaEpochPhotometry(params: GaiaEpochParams | null) {
  return useQuery({
    queryKey: ["gaia-epoch", params],
    queryFn: () => {
      const qs =
        "sourceId" in params!
          ? new URLSearchParams({ source_id: params!.sourceId })
          : new URLSearchParams({
              ra: String(params!.ra),
              dec: String(params!.dec),
            });
      return getJson<GaiaEpochResponse>(
        `/api/gaia-epoch?${qs}`,
        "Gaia epoch photometry"
      );
    },
    enabled: params !== null,
    // Gaia DR3 is frozen.
    staleTime: Infinity,
    retry: 1,
  });
}

export function usePs1Lightcurve(params: { ra: number; dec: number } | null) {
  return useQuery({
    queryKey: ["ps1-lightcurve", params?.ra, params?.dec],
    queryFn: () =>
      getJson<Ps1LightcurveResponse>(
        `/api/ps1-lightcurve?${new URLSearchParams({
          ra: String(params!.ra),
          dec: String(params!.dec),
        })}`,
        "Pan-STARRS light curve"
      ),
    enabled: params !== null,
    // PS1 DR2 is frozen.
    staleTime: Infinity,
    retry: 1,
  });
}
