"use client";

import { Grid } from "antd";
import Link from "next/link";
import { Fragment, type KeyboardEvent, useEffect, useRef } from "react";

import { BasketButton } from "@/app/components/basket/BasketButton";
import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";

import { CatalogMark } from "./CatalogMark";
import {
  type CatalogSummary,
  VERDICT_CLASSES,
  VERDICT_TEXT,
} from "./catalogSummary";
import { SeparationBar } from "./SeparationBar";
import { compassPoint, formatArcsec, type NearbySource } from "./shared";
import type { LinkedHighlight } from "./useLinkedHighlight";

/** Rows shown per catalog before "Show all". */
export const TOP_N = 3;

interface GroupedMatchesProps {
  groups: CatalogSummary[];
  highlight: LinkedHighlight;
  /** Catalogs listed in full rather than their top few. */
  expanded: Set<string>;
  /** Omitted when a single catalog is shown, always in full. */
  onToggleExpanded?: (slug: string) => void;
  /**
   * Approximate chance of an unrelated source as close as each catalog's
   * nearest match, as a formatted string; shown in the group header.
   */
  chanceBySlug?: Record<string, string>;
  /** Classes for the scrolling container (defaults fit layout C). */
  className?: string;
}

/** Scroll a row into view inside the table's own scroller, never the page. */
function scrollIntoScroller(scroller: HTMLElement | null, key: string) {
  const row = scroller?.querySelector<HTMLElement>(
    `tr[data-key="${CSS.escape(key)}"]`
  );
  if (!scroller || !row || scroller.scrollHeight <= scroller.clientHeight)
    return;
  const head = scroller.querySelector("thead")?.offsetHeight ?? 0;
  const top = row.offsetTop - head;
  const bottom = row.offsetTop + row.offsetHeight;
  if (top < scroller.scrollTop) scroller.scrollTo({ top, behavior: "smooth" });
  else if (bottom > scroller.scrollTop + scroller.clientHeight)
    scroller.scrollTo({
      top: bottom - scroller.clientHeight,
      behavior: "smooth",
    });
}

/**
 * One table, one row group per catalog: each sorted by separation and cut to
 * its nearest few, with an explicit line when a catalog has nothing.
 */
