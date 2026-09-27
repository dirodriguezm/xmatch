import { describe, expect, it } from "vitest";

import { buildGlossary, filterGlossary, groupByCatalog } from "./glossary";

const SOURCE = {
  gaia: {
    parallax: "Parallax (mas)",
    parallax_error: "Parallax uncertainty (mas)",
    ruwe: "Renormalised unit weight error",
  },
  allwise: { w1mpro: "W1 (3.4 µm) magnitude" },
};

describe("glossary", () => {
  const entries = buildGlossary(SOURCE);

  it("flattens every catalog field", () => {
    expect(entries).toHaveLength(4);
    expect(entries[0]).toEqual({
      catalog: "gaia",
      field: "parallax",
      description: "Parallax (mas)",
    });
  });

  it("filters on all terms, treating underscores as spaces", () => {
    expect(
      filterGlossary(entries, "parallax error").map((e) => e.field)
    ).toEqual(["parallax_error"]);
    expect(filterGlossary(entries, "W1").map((e) => e.field)).toEqual([
      "w1mpro",
    ]);
    expect(filterGlossary(entries, "allwise")).toHaveLength(1);
    expect(filterGlossary(entries, "  ")).toHaveLength(4);
    expect(filterGlossary(entries, "nothing-here")).toHaveLength(0);
  });

  it("groups by catalog in order", () => {
    const groups = groupByCatalog(entries);
    expect(groups.map(([c]) => c)).toEqual(["gaia", "allwise"]);
    expect(groups[0][1]).toHaveLength(3);
  });

  it("builds from the real descriptions", () => {
    expect(buildGlossary().length).toBeGreaterThan(50);
  });
});
