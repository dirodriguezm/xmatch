import { describe, expect, it } from "vitest";

import {
  gaiaSourceIdFromDesignation,
  parseGaiaEpochCsv,
  parseGaiaTapCone,
} from "./gaiaEpoch";

/** Real rows for Gaia DR3 2328619740476096384 (an RR Lyrae), trimmed. */
const CSV = `solution_id,source_id,transit_id,g_transit_time,g_transit_flux,g_transit_flux_error,g_transit_flux_over_error,g_transit_mag,bp_obs_time,bp_flux,bp_flux_error,bp_flux_over_error,bp_mag,rp_obs_time,rp_flux,rp_flux_error,rp_flux_over_error,rp_mag,variability_flag_g_reject,variability_flag_bp_reject,variability_flag_rp_reject,g_other_flags,bp_other_flags,rp_other_flags,rejected_by_photometry
375316653866487564,2328619740476096384,19904707766277818,1756.8137006432394,61379.645021363125,66.40325372862291,924.3469,13.717305935966815,1756.8139941058444,38115.55505560519,95.02268429563246,401.12057,13.885786594445392,1756.8140820892297,33731.73916407903,88.0843489264427,382.94815,13.427798667608632,false,false,false,5,0,0,false
375316653866487564,2328619740476096384,22755823142500184,1808.3702083380067,85639.92140745468,111.8572199076044,765.61816,13.355676216889247,,,,,,,,,,,false,true,true,4097,1024,16779264,false
375316653866487564,2328619740476096384,99999999999999999,1820.0,1,1,100,13.5,1820.0,1,1,100,13.6,1820.0,1,1,100,13.4,false,false,false,0,0,0,true
`;

describe("gaiaSourceIdFromDesignation", () => {
  it("keeps every digit of the 64-bit id", () => {
    expect(gaiaSourceIdFromDesignation("Gaia DR3 3836259600466717952")).toBe(
      "3836259600466717952"
    );
  });

  it("rejects anything else", () => {
    expect(gaiaSourceIdFromDesignation("1497p015_ac51-054642")).toBeNull();
    expect(gaiaSourceIdFromDesignation("Gaia DR2 123")).toBeNull();
  });
});

describe("parseGaiaEpochCsv", () => {
  const points = parseGaiaEpochCsv(CSV);

  it("expands each transit into G, BP and RP points", () => {
    const first = points.filter((p) => p.mjd! < 55197 + 1757);
    expect(first.map((p) => p.band).sort()).toEqual(["BP", "G", "RP"]);
  });

  it("converts Gaia time to MJD", () => {
    const g = points.find((p) => p.band === "G")!;
    expect(g.mjd).toBeCloseTo(1756.8137006432394 + 55197, 6);
  });

  it("derives magnitude errors from flux S/N", () => {
    const g = points.find((p) => p.band === "G")!;
    expect(g.magerr).toBeCloseTo(2.5 / Math.LN10 / 924.3469, 6);
  });

  it("drops rejected measurements and empty bands", () => {
    // Row 2: BP/RP rejected (and empty) → only G. Row 3: rejected entirely.
    expect(points).toHaveLength(4);
    expect(points.filter((p) => p.band === "G")).toHaveLength(2);
  });

  it("returns nothing for a non-CSV body", () => {
    expect(parseGaiaEpochCsv("<VOTABLE>error</VOTABLE>")).toEqual([]);
  });
});

describe("parseGaiaTapCone", () => {
  it("reads the exact source_id and the epoch flag", () => {
    expect(
      parseGaiaTapCone(
        "source_id,has_epoch_photometry,dist\n2328619740476096384,true,0.0001\n"
      )
    ).toEqual({ sourceId: "2328619740476096384", hasEpochPhotometry: true });
  });

  it("returns null for an empty cone", () => {
    expect(
      parseGaiaTapCone("source_id,has_epoch_photometry,dist\n")
    ).toBeNull();
  });
});