export function GroupedMatches({
  groups,
  highlight,
  expanded,
  onToggleExpanded,
  chanceBySlug,
  className = "self-start lg:max-h-[calc(100vh-230px)] lg:overflow-y-auto",
}: GroupedMatchesProps) {
  const scroller = useRef<HTMLDivElement>(null);
  // Phones drop the PA and basket columns (both are on the object page).
  const compact = Grid.useBreakpoint().sm === false;
  const cols = compact ? 4 : 6;
  const { hoveredKey, selectedKey, setHoveredKey, toggleSelected } = highlight;

  useEffect(() => {
    if (selectedKey) scrollIntoScroller(scroller.current, selectedKey);
  }, [selectedKey, expanded]);
  useEffect(() => {
    if (hoveredKey) scrollIntoScroller(scroller.current, hoveredKey);
  }, [hoveredKey]);

  const onRowKey = (e: KeyboardEvent, key: string) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      toggleSelected(key);
    }
  };

  const row = (s: NearbySource, radius: number | undefined, rank: number) => (
    <tr
      key={s.key}
      data-key={s.key}
      tabIndex={0}
      aria-selected={s.key === selectedKey}
      onMouseEnter={() => setHoveredKey(s.key)}
      onMouseLeave={() => setHoveredKey(null)}
      onFocus={() => setHoveredKey(s.key)}
      onBlur={() => setHoveredKey(null)}
      onClick={() => toggleSelected(s.key)}
      onKeyDown={(e) => onRowKey(e, s.key)}
      className={[
        "cursor-pointer border-t border-border/60 outline-none transition-colors focus-visible:ring-1 focus-visible:ring-primary focus-visible:ring-inset",
        s.key === selectedKey
          ? "bg-primary/15"
          : s.key === hoveredKey
            ? "bg-foreground/5"
            : "",
      ].join(" ")}
    >
      <td className="w-8 py-1.5 pl-2 text-right font-mono text-xs text-neutral-500">
        {rank}
      </td>
      <td className="max-w-0 py-1.5 pr-3 pl-2">
        <Link
          href={s.href}
          className="block truncate font-mono text-xs"
          onClick={(e) => e.stopPropagation()}
          title={s.objectId}
        >
          {s.objectId}
        </Link>
      </td>
      {/* Phones keep the numbers and drop the bar. */}
      <td className="py-1.5 pr-2 sm:pr-3 [&_svg]:hidden sm:[&_svg]:block">
        <SeparationBar
          source={s}
          scaleArcsec={radius ?? s.sepArcsec * 1.2}
          width={56}
        />
      </td>
      {!compact && (
        <td className="py-1.5 pr-3 text-right font-mono text-xs whitespace-nowrap">
          {s.pa.toFixed(0)}°{" "}
          <span className="inline-block min-w-[1.5rem] text-left text-neutral-400">
            {compassPoint(s.pa)}
          </span>
        </td>
      )}
      <td className="py-1.5 pr-2 text-right font-mono text-xs whitespace-nowrap">
        {s.photometry ? (
          <>
            {s.photometry.mag.toFixed(2)}
            <span className="ml-1 text-neutral-400">{s.photometry.band}</span>
          </>
        ) : (
          <span className="text-neutral-600">—</span>
        )}
      </td>
      {!compact && (
        <td className="w-9 py-1 pr-1 text-center">
          <span
            onClick={(e) => e.stopPropagation()}
            onKeyDown={(e) => e.stopPropagation()}
            role="presentation"
          >
            <BasketButton
              size="small"
              item={{
                objectId: s.objectId,
                catalog: s.slug,
                ra: s.ra,
                dec: s.dec,
              }}
            />
          </span>
        </td>
      )}
    </tr>
  );

  return (
    <div
      ref={scroller}
      className={`min-w-0 overflow-x-auto rounded-lg border border-border bg-surface ${className}`}
    >
      <table className="w-full table-fixed sm:min-w-[460px] border-collapse text-sm">
        <caption className="sr-only">
          Matches grouped by catalog, nearest first
        </caption>
        <thead className="sticky top-0 z-10 bg-surface-elevated text-xs text-neutral-400">
          <tr>
            <th scope="col" className="w-8 py-2 pl-2 text-right font-normal">
              #
            </th>
            <th scope="col" className="py-2 pl-2 text-left font-normal">
              Object ID
            </th>
            <th
              scope="col"
              className="w-[92px] py-2 text-left font-normal sm:w-[150px]"
            >
              Separation · σ
            </th>
            {!compact && (
              <th
                scope="col"
                className="w-[76px] py-2 pr-3 text-right font-normal"
              >
                PA
              </th>
            )}
            <th
              scope="col"
              className="w-[72px] py-2 pr-2 text-right font-normal sm:w-[78px]"
            >
              Mag
            </th>
            {!compact && (
              <th scope="col" className="w-9 py-2">
                <span className="sr-only">Basket</span>
              </th>
            )}
          </tr>
        </thead>
        {groups.map((g) => {
          const label = getSearchCatalogLabel(g.slug);
          const isOpen = expanded.has(g.slug);
          const shown = isOpen ? g.sources : g.sources.slice(0, TOP_N);
          const hidden = g.sources.length - shown.length;
          const r = g.radius !== undefined ? formatArcsec(g.radius) : "";
          return (
            <tbody key={g.slug} className="border-t-2 border-border">
              <tr>
                <th
                  scope="rowgroup"
                  colSpan={cols}
                  className="bg-foreground/[0.03] px-2 py-2 text-left font-normal"
                >
                  <span className="flex flex-wrap items-center gap-x-3 gap-y-1">
                    <span className="inline-flex items-center gap-2 text-sm font-medium text-foreground">
                      <CatalogMark slug={g.slug} size={11} />
                      {label}
                    </span>
                    <span className="text-xs text-neutral-400">
                      {g.sources.length}
                      {g.capped ? "+" : ""} within {r}
                      {g.capped && (
                        <span className="ml-1 text-amber-400">
                          (nearest 100 only)
                        </span>
                      )}
                    </span>
                    {chanceBySlug?.[g.slug] && (
                      <span
                        className="text-xs text-neutral-400"
                        title="Approximate chance that an unrelated source lies this close, from the density of sources in the outer half of the cone"
                      >
                        chance ≈ {chanceBySlug[g.slug]}
                      </span>
                    )}
                    {g.verdict !== "none" && (
                      <span
                        className={`rounded border px-1.5 py-0.5 text-[11px] leading-none ${VERDICT_CLASSES[g.verdict]}`}
                      >
                        {VERDICT_TEXT[g.verdict]}
                        {g.verdict === "ambiguous" &&
                          ` · ${g.within2Sigma} within 2σ`}
                      </span>
                    )}
                  </span>
                </th>
              </tr>
              {g.sources.length === 0 ? (
                <tr className="border-t border-border/60">
                  <td />
                  <td
                    colSpan={cols - 1}
                    className="py-2 pl-2 text-xs text-neutral-400"
                  >
                    None within {r} of the searched position
                  </td>
                </tr>
              ) : (
                <Fragment>
                  {shown.map((s, i) => row(s, g.radius, i + 1))}
                  {g.sources.length > TOP_N && onToggleExpanded && (
                    <tr className="border-t border-border/60">
                      <td />
                      <td colSpan={cols - 1} className="py-1 pl-1">
                        <button
                          type="button"
                          aria-expanded={isOpen}
                          onClick={() => onToggleExpanded(g.slug)}
                          className="cursor-pointer rounded px-1 py-0.5 text-xs text-primary hover:underline focus-visible:ring-1 focus-visible:ring-primary focus-visible:outline-none"
                        >
                          {isOpen
                            ? `Show nearest ${TOP_N} only`
                            : `Show all ${g.sources.length} ${label} (+${hidden})`}
                        </button>
                      </td>
                    </tr>
                  )}
                </Fragment>
              )}
            </tbody>
          );
        })}
      </table>
    </div>
  );
}
