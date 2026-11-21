/**
 * Upstream checks shown on the Status page. Most call an existing
 * same-origin proxy with cheap parameters (M31 is resolvable everywhere and
 * sits inside every survey footprint these services cover); sky image
 * surveys are probed through their mirrors, as Aladin loads them.
 */

export type CheckState = "up" | "degraded" | "down";

interface CheckBase {
  id: string;
  name: string;
  provider: string;
  usedFor: string;
  /**
   * Host that answers when all is well. When another host answers, the
   * page says a mirror is standing in.
   */
  primaryHost?: string;
  /** Some proxies cache upstream answers server-side; latency may reflect that. */
  cached?: boolean;
  /** Longer than the default probe timeout, for services slow by nature. */
  timeoutMs?: number;
}

export type UpstreamCheck = CheckBase &
  (
    | {
        kind: "proxy";
        /** Same-origin proxy URL to probe. */
        path: string;
      }
    | {
        kind: "hips";
        /** HiPS image survey ID, resolved to its mirrors via the MocServer. */
        hipsId: string;
      }
  );

const M31 = { ra: 10.684708, dec: 41.26875 };
const at = `ra=${M31.ra}&dec=${M31.dec}`;
/** M31's SIMBAD oid (the same on every SIMBAD mirror). */
const M31_OID = 1575544;
const ALASKY = "alasky.cds.unistra.fr";

export const UPSTREAM_CHECKS: UpstreamCheck[] = [
  {
    id: "conesearch",
    kind: "proxy",
    name: "XWave cone search",
    provider: "XWave API",
    usedFor: "Search results and counterparts",
    path: `/api/conesearch?${at}&radius=2&catalog=gaia&nneighbor=1`,
  },
  {
    id: "sesame",
    kind: "proxy",
    name: "Sesame name resolver",
    provider: "CDS",
    usedFor: "Turning object names into coordinates",
    path: "/api/sesame?name=M31",
    primaryHost: "cds.unistra.fr",
  },
  {
    id: "simbad",
    kind: "proxy",
    name: "SIMBAD",
    provider: "CDS",
    usedFor: "Object identity on the object page",
    path: `/api/simbad?${at}&radius=2`,
    primaryHost: "simbad.cds.unistra.fr",
    cached: true,
  },
  {
    id: "simbad-refs",
    kind: "proxy",
    name: "SIMBAD references",
    provider: "CDS",
    usedFor: "Literature tab on the object page",
    path: `/api/simbad-refs?oid=${M31_OID}&limit=1`,
    primaryHost: "simbad.cds.unistra.fr",
    cached: true,
  },
  {
    id: "vizier-sed",
    kind: "proxy",
    name: "VizieR SED",
    provider: "CDS",
    usedFor: "Published photometry in the SED",
    path: `/api/vizier-sed?${at}&radius=1`,
    primaryHost: "vizier.cds.unistra.fr",
    cached: true,
    // The CfA mirror takes 1–2 min for a position it hasn't cached.
    timeoutMs: 200_000,
  },
  {
    id: "galactic-dust",
    kind: "proxy",
    name: "Galactic dust reddening",
    provider: "NASA/IPAC IRSA",
    usedFor: "Extinction correction of the SED",
    path: `/api/galactic-dust?${at}`,
    cached: true,
  },
  {
    id: "hips-dss",
    kind: "hips",
    name: "DSS2 colour images",
    provider: "CDS HiPS and mirrors",
    usedFor: "Sky view on the object page",
    hipsId: "CDS/P/DSS2/color",
    primaryHost: ALASKY,
  },
  {
    id: "hips-ps1",
    kind: "hips",
    name: "Pan-STARRS DR1 colour images",
    provider: "CDS HiPS and mirrors",
    usedFor: "Sky image beside search results, north of −30°",
    hipsId: "CDS/P/PanSTARRS/DR1/color-i-r-g",
    primaryHost: ALASKY,
  },
  {
    id: "hips-skymapper",
    kind: "hips",
    name: "SkyMapper DR4 colour images",
    provider: "CDS HiPS",
    usedFor: "Sky image beside search results, south of −30° (else DSS2)",
    hipsId: "CDS/P/Skymapper/DR4/color",
    primaryHost: ALASKY,
  },
  {
    id: "ztf",
    kind: "proxy",
    name: "ZTF light curves",
    provider: "ALeRCE",
    usedFor: "ZTF time series on the object page",
    path: `/api/ztf-lightcurve?${at}&radius=1`,
  },
  {
    id: "ps1-lightcurve",
    kind: "proxy",
    name: "Pan-STARRS light curves",
    provider: "MAST (STScI)",
    usedFor: "Pan-STARRS time series on the object page",
    path: `/api/ps1-lightcurve?${at}`,
    cached: true,
  },
  {
    id: "crts",
    kind: "proxy",
    name: "Catalina (CRTS) light curves",
    provider: "Caltech",
    usedFor: "Catalina time series on the object page",
    path: `/api/crts-lightcurve?${at}`,
    cached: true,
  },
  {
    id: "gaia-epoch",
    kind: "proxy",
    name: "Gaia epoch photometry",
    provider: "ESA Gaia Archive",
    usedFor: "Gaia time series on the object page",
    path: `/api/gaia-epoch?${at}`,
    cached: true,
  },
  {
    id: "xmm-source",
    kind: "proxy",
    name: "5XMM-DR15 X-ray catalogue",
    provider: "ESA XMM-Newton Science Archive",
    usedFor: "X-ray source summary on the object page",
    path: `/api/xmm-source?${at}`,
    cached: true,
  },
  {
    id: "chandra-source",
    kind: "proxy",
    name: "Chandra Source Catalog",
    provider: "Chandra X-ray Center",
    usedFor: "X-ray source summary on the object page",
    path: `/api/chandra-source?${at}`,
    cached: true,
  },
  {
    id: "desi-target",
    kind: "proxy",
    name: "DESI targets",
    provider: "NOIRLab Astro Data Lab",
    usedFor: "DESI spectrum link on the object page",
    path: `/api/desi-target?${at}`,
    cached: true,
  },
];

/** The mirror answering for `check`, or null when it's the primary host. */
export function mirrorInUse(
  check: Pick<UpstreamCheck, "primaryHost">,
  host: string | null | undefined
): string | null {
  if (!host || !check.primaryHost) return null;
  return host === check.primaryHost ? null : host;
}

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
