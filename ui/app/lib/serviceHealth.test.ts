import { describe, expect, it } from "vitest";

import {
  getFailingServices,
  isClientError,
  isLocalRefusal,
  reportQueryError,
  reportServiceResult,
  serviceForQueryKey,
} from "./serviceHealth";

describe("serviceForQueryKey", () => {
  it("maps proxy-backed queries to their service", () => {
    expect(serviceForQueryKey(["simbad", { ra: 1, dec: 2 }])).toBe("SIMBAD");
    expect(serviceForQueryKey(["crts-lightcurve", 1, 2])).toBe(
      "Catalina light curves"
    );
  });

  it("counts the per-position sky image lookup as Sky images", () => {
    expect(serviceForQueryKey(["hips-at", "CDS/P/DSS2/color", 10, 41])).toBe(
      "Sky images"
    );
  });

  it("ignores queries that aren't external calls", () => {
    expect(serviceForQueryKey(["status-probe", "simbad"])).toBeNull();
    expect(serviceForQueryKey([42])).toBeNull();
    expect(serviceForQueryKey([])).toBeNull();
  });
});

describe("reportServiceResult", () => {
  it("tracks the latest outcome per service", () => {
    reportServiceResult("SIMBAD", false);
    reportServiceResult("SIMBAD", false);
    reportServiceResult("CRTS", false);
    expect(getFailingServices()).toEqual(["SIMBAD", "CRTS"]);

    reportServiceResult("SIMBAD", true);
    expect(getFailingServices()).toEqual(["CRTS"]);

    reportServiceResult(null, false);
    reportServiceResult("CRTS", true);
    expect(getFailingServices()).toEqual([]);
  });
});

describe("isClientError", () => {
  it("treats a 4xx as the service answering", () => {
    expect(isClientError({ status: 400 })).toBe(true);
    expect(isClientError(Object.assign(new Error("x"), { status: 422 }))).toBe(
      true
    );
  });

  it("treats 5xx, network failures and odd values as outages", () => {
    expect(isClientError({ status: 502 })).toBe(false);
    expect(isClientError(new TypeError("Failed to fetch"))).toBe(false);
    expect(isClientError({ status: "400" })).toBe(false);
    expect(isClientError(null)).toBe(false);
  });
});

describe("isLocalRefusal", () => {
  it("spots a request refused before it was sent", () => {
    expect(
      isLocalRefusal(
        Object.assign(new Error("x"), { status: 400, local: true })
      )
    ).toBe(true);
  });

  it("leaves real answers from the service alone", () => {
    expect(isLocalRefusal({ status: 400 })).toBe(false);
    expect(isLocalRefusal({ status: 400, local: false })).toBe(false);
    expect(isLocalRefusal(new TypeError("Failed to fetch"))).toBe(false);
    expect(isLocalRefusal(null)).toBe(false);
  });
});

describe("reportQueryError", () => {
  const conesearch = ["conesearch", 10, 41, "gaia"];

  it("keeps an outage flagged through a request refused before sending", () => {
    reportQueryError(
      conesearch,
      Object.assign(new Error("x"), { status: 502 })
    );
    expect(getFailingServices()).toEqual(["XWave API"]);

    // A radius over the limit never reaches the API, so it proves nothing.
    reportQueryError(
      conesearch,
      Object.assign(new Error("x"), { status: 400, local: true })
    );
    expect(getFailingServices()).toEqual(["XWave API"]);

    // The API itself answering 400 does show it is up.
    reportQueryError(
      conesearch,
      Object.assign(new Error("x"), { status: 400 })
    );
    expect(getFailingServices()).toEqual([]);
  });
});
