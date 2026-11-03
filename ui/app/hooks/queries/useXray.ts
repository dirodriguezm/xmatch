import { useQueries, useQuery } from "@tanstack/react-query";

import type { BatLightcurveResponse } from "@/app/api/bat-lightcurve/route";
import type { ChandraSourceResponse } from "@/app/api/chandra-source/route";
import type { FermiLcrResponse } from "@/app/api/fermi-lcr/route";
import type { HiligtResponse } from "@/app/api/hiligt/route";
import type { XmmSourceResponse } from "@/app/api/xmm-source/route";
import type { LcrCadence } from "@/app/lib/utils/gamma";
import { HILIGT_MISSIONS, type HiligtMission } from "@/app/lib/utils/hiligt";

async function getJson<T>(url: string, what: string): Promise<T> {
  const r = await fetch(url);
  if (!r.ok) {
    const body = await r.json().catch(() => ({}));
    throw new Error(body.error || `${what} lookup failed`);
  }
  return r.json();
}

type Position = { ra: number; dec: number };

/** Nearest 5XMM-DR15 source. */
export function useXmmSource(pos: Position | null) {
  return useQuery({
    queryKey: ["xmm-source", pos?.ra, pos?.dec],
    queryFn: () =>
      getJson<XmmSourceResponse>(
        `/api/xmm-source?${new URLSearchParams({
          ra: String(pos!.ra),
          dec: String(pos!.dec),
        })}`,
        "5XMM"
      ),
    enabled: pos !== null,
    // 5XMM-DR15 is frozen.
    staleTime: Infinity,
    retry: 1,
  });
}

/** Nearest Chandra Source Catalog 2.2 source. */
export function useChandraSource(pos: Position | null) {
  return useQuery({
    queryKey: ["chandra-source", pos?.ra, pos?.dec],
    queryFn: () =>
      getJson<ChandraSourceResponse>(
        `/api/chandra-source?${new URLSearchParams({
          ra: String(pos?.ra),
          dec: String(pos?.dec),
        })}`,
        "Chandra Source Catalog"
      ),
    enabled: pos !== null,
    // CSC 2.2 is frozen.
    staleTime: Infinity,
    retry: 1,
  });
}

/** Fermi-LAT light curve of a 4FGL source at one cadence. */
export function useFermiLightcurve(source: string | null, cadence: LcrCadence) {
  return useQuery({
    queryKey: ["fermi-lcr", source, cadence],
    queryFn: () =>
      getJson<FermiLcrResponse>(
        `/api/fermi-lcr?${new URLSearchParams({ source: source ?? "", cadence })}`,
        "Fermi LCR"
      ),
    enabled: source !== null,
    staleTime: 60 * 60 * 1000,
    retry: 1,
  });
}

/** Daily Swift-BAT monitor light curve of a source. */
export function useBatLightcurve(path: string | null) {
  return useQuery({
    queryKey: ["bat-lightcurve", path],
    queryFn: () =>
      getJson<BatLightcurveResponse>(
        `/api/bat-lightcurve?${new URLSearchParams({ path: path ?? "" })}`,
        "Swift-BAT monitor"
      ),
    enabled: path !== null,
    staleTime: 60 * 60 * 1000,
    retry: 1,
  });
}

export const HILIGT_MISSION_IDS = Object.keys(
  HILIGT_MISSIONS
) as HiligtMission[];

/** HILIGT fluxes / upper limits, one query per mission so each arrives alone. */
export function useHiligt(pos: Position | null) {
  return useQueries({
    queries: HILIGT_MISSION_IDS.map((mission) => ({
      queryKey: ["hiligt", mission, pos?.ra, pos?.dec],
      queryFn: () =>
        getJson<HiligtResponse>(
          `/api/hiligt?${new URLSearchParams({
            mission,
            ra: String(pos!.ra),
            dec: String(pos!.dec),
          })}`,
          HILIGT_MISSIONS[mission]
        ),
      enabled: pos !== null,
      staleTime: Infinity,
      retry: 1,
    })),
  });
}
