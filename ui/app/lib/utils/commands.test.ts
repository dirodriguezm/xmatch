import { describe, expect, it } from "vitest";

import { FOOTER_GROUPS, NAV_ITEMS } from "@/app/lib/constants/site";

import {
  buildNavCommands,
  buildSearchHref,
  cycleIndex,
  filterCommands,
  fuzzyScore,
  isMacPlatform,
  isTypingTarget,
  type PaletteCommand,
  paletteShortcutLabel,
} from "./commands";

describe("buildNavCommands", () => {
  const commands = buildNavCommands(NAV_ITEMS, FOOTER_GROUPS);

  it("has one command per distinct href, home first", () => {
    const hrefs = commands.map((c) => c.href);
    expect(hrefs[0]).toBe("/");
    expect(new Set(hrefs).size).toBe(hrefs.length);
    expect(hrefs).toContain("/about");
    expect(hrefs).toContain("/llms.txt");
  });

  it("keeps footer labels as keywords and flags external links", () => {
    const dev = commands.find((c) => c.href === "/developers")!;
    expect(dev.label).toBe("API");
    expect(dev.keywords).toContain("API playground");
    expect(dev.hint).toBe("G D");
    const gh = commands.find((c) => c.label === "GitHub")!;
    expect(gh.external).toBe(true);
  });
});

describe("fuzzyScore", () => {
  it("ranks exact > prefix > word start > substring > subsequence", () => {
    const exact = fuzzyScore("bulk", "bulk")!;
    const prefix = fuzzyScore("bul", "bulk cross-match")!;
    const word = fuzzyScore("cross", "bulk cross-match")!;
    const sub = fuzzyScore("ross", "bulk cross-match")!;
    const seq = fuzzyScore("bcm", "bulk cross-match")!;
    expect(exact).toBeGreaterThan(prefix);
    expect(prefix).toBeGreaterThan(word);
    expect(word).toBeGreaterThan(sub);
    expect(sub).toBeGreaterThan(seq);
    expect(seq).toBeGreaterThan(0);
  });

  it("returns null when characters are missing or out of order", () => {
    expect(fuzzyScore("xyz", "search")).toBeNull();
    expect(fuzzyScore("hcraes", "search")).toBeNull();
  });

  it("is case-insensitive and safe with regex characters", () => {
    expect(fuzzyScore("LLMS.TXT", "llms.txt")).toBe(1000);
    expect(fuzzyScore("(", "a (b)")).not.toBeNull();
  });
});

describe("filterCommands", () => {
  const cmds: PaletteCommand[] = [
    { id: "a", label: "Search", group: "Navigate" },
    { id: "b", label: "Bulk", group: "Navigate", keywords: ["csv upload"] },
    { id: "c", label: "API", group: "Navigate", keywords: ["playground"] },
  ];

  it("returns everything in order for an empty query", () => {
    expect(filterCommands(cmds, "  ")).toEqual(cmds);
  });

  it("matches keywords", () => {
    expect(filterCommands(cmds, "csv").map((c) => c.id)).toEqual(["b"]);
    expect(filterCommands(cmds, "play").map((c) => c.id)).toEqual(["c"]);
  });

  it("ranks label hits first", () => {
    expect(filterCommands(cmds, "s")[0].id).toBe("a");
  });
});

describe("buildSearchHref", () => {
  it("uses the landing page's query params", () => {
    const href = buildSearchHref(10.5, -3.25, [
      { catalog: "gaia", radius: 3, unit: "arcsec", enabled: true },
    ]);
    expect(href).toBe(
      "/search?ra=10.5&dec=-3.25&catalogRadii=gaia%3A3%3Aarcsec%3A1"
    );
  });

  it("defaults to every catalog", () => {
    const params = new URL(buildSearchHref(1, 2), "http://x").searchParams;
    expect(params.get("catalogRadii")).toContain("gaia:");
    expect(params.get("catalogRadii")).toContain("erosita:");
  });
});

describe("helpers", () => {
  it("cycles indices in both directions", () => {
    expect(cycleIndex(0, -1, 3)).toBe(2);
    expect(cycleIndex(2, 1, 3)).toBe(0);
    expect(cycleIndex(0, 1, 0)).toBe(0);
  });

  it("detects typing targets", () => {
    expect(isTypingTarget({ tagName: "input" })).toBe(true);
    expect(isTypingTarget({ tagName: "TEXTAREA" })).toBe(true);
    expect(isTypingTarget({ tagName: "DIV", isContentEditable: true })).toBe(
      true
    );
    expect(
      isTypingTarget({ tagName: "DIV", getAttribute: () => "textbox" })
    ).toBe(true);
    expect(isTypingTarget({ tagName: "BUTTON" })).toBe(false);
    expect(isTypingTarget(null)).toBe(false);
  });

  it("labels the shortcut per platform", () => {
    expect(isMacPlatform("MacIntel")).toBe(true);
    expect(isMacPlatform("Win32")).toBe(false);
    expect(paletteShortcutLabel(true)).toBe("⌘K");
    expect(paletteShortcutLabel(false)).toBe("Ctrl K");
  });
});
