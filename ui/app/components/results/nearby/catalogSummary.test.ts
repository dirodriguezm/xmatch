import { describe, expect, it } from "vitest";

import { listOrder, summarizeCatalogs } from "./catalogSummary";
import { combinedSigma, type NearbySource } from "./shared";

const src = (slug: string, sepArcsec: number): NearbySource => ({
  key: `${slug}:${sepArcsec}`,
  objectId: `${slug}-${sepArcsec}`,
  slug,
  catalog: slug,
  ra: 0,
  dec: 0,
  sepArcsec,
  pa: 0,
  east: 0,
  north: sepArcsec,
  sigma: combinedSigma(slug),
  href: "#",
});

describe("summarizeCatalogs", () => {
  const radii = { allwise: 3, gaia: 3, erosita: 20 };

  it("keeps every searched catalog, matched or not", () => {
    const groups = summarizeCatalogs([src("allwise", 0.22)], radii);
    expect(groups.map((g) => [g.slug, g.sources.length])).toEqual([
      ["allwise", 1],
      ["gaia", 0],
      ["erosita", 0],
    ]);
    expect(groups[1]?.verdict).toBe("none");
    expect(groups[1]?.radius).toBe(3);
  });

  it("marks only the catalogs whose search failed", () => {
    const groups = summarizeCatalogs([], radii, { gaia: "Timed out" });
    expect(groups.find((g) => g.slug === "gaia")?.error).toBe("Timed out");
    expect(groups.find((g) => g.slug === "erosita")).not.toHaveProperty(
      "error"
    );
  });
  it("marks a catalog still searching, so it never reads as 0 found", () => {
    const groups = summarizeCatalogs([src("allwise", 0.22)], radii, {}, [
      "gaia",
    ]);
    expect(groups.find((g) => g.slug === "gaia")?.pending).toBe(true);
    expect(groups.find((g) => g.slug === "erosita")).not.toHaveProperty(
      "pending"
    );
  });
});

describe("listOrder", () => {
  it("puts matches first, then failures, then empty catalogs", () => {
    const groups = summarizeCatalogs(
      [src("erosita", 4)],
      { allwise: 3, gaia: 3, erosita: 20 },
      { gaia: "" }
    );
    expect(listOrder(groups).map((g) => g.slug)).toEqual([
      "erosita",
      "gaia",
      "allwise",
    ]);
  });

  it("keeps a catalog being retried where its failure was", () => {
    const groups = summarizeCatalogs(
      [src("erosita", 4)],
      { allwise: 3, gaia: 3, erosita: 20 },
      {},
      ["gaia"]
    );
    expect(listOrder(groups).map((g) => g.slug)).toEqual([
      "erosita",
      "gaia",
      "allwise",
    ]);
  });
});
