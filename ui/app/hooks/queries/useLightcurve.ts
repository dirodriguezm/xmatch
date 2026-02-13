import { useQuery } from "@tanstack/react-query";

import { apiFetch } from "@/app/lib/api/client";
import type { components } from "@/types/xwave-api";

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

export interface LightcurveParams {
  ra: number;
  dec: number;
  radius: number;
  nneighbor?: number;
}

async function fetchLightcurve(
  params: LightcurveParams
): Promise<Lightcurve | null> {
  const searchParams = new URLSearchParams({
    ra: params.ra.toString(),
    dec: params.dec.toString(),
    radius: params.radius.toString(),
  });

  if (params.nneighbor !== undefined) {
    searchParams.set("nneighbor", params.nneighbor.toString());
  }

  return apiFetch<Lightcurve>(`/lightcurve?${searchParams}`);
}

export function useLightcurve(params: LightcurveParams | null) {
  return useQuery({
    queryKey: ["lightcurve", params],
    queryFn: () => fetchLightcurve(params!),
    enabled: params !== null,
  });
}
