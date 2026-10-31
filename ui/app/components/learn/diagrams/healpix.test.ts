import { describe, expect, it } from "vitest";

import {
  angToNest,
  childDigits,
  interleave,
  nestIndex,
  ringIndex,
  toRanges,
} from "./healpix";

// θ¹ Ori F, the worked example on the Learn page.
const RA = 83.81972;
const DEC = -5.39034;

describe("angToNest", () => {
  it("places θ¹ Ori F in equatorial face 5", () => {
    expect(angToNest(RA, DEC, 0).ipix).toBe(5);
  });

  it("matches the HEALPix level-12 index inside its Gaia DR3 source_id", () => {
    // Gaia source_id / 2^35 is the order-12 NESTED pixel.
    const sourceId = BigInt("3017364063331140224");
    expect(angToNest(RA, DEC, 12).ipix).toBe(
      Number(sourceId / BigInt(2) ** BigInt(35))
    );
  });

  it("is exact at order 18, past 32-bit range", () => {
    expect(angToNest(RA, DEC, 18).ipix).toBe(359697828043);
  });

  it("nests: dropping two bits per order gives the parent", () => {
    const p18 = angToNest(RA, DEC, 18).ipix;
    for (const order of [0, 1, 6, 12, 14]) {
      expect(angToNest(RA, DEC, order).ipix).toBe(
        Math.floor(p18 / 4 ** (18 - order))
      );
    }
  });
});

describe("childDigits", () => {
  it("rebuilds the index from the face and its child choices", () => {
    const pixel = angToNest(RA, DEC, 18);
    const rebuilt = childDigits(pixel, 18).reduce(
      (acc, d) => acc * 4 + d,
      pixel.face
    );
    expect(rebuilt).toBe(pixel.ipix);
  });
});

describe("pixel numbering", () => {
  it("interleaves ix on even bits and iy on odd bits", () => {
    expect(interleave(1, 0, 1)).toBe(1);
    expect(interleave(0, 1, 1)).toBe(2);
    expect(interleave(3, 3, 2)).toBe(15);
  });

  it("numbers base pixels the same in both schemes", () => {
    for (let f = 0; f < 12; f++) {
      expect(nestIndex(f, 0, 0, 0)).toBe(f);
      expect(ringIndex(f, 0, 0, 0)).toBe(f);
    }
  });

  it("gives every RING index exactly once", () => {
    const order = 2;
    const n = 2 ** order;
    const seen = new Set<number>();
    for (let f = 0; f < 12; f++)
      for (let ix = 0; ix < n; ix++)
        for (let iy = 0; iy < n; iy++) seen.add(ringIndex(f, ix, iy, order));
    expect(seen.size).toBe(12 * n * n);
    expect(Math.max(...seen)).toBe(12 * n * n - 1);
  });
});

describe("toRanges", () => {
  it("collapses consecutive indices", () => {
    expect(toRanges([1, 2, 3, 7, 9, 10])).toEqual([
      [1, 3],
      [7, 7],
      [9, 10],
    ]);
  });
});
