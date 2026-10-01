import { describe, expect, it } from "vitest";

import { compareLightcurveCatalogs } from "./lightcurve";

describe("compareLightcurveCatalogs", () => {
  it("orders infrared before optical, then by survey epoch", () => {
    expect(
      ["ztf", "gaia", "crts", "neowise", "ps1"].sort(compareLightcurveCatalogs)
    ).toEqual(["neowise", "crts", "ps1", "gaia", "ztf"]);
  });

  it("puts unknown surveys last, alphabetically", () => {
    expect(
      ["vlass", "ztf", "swift", "neowise"].sort(compareLightcurveCatalogs)
    ).toEqual(["neowise", "ztf", "swift", "vlass"]);
  });
});
