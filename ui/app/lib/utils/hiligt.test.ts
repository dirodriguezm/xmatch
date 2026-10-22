import { describe, expect, it } from "vitest";

import { isHiligtMission, isoToMjd, parseHiligt } from "./hiligt";

// Real records returned by ULSservice_passthru: an XMM slew detection of
// 3C 273, a ROSAT survey detection of 3C 273 (space-separated date), and a
// ROSAT survey upper limit at AO Com.
const SLEW_DETECTION = {
  _bck_counts: "2.3289474",
  _bkgimage: "NotSet",
  _crate: "26.144296",
  _crate_err: "1.8603947",
  _crate_flux: 3.3412410288e-11,
  _crate_flux_err: 2.3775844266e-12,
  _dec: "2.0524",
  _eef: "0.845",
  _ehigh: 2.0,
  _elow: 0.2,
  _end_date: "2012-07-16T23:57:27",
  _expmap: "P9230800003PNS003EXPMAP6055.FTZ",
  _exptime: "9.0834551",
  _filt: "Medium",
  _image: "P9230800003PNS003IMAGE_6055.FTZ",
  _instrum: "EPIC-pn",
  _label: "3C273",
  _mission: "XMM-Newton slew",
  _model: [3e20, "plaw", 2.0],
  _obsid: "9230800003",
  _ra: "187.2779",
  _src_counts: "203",
  _start_date: "2012-07-16T22:52:16",
  _status: "Ok",
  _ul: null,
  _ul_flux: " ",
  _ulsig: "2",
};
const ROSAT_DETECTION = {
  _bck_counts: "0.8843009720999999",
  _bkgimage: "NotSet",
  _crate: "7.8635",
  _crate_err: "0.1474",
  _crate_flux: 8.052224e-11,
  _crate_flux_err: 1.509376e-12,
  _dec: "2.0524",
  _eef: 1.0,
  _ehigh: 2.0,
  _elow: 0.2,
  _end_date: "1990-12-20 22:04:27",
  _expmap: "931734_exposure.fits",
  _exptime: "368.27",
  _filt: "Open",
  _image: "931734_image4.fits",
  _instrum: "ROSAT-PSPC",
  _label: "3C273",
  _mission: "ROSAT-Survey",
  _model: [3e20, "plaw", 2.0],
  _obsid: "931734",
  _ra: "187.2779",
  _src_counts: "2895.94",
  _start_date: "1990-12-19 01:15:26",
  _status: "Ok",
  _ul: null,
  _ul_flux: " ",
  _ulsig: "2",
};
const ROSAT_UPPER_LIMIT = {
  _bck_counts: "5.525",
  _bkgimage: "NotSet",
  _crate: null,
  _crate_err: null,
  _crate_flux: " ",
  _crate_flux_err: " ",
  _dec: "22.3950992",
  _eef: "0.58",
  _ehigh: 2.0,
  _elow: 0.2,
  _end_date: "1990-12-20T00:00:00",
  _expmap: "931332_exposure.fits",
  _exptime: "449.31927",
  _filt: "Open",
  _image: "931332_image4.fits",
  _instrum: "ROSAT-PSPC",
  _label: "AOCom",
  _mission: "ROSAT-Survey",
  _model: [3e20, "plaw", 2.0],
  _obsid: "WG931332P_N1",
  _ra: "189.0418645",
  _src_counts: "5",
  _start_date: "1990-11-30T00:00:00",
  _status: "Ok",
  _ul: "0.023794615",
  _ul_flux: 2.436568576e-13,
  _ulsig: "2",
};

describe("isoToMjd", () => {
  it("accepts both date formats HILIGT uses", () => {
    expect(isoToMjd("1858-11-17T00:00:00")).toBe(0);
    expect(isoToMjd("2000-01-01 12:00:00")).toBeCloseTo(51544.5, 6);
    expect(isoToMjd(" ")).toBeUndefined();
    expect(isoToMjd(null)).toBeUndefined();
  });
});

describe("isHiligtMission", () => {
  it("only allows the missions we query", () => {
    expect(isHiligtMission("XMMpnt")).toBe(true);
    expect(isHiligtMission("RosatSurvey")).toBe(true);
    expect(isHiligtMission("SwiftXRT")).toBe(false);
    expect(isHiligtMission("toString")).toBe(false);
  });
});

describe("parseHiligt", () => {
  it("reads detections with their band and mid-observation time", () => {
    const [p] = parseHiligt(JSON.stringify([SLEW_DETECTION]), "XMMslew");
    expect(p).toMatchObject({
      mission: "XMMslew",
      instrument: "EPIC-pn",
      obsid: "9230800003",
      band: "0.2-2",
      flux: 3.3412410288e-11,
      fluxErr: 2.3775844266e-12,
    });
    expect(p.upperLimit).toBeUndefined();
    // 2012-07-16 22:52 → 23:57 UTC
    expect(p.mjd).toBeCloseTo(56124.976, 2);
  });

  it("reads upper limits, treating blank fluxes as missing", () => {
    const [p] = parseHiligt(JSON.stringify([ROSAT_UPPER_LIMIT]), "RosatSurvey");
    expect(p.flux).toBeUndefined();
    expect(p.upperLimit).toBeCloseTo(2.4366e-13, 16);
    expect(p.sigma).toBe(2);
  });

  it("sorts by time and skips failed records", () => {
    const failed = { ...SLEW_DETECTION, _status: "Error" };
    const points = parseHiligt(
      JSON.stringify([SLEW_DETECTION, failed, ROSAT_DETECTION]),
      "RosatSurvey"
    );
    expect(points).toHaveLength(2);
    expect(points[0].mjd).toBeLessThan(points[1].mjd);
  });

  it("treats an empty or malformed body as no data", () => {
    expect(parseHiligt("", "XMMpnt")).toEqual([]);
    expect(parseHiligt("[]", "XMMpnt")).toEqual([]);
    expect(parseHiligt("<html>", "XMMpnt")).toEqual([]);
  });
});
