import { describe, expect, it, vi } from "vitest";

import {
  findRandomGaiaObject,
  firstConeHit,
  pickRandomFeatured,
  RANDOM_MAX_POSITIONS,
  RANDOM_RADII_ARCSEC,
  randomSkyPosition,
} from "./randomObject";

describe("randomSkyPosition", () => {
  it("maps the unit interval onto the full sphere", () => {
    expect(randomSkyPosition(() => 0)).toEqual({ ra: 0, dec: -90 });
    const mid = randomSkyPosition(() => 0.5);
    expect(mid.ra).toBe(180);
    expect(mid.dec).toBeCloseTo(0);
    expect(randomSkyPosition(() => 0.999999).dec).toBeGreaterThan(89);
  });
});

describe("pickRandomFeatured", () => {
  it("never indexes past the end", () => {
    expect(pickRandomFeatured(() => 0.9999999)).toBeDefined();
    expect(pickRandomFeatured(() => 0)).toBeDefined();
  });
});

describe("firstConeHit", () => {
  it("returns the closest row", () => {
    const payload = [
      {
        catalog: "gaia",
        data: [
          { id: "far", ra: 1, dec: 1, distance: 9 },
          { id: "near", ra: 2, dec: 2, distance: 1 },
        ],
      },
    ];
    expect(firstConeHit(payload)).toEqual({
      objectId: "near",
      catalog: "gaia",
      ra: 2,
      dec: 2,
    });
  });

  it("handles empty or odd payloads", () => {
    expect(firstConeHit(null)).toBeNull();
    expect(firstConeHit([])).toBeNull();
    expect(firstConeHit([{ data: [] }])).toBeNull();
  });
});

describe("findRandomGaiaObject", () => {
  const json = (body: unknown) =>
    Promise.resolve(new Response(JSON.stringify(body), { status: 200 }));

  it("grows the radius until a hit", async () => {
    const fetcher = vi
      .fn()
      .mockImplementationOnce(() => json([]))
      .mockImplementationOnce(() =>
        json([{ data: [{ id: "x", ra: 1, dec: 2, distance: 3 }] }])
      );
    const hit = await findRandomGaiaObject(fetcher, () => 0.3);
    expect(hit?.objectId).toBe("x");
    expect(fetcher).toHaveBeenCalledTimes(2);
    expect(fetcher.mock.calls[1][0]).toContain(
      `radius=${RANDOM_RADII_ARCSEC[1]}`
    );
  });

  it("gives up after every position and radius", async () => {
    const fetcher = vi.fn().mockImplementation(() => json([]));
    expect(await findRandomGaiaObject(fetcher, () => 0.3)).toBeNull();
    expect(fetcher).toHaveBeenCalledTimes(
      RANDOM_MAX_POSITIONS * RANDOM_RADII_ARCSEC.length
    );
  });
});
