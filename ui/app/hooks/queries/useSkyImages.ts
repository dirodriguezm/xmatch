import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";

import { hipsTilePath } from "@/app/lib/utils/healpix";

/**
 * How long each survey gets to answer before the next one is tried. Per
 * survey, not overall: a hung mirror (alaskybis during the CDS outage) must
 * not use up the time the fallback survey needs.
 */
const PROBE_TIMEOUT_MS = 4_000;

/** CDS registry that maps a HiPS ID to every server that hosts it. */
const MOC_SERVERS = [
  "https://alasky.cds.unistra.fr/MocServer/query",
  "https://alaskybis.cds.unistra.fr/MocServer/query",
];

type HipsRecord = Record<string, string>;

async function fetchHipsRecord(
  id: string,
  signal: AbortSignal
): Promise<HipsRecord | null> {
  const params = new URLSearchParams({
    ID: id,
    get: "record",
    fmt: "json",
    fields: "ID,hips_service_url*,hips_order,hips_tile_format",
  });
  for (const server of MOC_SERVERS) {
    try {
      const res = await fetch(`${server}?${params}`, { signal });
      if (!res.ok) continue;
      const records = (await res.json()) as HipsRecord[];
      return records[0] ?? null;
    } catch {
      if (signal.aborted) throw signal.reason;
    }
  }
  throw new Error("MocServer unreachable");
}

/** `hips_service_url`, `hips_service_url_1`, … in registry order. */
function mirrorUrls(record: HipsRecord): string[] {
  const urls = [record.hips_service_url];
  for (let i = 1; record[`hips_service_url_${i}`]; i++) {
    urls.push(record[`hips_service_url_${i}`]);
  }
  return urls.filter((u): u is string => !!u);
}

/** The mirror's host once it serves real `properties` (not an error page). */
async function probeMirror(url: string, signal: AbortSignal): Promise<string> {
  const res = await fetch(`${url.replace(/\/$/, "")}/properties`, { signal });
  if (!res.ok) throw new Error(`${url}: ${res.status}`);
  return new URL(url).host;
}

/**
 * Host of the fastest mirror serving HiPS `id` (roughly the one Aladin will
 * pick), or null when none does within the timeout or CDS's registry is
 * down. Throws only when `querySignal` aborts.
 */
export async function findHipsHost(
  id: string,
  querySignal: AbortSignal
): Promise<string | null> {
  const signal = AbortSignal.any([
    querySignal,
    AbortSignal.timeout(PROBE_TIMEOUT_MS),
  ]);
  try {
    const record = await fetchHipsRecord(id, signal);
    if (!record) return null;
    return await Promise.any(
      mirrorUrls(record).map((u) => probeMirror(u, signal))
    );
  } catch (error) {
    // Unmounted: let React Query drop it. Anything else: unavailable.
    if (querySignal.aborted) throw error;
    return null;
  }
}

/**
 * The first HiPS in `ids` that some mirror actually serves. Aladin picks the
 * mirror itself, but started against a survey nobody serves it shows a black
 * view (or a spinner forever) and throws errors nothing can catch, so the page
 * asks first. Rejects when none is available.
 */
