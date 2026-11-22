import { afterEach, describe, expect, it, vi } from "vitest";

import { isLocalRefusal } from "@/app/lib/serviceHealth";

import {
  ConeSearchError,
  fetchConeSearch,
  isRetryableConeError,
} from "./useConeSearch";

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("fetchConeSearch", () => {
  it("refuses a radius over the limit without a request", async () => {
    const fetchMock = vi.fn();
    vi.stubGlobal("fetch", fetchMock);
    const error = await fetchConeSearch({
      ra: 10,
      dec: 41,
      radius: 300,
    }).catch((e: unknown) => e);
    expect(fetchMock).not.toHaveBeenCalled();
    expect(error).toBeInstanceOf(ConeSearchError);
    // Never sent, so it must not count for or against the service.
    expect(isLocalRefusal(error)).toBe(true);
  });

  it("marks the service's own 4xx as sent", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => new Response("{}", { status: 400 }))
    );
    const error = await fetchConeSearch({ ra: 10, dec: 41, radius: 5 }).catch(
      (e: unknown) => e
    );
    expect(error).toMatchObject({ status: 400 });
    expect(isLocalRefusal(error)).toBe(false);
  });
});

describe("isRetryableConeError", () => {
  it("retries network errors and 5xx", () => {
    expect(isRetryableConeError(new TypeError("Failed to fetch"))).toBe(true);
    expect(isRetryableConeError(new ConeSearchError("x", 502))).toBe(true);
  });

  it("does not retry a request the service rejected", () => {
    expect(isRetryableConeError(new ConeSearchError("x", 400))).toBe(false);
    expect(isRetryableConeError(new ConeSearchError("x", 422, true))).toBe(
      false
    );
  });
});
