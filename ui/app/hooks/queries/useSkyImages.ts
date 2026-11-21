import { useQuery } from "@tanstack/react-query";

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
    fields: "ID,hips_service_url*",
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

/** Resolves when a mirror serves real `properties` (not an error page). */
async function probeMirror(url: string, signal: AbortSignal): Promise<void> {
  const res = await fetch(`${url.replace(/\/$/, "")}/properties`, { signal });
  if (!res.ok) throw new Error(`${url}: ${res.status}`);
}

/** False when no mirror serves `id` within the timeout (or CDS's registry is down). */
async function isHipsAvailable(
  id: string,
  querySignal: AbortSignal
): Promise<boolean> {
  const signal = AbortSignal.any([
    querySignal,
    AbortSignal.timeout(PROBE_TIMEOUT_MS),
  ]);
  try {
    const record = await fetchHipsRecord(id, signal);
    if (!record) return false;
    await Promise.any(mirrorUrls(record).map((u) => probeMirror(u, signal)));
    return true;
  } catch (error) {
    // Unmounted: let React Query drop it. Anything else: try the next survey.
    if (querySignal.aborted) throw error;
    return false;
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
        if (await isHipsAvailable(id, signal)) return id;
      }
      throw new Error("No HiPS mirror available");
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}
