/**
 * Upstream checks shown on the Status page. Each one calls an existing
 * same-origin proxy with cheap parameters (M31 is resolvable everywhere and
 * sits inside every survey footprint these services cover).
 */

export type CheckState = "up" | "degraded" | "down";

export interface UpstreamCheck {
  id: string;
  name: string;
  provider: string;
  usedFor: string;
  /** Same-origin proxy URL to probe. */
  path: string;
  /** Some proxies cache upstream answers server-side; latency may reflect that. */
  cached?: boolean;
}

const M31 = { ra: 10.684708, dec: 41.26875 };

export const UPSTREAM_CHECKS: UpstreamCheck[] = [
  {
    id: "conesearch",
    name: "XWave cone search",
    provider: "XWave API",
    usedFor: "Search results and counterparts",
    path: `/api/conesearch?ra=${M31.ra}&dec=${M31.dec}&radius=2&catalog=gaia&nneighbor=1`,
  },
  {
    id: "sesame",
    name: "Sesame name resolver",
    provider: "CDS",
    usedFor: "Turning object names into coordinates",
    path: "/api/sesame?name=M31",
  },
  {
    id: "simbad",
    name: "SIMBAD",
    provider: "CDS",
    usedFor: "Object identity on the object page",
    path: `/api/simbad?ra=${M31.ra}&dec=${M31.dec}&radius=2`,
    cached: true,
  },
  {
    id: "vizier-sed",
    name: "VizieR SED",
    provider: "CDS",
    usedFor: "Published photometry in the SED",
    path: `/api/vizier-sed?ra=${M31.ra}&dec=${M31.dec}&radius=1`,
    cached: true,
  },
  {
    id: "galactic-dust",
    name: "Galactic dust reddening",
    provider: "NASA/IPAC IRSA",
    usedFor: "Extinction correction of the SED",
    path: `/api/galactic-dust?ra=${M31.ra}&dec=${M31.dec}`,
    cached: true,
  },
  {
    id: "ztf",
    name: "ZTF light curves",
    provider: "ALeRCE",
    usedFor: "ZTF time series on the object page",
    path: `/api/ztf-lightcurve?ra=${M31.ra}&dec=${M31.dec}&radius=1`,
  },
  {
    id: "xmm-source",
    name: "5XMM-DR15 X-ray catalogue",
    provider: "ESA XMM-Newton Science Archive",
    usedFor: "X-ray source summary on the object page",
    path: `/api/xmm-source?ra=${M31.ra}&dec=${M31.dec}`,
    cached: true,
  },
];

/** Above this the service answers but is slow enough to notice. */
export const DEGRADED_LATENCY_MS = 4_000;

/** Classify a probe from its HTTP status (null = network error/timeout). */
export function classifyCheck(
  httpStatus: number | null,
  latencyMs: number | null
): CheckState {
  if (httpStatus === null) return "down";
  // 404 = "not found" answers (e.g. nothing at the position) still prove the
  // upstream is reachable; 5xx and 502s from the proxy mean it is not.
  const reachable =
    (httpStatus >= 200 && httpStatus < 300) || httpStatus === 404;
  if (!reachable) return "down";
  if (latencyMs !== null && latencyMs > DEGRADED_LATENCY_MS) return "degraded";
  return "up";
}

/** Worst state wins; an empty list counts as up. */
export function overallState(states: CheckState[]): CheckState {
  if (states.includes("down")) return "down";
  if (states.includes("degraded")) return "degraded";
  return "up";
}
