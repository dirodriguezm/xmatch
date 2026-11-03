import { describe, expect, it } from "vitest";

import {
  crtsTelescope,
  parseCrtsCsv,
  parseCrtsQueryPage,
} from "./crtsLightcurve";

const HEADER = "MasterID,Mag,Magerr,RA,Dec,MJD,Blend";

describe("parseCrtsQueryPage", () => {
  it("recognises uncovered areas", () => {
    const html =
      "<p>Photcat DB query</p><p>This area is not covered by CSS data.</p>";
    expect(parseCrtsQueryPage(html)).toEqual({
      status: "not_covered",
      csvUrl: null,
    });
  });

  it("recognises covered but empty areas", () => {
    const html =
      "This area is covered by the data release. However, no objects were found in the specified area .";
    expect(parseCrtsQueryPage(html).status).toBe("empty");
  });

  it("extracts the CSV link when rows were found", () => {
    const html =
      'There were 446 lines. <a href="http://nunuku.caltech.edu/DataRelease/upload/result_web_fileSbxUFV.csv">(right-mouse-click and save as to download )</a>';
    expect(parseCrtsQueryPage(html)).toEqual({
      status: "rows",
      csvUrl:
        "http://nunuku.caltech.edu/DataRelease/upload/result_web_fileSbxUFV.csv",
    });
  });
});

describe("crtsTelescope", () => {
  it("maps the MasterID prefix to the telescope", () => {
    expect(crtsTelescope("1015071022010")).toBe("CSS");
    expect(crtsTelescope("2015184009078")).toBe("MLS");
    expect(crtsTelescope("3015100003489")).toBe("SSS");
    expect(crtsTelescope("9000000000000")).toBe("CRTS");
  });
});

describe("parseCrtsCsv", () => {
  const RA = 206.1074685;
  const DEC = -15.9153253;

  it("merges every MasterID of the same star and bands by telescope", () => {
    const csv = [
      HEADER,
      `1015071022010,14.60,0.05,${RA}, ${DEC},53500.1,0`,
      `2015184009078,14.50,0.01,${RA}, ${DEC},53800.2,0`,
      `3015100003489,14.70,0.07,${RA}, ${DEC},53600.3,0`,
      // A neighbour 10″ away must not be merged in.
      `1015072099999,16.00,0.10,${RA}, ${DEC + 10 / 3600},53500.2,0`,
    ].join("\n");
    const lc = parseCrtsCsv(csv, RA, DEC);
    expect(lc.sources.map((s) => s.telescope).sort()).toEqual([
      "CSS",
      "MLS",
      "SSS",
    ]);
    expect(lc.points.map((p) => p.band)).toEqual(["CSS", "SSS", "MLS"]);
    expect(lc.saturated).toBe(false);
  });

  it("drops isolated spikes but keeps a dip seen all night", () => {
    const night = (mjd: number, mags: number[]) =>
      mags.map(
        (m, k) => `1121064053574,${m},0.05,${RA}, ${DEC},${mjd + k * 0.007},0`
      );
    const csv = [
      HEADER,
      ...night(53500.2, [14.1, 14.0, 14.1, 14.0]),
      ...night(53501.2, [14.1, 18.9, 14.0, 14.1]), // one spurious exposure
      ...night(53502.2, [15.6, 15.7, 15.6, 15.7]), // real eclipse
      ...night(53503.2, [14.0, 14.1, 14.0, 14.1]),
    ].join("\n");
    const lc = parseCrtsCsv(csv, RA, DEC);
    expect(lc.outliers).toBe(1);
    expect(lc.points).toHaveLength(15);
    expect(lc.points.some((p) => p.mag === 18.9)).toBe(false);
    expect(lc.points.filter((p) => p.mag! > 15.5)).toHaveLength(4);
  });

  it("flags saturated stars", () => {
    const csv = [
      HEADER,
      `1129055037278,11.6,0.05,${RA}, ${DEC},53527.1,0`,
      `1129055037278,11.7,0.05,${RA}, ${DEC},53527.2,0`,
    ].join("\n");
    expect(parseCrtsCsv(csv, RA, DEC).saturated).toBe(true);
  });

  it("returns nothing when no source lies within the match radius", () => {
    const csv = [HEADER, `1,14,0.05,${RA}, ${DEC + 5 / 3600},53500,0`].join(
      "\n"
    );
    expect(parseCrtsCsv(csv, RA, DEC).points).toEqual([]);
  });
});
