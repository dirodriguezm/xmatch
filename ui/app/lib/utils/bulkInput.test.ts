import { describe, expect, it } from "vitest";

import {
  BULK_EXAMPLE,
  bulkResultsToCsv,
  groupBulkResults,
  parseBulkInput,
  parseDec,
  parseRa,
  summarizeBulkResults,
} from "./bulkInput";

describe("parseRa / parseDec", () => {
  it("reads decimal degrees", () => {
    expect(parseRa("10.684708")).toBeCloseTo(10.684708, 6);
    expect(parseDec("-16.716")).toBeCloseTo(-16.716, 6);
  });

  it("reads sexagesimal with colons, spaces and letters", () => {
    expect(parseRa("06:45:08.917")).toBeCloseTo(101.287154, 5);
    expect(parseRa("06 45 08.917")).toBeCloseTo(101.287154, 5);
    expect(parseRa("06h45m08.917s")).toBeCloseTo(101.287154, 5);
    expect(parseDec("-16:42:58.02")).toBeCloseTo(-16.716117, 5);
    expect(parseDec("-16d42m58.02s")).toBeCloseTo(-16.716117, 5);
    expect(parseDec("+41°16′07.5″")).toBeCloseTo(41.26875, 5);
  });

  it("keeps the sign of a -00° declination", () => {
    expect(parseDec("-00:30:00")).toBeCloseTo(-0.5, 6);
  });

  it("rejects out-of-range values", () => {
    expect(parseRa("361")).toBeNull();
    expect(parseRa("25:00:00")).toBeNull();
    expect(parseDec("91")).toBeNull();
    expect(parseDec("+10:61:00")).toBeNull();
    expect(parseRa("abc")).toBeNull();
  });
});

describe("parseBulkInput", () => {
  it("parses the built-in example completely", () => {
    const { targets, errors, header } = parseBulkInput(BULK_EXAMPLE);
    expect(errors).toEqual([]);
    expect(header).toEqual(["name", "ra", "dec"]);
    expect(targets).toHaveLength(8);
    expect(targets[0]).toMatchObject({ index: 0, name: "M31 nucleus" });
    expect(targets[3].ra).toBeCloseTo(101.287154, 5);
  });

  it("maps header columns case-insensitively in any order", () => {
    const { targets, errors } = parseBulkInput(
      "DEJ2000\tObject\tRAJ2000\tmag\n41.27\tM31\t10.68\t3.4"
    );
    expect(errors).toEqual([]);
    expect(targets[0]).toMatchObject({ name: "M31", ra: 10.68, dec: 41.27 });
  });

  it("guesses layout without a header", () => {
    const { targets } = parseBulkInput(
      ["10.68 41.27", "Vega,279.23,38.78", "88.79,7.41,Betelgeuse"].join("\n")
    );
    // Delimiter comes from the first line (whitespace); commas then fail.
    expect(targets[0]).toMatchObject({ ra: 10.68, dec: 41.27, name: "" });
  });

  it("handles comma-separated name-first and ra-first rows", () => {
    const { targets, errors } = parseBulkInput(
      "Vega,279.23,38.78\n88.79,7.41,Betelgeuse"
    );
    expect(errors).toEqual([]);
    expect(targets[0]).toMatchObject({ name: "Vega", ra: 279.23 });
    expect(targets[1]).toMatchObject({ name: "Betelgeuse", ra: 88.79 });
  });

  it("handles whitespace names and space-separated sexagesimal", () => {
    const { targets, errors } = parseBulkInput(
      "Proxima Cen 14 29 42.946 -62 40 46.17\n10.68 41.27 M31 nucleus"
    );
    expect(errors).toEqual([]);
    expect(targets[0].name).toBe("Proxima Cen");
    expect(targets[0].dec).toBeCloseTo(-62.679492, 5);
    expect(targets[1]).toMatchObject({ name: "M31 nucleus", ra: 10.68 });
  });

  it("reports per-line errors with 1-based line numbers and skips comments", () => {
    const { targets, errors } = parseBulkInput(
      "# my list\nname,ra,dec\nA,10,20\nB,400,20\nC,10,-95\n\nD,1,2"
    );
    expect(targets.map((t) => t.name)).toEqual(["A", "D"]);
    expect(targets.map((t) => t.index)).toEqual([0, 1]);
    expect(errors.map((e) => e.line)).toEqual([4, 5]);
    expect(errors[0].message).toMatch(/RA/);
    expect(errors[1].message).toMatch(/Dec/);
  });

  it("honours quoted CSV fields", () => {
    const { targets } = parseBulkInput('name,ra,dec\n"NGC 1, core",10,20');
    expect(targets[0].name).toBe("NGC 1, core");
  });

  it("returns nothing for empty input", () => {
    expect(parseBulkInput("  \n# only a comment\n")).toEqual({
      targets: [],
      errors: [],
      header: null,
    });
  });
});

describe("result shaping", () => {
  const { targets } = parseBulkInput("A,10,20\nB,30,40\nC,50,60");
  const response = [
    {
      catalog: "allwise",
      data: [{ id: "w2", ra: 10, dec: 20, cat: "allwise", distance: 3.2 }],
      index: 0,
    },
    {
      catalog: "erosita",
      data: [{ id: "e1", ra: 10, dec: 20, cat: "erosita", distance: 1.1 }],
      index: 0,
    },
    {
      catalog: "allwise",
      data: [{ id: "w9", ra: 50, dec: 60, cat: "allwise", distance: 0.4 }],
      index: 2,
    },
  ];

  it("groups by input index and sorts by separation", () => {
    const rows = groupBulkResults(targets, response);
    expect(rows[0].matches.map((m) => m.id)).toEqual(["e1", "w2"]);
    expect(rows[1].matches).toEqual([]);
    expect(rows[2].matches[0].separation).toBe(0.4);
  });

  it("summarizes matched / unmatched / per-catalog", () => {
    const summary = summarizeBulkResults(groupBulkResults(targets, response));
    expect(summary).toEqual({
      inputs: 3,
      matched: 2,
      unmatched: 1,
      totalMatches: 3,
      perCatalog: { allwise: 2, erosita: 1 },
    });
  });

  it("writes a CSV line per match plus one per unmatched input", () => {
    const csv = bulkResultsToCsv(groupBulkResults(targets, response));
    const lines = csv.split("\n");
    expect(lines[0]).toMatch(/^input_row,input_name/);
    expect(lines).toHaveLength(1 + 2 + 1 + 1);
    expect(lines[3]).toBe("2,B,30,40,,,,,");
  });
});
