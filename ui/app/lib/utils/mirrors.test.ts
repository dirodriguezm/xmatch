import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { fetchFromMirrors, resetPreferredMirrors } from "./mirrors";

const BASES = ["https://a.test", "https://b.test", "https://c.test"];

describe("fetchFromMirrors", () => {
  const fetchMock = vi.fn<typeof fetch>();

  beforeEach(() => {
    resetPreferredMirrors();
    vi.stubGlobal("fetch", fetchMock);
    vi.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    fetchMock.mockReset();
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it("returns the first mirror when it answers OK", async () => {
    fetchMock.mockResolvedValueOnce(new Response("ok"));

    const { response, host } = await fetchFromMirrors(
      BASES,
      "/q?x=1",
      {},
      1000
    );

    expect(await response.text()).toBe("ok");
    expect(host).toBe("a.test");
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(fetchMock.mock.calls[0]?.[0]).toBe("https://a.test/q?x=1");
  });

  it("skips mirrors that error or answer non-2xx", async () => {
    fetchMock
      .mockResolvedValueOnce(new Response("down", { status: 404 }))
      .mockRejectedValueOnce(new TypeError("fetch failed"))
      .mockResolvedValueOnce(new Response("mirror"));

    const { response, host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

    expect(await response.text()).toBe("mirror");
    expect(host).toBe("c.test");
    expect(fetchMock.mock.calls.map((c) => c[0])).toEqual([
      "https://a.test/q",
      "https://b.test/q",
      "https://c.test/q",
    ]);
  });

  it("returns the last HTTP response when no mirror is OK", async () => {
    fetchMock
      .mockResolvedValueOnce(new Response("", { status: 404 }))
      .mockResolvedValueOnce(new Response("", { status: 503 }))
      .mockRejectedValueOnce(new TypeError("fetch failed"));

    const { response, host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

    expect(response.status).toBe(503);
    expect(host).toBe("b.test");
  });

  it("throws when every mirror fails without answering", async () => {
    fetchMock.mockRejectedValue(new TypeError("fetch failed"));

    await expect(fetchFromMirrors(BASES, "/q", {}, 1000)).rejects.toThrow(
      "fetch failed"
    );
  });

  it("reports the host the request was redirected to", async () => {
    const redirected = new Response("ok");
    Object.defineProperty(redirected, "url", {
      value: "https://mirror.test/q",
    });
    fetchMock.mockResolvedValueOnce(redirected);

    const { host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

    expect(host).toBe("mirror.test");
  });

  it("gives each mirror its own timeout signal", async () => {
    fetchMock.mockResolvedValue(new Response("", { status: 500 }));

    await fetchFromMirrors(BASES, "/q", {}, 1000);

    const signals = fetchMock.mock.calls.map((c) => c[1]?.signal);
    expect(new Set(signals).size).toBe(3);
  });

  it("uses a mirror's own timeout over the default", async () => {
    const timeout = vi.spyOn(AbortSignal, "timeout");
    fetchMock.mockResolvedValue(new Response("", { status: 500 }));

    await fetchFromMirrors(
      ["https://a.test", { base: "https://b.test", timeoutMs: 90_000 }],
      "/q",
      {},
      1000
    );

    expect(timeout.mock.calls.map((c) => c[0])).toEqual([1000, 90_000]);
  });

  describe("after the primary fails", () => {
    const urls = () => fetchMock.mock.calls.map((c) => c[0]);

    beforeEach(async () => {
      vi.useFakeTimers({ toFake: ["Date"] });
      fetchMock
        .mockResolvedValueOnce(new Response("", { status: 404 }))
        .mockResolvedValueOnce(new Response("ok"));
      await fetchFromMirrors(BASES, "/q", {}, 1000);
      fetchMock.mockClear();
    });

    it("tries the mirror that answered first", async () => {
      fetchMock.mockResolvedValueOnce(new Response("ok"));

      const { host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

      expect(host).toBe("b.test");
      expect(urls()).toEqual(["https://b.test/q"]);
    });

    it("gives the primary another chance after 10 minutes", async () => {
      vi.setSystemTime(Date.now() + 10 * 60_000 + 1);
      fetchMock.mockResolvedValueOnce(new Response("ok"));

      const { host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

      expect(host).toBe("a.test");
    });

    it("falls back to the others when that mirror fails too", async () => {
      fetchMock
        .mockRejectedValueOnce(new TypeError("fetch failed"))
        .mockResolvedValueOnce(new Response("ok"));

      const { host } = await fetchFromMirrors(BASES, "/q", {}, 1000);

      expect(host).toBe("a.test");
      expect(urls()).toEqual(["https://b.test/q", "https://a.test/q"]);
    });
  });
});
