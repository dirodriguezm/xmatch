import { describe, expect, it } from "vitest";

import { buildSedPoints, type PhotometrySource } from "./sed";

const source = (
  catalog: string,
  record: Record<string, unknown>,
  isSelf = true
): PhotometrySource => ({
  catalog,
  record,
  isSelf,
  separationArcsec: isSelf ? 0 : 0.4,
});

describe("buildSedPoints", () => {
  it("converts a Vega magnitude to νFν with the band zero point", () => {
    // W1 = 0 → F_ν = 309.54 Jy; νFν = F_ν · 1e−23 · c/λ
    const [p] = buildSedPoints([
      source("allwise", { w1mpro: 0, w1sigmpro: 0.02 }),
    ]);
    const expected = 309.54 * 1e-23 * (2.99792458e14 / 3.3526);
    expect(p.band).toBe("W1");
    expect(p.nuFnu / expected).toBeCloseTo(1, 6);
    expect(p.upperLimit).toBe(false);
  });

  it("drops 5 magnitudes per factor of 100 in flux", () => {
    const [bright] = buildSedPoints([
      source("allwise", { w2mpro: 10, w2sigmpro: 0.05 }),
    ]);
    const [faint] = buildSedPoints([
      source("allwise", { w2mpro: 15, w2sigmpro: 0.05 }),
    ]);
    expect(bright.nuFnu / faint.nuFnu).toBeCloseTo(100, 6);
  });

  it("marks AllWISE magnitudes without an uncertainty as upper limits", () => {
    const points = buildSedPoints([
      source("allwise", { w3mpro: 12.06, w3sigmpro: null }),
    ]);
    expect(points).toHaveLength(1);
    expect(points[0].upperLimit).toBe(true);
    expect(points[0].nuFnuErr).toBeUndefined();
  });

  it("treats Gaia magnitudes of exactly 0 as missing", () => {
    const points = buildSedPoints([
      source("gaia", {
        phot_g_mean_mag: 18.68,
        phot_g_mean_flux_over_error: 484,
        phot_bp_mean_mag: 0,
        phot_rp_mean_mag: 0,
      }),
    ]);
    expect(points.map((p) => p.band)).toEqual(["G"]);
  });

  it("derives the Gaia magnitude error from flux S/N", () => {
    const [g] = buildSedPoints([
      source("gaia", {
        phot_g_mean_mag: 18,
        phot_g_mean_flux_over_error: 100,
      }),
    ]);
    // σ_m = 2.5 / ln10 / (S/N)
    expect(g.magErr).toBeCloseTo(2.5 / Math.LN10 / 100, 6);
  });

  it("merges sources and sorts by wavelength, keeping provenance", () => {
    const points = buildSedPoints([
      source("allwise", { w1mpro: 17.08, w1sigmpro: 0.13 }),
      source(
        "gaia",
        { phot_g_mean_mag: 18.68, phot_g_mean_flux_over_error: 484 },
        false
      ),
    ]);
    expect(points.map((p) => p.band)).toEqual(["G", "W1"]);
    expect(points[0].isSelf).toBe(false);
    expect(points[0].separationArcsec).toBe(0.4);
  });

  it("skips null and sentinel values", () => {
    const points = buildSedPoints([
      source("allwise", { w1mpro: null, w2mpro: -999, j_m_2mass: undefined }),
    ]);
    expect(points).toHaveLength(0);
  });
});
