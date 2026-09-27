import {
  parseSesame,
  type ApiCatalogResult,
  type CatalogOrAll,
  type SesameResult,
} from "./helpers.js";

export interface ClientConfig {
  apiUrl: string;
  sesameUrl: string;
  timeoutMs: number;
  userAgent: string;
}

export function configFromEnv(env: NodeJS.ProcessEnv = process.env): ClientConfig {
  const timeout = Number(env.XWAVE_TIMEOUT_MS);
  return {
    apiUrl: (env.XWAVE_API_URL || "https://xwave-astro.udp.cl/v1").replace(/\/+$/, ""),
    sesameUrl: (env.XWAVE_SESAME_URL || "https://cds.unistra.fr/cgi-bin/nph-sesame/-oI/A").replace(/\/+$/, ""),
    timeoutMs: Number.isFinite(timeout) && timeout > 0 ? timeout : 30_000,
    userAgent: "xwave-mcp/0.1.0",
  };
}

export class ApiError extends Error {
  constructor(message: string, readonly status?: number) {
    super(message);
    this.name = "ApiError";
  }
}

/** fetch with an AbortController timeout. Returns null body for 204. */
async function request(
  cfg: ClientConfig,
  url: string,
  init: RequestInit = {},
): Promise<{ status: number; text: string | null }> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), cfg.timeoutMs);
  try {
    const res = await fetch(url, {
      ...init,
      signal: controller.signal,
      headers: { "User-Agent": cfg.userAgent, Accept: "application/json, text/plain", ...(init.headers ?? {}) },
    });
    if (res.status === 204) return { status: 204, text: null };
    const text = await res.text();
    if (!res.ok) {
      throw new ApiError(`${describeUrl(url)} returned HTTP ${res.status}: ${extractError(text)}`, res.status);
    }
    return { status: res.status, text };
  } catch (err) {
    if (err instanceof ApiError) throw err;
    if ((err as Error).name === "AbortError") {
      throw new ApiError(`${describeUrl(url)} timed out after ${cfg.timeoutMs} ms. Try a smaller radius or fewer positions.`);
    }
    throw new ApiError(`Could not reach ${describeUrl(url)}: ${(err as Error).message}`);
  } finally {
    clearTimeout(timer);
  }
}

function describeUrl(url: string): string {
  try {
    const u = new URL(url);
    return `${u.host}${u.pathname}`;
  } catch {
    return url;
  }
}

function extractError(text: string): string {
  try {
    const j = JSON.parse(text) as Record<string, unknown>;
    if (typeof j.Reason === "string") return `${j.Field ? `${j.Field}: ` : ""}${j.Reason}`;
    if (typeof j.error === "string") return j.error;
    if (typeof j.message === "string") return j.message;
  } catch {
    /* not JSON */
  }
  return text.slice(0, 300) || "(empty body)";
}

function parseJson<T>(text: string | null, fallback: T): T {
  if (text === null || text.trim() === "") return fallback;
  return JSON.parse(text) as T;
}

export class XWaveClient {
  constructor(readonly cfg: ClientConfig) {}

  async resolveName(name: string): Promise<SesameResult | null> {
    const { text } = await request(this.cfg, `${this.cfg.sesameUrl}?${encodeURIComponent(name.trim())}`);
    return text ? parseSesame(text, name.trim()) : null;
  }

  /** radius in arcsec. 204 → []. */
  async coneSearch(ra: number, dec: number, radiusArcsec: number, catalog: CatalogOrAll, nneighbor: number): Promise<ApiCatalogResult[]> {
    const q = new URLSearchParams({
      ra: String(ra),
      dec: String(dec),
      radius: String(radiusArcsec),
      catalog,
      nneighbor: String(nneighbor),
    });
    const { text } = await request(this.cfg, `${this.cfg.apiUrl}/conesearch?${q}`);
    return parseJson<ApiCatalogResult[]>(text, []);
  }

  /**
   * Bulk cone search. The radius is sent in ARCSECONDS: the Go struct comment
   * says degrees, but the live service treats it as arcsec (verified: 5 → match
   * at 0.22", 0.0014 → 204).
   */
  async bulkConeSearch(ra: number[], dec: number[], radiusArcsec: number, catalog: CatalogOrAll, nneighbor: number): Promise<ApiCatalogResult[]> {
    const { text } = await request(this.cfg, `${this.cfg.apiUrl}/bulk-conesearch`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ ra, dec, radius: radiusArcsec, catalog, nneighbor }),
    });
    return parseJson<ApiCatalogResult[]>(text, []);
  }

  async metadata(id: string, catalog: string): Promise<Record<string, unknown> | null> {
    const q = new URLSearchParams({ id, catalog });
    const { text } = await request(this.cfg, `${this.cfg.apiUrl}/metadata?${q}`);
    return parseJson<Record<string, unknown> | null>(text, null);
  }

  async lightcurve(ra: number, dec: number, radiusArcsec: number, nneighbor: number): Promise<LightcurveResponse> {
    const q = new URLSearchParams({ ra: String(ra), dec: String(dec), radius: String(radiusArcsec), nneighbor: String(nneighbor) });
    const { text } = await request(this.cfg, `${this.cfg.apiUrl}/lightcurve?${q}`);
    return parseJson<LightcurveResponse>(text, { detections: [], non_detections: [], forced_photometry: [] });
  }
}

export interface LightcurveEntry {
  catalog: string;
  id: string;
  object_id: string;
  mjd: number;
  mag: number;
  magerr: number;
  data?: Record<string, unknown>;
}

export interface LightcurveResponse {
  detections: LightcurveEntry[] | null;
  non_detections: LightcurveEntry[] | null;
  forced_photometry: LightcurveEntry[] | null;
}