export function useAvailableHips(ids: readonly string[]) {
  return useQuery({
    queryKey: ["available-hips", ...ids],
    queryFn: async ({ signal }) => {
      for (const id of ids) {
        if (await findHipsHost(id, signal)) return id;
      }
      throw new Error("No HiPS mirror available");
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}

/** File extension of a HiPS `hips_tile_format` (first listed wins). */
const TILE_EXTENSIONS: Record<string, string> = {
  jpeg: "jpg",
  png: "png",
  webp: "webp",
};

/** A survey and the mirror that serves it at a given position. */
export interface HipsAtPosition {
  id: string;
  /** Mirror base URL, to hand Aladin instead of the ID. */
  url: string;
  /** The tile under the position on that mirror, which loaded. */
  tileUrl: string;
}

/**
 * The first HiPS in `ids` with a mirror that serves the deepest tile under
 * (ra, dec), and that mirror's URL. A mirror can answer `properties` and
 * still lack part of the sky: MAST's Pan-STARRS copy has no tiles deeper
 * than order 7 around Orion, which left Aladin black there. Aladin gets the
 * URL rather than the ID, since given an ID it picks a mirror itself and can
 * land on one that 404s (alasky during the CDS outage).
 *
 * Null when the position isn't covered: some server said so (a mirror's 4xx
 * for the tile, or no usable registry record). Throws only when nothing gave
 * an answer, every registry and mirror request failing or timing out, so a
 * gap in the sky doesn't read as an outage in the header.
 */
export async function findHipsAt(
  ids: readonly string[],
  ra: number,
  dec: number,
  querySignal: AbortSignal
): Promise<HipsAtPosition | null> {
  let answered = false;
  for (const id of ids) {
    const signal = AbortSignal.any([
      querySignal,
      AbortSignal.timeout(PROBE_TIMEOUT_MS),
    ]);
    try {
      const record = await fetchHipsRecord(id, signal);
      const order = Number(record?.hips_order);
      const format = record?.hips_tile_format?.split(" ")[0] ?? "jpeg";
      const ext = TILE_EXTENSIONS[format];
      if (!record || !order || !ext) {
        answered = true;
        continue;
      }
      const path = hipsTilePath(order, ra, dec, ext);
      // First mirror to serve the tile wins; the others are cut off.
      const race = new AbortController();
      const url = await Promise.any(
        mirrorUrls(record).map(async (u) => {
          const base = u.replace(/\/$/, "");
          const res = await fetch(`${base}/${path}`, {
            signal: AbortSignal.any([signal, race.signal]),
          });
          if (!res.ok) {
            // A 4xx is the mirror up and lacking this patch of sky.
            if (res.status < 500) answered = true;
            throw new Error(`${base}: ${res.status}`);
          }
          return base;
        })
      );
      race.abort();
      return { id, url, tileUrl: `${url}/${path}` };
    } catch (error) {
      // Unmounted: let React Query drop it. Anything else: next survey.
      if (querySignal.aborted) throw error;
    }
  }
  if (answered) return null;
  throw new Error("No sky image server answered");
}

/** `findHipsAt` as a query; `data` is null where no survey covers the sky. */
export function useHipsAt(ids: readonly string[], ra: number, dec: number) {
  return useQuery({
    queryKey: ["hips-at", ...ids, ra, dec],
    queryFn: ({ signal }) => findHipsAt(ids, ra, dec, signal),
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}

/** After this long without the tile, say so rather than spin silently. */
const TILE_SLOW_MS = 12_000;
/** After this long, stop covering the view: the image may be fine. */
const TILE_GIVE_UP_MS = 30_000;

/** Where a probe of one tile stands. */
export type TileProbe = "pending" | "loaded" | "failed";

/**
 * What the view should say about its tiles: still `loading`, `slow`,
 * `loaded`, `failed`, or `unknown` once the probe has hung so long that
 * covering the view would hide an image that may well be there.
 */
export type TileState = "loading" | "slow" | "loaded" | "failed" | "unknown";

export function tileState(probe: TileProbe, elapsedMs: number): TileState {
  if (probe !== "pending") return probe;
  if (elapsedMs >= TILE_GIVE_UP_MS) return "unknown";
  return elapsedMs >= TILE_SLOW_MS ? "slow" : "loading";
}

/**
 * Whether the survey's tile at the target loads, from an image probe of it.
 * Aladin reports no tile progress, and resource timing can't stand in: a
 * tile served from the image cache (a remount over the same view) records
 * no entry, nor does any once the buffer is full, and outside Chromium a
 * 404 looks like a success. An image's load or error event fires in every
 * browser, from cache too. Aladin requests tiles with the same CORS mode,
 * so the probe and its tiles share the HTTP cache.
 */
export function useTileState(tileUrl: string): TileState {
  const [probe, setProbe] = useState<TileProbe>("pending");
  const [elapsed, setElapsed] = useState(0);
  useEffect(() => {
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => setProbe("loaded");
    img.onerror = () => setProbe("failed");
    img.src = tileUrl;
    const timers = [TILE_SLOW_MS, TILE_GIVE_UP_MS].map((ms) =>
      setTimeout(() => setElapsed(ms), ms)
    );
    return () => {
      img.onload = null;
      img.onerror = null;
      timers.forEach(clearTimeout);
    };
  }, [tileUrl]);
  return tileState(probe, elapsed);
}
