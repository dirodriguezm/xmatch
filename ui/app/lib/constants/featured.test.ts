import { describe, expect, it } from "vitest";

import { FEATURED_OBJECTS, isoWeek, objectOfTheWeek } from "./featured";

describe("isoWeek", () => {
  it("matches known ISO weeks", () => {
    expect(isoWeek(new Date(Date.UTC(2026, 0, 1)))).toBe(1); // Thu
    expect(isoWeek(new Date(Date.UTC(2021, 0, 3)))).toBe(53); // Sun
    expect(isoWeek(new Date(Date.UTC(2024, 11, 30)))).toBe(1);
    expect(isoWeek(new Date(Date.UTC(2026, 8, 26)))).toBe(39);
  });
});

describe("objectOfTheWeek", () => {
  it("is stable within a week and rotates across weeks", () => {
    const mon = objectOfTheWeek(new Date(Date.UTC(2026, 8, 21)));
    const sun = objectOfTheWeek(new Date(Date.UTC(2026, 8, 27)));
    const next = objectOfTheWeek(new Date(Date.UTC(2026, 8, 28)));
    expect(mon).toBe(sun);
    expect(next).not.toBe(mon);
  });

  it("has unique slugs and valid coordinates", () => {
    const slugs = new Set(FEATURED_OBJECTS.map((o) => o.slug));
    expect(slugs.size).toBe(FEATURED_OBJECTS.length);
    for (const o of FEATURED_OBJECTS) {
      expect(o.ra).toBeGreaterThanOrEqual(0);
      expect(o.ra).toBeLessThan(360);
      expect(Math.abs(o.dec)).toBeLessThanOrEqual(90);
    }
  });
});
