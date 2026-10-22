import { NextResponse } from "next/server";

import { API_ORIGIN } from "@/app/lib/constants/site";

/**
 * Liveness of the XWave API: pings `${API_ORIGIN}/ping` (which answers
 * "pong") and reports how long it took. Never cached.
 */

export const dynamic = "force-dynamic";

const PING_URL = `${API_ORIGIN}/ping`;
const TIMEOUT_MS = 5_000;

export interface HealthResponse {
  status: "up" | "down";
  latencyMs: number | null;
  checkedAt: string;
  error?: string;
}

export async function GET() {
  const started = performance.now();
  let body: HealthResponse;
  try {
    const resp = await fetch(PING_URL, {
      cache: "no-store",
      signal: AbortSignal.timeout(TIMEOUT_MS),
    });
    const text = (await resp.text()).trim();
    const latencyMs = Math.round(performance.now() - started);
    body =
      resp.ok && text === "pong"
        ? { status: "up", latencyMs, checkedAt: new Date().toISOString() }
        : {
            status: "down",
            latencyMs,
            checkedAt: new Date().toISOString(),
            error: `Unexpected response (HTTP ${resp.status})`,
          };
  } catch (err) {
    const timedOut = err instanceof Error && err.name === "TimeoutError";
    body = {
      status: "down",
      latencyMs: null,
      checkedAt: new Date().toISOString(),
      error: timedOut
        ? `No reply within ${TIMEOUT_MS / 1000} s`
        : "Could not reach the API",
    };
  }
  return NextResponse.json(body, {
    headers: { "Cache-Control": "no-store, max-age=0" },
  });
}
