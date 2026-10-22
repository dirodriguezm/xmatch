import { useQueries, useQuery } from "@tanstack/react-query";

import type { HiligtResponse } from "@/app/api/hiligt/route";
import type { XmmSourceResponse } from "@/app/api/xmm-source/route";
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
