import { afterEach, describe, expect, it, vi } from "vitest";

import { findHipsAt, tileState } from "./useSkyImages";

afterEach(() => {
  vi.unstubAllGlobals();
});

const PS1 = "CDS/P/PanSTARRS/DR1/color-i-r-g";
const DSS = "CDS/P/DSS2/color";
const MIRRORS: Record<string, string[]> = {
  [PS1]: ["https://a.example/ps1", "https://b.example/ps1"],
  [DSS]: ["https://a.example/dss"],
};

/**
 * A registry listing MIRRORS, and tile servers answering per mirror:
 * a status, or "down" for a network failure.
 */
function stubServers(tiles: Record<string, number | "down">) {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (input: string) => {
      const url = new URL(input);
      if (url.pathname.endsWith("/MocServer/query")) {
        const urls = MIRRORS[url.searchParams.get("ID") ?? ""] ?? [];
        const record: Record<string, string> = {
          hips_order: "11",
          hips_tile_format: "jpeg",
        };
        urls.forEach((u, i) => {
          record[i ? `hips_service_url_${i}` : "hips_service_url"] = u;
        });
        return Response.json(urls.length ? [record] : []);
      }
      const base = Object.keys(tiles).find((b) => input.startsWith(b));
      const answer = base ? tiles[base] : 404;
      if (answer === "down") throw new TypeError("Failed to fetch");
      return new Response("", { status: answer });
    })
  );
}

const find = () =>
  findHipsAt([PS1, DSS], 83.8220833, -5.3911111, new AbortController().signal);

describe("findHipsAt", () => {
  it("returns the mirror serving the tile, and that tile", async () => {
    stubServers({ "https://a.example/ps1": 404, "https://b.example/ps1": 200 });
    expect(await find()).toEqual({
      id: PS1,
      url: "https://b.example/ps1",
      tileUrl: "https://b.example/ps1/Norder11/Dir21950000/Npix21954213.jpg",
    });
  });

  it("falls back to the next survey where the first has a gap", async () => {
    stubServers({ "https://a.example/dss": 200 });
    expect(await find()).toMatchObject({ id: DSS });
  });

  it("reports a gap in the sky as uncovered, not as a failure", async () => {
    stubServers({});
    expect(await find()).toBeNull();
  });

  it("fails only when no server answered", async () => {
    stubServers({
      "https://a.example/ps1": "down",
      "https://b.example/ps1": 503,
      "https://a.example/dss": "down",
    });
    await expect(find()).rejects.toThrow();
  });

  it("fails when the registry is unreachable", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => {
        throw new TypeError("Failed to fetch");
      })
    );
    await expect(find()).rejects.toThrow();
  });
});

describe("tileState", () => {
  it("follows the probe once it settles", () => {
    expect(tileState("loaded", 0)).toBe("loaded");
    expect(tileState("failed", 0)).toBe("failed");
    expect(tileState("loaded", 60_000)).toBe("loaded");
  });

  it("calls a pending probe slow, then stops claiming anything", () => {
    expect(tileState("pending", 0)).toBe("loading");
    expect(tileState("pending", 12_000)).toBe("slow");
    expect(tileState("pending", 30_000)).toBe("unknown");
  });
});
