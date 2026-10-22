import { describe, expect, it } from "vitest";

import { buildChangelogRss, CHANGELOG, rssDate } from "./changelog";

describe("CHANGELOG", () => {
  it("is sorted newest first with unique ids", () => {
    const dates = CHANGELOG.map((e) => e.date);
    expect([...dates].sort().reverse()).toEqual(dates);
    expect(new Set(CHANGELOG.map((e) => e.id)).size).toBe(CHANGELOG.length);
  });
});

describe("buildChangelogRss", () => {
  const xml = buildChangelogRss(
    [
      {
        id: "a",
        date: "2026-09-26",
        title: "Fish & <chips>",
        summary: "s",
        tags: ["New"],
        items: ["one"],
      },
    ],
    "https://example.org/"
  );

  it("produces an RSS 2.0 channel with escaped items", () => {
    expect(xml).toMatch(/^<\?xml version="1.0"/);
    expect(xml).toContain('<rss version="2.0"');
    expect(xml).toContain("<title>Fish &amp; &lt;chips&gt;</title>");
    expect(xml).toContain("<link>https://example.org/changelog#a</link>");
    expect(xml).toContain("<category>New</category>");
    expect(xml).toContain("&lt;li&gt;one&lt;/li&gt;");
    expect(xml).not.toContain("<li>");
  });

  it("uses RFC 822 dates", () => {
    expect(rssDate("2026-09-26")).toBe("Sat, 26 Sep 2026 12:00:00 GMT");
  });
});
