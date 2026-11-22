import { useSyncExternalStore } from "react";

/**
 * Services the app has seen failing in this session, from the queries it
 * already runs (no extra requests). Feeds the header's status indicator;
 * the Status page runs real checks.
 */

/** First query-key element → service name, for queries behind a proxy. */
const SERVICE_BY_QUERY_KEY: Record<string, string> = {
  conesearch: "XWave API",
  counterpart: "XWave API",
  neighbors: "XWave API",
  metadata: "XWave API",
  lightcurve: "XWave API",
  "bulk-metadata": "XWave API",
  simbad: "SIMBAD",
  "simbad-refs": "SIMBAD references",
  "vizier-sed": "VizieR SED",
  "galactic-reddening": "Galactic dust reddening",
  "available-hips": "Sky images",
  "hips-at": "Sky images",
  "ztf-lightcurve": "ZTF light curves",
  "ps1-lightcurve": "Pan-STARRS light curves",
  "crts-lightcurve": "Catalina light curves",
  "gaia-epoch": "Gaia epoch photometry",
  "xmm-source": "XMM-Newton catalogue",
  "chandra-source": "Chandra Source Catalog",
  "desi-target": "DESI targets",
  "desi-spectrum": "DESI spectra",
  "fermi-lcr": "Fermi LAT light curves",
  "bat-lightcurve": "Swift BAT light curves",
  hiligt: "HILIGT upper limits",
};

/** Service name for a query key, or null when it isn't an external call. */
export function serviceForQueryKey(
  queryKey: readonly unknown[]
): string | null {
  const head = queryKey[0];
  return typeof head === "string" ? (SERVICE_BY_QUERY_KEY[head] ?? null) : null;
}

let failing: readonly string[] = [];
const listeners = new Set<() => void>();

function emit(next: readonly string[]) {
  failing = next;
  listeners.forEach((notify) => notify());
}

/**
 * Whether a failed query still shows the service is up: a 4xx means it
 * answered and refused the request (out-of-range RA, radius too large), which
 * is the user's input, not an outage. Duck-typed so ApiError, ConeSearchError
 * and any other error carrying an HTTP status all count.
 */
export function isClientError(error: unknown): boolean {
  if (typeof error !== "object" || error === null || !("status" in error)) {
    return false;
  }
  const { status } = error;
  return typeof status === "number" && status >= 400 && status < 500;
}

/**
 * Whether a failed query never reached the service: refused in the browser
 * before sending (a radius over the limit). It says nothing either way about
 * the service, so it must neither raise nor clear a failing flag.
 */
export function isLocalRefusal(error: unknown): boolean {
  return (
    typeof error === "object" &&
    error !== null &&
    "local" in error &&
    error.local === true
  );
}

/** Record the outcome of a call to `service`; the latest outcome wins. */
export function reportServiceResult(service: string | null, ok: boolean) {
  if (!service) return;
  const isFailing = failing.includes(service);
  if (ok && isFailing) emit(failing.filter((s) => s !== service));
  else if (!ok && !isFailing) emit([...failing, service]);
}

/**
 * Record a failed query. A 4xx is the service answering, so it doesn't count
 * as an outage; a request refused before sending says nothing either way.
 */
export function reportQueryError(queryKey: readonly unknown[], error: unknown) {
  if (isLocalRefusal(error)) return;
  reportServiceResult(serviceForQueryKey(queryKey), isClientError(error));
}

/** Current failing services (what the hook below renders). */
export function getFailingServices(): readonly string[] {
  return failing;
}

function subscribe(onChange: () => void) {
  listeners.add(onChange);
  return () => listeners.delete(onChange);
}

const EMPTY: readonly string[] = [];

/** Names of the services whose latest call in this session failed. */
export function useFailingServices(): readonly string[] {
  return useSyncExternalStore(subscribe, getFailingServices, () => EMPTY);
}
