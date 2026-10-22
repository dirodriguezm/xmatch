/**
 * "Random object": a uniformly random sky position, then a Gaia cone search
 * with a growing radius until something turns up.
 */

import {
  FEATURED_OBJECTS,
  type FeaturedObject,
} from "@/app/lib/constants/featured";
import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

/** Radii tried in turn, in arcsec; capped by what the backend answers fast. */
export const RANDOM_RADII_ARCSEC = [10, 30, 60, MAX_RADIUS_ARCSEC].filter(
  (r) => r <= MAX_RADIUS_ARCSEC
);

/** Sky positions tried before giving up and falling back to a featured one. */
export const RANDOM_MAX_POSITIONS = 3;

export interface RandomHit {
  objectId: string;
  catalog: string;
  ra: number;
  dec: number;
}

/**
 * Uniform point on the celestial sphere (equal probability per steradian, so
 * declination is drawn via arcsin rather than uniformly).
 */
export function randomSkyPosition(rand: () => number = Math.random): {
  ra: number;
  dec: number;
} {
  const ra = rand() * 360;
  const dec = (Math.asin(2 * rand() - 1) * 180) / Math.PI;
  return { ra, dec };
}

export function pickRandomFeatured(
  rand: () => number = Math.random,
  list: FeaturedObject[] = FEATURED_OBJECTS
): FeaturedObject {
  return list[Math.min(list.length - 1, Math.floor(rand() * list.length))];
}

interface ConeRow {
  id: string;
  ra: number;
  dec: number;
  distance?: number;
}

/** First (closest) object in a conesearch proxy response, if any. */
export function firstConeHit(
  payload: unknown,
  catalog = "gaia"
): RandomHit | null {
  if (!Array.isArray(payload)) return null;
  for (const group of payload) {
    const rows = (group as { data?: ConeRow[] })?.data;
    if (!Array.isArray(rows) || rows.length === 0) continue;
    const best = [...rows].sort(
      (a, b) => (a.distance ?? 0) - (b.distance ?? 0)
    )[0];
    if (typeof best?.id !== "string") continue;
    return { objectId: best.id, catalog, ra: best.ra, dec: best.dec };
  }
  return null;
}

/** Search random positions with growing radii; null if nothing was found. */
export async function findRandomGaiaObject(
  fetcher: typeof fetch = fetch,
  rand: () => number = Math.random
): Promise<RandomHit | null> {
  for (let attempt = 0; attempt < RANDOM_MAX_POSITIONS; attempt++) {
    const { ra, dec } = randomSkyPosition(rand);
    for (const radius of RANDOM_RADII_ARCSEC) {
      const params = new URLSearchParams({
        ra: ra.toFixed(6),
        dec: dec.toFixed(6),
        radius: String(radius),
        catalog: "gaia",
        nneighbor: "1",
      });
      try {
        const res = await fetcher(`/api/conesearch?${params}`);
        if (!res.ok) break;
        const hit = firstConeHit(await res.json());
        if (hit) return hit;
      } catch {
        break;
      }
    }
  }
  return null;
}
