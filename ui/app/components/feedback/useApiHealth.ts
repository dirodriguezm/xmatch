"use client";

import { useQuery } from "@tanstack/react-query";

import type { HealthResponse } from "@/app/api/health/route";

export const HEALTH_POLL_MS = 30_000;

async function fetchHealth(): Promise<HealthResponse> {
  const resp = await fetch("/api/health", { cache: "no-store" });
  if (!resp.ok) throw new Error(`Health check failed (HTTP ${resp.status})`);
  return resp.json();
}

/** XWave API liveness from /api/health, re-polled every 30 s. */
export function useApiHealth(enabled = true) {
  return useQuery({
    queryKey: ["api-health"],
    queryFn: fetchHealth,
    enabled,
    refetchInterval: HEALTH_POLL_MS,
    refetchIntervalInBackground: false,
    staleTime: 0,
    retry: false,
  });
}
