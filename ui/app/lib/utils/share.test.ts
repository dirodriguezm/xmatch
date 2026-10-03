import { describe, expect, it } from "vitest";

import {
  absoluteUrl,
  buildEmbedHtml,
  buildShareTargets,
  catalogFromObjectId,
  formatRaDec,
  pickPhotometry,
} from "./share";

describe("absoluteUrl", () => {
  it("joins origin and path without doubled slashes", () => {
    expect(absoluteUrl("/object/x", "https://a.b/")).toBe(
      "https://a.b/object/x"
    );
    expect(absoluteUrl("search?ra=1", "https://a.b")).toBe(
      "https://a.b/search?ra=1"
    );
  });
});

describe("buildEmbedHtml", () => {
  it("produces a lazy full-width iframe with escaped attributes", () => {
    const html = buildEmbedHtml("https://a.b/embed?x=1&y=2", 'Obj "1"');
    expect(html).toContain('src="https://a.b/embed?x=1&amp;y=2"');
    expect(html).toContain('title="Obj &quot;1&quot; on XWave"');
    expect(html).toContain('width="100%"');
    expect(html).toContain('height="420"');
    expect(html).toContain('loading="lazy"');
  });
});

describe("buildShareTargets", () => {
  it("encodes the url into each target", () => {
    const targets = buildShareTargets("https://a.b/o?c=1&d=2", "Gaia DR3 1");
    expect(targets.map((t) => t.id)).toEqual(["x", "bluesky", "email"]);
    const x = new URL(targets[0].href);
    expect(x.searchParams.get("url")).toBe("https://a.b/o?c=1&d=2");
    const bsky = new URL(targets[1].href);
    expect(bsky.searchParams.get("text")).toBe(
      "Gaia DR3 1 on XWave https://a.b/o?c=1&d=2"
    );
    expect(targets[2].href).toMatch(/^mailto:\?subject=/);
  });
});

describe("catalogFromObjectId", () => {
  it("recognises catalog prefixes", () => {
    expect(catalogFromObjectId("Gaia DR3 381267755864434048")?.slug).toBe(
      "gaia"
    );
    expect(catalogFromObjectId("J004243.11+411609.2")?.slug).toBe("allwise");
    expect(catalogFromObjectId("1eRASS J004243.1+411609")?.slug).toBe(
      "erosita"
    );
    expect(catalogFromObjectId("M31")).toBeNull();
  });
});

describe("formatRaDec", () => {
  it("signs declination", () => {
    expect(formatRaDec(10.5, -3.25, 2)).toBe("RA 10.50°, Dec −3.25°");
    expect(formatRaDec(10.5, 3.25, 2)).toBe("RA 10.50°, Dec +3.25°");
  });
});

describe("pickPhotometry", () => {
  it("picks known bands and skips zero-filled values", () => {
    const rows = pickPhotometry("allwise", {
      w1mpro: 12.1,
      w1sigmpro: 0.02,
      w2mpro: 0,
      j_m_2mass: 13.4,
    });
    expect(rows).toEqual([
      { band: "W1", value: 12.1, error: 0.02, unit: "mag" },
      { band: "J", value: 13.4, error: undefined, unit: "mag" },
    ]);
  });

  it("falls back to flux fields for unknown catalogs", () => {
    const rows = pickPhotometry("erosita", {
      ml_flux_1: 1e-14,
      ml_flux_err_1: 1e-15,
      ra: 1,
    });
    expect(rows.map((r) => r.band)).toEqual(["ml_flux_1"]);
  });

  it("handles missing records", () => {
    expect(pickPhotometry("gaia", null)).toEqual([]);
  });
});
