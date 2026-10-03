/**
 * Helpers for share links, embeds and Open Graph cards.
 */

import { API_BASE_URL } from "@/app/lib/api/client";
import { SITE_NAME, SITE_URL } from "@/app/lib/constants/site";

/** Join an app-relative path onto an origin, without doubling slashes. */
export function absoluteUrl(path: string, origin: string = SITE_URL): string {
  const base = origin.replace(/\/+$/, "");
  const rel = path.startsWith("/") ? path : `/${path}`;
  return `${base}${rel}`;
}

function escapeHtmlAttr(value: string): string {
  return value
    .replace(/&/g, "&amp;")
    .replace(/"/g, "&quot;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}

/** Copy-paste `<iframe>` snippet for an embeddable view. */
export function buildEmbedHtml(
  src: string,
  title: string,
  height = 420
): string {
  return `<iframe src="${escapeHtmlAttr(src)}" title="${escapeHtmlAttr(`${title} on ${SITE_NAME}`)}" width="100%" height="${height}" loading="lazy" style="border:0;border-radius:8px" allowfullscreen></iframe>`;
}

export interface ShareTarget {
  id: "x" | "bluesky" | "email";
  label: string;
  href: string;
}

/** Intent URLs for posting a link to X, Bluesky or email. */
export function buildShareTargets(url: string, title: string): ShareTarget[] {
  const text = `${title} on ${SITE_NAME}`;
  return [
    {
      id: "x",
      label: "X",
      href: `https://x.com/intent/post?${new URLSearchParams({ text, url })}`,
    },
    {
      id: "bluesky",
      label: "Bluesky",
      href: `https://bsky.app/intent/compose?${new URLSearchParams({
        text: `${text} ${url}`,
      })}`,
    },
    {
      id: "email",
      label: "Email",
      href: `mailto:?subject=${encodeURIComponent(text)}&body=${encodeURIComponent(url)}`,
    },
  ];
}

export interface CatalogGuess {
  /** Search-catalog slug: `gaia`, `allwise`, `erosita`. */
  slug: "gaia" | "allwise" | "erosita";
  /** Display label including release, e.g. "Gaia DR3". */
  label: string;
}

/** Infer the catalog from an object identifier's prefix, when it has one. */
export function catalogFromObjectId(objectId: string): CatalogGuess | null {
  const id = objectId.trim();
  if (/^gaia\s*dr3\b/i.test(id)) return { slug: "gaia", label: "Gaia DR3" };
  if (/^1eRASS\b/i.test(id)) return { slug: "erosita", label: "eROSITA" };
  if (/^J\d{6}\.\d{2}[+-]\d{6}\.\d$/.test(id))
    return { slug: "allwise", label: "AllWISE" };
  return null;
}

const CATALOG_LABELS: Record<string, string> = {
  gaia: "Gaia DR3",
  allwise: "AllWISE",
  erosita: "eROSITA",
};

/** Display label for a catalog slug, falling back to the slug itself. */
export function catalogLabel(slug: string): string {
  return CATALOG_LABELS[slug] ?? slug;
}

/** Short "RA 10.67964°, Dec +41.26972°" string for cards and descriptions. */
export function formatRaDec(ra: number, dec: number, decimals = 5): string {
  const sign = dec >= 0 ? "+" : "−";
  return `RA ${ra.toFixed(decimals)}°, Dec ${sign}${Math.abs(dec).toFixed(decimals)}°`;
}

export interface PhotometryRow {
  band: string;
  value: number;
  error?: number;
  unit: string;
}

const PHOTOMETRY_FIELDS: Record<
  string,
  { band: string; field: string; errField?: string; unit: string }[]
> = {
  gaia: [
    { band: "G", field: "phot_g_mean_mag", unit: "mag" },
    { band: "BP", field: "phot_bp_mean_mag", unit: "mag" },
    { band: "RP", field: "phot_rp_mean_mag", unit: "mag" },
  ],
  allwise: [
    { band: "W1", field: "w1mpro", errField: "w1sigmpro", unit: "mag" },
    { band: "W2", field: "w2mpro", errField: "w2sigmpro", unit: "mag" },
    { band: "W3", field: "w3mpro", errField: "w3sigmpro", unit: "mag" },
    { band: "W4", field: "w4mpro", errField: "w4sigmpro", unit: "mag" },
    { band: "J", field: "j_m_2mass", errField: "j_msig_2mass", unit: "mag" },
    { band: "H", field: "h_m_2mass", errField: "h_msig_2mass", unit: "mag" },
    { band: "Ks", field: "k_m_2mass", errField: "k_msig_2mass", unit: "mag" },
  ],
};

function finiteNonZero(v: unknown): number | undefined {
  return typeof v === "number" && Number.isFinite(v) && v !== 0 ? v : undefined;
}

/**
 * Key photometry for a metadata record. Zero values are treated as missing
 * (the API fills absent measurements with 0). Catalogs without a known band
 * list (eROSITA) fall back to any `*flux*` fields.
 */
export function pickPhotometry(
  catalog: string,
  record: Record<string, unknown> | null | undefined
): PhotometryRow[] {
  if (!record) return [];
  const fields = PHOTOMETRY_FIELDS[catalog];
  if (fields) {
    return fields.flatMap(({ band, field, errField, unit }) => {
      const value = finiteNonZero(record[field]);
      if (value === undefined) return [];
      const error = errField ? finiteNonZero(record[errField]) : undefined;
      return [{ band, value, error, unit }];
    });
  }
  return Object.keys(record)
    .filter((k) => /flux/i.test(k) && !/err/i.test(k))
    .flatMap((k) => {
      const value = finiteNonZero(record[k]);
      return value === undefined ? [] : [{ band: k, value, unit: "erg/s/cm²" }];
    })
    .slice(0, 6);
}

/**
 * Server-side lookup of an object's metadata row for OG cards and page
 * metadata. Never throws: returns null on timeout, API error or missing
 * coordinates so previews degrade to the id alone.
 */
export async function fetchObjectSummary(
  objectId: string,
  catalog: string,
  timeoutMs = 2000
): Promise<{
  ra: number;
  dec: number;
  record: Record<string, unknown>;
} | null> {
  try {
    const params = new URLSearchParams({ id: objectId, catalog });
    const res = await fetch(`${API_BASE_URL}/metadata?${params}`, {
      signal: AbortSignal.timeout(timeoutMs),
      next: { revalidate: 86400 },
    });
    if (!res.ok || res.status === 204) return null;
    const record = (await res.json()) as Record<string, unknown> | null;
    if (
      !record ||
      typeof record.ra !== "number" ||
      typeof record.dec !== "number"
    )
      return null;
    return { ra: record.ra, dec: record.dec, record };
  } catch {
    return null;
  }
}
