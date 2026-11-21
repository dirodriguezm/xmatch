import { describe, expect, it } from "vitest";

import {
  getFailingServices,
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
