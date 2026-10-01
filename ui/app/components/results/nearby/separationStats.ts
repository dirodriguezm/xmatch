import { CONE_SEARCH_MAX_NEIGHBORS } from "@/app/lib/constants/search";

import type { NearbySource } from "./shared";

/** At or below this many matches the filter tools start collapsed. */
export const SPARSE_MAX = 20;

/** Annulus sources needed before a field density is worth quoting. */
const MIN_FIELD_SOURCES = 5;

/**
 * Approximate surface density of unrelated sources around the target, from
 * the outer half (in radius) of a catalog's cone, where a true counterpart is
 * unlikely to sit. When the catalog hit the nearest-N cap the cone is only
 * complete out to the farthest returned source, so that is used instead.
 */
export interface FieldDensity {
  /** Sources per square arcsecond. */
  perSqArcsec: number;
  /** Radius out to which the cone is complete, in arcseconds. */
  completeTo: number;
  /** Sources the estimate is based on. */
  annulusCount: number;
  capped: boolean;
}

export function fieldDensity(
  sorted: NearbySource[],
  radius: number | undefined
): FieldDensity | null {
  if (sorted.length === 0) return null;
  const capped = sorted.length >= CONE_SEARCH_MAX_NEIGHBORS;
  const farthest = sorted[sorted.length - 1].sepArcsec;
  const completeTo = capped || radius === undefined ? farthest : radius;
  const inner = completeTo / 2;
  const annulusCount = sorted.filter(
    (s) => s.sepArcsec > inner && s.sepArcsec <= completeTo
  ).length;
  if (annulusCount < MIN_FIELD_SOURCES) return null;
  const area = Math.PI * (completeTo ** 2 - inner ** 2);
  return { perSqArcsec: annulusCount / area, completeTo, annulusCount, capped };
}

/** Expected unrelated sources with lo < separation ≤ hi. */
export function expectedByChance(d: FieldDensity, lo: number, hi: number) {
  const a = Math.min(lo, d.completeTo);
  const b = Math.min(hi, d.completeTo);
  return Math.max(0, d.perSqArcsec * Math.PI * (b * b - a * a));
}

/** Chance of at least one unrelated source within r (Poisson). */
export function chanceWithin(d: FieldDensity, r: number) {
  return 1 - Math.exp(-d.perSqArcsec * Math.PI * r * r);
}

export function formatChance(p: number): string {
  if (p < 0.001) return "<0.1%";
  if (p < 0.1) return `${Number((p * 100).toFixed(1))}%`;
  return `${Math.round(p * 100)}%`;
}

export interface CatalogGroup {
  slug: string;
  /** Matches of this catalog, nearest first. */
  sources: NearbySource[];
  radius: number | undefined;
  density: FieldDensity | null;
}

export function groupByCatalog(
  catalogs: string[],
  sources: NearbySource[],
  radii: Record<string, number>
): CatalogGroup[] {
  return catalogs.map((slug) => {
    const list = sources
      .filter((s) => s.slug === slug)
      .sort((a, b) => a.sepArcsec - b.sepArcsec);
    return {
      slug,
      sources: list,
      radius: radii[slug],
      density: fieldDensity(list, radii[slug]),
    };
  });
}

export interface SeparationFilterState {
  /** Separation range in arcseconds, or null for everything. */
  range: [number, number] | null;
  /** Keep only the nearest N per catalog (within the range), or null. */
  nearest: number | null;
}

export const NO_FILTER: SeparationFilterState = { range: null, nearest: null };

export const isFiltered = (f: SeparationFilterState) =>
  f.range !== null || f.nearest !== null;

export function applyFilter(
  sources: NearbySource[],
  f: SeparationFilterState
): NearbySource[] {
  const inRange = f.range
    ? sources.filter(
        (s) => s.sepArcsec >= f.range![0] && s.sepArcsec <= f.range![1]
      )
    : sources;
  if (f.nearest === null) return inRange;
  const seen = new Map<string, number>();
  // Sources arrive nearest first, so the first N of each catalog are kept.
  return [...inRange]
    .sort((a, b) => a.sepArcsec - b.sepArcsec)
    .filter((s) => {
      const n = (seen.get(s.slug) ?? 0) + 1;
      seen.set(s.slug, n);
      return n <= f.nearest!;
    });
}

const NICE = [
  0.5, 1, 2, 3, 4, 5, 6, 8, 10, 12, 15, 20, 25, 30, 40, 50, 60, 90, 120, 180,
  300,
];

/** A round map radius at or above r. */
export function niceExtent(r: number): number {
  return NICE.find((n) => n >= r) ?? Math.ceil(r / 60) * 60;
}

/** How far a catalog's matches can be trusted to be complete. */
export function groupExtent(g: CatalogGroup): number {
  const farthest = g.sources[g.sources.length - 1]?.sepArcsec ?? 0;
  if (g.sources.length >= CONE_SEARCH_MAX_NEIGHBORS) return farthest;
  return g.radius ?? farthest;
}

/**
 * Right edge of the separation axis: the widest cone of a catalog that
 * returned anything — an empty 60″ eROSITA cone would squash 3″ matches.
 */
export function separationAxisMax(groups: CatalogGroup[]): number {
  const withMatches = groups.filter((g) => g.sources.length > 0);
  return Math.max(1, ...withMatches.map(groupExtent));
}
