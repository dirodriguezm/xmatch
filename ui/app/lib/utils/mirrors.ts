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
 * GET `path` from each base URL in turn and return the first OK response.
 * A mirror is skipped when it errors, times out (`timeoutMs` each, so one
 * hung server can't eat the others' time) or answers non-2xx. When none is
 * OK, returns the last HTTP response if there was one, else throws.
 */
export async function fetchFromMirrors(
  bases: readonly string[],
  path: string,
  init: RequestInit,
  timeoutMs: number
): Promise<Response> {
  let lastResponse: Response | undefined;
  let lastError: unknown = new Error("No mirrors given");

  for (const base of bases) {
    try {
      const resp = await fetch(`${base}${path}`, {
        ...init,
        signal: AbortSignal.timeout(timeoutMs),
      });
      if (resp.ok) return resp;
      lastResponse = resp;
      console.warn(`Mirror ${base} answered ${resp.status}, trying next`);
    } catch (error) {
      lastError = error;
      console.warn(`Mirror ${base} failed, trying next:`, error);
    }
  }

  if (lastResponse) return lastResponse;
  throw lastError;
}
