import { CONE_SEARCH_MAX_NEIGHBORS } from "@/app/lib/constants/search";

import { catalogsOf, type NearbySource, sigmaRatio } from "./shared";

/**
 * How clear-cut the nearest match of one catalog is, judged only from the
 * typical (survey-wide) positional error — never per-source errors.
 *  - none: nothing within the search radius
 *  - unique: nearest within 2σ, no other source within 2σ
 *  - ambiguous: two or more sources within 2σ
 *  - offset: even the nearest lies beyond 2σ
 */
export type MatchVerdict = "none" | "unique" | "ambiguous" | "offset";

export interface CatalogSummary {
  slug: string;
  /** Search radius in arcseconds, if the catalog was searched. */
  radius?: number;
  /** Matches of this catalog, nearest first. */
  sources: NearbySource[];
  best?: NearbySource;
  /** Sources within 2σ of the searched position (typical error). */
  within2Sigma: number;
  verdict: MatchVerdict;
  /** The cone search stopped at its neighbour limit. */
  capped: boolean;
}

export function summarizeCatalogs(
  sources: NearbySource[],
  radii: Record<string, number>
): CatalogSummary[] {
  return catalogsOf(radii, sources).map((slug) => {
    const own = sources
      .filter((s) => s.slug === slug)
      .sort((a, b) => a.sepArcsec - b.sepArcsec);
    const within2Sigma = own.filter((s) => sigmaRatio(s) <= 2).length;
    const verdict: MatchVerdict =
      own.length === 0
        ? "none"
        : within2Sigma === 0
          ? "offset"
          : within2Sigma === 1
            ? "unique"
            : "ambiguous";
    return {
      slug,
      radius: radii[slug],
      sources: own,
      best: own[0],
      within2Sigma,
      verdict,
      capped: own.length >= CONE_SEARCH_MAX_NEIGHBORS,
    };
  });
}

export const VERDICT_TEXT: Record<MatchVerdict, string> = {
  none: "No match",
  unique: "Unique within 2σ",
  ambiguous: "Ambiguous",
  offset: "None within 2σ",
};

export const VERDICT_CLASSES: Record<MatchVerdict, string> = {
  none: "border-border text-neutral-400",
  unique: "border-green-500/40 bg-green-500/10 text-green-300",
  ambiguous: "border-amber-500/40 bg-amber-500/10 text-amber-300",
  offset: "border-border bg-foreground/5 text-neutral-300",
};

/** Caveat shown next to every σ verdict. */
export const SIGMA_CAVEAT =
  "σ uses each survey's typical positional error combined with 0.3″ for the searched position — not per-source errors — so treat the verdicts as a guide.";
