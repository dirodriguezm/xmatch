"use client";

import { useSearchParams } from "next/navigation";
import { useMemo } from "react";

import {
  radiiFromParam,
  toNearbySources,
} from "@/app/components/results/nearby";
import {
  catalogsOf,
  type NearbySource,
} from "@/app/components/results/nearby/shared";
import { useLinkedHighlight } from "@/app/components/results/nearby/useLinkedHighlight";
import type { EnrichedResult } from "@/app/components/results/ResultsTable";

/** Matches as sky offsets, the searched radii and a shared hover/selection. */
export function useNearbyResults(
  data: EnrichedResult[],
  target: { ra: number; dec: number } | null
) {
  const params = useSearchParams();
  const catalogRadii = params.get("catalogRadii") ?? "";
  const radii = useMemo(() => radiiFromParam(catalogRadii), [catalogRadii]);
  const sources = useMemo<NearbySource[]>(
    () => (target ? toNearbySources(data, target) : []),
    [data, target]
  );
  const catalogs = useMemo(() => catalogsOf(radii, sources), [radii, sources]);
  const highlight = useLinkedHighlight();
  return { radii, sources, catalogs, highlight };
}
