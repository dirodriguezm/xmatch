import { describe, expect, it } from "vitest";

import {
  buildFeedbackIssue,
  buildFlagIssue,
  buildIssueUrl,
  findFlag,
  FLAG_STORAGE_KEY,
  readFlags,
  saveFlag,
} from "./githubIssue";

function memoryStorage(initial: Record<string, string> = {}) {
  const data = { ...initial };
  return {
    getItem: (k: string) => (k in data ? data[k] : null),
    setItem: (k: string, v: string) => {
      data[k] = v;
    },
    data,
  };
}

describe("buildIssueUrl", () => {
  it("encodes title, body and labels as query params", () => {
    const url = new URL(
      buildIssueUrl(
        { title: "A & B", body: "line 1\nline 2", labels: ["bug", "data"] },
        "https://github.com/o/r"
      )
    );
    expect(url.origin + url.pathname).toBe("https://github.com/o/r/issues/new");
    expect(url.searchParams.get("title")).toBe("A & B");
    expect(url.searchParams.get("body")).toBe("line 1\nline 2");
    expect(url.searchParams.get("labels")).toBe("bug,data");
  });

  it("omits labels when none are given", () => {
    const url = new URL(buildIssueUrl({ title: "t", body: "b" }));
    expect(url.searchParams.has("labels")).toBe(false);
  });
});

describe("buildFlagIssue", () => {
  it("includes object, catalog, coordinates, reason and page", () => {
    const issue = buildFlagIssue({
      objectId: "123",
      catalog: "gaia",
      ra: 10.1234567,
      dec: -5.5,
      reason: "blended",
      note: "  two stars  ",
      pageUrl: "https://x/object/123",
    });
    expect(issue.title).toBe(
      "[Match report] Blended / confused source: gaia 123"
    );
    expect(issue.body).toContain("`123`");
    expect(issue.body).toContain("10.123457, -5.500000");
    expect(issue.body).toContain("https://x/object/123");
    expect(issue.body).toContain("two stars");
    expect(issue.labels).toEqual(["match-report"]);
  });

  it("truncates long notes", () => {
    const issue = buildFlagIssue({
      objectId: "1",
      catalog: "gaia",
      ra: 0,
      dec: 0,
      reason: "other",
      note: "x".repeat(5000),
    });
    expect(issue.body.length).toBeLessThan(2400);
  });
});

describe("buildFeedbackIssue", () => {
  it("prefixes the title and maps labels by kind", () => {
    const issue = buildFeedbackIssue({
      kind: "bug",
      title: " Search hangs ",
      description: "It hangs",
      steps: "1. search",
    });
    expect(issue.title).toBe("[Bug report] Search hangs");
    expect(issue.labels).toEqual(["bug"]);
    expect(issue.body).toContain("### Steps to reproduce");
  });
});

describe("flag storage", () => {
  it("round-trips and replaces previous flags for the same object", () => {
    const s = memoryStorage();
    saveFlag(s, "Gaia", "1", "other", new Date("2026-01-01T00:00:00Z"));
    saveFlag(s, "gaia", "1", "blended", new Date("2026-01-02T00:00:00Z"));
    saveFlag(s, "allwise", "2", "artefact");
    expect(readFlags(s)).toHaveLength(2);
    expect(findFlag(s, "gaia", "1")?.reason).toBe("blended");
    expect(findFlag(s, "erosita", "1")).toBeUndefined();
  });

  it("tolerates corrupt or missing storage", () => {
    expect(readFlags(undefined)).toEqual([]);
    expect(readFlags(memoryStorage({ [FLAG_STORAGE_KEY]: "{bad" }))).toEqual(
      []
    );
    expect(readFlags(memoryStorage({ [FLAG_STORAGE_KEY]: "{}" }))).toEqual([]);
  });
});
