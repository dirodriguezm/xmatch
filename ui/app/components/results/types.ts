import type { MagValue } from "@/app/lib/constants/photometry";

export interface CrossmatchResult {
  key: string;
  objectId: string;
  ra: number;
  dec: number;
  angularDistance: number;
  catalog: string;
  ipix?: number;
  /**
   * Canonical request slug, used only as a join key for photometry. Optional
   * because the object detail page builds a CrossmatchResult by hand.
   */
  catalogSlug?: string;
}

/** A search result once photometry has been merged in. */
export interface EnrichedResult extends CrossmatchResult {
  photometry?: MagValue;
}
