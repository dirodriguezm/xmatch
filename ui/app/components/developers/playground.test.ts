import { describe, expect, it } from "vitest";

import { requestUrl } from "@/app/lib/utils/snippets";

import {
  buildPlaygroundRequest,
  DEFAULT_STATE,
  formatBytes,
  parsePositions,
  proxyUrl,
  summarizeResult,
} from "./playground";

describe("parsePositions", () => {
  it("parses degrees and sexagesimal, skipping comments and blanks", () => {
    const p = parsePositions("10.5 -3.25 # a\n\n05:34:31.9 +22:00:52\nfoo bar");
    expect(p.ra[0]).toBe(10.5);
    expect(p.dec[0]).toBe(-3.25);
    expect(p.ra[1]).toBeCloseTo(83.63292, 4);
    expect(p.invalid).toEqual([4]);
  });
});

describe("buildPlaygroundRequest", () => {
  it("builds a cone search with optional getMetadata", () => {
    const req = buildPlaygroundRequest("conesearch", {
      ...DEFAULT_STATE,
      getMetadata: true,
    });
    expect(requestUrl(req)).toBe(
      "https://xwave-astro.udp.cl/v1/conesearch?ra=10.6847&dec=41.269&radius=10&catalog=all&nneighbor=5&getMetadata=true"
    );
  });

  it("adds the light-curve catalog only when filtered", () => {
    const all = buildPlaygroundRequest("lightcurve", DEFAULT_STATE);
    expect(requestUrl(all)).not.toContain("catalog");
    const ztf = buildPlaygroundRequest("lightcurve", {
      ...DEFAULT_STATE,
      lightcurveCatalog: "ztf",
    });
    expect(requestUrl(ztf)).toContain("catalog=ztf");
  });

  it("builds a bulk POST body from the positions textarea", () => {
    const req = buildPlaygroundRequest("bulk", DEFAULT_STATE);
    expect(req.method).toBe("POST");
    if (req.method !== "POST") return;
    expect((req.body.ra as number[]).length).toBe(3);
    expect(req.body.radius).toBe(5);
  });
});

describe("proxyUrl", () => {
  it("maps requests to the Next.js proxy routes", () => {
    expect(proxyUrl(buildPlaygroundRequest("metadata", DEFAULT_STATE))).toBe(
      "/api/metadata?id=0098p408_ac51-043708&catalog=allwise"
    );
    expect(proxyUrl(buildPlaygroundRequest("bulk", DEFAULT_STATE))).toBe(
      "/api/bulk-conesearch"
    );
  });
});

describe("formatting", () => {
  it("formats sizes", () => {
    expect(formatBytes(500)).toBe("500 B");
    expect(formatBytes(2048)).toBe("2.0 KB");
  });

  it("summarizes result shapes", () => {
    expect(summarizeResult([{ data: [1, 2] }, { data: [3] }])).toBe(
      "2 groups, 3 rows"
    );
    expect(summarizeResult({ detections: [1], non_detections: [], x: 1 })).toBe(
      "1 detections, 0 non detections"
    );
    expect(summarizeResult({ a: 1, b: 2 })).toBe("2 fields");
    expect(summarizeResult("x")).toBeNull();
  });
});
