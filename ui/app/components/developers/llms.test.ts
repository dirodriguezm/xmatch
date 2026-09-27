import { describe, expect, it } from "vitest";

import { buildLlmsFullTxt, buildLlmsTxt, LLMS_EXAMPLES } from "./llms";

describe("llms.txt", () => {
  const txt = buildLlmsTxt();

  it("follows the llmstxt.org layout", () => {
    const lines = txt.split("\n");
    expect(lines[0]).toMatch(/^# XWave$/);
    expect(lines[2]).toMatch(/^> /);
    expect(txt).toContain("## API endpoints");
    expect(txt).toContain("## Optional");
  });

  it("lists every endpoint and catalog", () => {
    for (const path of [
      "/conesearch",
      "/metadata",
      "/lightcurve",
      "/bulk-conesearch",
      "/bulk-metadata",
    ]) {
      expect(txt).toContain(path);
    }
    expect(txt).toContain("catalog=gaia");
    expect(txt).toContain("[AllWISE Source Catalog]");
    expect(txt).not.toContain("~~");
  });

  it("uses public API example URLs", () => {
    expect(LLMS_EXAMPLES.conesearch).toMatch(
      /^https:\/\/xwave-astro\.udp\.cl\/v1\/conesearch\?ra=/
    );
  });
});

describe("llms-full.txt", () => {
  it("documents every endpoint section", () => {
    const full = buildLlmsFullTxt();
    expect(full).toContain("## GET /conesearch");
    expect(full).toContain("## POST /bulk-conesearch");
    expect(full).toContain("## GET /metadata");
    expect(full).toContain("## POST /bulk-metadata");
    expect(full).toContain("## GET /lightcurve");
    expect(full).not.toContain("mcp");
  });
});
