import { CONE_SEARCH_MAX_NEIGHBORS } from "@/app/lib/constants/search";

import type { NearbySource } from "./shared";

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

/** Chance of at least one unrelated source within r (Poisson). */
export function chanceWithin(d: FieldDensity, r: number) {
  return 1 - Math.exp(-d.perSqArcsec * Math.PI * r * r);
}

export function formatChance(p: number): string {
  if (p < 0.001) return "<0.1%";
  if (p < 0.1) return `${Number((p * 100).toFixed(1))}%`;
  return `${Math.round(p * 100)}%`;
}
