/**
 * CDS services with their official mirrors, tried in order. Strasbourg first;
 * CfA (Harvard) and IDIA (South Africa) host copies that CDS points users to
 * when Strasbourg is down (as during the November 2026 incident). SIMBAD
 * object ids (oid) are the same on every mirror, so a match from one can be
 * followed up on another.
 */
export const SESAME_MIRRORS = [
  "https://cds.unistra.fr/cgi-bin/nph-sesame",
  "https://vizier.cfa.harvard.edu/viz-bin/nph-sesame",
  "https://vizier.idia.ac.za/viz-bin/nph-sesame",
] as const;

export const SIMBAD_TAP_MIRRORS = [
  "https://simbad.cds.unistra.fr/simbad/sim-tap/sync",
  "https://simbad.cfa.harvard.edu/simbad/sim-tap/sync",
  "https://simbad.idia.ac.za/simbad/sim-tap/sync",
] as const;

/**
 * VizieR's photometry (SED) service. Only CfA serves it in full: IDIA
 * answers 500, China-VO holds only small catalogues (IRAS rows only). CfA
 * needs 1–2 minutes for a position it hasn't cached (~1 s afterwards); the
 * SED proxy retries its interim answers.
 */
export const VIZIER_SED_MIRRORS: readonly Mirror[] = [
  { base: "https://vizier.cds.unistra.fr/viz-bin/sed", timeoutMs: 20_000 },
  // CfA's gateway answers 504 after 60 s while it keeps computing.
  { base: "https://vizier.cfa.harvard.edu/viz-bin/sed", timeoutMs: 65_000 },
];

/**
 * Response header naming the host that actually answered a proxied request,
 * so the status page can tell when a mirror is standing in for the primary.
 */
export const UPSTREAM_HOST_HEADER = "x-upstream-host";

export interface MirrorResponse {
  response: Response;
  /** Host that answered, after redirects (CDS may redirect to a mirror). */
  host: string;
}

/** A base URL, or one with its own timeout (e.g. a slow but working mirror). */
export type Mirror = string | { base: string; timeoutMs: number };

/**
 * How long a working mirror is tried first after the primary failed, before
 * the primary gets another chance. Saves every request from re-discovering
 * the outage, which costs a full timeout when the primary hangs.
 */
const PREFER_MIRROR_MS = 10 * 60_000;

/** Per mirror list: the mirror to try first, and until when. */
const preferred = new Map<string, { base: string; until: number }>();

/** Forget remembered mirrors (tests). */
export function resetPreferredMirrors() {
  preferred.clear();
}

const baseOf = (m: Mirror) => (typeof m === "string" ? m : m.base);

/**
 * GET `path` from each mirror in turn and return the first OK response.
 * A mirror is skipped when it errors, times out (its own `timeoutMs`, else
 * the default, so one hung server can't eat the others' time) or answers
 * non-2xx. When none is OK, returns the last HTTP response if there was one,
 * else throws.
 *
 * After the primary (first) mirror fails and another answers, that one is
 * tried first for the next {@link PREFER_MIRROR_MS}; then the primary again.
 */
export async function fetchFromMirrors(
  mirrors: readonly Mirror[],
  path: string,
  init: RequestInit,
  timeoutMs: number
): Promise<MirrorResponse> {
  const key = mirrors.map(baseOf).join("|");
  const primary = mirrors[0] && baseOf(mirrors[0]);
  const remembered = preferred.get(key);
  const first =
    remembered && remembered.until > Date.now()
      ? mirrors.find((m) => baseOf(m) === remembered.base)
      : undefined;
  const order = first
    ? [first, ...mirrors.filter((m) => m !== first)]
    : mirrors;

  let last: MirrorResponse | undefined;
  let lastError: unknown = new Error("No mirrors given");

  for (const mirror of order) {
    const base = baseOf(mirror);
    try {
      const response = await fetch(`${base}${path}`, {
        ...init,
        signal: AbortSignal.timeout(
          typeof mirror === "string" ? timeoutMs : mirror.timeoutMs
        ),
      });
      const result = { response, host: new URL(response.url || base).host };
      if (response.ok) {
        if (base === primary) preferred.delete(key);
        else preferred.set(key, { base, until: Date.now() + PREFER_MIRROR_MS });
        return result;
      }
      last = result;
      console.warn(`Mirror ${base} answered ${response.status}, trying next`);
    } catch (error) {
      lastError = error;
      console.warn(`Mirror ${base} failed, trying next:`, error);
    }
    if (base === remembered?.base) preferred.delete(key);
  }

  if (last) return last;
  throw lastError;
}
