"use client";

import type { KeyboardEvent } from "react";

import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";

import { MarkPath } from "./CatalogMark";
import {
  agreement,
  AGREEMENT_LABEL,
  compassPoint,
  fillClass,
  formatArcsec,
  type NearbySource,
  shapeOf,
  sigmaRatio,
  strokeClass,
} from "./shared";

interface OffsetMapProps {
  sources: NearbySource[];
  /** Search radius per catalog slug, drawn as a dashed circle each. */
  radii: Record<string, number>;
  /** Radius of the outer edge of the map, in arcseconds. */
  extentArcsec: number;
  hoveredKey?: string | null;
  selectedKey?: string | null;
  onHover?: (key: string | null) => void;
  onSelect?: (key: string) => void;
  /** Size of the drawing, in SVG units (and default pixels). */
  size?: number;
  /** Shrink to a thumbnail: no labels, no interaction. */
  compact?: boolean;
  /** Shade each source's 1σ positional uncertainty. */
  showErrors?: boolean;
  className?: string;
}

const NICE_STEPS = [0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 60];

/** A ring spacing that gives two to five rings inside the extent. */
function ringStep(extent: number): number {
  return NICE_STEPS.find((s) => extent / s <= 5) ?? 60;
}

/**
 * Offsets of the matches from the searched position: north up, east left,
 * as on the object page's NeighborMap, plus each catalog's search radius and
 * the 1σ error of every source.
 */
