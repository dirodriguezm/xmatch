import { describe, expect, it } from "vitest";

import { parsePs1DetectionsCsv } from "./ps1Lightcurve";

/** Real detections of PS1 78000043968034076, plus a synthetic neighbour. */
const CSV = `objID,obsTime,filterID,psfFlux,psfFluxErr,ra,dec
78000043968034076,55447.4966518,2,0.004172800108790398,6.930919880687725e-06,4.39682901,-24.99687932
78000043968034076,55457.4638563,3,0.003500829916447401,6.680390015390003e-06,4.39683257,-24.99688495
78000043968034076,55447.4085872,5,-1e-06,6.87e-06,4.39682413,-24.99687778
99999999999999999,55450.0,1,0.001,1e-05,4.39700000,-24.99700000
`;

describe("parsePs1DetectionsCsv", () => {
  const { objId, points } = parsePs1DetectionsCsv(
    CSV,
    4.396816816814972,
    -24.99689333457813
  );

  it("keeps the 64-bit objID exactly", () => {
    // As a JS number this would round to 78000043968034080.
    expect(objId).toBe("78000043968034076");
  });

  it("keeps only the object nearest the target", () => {
    expect(points.every((p) => p.band !== "g")).toBe(true);
  });

  it("converts Jy to AB magnitudes", () => {
    const r = points.find((p) => p.band === "r")!;
    expect(r.mag).toBeCloseTo(-2.5 * Math.log10(0.004172800108790398) + 8.9, 6);
    expect(r.magerr).toBeCloseTo(
      ((2.5 / Math.LN10) * 6.930919880687725e-6) / 0.004172800108790398,
      6
    );
  });

  it("drops non-positive fluxes and sorts by time", () => {
    expect(points.map((p) => p.band)).toEqual(["r", "i"]);
  });

  it("returns nothing when columns are missing", () => {
    expect(parsePs1DetectionsCsv("foo,bar\n1,2\n", 0, 0)).toEqual({
      objId: null,
      points: [],
    });
  });
});
