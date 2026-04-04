import { useQuery } from "@tanstack/react-query";

import type { components } from "@/types/xwave-api";

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

export interface ZtfLightcurveParams {
  ra: number;
  dec: number;
  radius?: number;
}

async function fetchZtfLightcurve(
  params: ZtfLightcurveParams
): Promise<Lightcurve | null> {
  const searchParams = new URLSearchParams({
    ra: params.ra.toString(),
    dec: params.dec.toString(),
  });

  if (params.radius !== undefined) {
    searchParams.set("radius", params.radius.toString());
  }

  const response = await fetch(`/api/ztf-lightcurve?${searchParams}`);

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || "Failed to fetch ZTF lightcurve data");
  }

  return response.json();
}

export function useZtfLightcurve(params: ZtfLightcurveParams | null) {
  return useQuery({
    queryKey: ["ztf-lightcurve", params],
    queryFn: () => fetchZtfLightcurve(params!),
    enabled: params !== null,
  });
}
