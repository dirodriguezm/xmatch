import { describe, expect, it } from "vitest";

import {
  dereddenFactor,
  dereddeningBlocker,
  extinctionMag,
  extinctionRatio,
} from "./extinction";

describe("extinctionRatio (CCM89 + O'Donnell94, R_V = 3.1)", () => {
  it("is 1 in V by definition", () => {
    expect(extinctionRatio(0.5495)).toBeCloseTo(1, 2);
  });

  it("gives A_B/A_V = (R_V + 1)/R_V in B", () => {
    expect(extinctionRatio(0.44)).toBeCloseTo(4.1 / 3.1, 1);
  });

  it("matches the published near-IR values", () => {
    expect(extinctionRatio(2.2)).toBeCloseTo(0.114, 2); // K
    expect(extinctionRatio(1.25)).toBeCloseTo(0.282, 2); // J
  });

  it("peaks at the 2175 Å bump relative to neighbouring UV", () => {
    const bump = extinctionRatio(0.2175)!;
    expect(bump).toBeGreaterThan(extinctionRatio(0.26)!);
    expect(bump).toBeGreaterThan(extinctionRatio(0.19)!);
  });

  it("returns null shortward of 0.1 µm", () => {
    expect(extinctionRatio(0.09)).toBeNull();
  });
});

describe("dereddening", () => {
  it("scales with E(B−V)", () => {
    expect(extinctionMag(0.5495, 0.1)).toBeCloseTo(0.31, 2);
  });

  it("brightens the flux by 10^(0.4 A)", () => {
    expect(dereddenFactor(0.5495, 0.1)).toBeCloseTo(10 ** (0.4 * 0.31), 2);
  });

  it("leaves out-of-range wavelengths untouched", () => {
    expect(dereddenFactor(0.05, 0.5)).toBe(1);
  });
});

describe("dereddeningBlocker", () => {
  it("allows typical high-latitude fields", () => {
    expect(dereddeningBlocker(0.0163, 42.1)).toBeNull();
  });

  it("refuses the Galactic plane", () => {
    expect(dereddeningBlocker(0.3, -2)).toMatch(/Galactic plane/);
  });

  it("refuses map values too large to trust (Sgr A*: E(B−V) ≈ 85)", () => {
    expect(dereddeningBlocker(84.7, 30)).toMatch(/too high/);
  });
});
