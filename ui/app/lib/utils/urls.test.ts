import { describe, expect, it } from "vitest";

import { buildAdsObjectUrl } from "./urls";

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
