import { describe, expect, it } from "vitest";

import { ang2pixNest, hipsTilePath } from "./healpix";

describe("ang2pixNest", () => {
  it("puts order-0 positions on the right base face", () => {
    expect(ang2pixNest(0, 0, 90)).toBe(0);
    expect(ang2pixNest(0, 0, 0)).toBe(4);
    expect(ang2pixNest(0, 45, 60)).toBe(0);
    expect(ang2pixNest(0, 45, -60)).toBe(8);
    expect(ang2pixNest(0, 180, -90)).toBe(10);
  });

  it("matches the tiles Aladin requests around Orion", () => {
    // Seen in the browser for the view centred on 83.8220833, -5.3911111.
    expect(ang2pixNest(3, 83.8220833, -5.3911111)).toBe(334);
    expect(ang2pixNest(8, 83.8220833, -5.3911111)).toBe(343034);
    expect(ang2pixNest(11, 83.8220833, -5.3911111)).toBe(21954213);
  });

  it("nests: a cell's parent is its index shifted by two bits per order", () => {
    const deep = ang2pixNest(11, 10.6847, 41.269);
    expect(ang2pixNest(9, 10.6847, 41.269)).toBe(Math.floor(deep / 16));
  });
});

describe("hipsTilePath", () => {
  it("groups tiles in directories of 10 000", () => {
    expect(hipsTilePath(11, 83.8220833, -5.3911111, "jpg")).toBe(
      "Norder11/Dir21950000/Npix21954213.jpg"
    );
  });
});
