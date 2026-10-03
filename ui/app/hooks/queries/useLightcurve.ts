import { useQuery } from "@tanstack/react-query";

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

  const response = await fetch(`/api/lightcurve?${searchParams}`);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || "Failed to fetch lightcurve data");
  }

  return response.json();
}

export function useLightcurve(params: LightcurveParams | null) {
  return useQuery({
    queryKey: ["lightcurve", params],
    queryFn: () => fetchLightcurve(params!),
    enabled: params !== null,
  });
}
