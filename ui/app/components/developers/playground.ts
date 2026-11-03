/**
 * Pure helpers for the /developers API playground: endpoint definitions,
 * form state → ApiRequest, and response formatting.
 */

import { parseCoordinates } from "@/app/lib/utils/coordinates";
import {
  type ApiRequest,
  bulkConeSearchRequest,
  coneSearchRequest,
  lightcurveRequest,
  metadataRequest,
} from "@/app/lib/utils/snippets";

export type EndpointKey = "conesearch" | "metadata" | "lightcurve" | "bulk";

export interface EndpointInfo {
  key: EndpointKey;
  label: string;
  method: "GET" | "POST";
  path: string;
  description: string;
}

export const ENDPOINTS: EndpointInfo[] = [
  {
    key: "conesearch",
    label: "Cone search",
    method: "GET",
    path: "/conesearch",
    description: "Sources within a radius of one position, grouped by catalog.",
  },
  {
    key: "metadata",
    label: "Metadata",
    method: "GET",
    path: "/metadata",
    description: "The full catalog row for one source id.",
  },
  {
    key: "lightcurve",
    label: "Light curve",
    method: "GET",
    path: "/lightcurve",
    description: "Time-series photometry (ZTF, NEOWISE) near a position.",
  },
  {
    key: "bulk",
    label: "Bulk cone search",
    method: "POST",
    path: "/bulk-conesearch",
    description: "Many positions in one request with a shared radius.",
  },
];

export interface PlaygroundState {
  ra: number;
  dec: number;
  /** Arcseconds. */
  radius: number;
  catalog: string;
  nneighbor: number;
  getMetadata: boolean;
  metadataId: string;
  metadataCatalog: string;
  lightcurveRadius: number;
  lightcurveCatalog: string;
  /** One "ra dec" (degrees or sexagesimal) per line. */
  bulkPositions: string;
  bulkRadius: number;
}

/** M31, plus a known AllWISE source at its core. */
export const DEFAULT_STATE: PlaygroundState = {
  ra: 10.6847,
  dec: 41.269,
  radius: 10,
  catalog: "all",
  nneighbor: 5,
  getMetadata: false,
  metadataId: "0098p408_ac51-043708",
  metadataCatalog: "allwise",
  lightcurveRadius: 5,
  lightcurveCatalog: "all",
  bulkPositions: [
    "10.6847 41.2690  # M31",
    "83.6331 22.0145  # Crab Nebula",
    "05:34:31.9 +22:00:52  # Crab Pulsar",
  ].join("\n"),
  bulkRadius: 5,
};

export interface ParsedPositions {
  ra: number[];
  dec: number[];
  /** 1-based line numbers that could not be parsed. */
  invalid: number[];
}

/** Parse one position per line; `#` starts a comment, blank lines are skipped. */
export function parsePositions(text: string): ParsedPositions {
  const out: ParsedPositions = { ra: [], dec: [], invalid: [] };
  text.split(/\r?\n/).forEach((raw, i) => {
    const line = raw.replace(/#.*$/, "").trim();
    if (!line) return;
    const coords = parseCoordinates(line);
    if (!coords) {
      out.invalid.push(i + 1);
      return;
    }
    out.ra.push(round(coords.ra));
    out.dec.push(round(coords.dec));
  });
  return out;
}

function round(v: number, digits = 6): number {
  const f = 10 ** digits;
  return Math.round(v * f) / f;
}

/** Build the request the form describes. */
export function buildPlaygroundRequest(
  endpoint: EndpointKey,
  s: PlaygroundState
): ApiRequest {
  switch (endpoint) {
    case "conesearch": {
      const req = coneSearchRequest({
        ra: s.ra,
        dec: s.dec,
        radius: s.radius,
        catalog: s.catalog,
        nneighbor: s.nneighbor,
      });
      if (s.getMetadata && req.method === "GET")
        req.params.getMetadata = "true";
      return req;
    }
    case "metadata":
      return metadataRequest(s.metadataId.trim(), s.metadataCatalog);
    case "lightcurve": {
      const req = lightcurveRequest({
        ra: s.ra,
        dec: s.dec,
        radius: s.lightcurveRadius,
      });
      if (s.lightcurveCatalog !== "all" && req.method === "GET") {
        req.params.catalog = s.lightcurveCatalog;
      }
      return req;
    }
    case "bulk": {
      const { ra, dec } = parsePositions(s.bulkPositions);
      return bulkConeSearchRequest({
        ra,
        dec,
        radius: s.bulkRadius,
        catalog: s.catalog,
        nneighbor: s.nneighbor,
      });
    }
  }
}

/** Same request routed through this app's Next.js proxy (fallback for CORS). */
export function proxyUrl(req: ApiRequest): string {
  if (req.method === "POST") return `/api${req.path}`;
  const qs = new URLSearchParams(
    Object.entries(req.params).map(([k, v]) => [k, String(v)])
  );
  return `/api${req.path}?${qs}`;
}

/** Human-readable byte size. */
export function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(2)} MB`;
}

/** Number of top-level results, for the status line. */
export function summarizeResult(data: unknown): string | null {
  if (Array.isArray(data)) {
    const rows = data.reduce<number>(
      (n, g) =>
        n +
        (g &&
        typeof g === "object" &&
        Array.isArray((g as { data?: unknown }).data)
          ? (g as { data: unknown[] }).data.length
          : 1),
      0
    );
    return `${data.length} group${data.length === 1 ? "" : "s"}, ${rows} row${rows === 1 ? "" : "s"}`;
  }
  if (data && typeof data === "object") {
    const obj = data as Record<string, unknown>;
    const arrays = Object.entries(obj).filter(([, v]) => Array.isArray(v));
    if (arrays.length > 0) {
      return arrays
        .map(([k, v]) => `${(v as unknown[]).length} ${k.replace(/_/g, " ")}`)
        .join(", ");
    }
    return `${Object.keys(obj).length} fields`;
  }
  return null;
}
