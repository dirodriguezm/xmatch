import { describe, expect, it } from "vitest";

import { parseXmmSourceCsv, xmmConeQuery } from "./xmm5";

// Real rows for 3C 273 from xsa.v_epic_source_cat: the stacked summary row
// and one of its 42 per-observation rows (no n_obs).
const HEADER =
  "srcid,iauname,ra,dec,radec_err,ep_flux,ep_flux_err,ep_1_flux,ep_1_flux_err,ep_2_flux,ep_2_flux_err,ep_3_flux,ep_3_flux_err,ep_4_flux,ep_4_flux_err,ep_5_flux,ep_5_flux_err,ep_hr1,ep_hr1_err,ep_hr2,ep_hr2_err,ep_hr3,ep_hr3_err,ep_hr4,ep_hr4_err,n_obs,mjd_first,mjd_last,var_flag,var_prob,fvar,fvar_err,stack_gamma,stack_nh,extent,sum_flag,gaiadr3_source_id,gaia_match_prob,wise_name,wise_match_prob,info_counterparts";
const STACKED =
  "3011277010100001,5XMM J122906.6+020308,187.27784189625316,2.052419894937612,1.290723230340518E-5,1.4098347489444052E-10,2.824588571934475E-14,1.4743236145808858E-11,3.193960920803941E-15,1.8617697661316157E-11,2.5458077096749245E-15,2.1117602458375906E-11,4.398705199046152E-15,3.346436322493318E-11,9.58923373784457E-15,5.3689525048028486E-11,2.5750840058924064E-14,0.14149458706378937,1.313964748987928E-4,-0.01713445968925953,1.3261608546599746E-4,-0.2743474841117859,1.652771170483902E-4,-0.41701826453208923,2.3090494505595416E-4,43,51708.3806018519,60317.4226273148,true,0.0,0.0195704735815525,0.00208212924189866,1.738886833190918,9.999999980506448E18,0.0,1,3700386905605055360,0.9702302123216551,J122906.69+020308.6,0.9950495841005519,https://xmm-catalog.irap.omp.eu/source/3011277010100001";
const PER_OBS =
  "3011277010100001,5XMM J122906.6+020308,187.27784189625316,2.052419894937612,1.290723230340518E-5,2.1905477431971576E-10,6.416704190562172E-13,2.2038045000005724E-11,7.38610222788226E-14,2.916151001230993E-11,9.568112632491604E-14,3.3611984723291854E-11,1.0470553731770776E-13,5.127278809857749E-11,2.272887010979635E-13,8.295743259001398E-11,5.767693180697142E-13,0.09071991592645645,0.002326738787814975,-0.021872377023100853,0.002239842899143696,-0.31300613284111023,0.0024441112764179707,-0.40247538685798645,0.0034507683012634516,,52259.6540856482,52259.7212384259,false,,,,1.738886833190918,9.999999980506448E18,0.0,1,,,,,";

const RA = 187.2779;
const DEC = 2.0524;

describe("xmmConeQuery", () => {
  it("cones the stacked catalogue only", () => {
    const q = xmmConeQuery(RA, DEC);
    expect(q).toContain("FROM xsa.v_epic_source_cat");
    expect(q).toContain("CIRCLE('ICRS', 187.2779, 2.0524,");
    expect(q).toContain("AND n_obs IS NOT NULL");
    expect(q).not.toMatch(/class/);
  });
});

describe("parseXmmSourceCsv", () => {
  it("parses the stacked row of 3C 273", () => {
    const s = parseXmmSourceCsv([HEADER, STACKED].join("\n"), RA, DEC)!;
    expect(s.name).toBe("5XMM J122906.6+020308");
    expect(s.separationArcsec).toBeLessThan(1);
    expect(s.flux?.value).toBeCloseTo(1.4098e-10, 13);
    expect(s.bandFluxes).toHaveLength(5);
    expect(s.bandFluxes[0]).toMatchObject({ label: "0.2–0.5 keV" });
    expect(s.hardnessRatios.map((h) => h.label)).toEqual([
      "HR1",
      "HR2",
      "HR3",
      "HR4",
    ]);
    expect(s.observations).toBe(43);
    expect(s.variable).toBe(true);
    expect(s.gamma).toBeCloseTo(1.739, 3);
    expect(s.gaia).toEqual({
      sourceId: "3700386905605055360",
      probability: expect.closeTo(0.97, 2),
    });
    expect(s.wise?.name).toBe("J122906.69+020308.6");
    expect(s.infoUrl).toBe(
      "https://xmm-catalog.irap.omp.eu/source/3011277010100001"
    );
  });

  it("ignores per-observation rows", () => {
    expect(parseXmmSourceCsv([HEADER, PER_OBS].join("\n"), RA, DEC)).toBeNull();
    const s = parseXmmSourceCsv(
      [HEADER, PER_OBS, STACKED].join("\n"),
      RA,
      DEC
    )!;
    expect(s.observations).toBe(43);
  });

  it("returns null when nothing lies within the match radius", () => {
    expect(
      parseXmmSourceCsv([HEADER, STACKED].join("\n"), RA, DEC + 10 / 3600)
    ).toBeNull();
    expect(parseXmmSourceCsv(HEADER, RA, DEC)).toBeNull();
    expect(parseXmmSourceCsv("", RA, DEC)).toBeNull();
  });
});
