import { describe, expect, it } from "vitest";

import {
  buildSnippet,
  bulkConeSearchRequest,
  coneSearchRequest,
  requestUrl,
} from "./snippets";

describe("snippets", () => {
  const cone = coneSearchRequest({
    ra: 10.5,
    dec: -3.25,
    radius: 2,
    catalog: "gaia",
  });

  it("builds the public API URL for GET requests", () => {
    expect(requestUrl(cone)).toBe(
      "https://xwave-astro.udp.cl/v1/conesearch?ra=10.5&dec=-3.25&radius=2&catalog=gaia"
    );
  });

  it("renders every language for a GET request", () => {
    expect(buildSnippet(cone, "curl")).toContain('--data-urlencode "radius=2"');
    expect(buildSnippet(cone, "python")).toContain('"catalog": "gaia",');
    expect(buildSnippet(cone, "javascript")).toContain('"dec": "-3.25"');
  });

  it("renders POST bodies with query parameters and Python literals", () => {
    const bulk = bulkConeSearchRequest({ ra: [1, 2], dec: [3, 4], radius: 5 });
    const py = buildSnippet(bulk, "python");
    expect(py).toContain('"ra": [1, 2],');
    expect(py).toContain('"radius": 5,');
    expect(py).toContain("requests.post");
    expect(buildSnippet(bulk, "curl")).toContain("nneighbor=1");
    expect(requestUrl(bulk)).toBe(
      "https://xwave-astro.udp.cl/v1/bulk-conesearch?radius=5&catalog=all&nneighbor=1"
    );
  });
});
