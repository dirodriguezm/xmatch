import { describe, expect, it } from "vitest";

import { buildAdsObjectUrl, buildObjectUrl } from "./urls";

describe("buildAdsObjectUrl", () => {
  it("searches ADS by object name, newest first", () => {
    expect(buildAdsObjectUrl("3C 273")).toBe(
      "https://ui.adsabs.harvard.edu/search/q=object%3A%223C%20273%22&sort=date%20desc"
    );
  });

  it("escapes SIMBAD names with special characters", () => {
    expect(buildAdsObjectUrl("* alf Ori")).toContain(
      "object%3A%22*%20alf%20Ori%22"
    );
  });
});

describe("buildObjectUrl", () => {
  it("encodes the id and catalog", () => {
    expect(buildObjectUrl("Gaia DR3 1", "gaia")).toBe(
      "/object/Gaia%20DR3%201?catalog=gaia"
    );
  });

  it("tags links from the bulk page", () => {
    expect(buildObjectUrl("x", "allwise", "bulk")).toBe(
      "/object/x?catalog=allwise&from=bulk"
    );
  });
});
