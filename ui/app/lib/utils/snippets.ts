/**
 * Copy-as-code: turn an XWave request into equivalent curl / Python / JS
 * snippets that call the public REST API directly (not the Next.js proxy).
 */

import { API_BASE_URL } from "@/app/lib/api/client";

export type SnippetLanguage = "curl" | "python" | "javascript";

export const SNIPPET_LANGUAGES: { key: SnippetLanguage; label: string }[] = [
  { key: "curl", label: "curl" },
  { key: "python", label: "Python" },
  { key: "javascript", label: "JavaScript" },
];

export type ApiRequest =
  | {
      method: "GET";
      /** Path under /v1, e.g. "/conesearch". */
      path: string;
      params: Record<string, string | number>;
    }
  | {
      method: "POST";
      path: string;
      body: Record<string, unknown>;
    };

export function coneSearchRequest(opts: {
  ra: number;
  dec: number;
  /** Radius in arcsec. */
  radius: number;
  catalog?: string;
  nneighbor?: number;
}): ApiRequest {
  const params: Record<string, string | number> = {
    ra: opts.ra,
    dec: opts.dec,
    radius: opts.radius,
    catalog: opts.catalog ?? "all",
  };
  if (opts.nneighbor != null) params.nneighbor = opts.nneighbor;
  return { method: "GET", path: "/conesearch", params };
}

export function metadataRequest(id: string, catalog: string): ApiRequest {
  return { method: "GET", path: "/metadata", params: { id, catalog } };
}

export function lightcurveRequest(opts: {
  ra: number;
  dec: number;
  radius: number;
}): ApiRequest {
  return {
    method: "GET",
    path: "/lightcurve",
    params: { ra: opts.ra, dec: opts.dec, radius: opts.radius },
  };
}

export function bulkConeSearchRequest(opts: {
  ra: number[];
  dec: number[];
  radius: number;
  catalog?: string;
  nneighbor?: number;
}): ApiRequest {
  return {
    method: "POST",
    path: "/bulk-conesearch",
    body: {
      ra: opts.ra,
      dec: opts.dec,
      radius: opts.radius,
      catalog: opts.catalog ?? "all",
      nneighbor: opts.nneighbor ?? 1,
    },
  };
}

function queryString(params: Record<string, string | number>): string {
  return new URLSearchParams(
    Object.entries(params).map(([k, v]) => [k, String(v)])
  ).toString();
}

export function requestUrl(req: ApiRequest): string {
  const base = `${API_BASE_URL}${req.path}`;
  return req.method === "GET" ? `${base}?${queryString(req.params)}` : base;
}

function pyLiteral(value: unknown): string {
  if (value === null || value === undefined) return "None";
  if (typeof value === "boolean") return value ? "True" : "False";
  if (typeof value === "number") return String(value);
  if (typeof value === "string") return JSON.stringify(value);
  if (Array.isArray(value)) return `[${value.map(pyLiteral).join(", ")}]`;
  return `{${Object.entries(value as Record<string, unknown>)
    .map(([k, v]) => `${JSON.stringify(k)}: ${pyLiteral(v)}`)
    .join(", ")}}`;
}

function pyDict(obj: Record<string, unknown>, indent = "    "): string {
  const lines = Object.entries(obj).map(
    ([k, v]) => `${indent}${JSON.stringify(k)}: ${pyLiteral(v)},`
  );
  return `{\n${lines.join("\n")}\n}`;
}

function jsObject(obj: Record<string, unknown>): string {
  return JSON.stringify(obj, null, 2);
}

export function buildSnippet(req: ApiRequest, lang: SnippetLanguage): string {
  const base = `${API_BASE_URL}${req.path}`;
  switch (lang) {
    case "curl":
      return req.method === "GET"
        ? `curl -G "${base}" \\\n${Object.entries(req.params)
            .map(([k, v]) => `  --data-urlencode "${k}=${v}"`)
            .join(" \\\n")}`
        : `curl -X POST "${base}" \\\n  -H "Content-Type: application/json" \\\n  -d '${JSON.stringify(req.body)}'`;
    case "python":
      return req.method === "GET"
        ? `import requests

params = ${pyDict(req.params)}
resp = requests.get("${base}", params=params, timeout=60)
resp.raise_for_status()
data = resp.json() if resp.status_code != 204 else []
print(data)`
        : `import requests

payload = ${pyDict(req.body)}
resp = requests.post("${base}", json=payload, timeout=120)
resp.raise_for_status()
data = resp.json() if resp.status_code != 204 else []
print(data)`;
    case "javascript":
      return req.method === "GET"
        ? `const params = new URLSearchParams(${jsObject(
            Object.fromEntries(
              Object.entries(req.params).map(([k, v]) => [k, String(v)])
            )
          )});
const resp = await fetch(\`${base}?\${params}\`);
const data = resp.status === 204 ? [] : await resp.json();
console.log(data);`
        : `const resp = await fetch("${base}", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify(${jsObject(req.body)}),
});
const data = resp.status === 204 ? [] : await resp.json();
console.log(data);`;
  }
}
