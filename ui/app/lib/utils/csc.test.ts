import { describe, expect, it } from "vitest";

import { cscConeUrl, parseCscCone, parseVoTable } from "./csc";

const field = (name: string) =>
  `<FIELD datatype="double" name="${name}">\n</FIELD>`;

const COLUMNS = [
  "name",
  "ra",
  "dec",
  "err_ellipse_r0",
  "conf_flag",
  "extent_flag",
  "sat_src_flag",
  "flux_aper_b",
  "flux_aper_lolim_b",
  "flux_aper_hilim_b",
  "significance",
  "hard_hm",
  "hard_hm_lolim",
  "hard_hm_hilim",
  "var_intra_index_b",
  "var_inter_index_b",
];

const row = (cells: string[]) =>
  `<TR>\n${cells.map((c) => `<TD>${c}</TD>`).join("\n")}\n</TR>`;

// Trimmed from the CSC 2.2 cone search around 3C 273, plus a farther source.
const VOTABLE = `<VOTABLE version="1.2"><RESOURCE type="results"><TABLE name="RA:187.2779158">
${COLUMNS.map(field).join("\n")}
<DATA><TABLEDATA>
${row([
  "2CXO J122906.6+020308",
  "187.2778492039563",
  "2.0524019046366413",
  "0.6179603405207067",
  "F",
  "T",
  "T",
  "1.4168050808297E-11",
  "1.4094782486356E-11",
  "1.4241319130239E-11",
  "287.86547845913327",
  "0.34978138663335",
  "0.34540911930044",
  "0.35540287320425",
  "0",
  "8",
])}
${row([
  "2CXO J122906.7+020309",
  "187.2780",
  "2.0530",
  "",
  "F",
  "F",
  "F",
  "",
  "",
  "",
  "3.1",
  "",
  "",
  "",
  "-2147483648",
  "-2147483648",
])}
</TABLEDATA></DATA></TABLE></RESOURCE></VOTABLE>`;

const RA = 187.2779158387138;
const DEC = 2.052388668185861;

describe("parseVoTable", () => {
  it("maps cells to field names, empty cells to empty strings", () => {
    const rows = parseVoTable(VOTABLE);
    expect(rows).toHaveLength(2);
    expect(rows[0]?.name).toBe("2CXO J122906.6+020308");
    expect(rows[1]?.flux_aper_b).toBe("");
  });
});

describe("parseCscCone", () => {
  it("picks the nearest source and reads its values and flags", () => {
    const s = parseCscCone(VOTABLE, RA, DEC);
    expect(s?.name).toBe("2CXO J122906.6+020308");
    expect(s?.separationArcsec).toBeLessThan(0.3);
    expect(s?.flux?.value).toBeCloseTo(1.4168e-11, 14);
    expect(s?.hardHm?.lo).toBeCloseTo(0.3454, 4);
    expect(s?.varInter).toBe(8);
    expect(s?.varIntra).toBe(0);
    expect(s?.extended).toBe(true);
    expect(s?.saturated).toBe(true);
    expect(s?.confused).toBe(false);
  });

  it("treats the integer null as missing", () => {
    const s = parseCscCone(VOTABLE, 187.278, 2.053);
    expect(s?.name).toBe("2CXO J122906.7+020309");
    expect(s?.varInter).toBeUndefined();
    expect(s?.flux).toBeUndefined();
  });

  it("returns null when nothing is within the match radius", () => {
    expect(parseCscCone(VOTABLE, RA + 0.01, DEC)).toBeNull();
  });
});

describe("cscConeUrl", () => {
  it("searches twice the match radius, in degrees", () => {
    const url = new URL(cscConeUrl(RA, DEC));
    expect(url.searchParams.get("SR")).toBe(String(6 / 3600));
  });
});
