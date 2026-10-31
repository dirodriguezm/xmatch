import { afterEach, describe, expect, it, vi } from "vitest";

import { ApiError, apiFetch } from "./client";
import {
  bulkConeSearchResponseSchema,
  coneSearchResponseSchema,
  lightcurveResponseSchema,
} from "./schemas";

function mockFetch(status: number, body?: unknown) {
  const text = body === undefined ? "" : JSON.stringify(body);
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(status === 204 ? null : text, { status }))
  );
}

afterEach(() => vi.unstubAllGlobals());

// Trimmed from a real /conesearch answer around M31.
const CONE = [
  {
    catalog: "gaia",
    data: [
      {
        id: "Gaia DR3 381266999950756352",
        ipix: 45450570725,
        ra: 10.684112026261536,
        dec: 41.266847862838205,
        cat: "gaia",
        distance: 7.90936140098843,
      },
    ],
    index: 0,
    total: 4,
    total_in_catalog: 4,
  },
];

describe("apiFetch", () => {
  it("returns the parsed body and keeps unknown fields", async () => {
    mockFetch(200, CONE);
    const result = await apiFetch("/conesearch?ra=1", coneSearchResponseSchema);
    expect(result).toEqual(CONE);
  });

  it("returns null on 204", async () => {
    mockFetch(204);
    await expect(
      apiFetch("/conesearch", coneSearchResponseSchema)
    ).resolves.toBeNull();
  });

  it("throws a 502 naming the endpoint and field when the shape is off", async () => {
    mockFetch(200, [{ catalog: "gaia", data: [{ id: 42, ra: 1, dec: 2 }] }]);
    const error = await apiFetch(
      "/conesearch?ra=1",
      coneSearchResponseSchema
    ).catch((e: unknown) => e);
    expect(error).toBeInstanceOf(ApiError);
    expect((error as ApiError).status).toBe(502);
    expect((error as ApiError).message).toContain("/conesearch:");
    expect((error as ApiError).message).not.toContain("ra=1");
    expect((error as ApiError).message).toContain("data[0].id");
  });

  it("formats the service's validation error", async () => {
    mockFetch(400, {
      ErrValue: "x",
      Field: "RA",
      Reason: "Could not parse float.",
    });
    await expect(
      apiFetch("/conesearch", coneSearchResponseSchema)
    ).rejects.toMatchObject({
      status: 400,
      message: "RA: Could not parse float.",
    });
  });

  it("passes a plain-text error through", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => new Response("upstream down", { status: 503 }))
    );
    await expect(
      apiFetch("/conesearch", coneSearchResponseSchema)
    ).rejects.toMatchObject({ status: 503, message: "upstream down" });
  });
});

describe("schemas", () => {
  it("turns a null match list into an empty one", () => {
    const parsed = coneSearchResponseSchema.parse([
      { catalog: "gaia", data: null },
    ]);
    expect(parsed[0]?.data).toEqual([]);
  });

  it("requires the input row index in bulk answers", () => {
    expect(
      bulkConeSearchResponseSchema.safeParse([{ catalog: "gaia", data: [] }])
        .success
    ).toBe(false);
  });

  it("accepts non-detections without a magnitude", () => {
    const parsed = lightcurveResponseSchema.parse({
      detections: [
        { catalog: "neowise", mjd: 59211.2, mag: 3.87, magerr: 0.04 },
      ],
      non_detections: [{ mjd: 59300, diffmaglim: 20.1 }],
      forced_photometry: null,
    });
    expect(parsed.non_detections).toHaveLength(1);
    expect(parsed.forced_photometry).toEqual([]);
  });
});
