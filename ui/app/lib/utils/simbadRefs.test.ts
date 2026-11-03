import { describe, expect, it } from "vitest";

import {
  adsAbstractUrl,
  parseSimbadRefs,
  simbadRefsCountQuery,
  simbadRefsQuery,
} from "./simbadRefs";

describe("simbadRefsQuery", () => {
  it("lists an object's references newest first", () => {
    const q = simbadRefsQuery(1940765, 25);
    expect(q).toContain("SELECT TOP 25");
    expect(q).toContain("WHERE h.oidref = 1940765");
    expect(q).toContain("ORDER BY yr DESC");
    expect(simbadRefsCountQuery(1940765)).toContain("oidref = 1940765");
  });

  it("puts the papers about fewest objects first when focused", () => {
    const q = simbadRefsQuery(1940765, 25, { sort: "focused" });
    expect(q).toContain("r.nbobject");
    expect(q).toContain("ORDER BY r.nbobject ASC, yr DESC");
  });

  it("filters by year in both the list and the count", () => {
    const opts = { fromYear: 2010, toYear: 2020 };
    const where = 'r."year" >= 2010 AND r."year" < 2020';
    expect(simbadRefsQuery(1940765, 25, opts)).toContain(where);
    expect(simbadRefsCountQuery(1940765, opts)).toContain(where);
    expect(simbadRefsCountQuery(1940765, { toYear: 1990 })).toContain(
      'r."year" < 1990'
    );
  });
});

describe("parseSimbadRefs", () => {
  it("parses TAP JSON rows (real 3C 273 reference)", () => {
    expect(
      parseSimbadRefs([
        [
          "2026RNAAS..10..192D",
          "Geometric Masking in AGN Jets and its Implications for Unification and Blazar Physics.",
          2026,
          "RNAAS",
          "10.3847/2515-5172/ae8995",
          1,
        ],
        ["", "no bibcode", 2020, "X", null],
        "not a row",
      ])
    ).toEqual([
      {
        bibcode: "2026RNAAS..10..192D",
        title:
          "Geometric Masking in AGN Jets and its Implications for Unification and Blazar Physics",
        year: 2026,
        journal: "RNAAS",
        doi: "10.3847/2515-5172/ae8995",
        objects: 1,
      },
    ]);
  });

  it("tolerates missing data", () => {
    expect(parseSimbadRefs(undefined)).toEqual([]);
    expect(
      parseSimbadRefs([["1999A&A...1..1X", null, null, null, null]])
    ).toEqual([{ bibcode: "1999A&A...1..1X" }]);
  });
});

describe("adsAbstractUrl", () => {
  it("escapes bibcodes with ampersands", () => {
    expect(adsAbstractUrl("2023A&A...674A...1G")).toBe(
      "https://ui.adsabs.harvard.edu/abs/2023A%26A...674A...1G/abstract"
    );
  });
});
