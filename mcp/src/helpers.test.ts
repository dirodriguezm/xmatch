import { describe, expect, it } from "vitest";
import {
  ValidationError,
  flattenMatches,
  groupBulkByIndex,
  objectUrl,
  parseSesame,
  summarizeLightcurve,
  summarizeMatches,
  validatePosition,
  validateRadius,
} from "./helpers.js";

const SESAME_BETELGEUSE = `# Betelgeuse\t#Q8299282
#=Sc=Simbad (CDS, via client/server):    1   134ms
%@ @843879
%C.0 s*r
%J 88.79293899 +7.40706400 = 05 55 10.305  +07 24 25.43
%J.E [9.04 5.72 90] A 2007A&A...474..653V
%I.0 * alf Ori
%I NAME Betelgeuse
#B 1803
`;

const SESAME_MISS = `# notarealobjectxyz\t#Q8299283
#!Sc=Simbad (CDS, via client/server): ***  The catalogue name is unknown: notarealobjectxyz
#=N=NED
# (nothing found) #
`;

describe("parseSesame", () => {
  it("parses the %J line and identity fields", () => {
    const r = parseSesame(SESAME_BETELGEUSE, "Betelgeuse");
    expect(r).toMatchObject({ name: "Betelgeuse", ra: 88.79293899, dec: 7.407064, mainId: "* alf Ori", objectType: "s*r", source: "Simbad" });
  });
  it("returns null when nothing is found", () => {
    expect(parseSesame(SESAME_MISS, "x")).toBeNull();
  });
  it("handles negative declinations", () => {
    expect(parseSesame("%J 83.8 -5.39 = x\n", "Orion")?.dec).toBe(-5.39);
  });
});

describe("validation", () => {
  it("accepts valid positions", () => {
    expect(() => validatePosition(10.68, 41.27, 5)).not.toThrow();
    expect(() => validatePosition(0, -90, 120)).not.toThrow();
  });
  it("rejects out-of-range values with helpful messages", () => {
    expect(() => validatePosition(361, 0)).toThrow(/between 0 and 360/);
    expect(() => validatePosition(10, 91)).toThrow(/between -90 and 90/);
    expect(() => validateRadius(121)).toThrow(/at most 120/);
    expect(() => validateRadius(0)).toThrow(ValidationError);
    expect(() => validateRadius(Number.NaN)).toThrow(ValidationError);
  });
});

describe("objectUrl", () => {
  it("encodes ids with spaces", () => {
    expect(objectUrl("Gaia DR3 123", "gaia")).toBe("https://xwave-rho.vercel.app/object/Gaia%20DR3%20123?catalog=gaia");
  });
  it("respects a custom base without doubling slashes", () => {
    expect(objectUrl("a+b", "erosita", "http://x/")).toBe("http://x/object/a%2Bb?catalog=erosita");
  });
});

const BULK = [
  { catalog: "allwise", data: [{ id: "w1", ra: 1, dec: 2, cat: "allwise", distance: 0.21934 }], index: 0 },
  { catalog: "gaia", data: [{ id: "g1", ra: 1, dec: 2, cat: "gaia", distance: 7.9 }], index: 0 },
  { catalog: "erosita", data: [{ id: "e1", ra: 3, dec: 4, cat: "erosita", distance: 1.5 }], index: 2 },
  { catalog: "allwise", data: null, index: 1 },
];

describe("match flattening", () => {
  it("sorts by distance and adds urls", () => {
    const m = flattenMatches(BULK);
    expect(m.map((x) => x.id)).toEqual(["w1", "e1", "g1"]);
    expect(m[0]!.distance_arcsec).toBe(0.219);
    expect(m[0]!.url).toContain("/object/w1?catalog=allwise");
  });
  it("groups bulk results by input index", () => {
    const g = groupBulkByIndex(BULK, 3);
    expect(g.map((x) => x.map((m) => m.id))).toEqual([["w1", "g1"], [], ["e1"]]);
  });
  it("summarises empty and non-empty results", () => {
    expect(summarizeMatches([])).toBe("No sources found.");
    expect(summarizeMatches(flattenMatches(BULK))).toMatch(/^3 source\(s\): AllWISE 1, eROSITA eRASS1 1, Gaia DR3 1/);
  });
});

describe("summarizeLightcurve", () => {
  it("computes per-survey ranges and median", () => {
    const s = summarizeLightcurve([
      { catalog: "neowise", mjd: 3, mag: 7, magerr: 0.1 },
      { catalog: "neowise", mjd: 1, mag: 5, magerr: 0.1 },
      { catalog: "neowise", mjd: 2, mag: 6, magerr: 0.1 },
      { catalog: "neowise", mjd: 4, mag: 8, magerr: 0.1 },
    ]);
    expect(s).toEqual([{ catalog: "neowise", n: 4, mjd_min: 1, mjd_max: 4, mag_min: 5, mag_max: 8, mag_median: 6.5 }]);
  });
});