export function OffsetMap({
  sources,
  radii,
  extentArcsec,
  hoveredKey = null,
  selectedKey = null,
  onHover,
  onSelect,
  size = 300,
  compact = false,
  showErrors = true,
  className = "",
}: OffsetMapProps) {
  const half = size / 2;
  const pad = compact ? 3 : 18;
  const scale = (half - pad) / extentArcsec;
  const step = ringStep(extentArcsec);
  const rings = Array.from(
    { length: Math.floor(extentArcsec / step + 1e-9) },
    (_, i) => (i + 1) * step
  );
  const activeKey = hoveredKey ?? selectedKey;
  const interactive = !compact && (onHover || onSelect);

  // DOM order stays fixed so keyboard focus is not lost when a mark becomes
  // active; the active mark is drawn a second time on top instead.
  const activeSource = compact
    ? undefined
    : sources.find((s) => s.key === activeKey);

  const place = (s: NearbySource) => {
    // Sources beyond the extent (zoomed in) sit on the edge, faded.
    const k = Math.min(1, extentArcsec / Math.max(s.sepArcsec, 1e-9));
    return {
      x: -s.east * k * scale,
      y: -s.north * k * scale,
      outside: k < 1,
    };
  };

  const onKey = (e: KeyboardEvent, key: string) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      onSelect?.(key);
    }
  };

  return (
    <svg
      viewBox={`${-half} ${-half} ${size} ${size}`}
      width={size}
      height={size}
      className={`shrink-0 text-foreground ${className}`}
      role={interactive ? "group" : "img"}
      aria-label={`Sky offsets of ${sources.length} matches from the searched position, north up, east left`}
    >
      <circle
        r={half - pad}
        className="fill-surface"
        stroke="currentColor"
        strokeOpacity={0.12}
      />
      {rings.map((r) => (
        <g key={r}>
          <circle
            r={r * scale}
            fill="none"
            stroke="currentColor"
            strokeOpacity={0.12}
            strokeDasharray="2 3"
          />
          {!compact && size >= 200 && (
            <text
              x={r * scale * Math.SQRT1_2 + 2}
              y={r * scale * Math.SQRT1_2 + 9}
              fontSize={9}
              fill="currentColor"
              fillOpacity={0.4}
            >
              {formatArcsec(r)}
            </text>
          )}
        </g>
      ))}

      {/* Search radius per catalog */}
      {/* Labelled in the legend: equal radii would print on top of each other. */}
      {Object.entries(radii).map(([slug, r]) =>
        r <= extentArcsec ? (
          <g key={slug}>
            <circle
              r={r * scale}
              fill="none"
              className={strokeClass(slug)}
              strokeOpacity={0.6}
              strokeWidth={compact ? 1 : 1.25}
              strokeDasharray="5 3"
            />
          </g>
        ) : null
      )}

      {!compact && (
        <>
          <text
            y={-half + 11}
            textAnchor="middle"
            fontSize={10}
            fill="currentColor"
            fillOpacity={0.6}
          >
            N
          </text>
          <text
            x={-half + 4}
            y={4}
            fontSize={10}
            fill="currentColor"
            fillOpacity={0.6}
          >
            E
          </text>
        </>
      )}

      {/* The searched position */}
      <g stroke="currentColor" strokeOpacity={0.8}>
        <line x1={-8} x2={-3} />
        <line x1={3} x2={8} />
        <line y1={-8} y2={-3} />
        <line y1={3} y2={8} />
      </g>

      {showErrors &&
        sources.map((s) => {
          const { x, y, outside } = place(s);
          const r = s.sigma * scale;
          if (outside || r < 2) return null;
          return (
            <circle
              key={`err-${s.key}`}
              cx={x}
              cy={y}
              r={r}
              className={`${fillClass(s.slug)} ${strokeClass(s.slug)}`}
              fillOpacity={s.key === activeKey ? 0.22 : 0.08}
              strokeOpacity={0.35}
              pointerEvents="none"
            />
          );
        })}

      {sources.map((s) => {
        const { x, y, outside } = place(s);
        const active = s.key === activeKey;
        const selected = s.key === selectedKey;
        const dim = activeKey !== null && !active;
        const r = compact ? 2.5 : active ? 6.5 : 5;
        const label = `${getSearchCatalogLabel(s.slug)} ${s.objectId}, ${formatArcsec(s.sepArcsec)} ${compassPoint(s.pa)}, ${sigmaRatio(s).toFixed(1)}σ (${AGREEMENT_LABEL[agreement(s)]})`;
        return (
          <g
            key={s.key}
            className={
              interactive
                ? "cursor-pointer outline-none focus-visible:opacity-100"
                : undefined
            }
            tabIndex={interactive ? 0 : undefined}
            role={interactive ? "button" : undefined}
            aria-label={interactive ? label : undefined}
            aria-pressed={interactive ? selected : undefined}
            onMouseEnter={interactive ? () => onHover?.(s.key) : undefined}
            onMouseLeave={interactive ? () => onHover?.(null) : undefined}
            onFocus={interactive ? () => onHover?.(s.key) : undefined}
            onBlur={interactive ? () => onHover?.(null) : undefined}
            onClick={interactive ? () => onSelect?.(s.key) : undefined}
            onKeyDown={interactive ? (e) => onKey(e, s.key) : undefined}
          >
            {/* Larger invisible hit area for small marks */}
            {interactive && <circle cx={x} cy={y} r={10} fill="transparent" />}
            {selected && !compact && (
              <circle
                cx={x}
                cy={y}
                r={11}
                fill="none"
                stroke="currentColor"
                strokeWidth={1.5}
              />
            )}
            <MarkPath
              shape={shapeOf(s.slug)}
              r={r}
              x={x}
              y={y}
              className={fillClass(s.slug)}
              fillOpacity={outside ? 0.4 : dim ? 0.35 : 0.95}
              stroke={active && !compact ? "white" : "black"}
              strokeWidth={active && !compact ? 1.5 : 0.75}
            />
            {!compact && <title>{label}</title>}
          </g>
        );
      })}
      {activeSource && (
        <MarkPath
          shape={shapeOf(activeSource.slug)}
          r={6.5}
          x={place(activeSource).x}
          y={place(activeSource).y}
          className={`${fillClass(activeSource.slug)} pointer-events-none`}
          stroke="white"
          strokeWidth={1.5}
        />
      )}
    </svg>
  );
}
