/**
 * Pure logic behind the ⌘K command palette and global keyboard shortcuts:
 * building the command list, fuzzy matching, and key handling helpers.
 */

import { buildDefaultCatalogConfigs } from "@/app/lib/constants/catalogs";
import {
  type CatalogRadiusConfig,
  encodeCatalogRadii,
} from "@/app/lib/constants/search";
import type { NavItem } from "@/app/lib/constants/site";
import { isExternalHref } from "@/app/lib/constants/site";

export type CommandGroup = "Search" | "Navigate" | "Actions";

export interface PaletteCommand {
  id: string;
  label: string;
  group: CommandGroup;
  /** Destination for navigation commands. */
  href?: string;
  /** Opens in a new tab instead of client-side navigation. */
  external?: boolean;
  /** Extra words that should match, e.g. "api" for the playground. */
  keywords?: string[];
  /** Right-aligned hint, e.g. a shortcut or the path. */
  hint?: string;
}

/** "g then <key>" navigation shortcuts. */
export const GO_SHORTCUTS: Record<string, { href: string; label: string }> = {
  h: { href: "/", label: "Home" },
  s: { href: "/search", label: "Search" },
  b: { href: "/bulk", label: "Bulk" },
  e: { href: "/explore", label: "Explore" },
  d: { href: "/developers", label: "Developers" },
  l: { href: "/learn", label: "Learn" },
};

/** How long after "g" the second key still counts, in ms. */
export const GO_SEQUENCE_TIMEOUT_MS = 1200;

const EXTRA_KEYWORDS: Record<string, string[]> = {
  "/": ["home", "landing", "start"],
  "/search": ["cone", "crossmatch", "cross-match", "find"],
  "/bulk": ["upload", "csv", "list", "batch"],
  "/developers": ["api", "playground", "curl", "python", "code", "docs"],
  "/learn": ["faq", "help", "tutorial", "guide"],
  "/llms.txt": ["llm", "ai", "agent"],
  "/status": ["uptime", "health"],
};

/**
 * One navigation command per distinct page, header items first. Labels from
 * the header win; footer labels are kept as keywords so both match.
 */
export function buildNavCommands(
  navItems: NavItem[],
  footerGroups: { title: string; items: NavItem[] }[]
): PaletteCommand[] {
  const byHref = new Map<string, PaletteCommand>();
  const add = (item: NavItem) => {
    const existing = byHref.get(item.href);
    if (existing) {
      if (
        item.label !== existing.label &&
        !existing.keywords?.includes(item.label)
      ) {
        existing.keywords = [...(existing.keywords ?? []), item.label];
      }
      return;
    }
    const goKey = Object.entries(GO_SHORTCUTS).find(
      ([, v]) => v.href === item.href
    )?.[0];
    const external = isExternalHref(item.href);
    byHref.set(item.href, {
      id: `nav:${item.href}`,
      label: item.label,
      group: "Navigate",
      href: item.href,
      external,
      keywords: [...(EXTRA_KEYWORDS[item.href] ?? [])],
      hint: goKey ? `G ${goKey.toUpperCase()}` : external ? "↗" : item.href,
    });
  };

  add({ href: "/", label: "Home" });
  navItems.forEach(add);
  footerGroups.forEach((g) => g.items.forEach(add));
  return [...byHref.values()];
}

/**
 * Fuzzy score of `query` against `text`: higher is better, `null` if the
 * query characters do not appear in order. Prefix and word-start matches
 * beat scattered subsequences.
 */
export function fuzzyScore(query: string, text: string): number | null {
  const q = query.trim().toLowerCase();
  const t = text.toLowerCase();
  if (!q) return 0;

  if (t === q) return 1000;
  if (t.startsWith(q)) return 800 - t.length;
  const wordIdx = t.search(
    new RegExp(`\\b${q.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}`)
  );
  if (wordIdx >= 0) return 600 - wordIdx;
  const idx = t.indexOf(q);
  if (idx >= 0) return 400 - idx;

  // Ordered subsequence: reward consecutive runs, penalise gaps.
  let score = 200;
  let ti = 0;
  let prev = -2;
  for (const ch of q) {
    if (ch === " ") continue;
    const found = t.indexOf(ch, ti);
    if (found < 0) return null;
    score += found === prev + 1 ? 5 : -(found - ti);
    prev = found;
    ti = found + 1;
  }
  return Math.max(score, 1);
}

/** Best score across a command's label and keywords. */
export function scoreCommand(
  query: string,
  command: PaletteCommand
): number | null {
  const candidates = [command.label, ...(command.keywords ?? [])];
  let best: number | null = null;
  candidates.forEach((text, i) => {
    const s = fuzzyScore(query, text);
    // Keyword hits rank slightly below equal label hits.
    const adjusted = s == null ? null : i === 0 ? s : s - 50;
    if (adjusted != null && (best == null || adjusted > best)) best = adjusted;
  });
  return best;
}

/** Filter and rank commands; an empty query keeps the original order. */
export function filterCommands(
  commands: PaletteCommand[],
  query: string
): PaletteCommand[] {
  if (!query.trim()) return commands;
  return commands
    .map((c, i) => ({ c, i, s: scoreCommand(query, c) }))
    .filter(
      (x): x is { c: PaletteCommand; i: number; s: number } => x.s != null
    )
    .sort((a, b) => b.s - a.s || a.i - b.i)
    .map((x) => x.c);
}

/** `/search` URL for a position, using the same params as the landing page. */
export function buildSearchHref(
  ra: number,
  dec: number,
  configs: CatalogRadiusConfig[] = buildDefaultCatalogConfigs()
): string {
  const params = new URLSearchParams({
    ra: ra.toString(),
    dec: dec.toString(),
    catalogRadii: encodeCatalogRadii(configs),
  });
  return `/search?${params.toString()}`;
}

/** Wrap the list index in either direction. */
export function cycleIndex(current: number, delta: number, length: number) {
  if (length <= 0) return 0;
  return (((current + delta) % length) + length) % length;
}

interface TargetLike {
  tagName?: string;
  isContentEditable?: boolean;
  getAttribute?: (name: string) => string | null;
}

/** True when keystrokes belong to a text field, not to global shortcuts. */
export function isTypingTarget(target: TargetLike | null | undefined): boolean {
  if (!target) return false;
  if (target.isContentEditable) return true;
  const tag = target.tagName?.toUpperCase();
  if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return true;
  const role = target.getAttribute?.("role");
  return role === "textbox" || role === "combobox" || role === "searchbox";
}

export function isMacPlatform(platform: string): boolean {
  return /mac|iphone|ipad|ipod/i.test(platform);
}

/** "⌘K" on Apple platforms, "Ctrl K" elsewhere. */
export function paletteShortcutLabel(isMac: boolean): string {
  return isMac ? "⌘K" : "Ctrl K";
}
