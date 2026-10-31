import { describe, expect, it } from "vitest";

import {
  binPoints,
  findGammaAssociations,
  metToMjd,
  parseBatLightcurve,
  parseLcrLightcurve,
} from "./gamma";

describe("metToMjd", () => {
  it("starts Fermi time at MJD 51910 (2001-01-01)", () => {
    expect(metToMjd(0)).toBe(51910);
    expect(metToMjd(86_400)).toBe(51911);
  });
});

describe("parseLcrLightcurve", () => {
  // Shape of the repository's lightCurveData JSON, values from 3C 273.
  const json = {
    ts: [
      [239889601, 82.29],
      [326376001, 1.2],
    ],
    flux: [[239889601, 3.03e-7]],
    flux_error: [[239889601, 2.5e-7, 3.56e-7]],
    flux_upper_limits: [[326376001, 1.16e-7]],
    fit_convergence: [
      [239889601, 0],
      [326376001, 0],
    ],
  };

  it("keeps detections with their bounds and limits as limits, in time order", () => {
    const points = parseLcrLightcurve(json);
    expect(points).toHaveLength(2);
    expect(points[0]).toMatchObject({
      value: 3.03e-7,
      lo: 2.5e-7,
      hi: 3.56e-7,
    });
    expect(points[0]?.mjd).toBeCloseTo(51910 + 239889601 / 86_400, 6);
    expect(points[1]).toMatchObject({ value: 1.16e-7, upperLimit: true });
  });

  it("drops bins whose fit did not converge", () => {
    const points = parseLcrLightcurve({
      ...json,
      fit_convergence: [
        [239889601, 156],
        [326376001, 0],
      ],
    });
    expect(points.map((p) => p.upperLimit)).toEqual([true]);
  });

  it("returns nothing for an empty or malformed answer", () => {
    expect(parseLcrLightcurve(null)).toEqual([]);
    expect(parseLcrLightcurve({})).toEqual([]);
  });
});

describe("parseBatLightcurve", () => {
  const text = `# Swift-BAT daily average light curve for 3C 273
#     TIME           RATE          ERROR      YEAR          DAY     STAT_ERR        SYS_ERR   DATA_FLAG   TIMEDEL_EXPO     TIMEDEL_CODED  TIMEDEL_DITHERED
      53416      -0.00039373     0.00079306   2005          46     0.00067889    -0.00001378      0   19104.00000000  6585.40625000     0.00000000
      53417      -0.00050299     0.00057838   2005          47     0.00049496    -0.00001760      1   14048.00000000  9546.75000000     0.00000000
      53419       0.00295081     0.00097062   2005          49     0.00082630     0.00010328      0   14528.00000000  3884.00000000  3264.00000000
`;

  it("reads good days and drops flagged ones", () => {
    const points = parseBatLightcurve(text);
    expect(points.map((p) => p.mjd)).toEqual([53416, 53419]);
    expect(points[1]?.lo).toBeCloseTo(0.00295081 - 0.00097062, 8);
  });
});

describe("binPoints", () => {
  it("weights by inverse variance within each bin", () => {
    const points = [
      { mjd: 0, value: 1, lo: 0, hi: 2 }, // σ = 1
      { mjd: 1, value: 4, lo: 2, hi: 6 }, // σ = 2
      { mjd: 40, value: 5, lo: 4, hi: 6 },
    ];
    const bins = binPoints(points, 30);
    expect(bins).toHaveLength(2);
    // (1·1 + 4·¼) / (1 + ¼) = 1.6, σ = 1/√1.25
    expect(bins[0]?.value).toBeCloseTo(1.6, 10);
    expect((bins[0]?.hi ?? NaN) - (bins[0]?.value ?? NaN)).toBeCloseTo(
      1 / Math.sqrt(1.25),
      10
    );
    expect(bins[1]?.value).toBe(5);
  });
});

describe("findGammaAssociations", () => {
  it("links 3C 273 to its 4FGL source and BAT monitor entry", () => {
    const g = findGammaAssociations(187.2779158387138, 2.052388668185861);
    expect(g.fermi?.name).toBe("4FGL J1229.0+0202");
    expect(g.fermi?.assoc).toBe("3C 273");
    expect(g.bat?.path).toBe("weak/3C273");
  });

  it("finds nothing for a blank patch of sky", () => {
    expect(findGammaAssociations(10.6847, 41.269)).toEqual({});
  });
});
