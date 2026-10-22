import { useQuery } from "@tanstack/react-query";

export interface DesiSpectrumResult {
  found: boolean;
  sparclId?: string;
  targetid?: string;
  spectype?: string | null;
  redshift?: number | null;
  wavelength?: number[];
  flux?: number[];
  model?: number[] | null;
  ivar?: number[] | null;
}

export interface DesiSpectrumParams {
  targetid: string;
}

async function fetchDesiSpectrum(
  p: DesiSpectrumParams
): Promise<DesiSpectrumResult> {
  const qs = new URLSearchParams({ targetid: p.targetid });
  const r = await fetch(`/api/desi-spectrum?${qs}`);
  if (!r.ok) throw new Error("DESI spectrum lookup failed");
  return r.json();
}

export function useDesiSpectrum(targetid: string | null) {
  return useQuery({
    queryKey: ["desi-spectrum", targetid],
    queryFn: () => fetchDesiSpectrum({ targetid: targetid! }),
    enabled: !!targetid,
    staleTime: 24 * 60 * 60 * 1000,
    // Ride out transient SPARCL blips (e.g. a brief upstream block) without a
    // manual reload. A `{found:false}` response is a non-error 200, so it isn't
    // retried — only thrown errors (non-2xx) are.
    retry: 2,
    retryDelay: (n) => Math.min(1000 * 2 ** n, 8000),
  });
}
