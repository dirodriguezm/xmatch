import type { MagValue } from "@/app/lib/constants/photometry";
import {
  convertRadiusToArcsec,
  decodeCatalogRadii,
} from "@/app/lib/constants/search";
import { positionAngle } from "@/app/lib/utils/coordinates";
import { buildObjectUrl } from "@/app/lib/utils/urls";

import type { EnrichedResult } from "../ResultsTable";

/** A match placed on the sky relative to the searched position. */
export interface NearbySource {
  key: string;
  objectId: string;
  /** Request slug (`allwise`, `gaia`, `erosita`). */
  slug: string;
  /** Catalog value the API returned, used for object-page links. */
  catalog: string;
  ra: number;
  dec: number;
  sepArcsec: number;
  /** Degrees east of north. */
  pa: number;
  /** Offset east (+) and north (+) of the target, in arcseconds. */
  east: number;
  north: number;
  /** Combined 1σ positional uncertainty of the pair, in arcseconds. */
  sigma: number;
  photometry?: MagValue;
  href: string;
}

/**
 * Typical 1σ astrometric uncertainty per catalog, in arcseconds. The cone
 * search does not return per-source errors, so these are survey-wide medians:
 * Gaia DR3 (incl. proper motion to J2000 for most sources), AllWISE (S/N≈10),
 * eRASS1 (median POS_ERR).
 */
export const CATALOG_SIGMA_ARCSEC: Record<string, number> = {
  gaia: 0.1,
  allwise: 0.3,
  erosita: 4.5,
};

/** Assumed 1σ uncertainty of the searched position itself. */
export const TARGET_SIGMA_ARCSEC = 0.3;

const DEFAULT_SIGMA = 1;

export function combinedSigma(slug: string): number {
  const cat = CATALOG_SIGMA_ARCSEC[slug] ?? DEFAULT_SIGMA;
  return Math.hypot(cat, TARGET_SIGMA_ARCSEC);
}

export type Agreement = "consistent" | "marginal" | "offset";

/** How the separation compares with the positional uncertainty. */
export function agreement(source: NearbySource): Agreement {
  const n = source.sepArcsec / source.sigma;
  if (n <= 2) return "consistent";
  if (n <= 3) return "marginal";
  return "offset";
}

export const AGREEMENT_LABEL: Record<Agreement, string> = {
  consistent: "within 2σ",
  marginal: "2–3σ",
  offset: "> 3σ",
};

export const AGREEMENT_TEXT_CLASSES: Record<Agreement, string> = {
  consistent: "text-green-400",
  marginal: "text-amber-400",
  offset: "text-neutral-400",
};

export const sigmaRatio = (s: NearbySource) => s.sepArcsec / s.sigma;

export function toNearbySources(
  data: EnrichedResult[],
  target: { ra: number; dec: number }
): NearbySource[] {
  return data
    .map((r) => {
      const slug = (r.catalogSlug ?? r.catalog).toLowerCase();
      const pa = positionAngle(target.ra, target.dec, r.ra, r.dec);
      const rad = (pa * Math.PI) / 180;
      return {
        key: `${slug}:${r.key}`,
        objectId: r.objectId,
        slug,
        catalog: r.catalog,
        ra: r.ra,
        dec: r.dec,
        sepArcsec: r.angularDistance,
        pa,
        east: r.angularDistance * Math.sin(rad),
        north: r.angularDistance * Math.cos(rad),
        sigma: combinedSigma(slug),
        photometry: r.photometry,
        href: buildObjectUrl(r.objectId, r.catalog),
      };
    })
    .sort((a, b) => a.sepArcsec - b.sepArcsec);
}

/** Search radius per enabled catalog slug, in arcseconds, from the URL value. */
export function radiiFromParam(catalogRadii: string): Record<string, number> {
  return Object.fromEntries(
    decodeCatalogRadii(catalogRadii)
      .filter((c) => c.enabled)
      .map((c) => [
        c.catalog.toLowerCase(),
        convertRadiusToArcsec(c.radius, c.unit),
      ])
  );
}

/** Catalogs in a stable order: searched ones first, then any extra in the data. */
export function catalogsOf(
  radii: Record<string, number>,
  sources: NearbySource[]
): string[] {
  return [...new Set([...Object.keys(radii), ...sources.map((s) => s.slug)])];
}

/** Largest search radius, falling back to the farthest source. */
export function maxRadius(
  radii: Record<string, number>,
  sources: NearbySource[]
): number {
  const r = Math.max(0, ...Object.values(radii));
  const s = Math.max(0, ...sources.map((x) => x.sepArcsec));
  return Math.max(r, s, 1);
}

// Tailwind needs full literal class names, so these live as lookup tables.
export const CATALOG_STROKE_CLASSES: Record<string, string> = {
  allwise: "stroke-purple-500",
  gaia: "stroke-blue-500",
  erosita: "stroke-pink-500",
};
export const CATALOG_FILL: Record<string, string> = {
  allwise: "fill-purple-500",
  gaia: "fill-blue-500",
  erosita: "fill-pink-500",
};
export const CATALOG_BG: Record<string, string> = {
  allwise: "bg-purple-500",
  gaia: "bg-blue-500",
  erosita: "bg-pink-500",
};

export const strokeClass = (slug: string) =>
  CATALOG_STROKE_CLASSES[slug] ?? "stroke-neutral-400";
export const fillClass = (slug: string) =>
  CATALOG_FILL[slug] ?? "fill-neutral-400";
export const bgClass = (slug: string) => CATALOG_BG[slug] ?? "bg-neutral-400";

/** Shape per catalog, so colour is never the only cue. */
export type MarkShape = "circle" | "square" | "diamond" | "triangle";
export const CATALOG_SHAPE: Record<string, MarkShape> = {
  gaia: "circle",
  allwise: "square",
  erosita: "diamond",
};
export const shapeOf = (slug: string): MarkShape =>
  CATALOG_SHAPE[slug] ?? "triangle";

export const COMPASS = ["N", "NE", "E", "SE", "S", "SW", "W", "NW"];
export const compassPoint = (pa: number) => COMPASS[Math.round(pa / 45) % 8];

/** Arcseconds with two decimals (one from 10″), trailing zeros dropped. */
export function formatArcsec(v: number): string {
  return `${Number(v.toFixed(v >= 10 ? 1 : 2))}″`;
}

/**
 * Scroll a table row (antd sets `data-row-key`) into view inside the table's
 * own scroller only — never the page, so hovering the map cannot make the
 * window jump.
 */
export function scrollRowIntoView(container: HTMLElement | null, key: string) {
  const row = container?.querySelector<HTMLElement>(
    `tr[data-row-key="${CSS.escape(key)}"]`
  );
  const body = row?.closest<HTMLElement>(".ant-table-body");
  if (!row || !body) return;
  const top = row.offsetTop;
  const bottom = top + row.offsetHeight;
  if (top < body.scrollTop) body.scrollTo({ top, behavior: "smooth" });
  else if (bottom > body.scrollTop + body.clientHeight)
    body.scrollTo({ top: bottom - body.clientHeight, behavior: "smooth" });
}
