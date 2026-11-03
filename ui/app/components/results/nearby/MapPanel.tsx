"use client";

import { Flex, Segmented, Typography } from "antd";
import { useMemo, useState } from "react";

import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";
import { CONE_SEARCH_MAX_NEIGHBORS } from "@/app/lib/constants/search";

import { CatalogMark } from "./CatalogMark";
import { OffsetMap } from "./OffsetMap";
import { formatArcsec, maxRadius, type NearbySource } from "./shared";
import type { LinkedHighlight } from "./useLinkedHighlight";

const { Text } = Typography;

interface MapPanelProps {
  sources: NearbySource[];
  radii: Record<string, number>;
  highlight: LinkedHighlight;
  size?: number;
  className?: string;
  /** Sizing of the drawing itself; it scales down to fit narrow screens. */
  mapClassName?: string;
}

/**
 * Offset map with a legend and a zoom to each catalog's radius — a 3″ Gaia
 * match is a speck on a 20″ eROSITA map otherwise.
 */
export function MapPanel({
  sources,
  radii,
  highlight,
  size = 320,
  className = "",
  mapClassName = "h-auto w-full max-w-[320px]",
}: MapPanelProps) {
  const outer = maxRadius(radii, sources);
  // Each catalog's radius, plus a tight "fit" level when every match sits
  // well inside the smallest one (e.g. 100 Gaia sources within 13″ of a 30″
  // search) — otherwise they'd bunch up in the middle.
  const farthest = Math.max(0, ...sources.map((s) => s.sepArcsec));
  const zoomLevels = useMemo(() => {
    const levels = [...Object.values(radii), outer].filter((r) => r > 0);
    const smallest = Math.min(...levels);
    const fit = Math.ceil(farthest * 1.1);
    if (farthest > 0 && fit < smallest * 0.7) levels.push(fit);
    return [...new Set(levels)].sort((a, b) => a - b);
  }, [radii, outer, farthest]);
  // Open on the tightest zoom that still shows every match: with Gaia at 3″
  // and eROSITA at 20″ but only Gaia matches, 20″ would bunch them up.
  const fitting = useMemo(
    () => zoomLevels.find((r) => r >= farthest) ?? outer,
    [zoomLevels, farthest, outer]
  );
  const [extent, setExtent] = useState<number | null>(null);
  const current = extent && zoomLevels.includes(extent) ? extent : fitting;
  const catalogs = [...new Set(sources.map((s) => s.slug))];

  return (
    <Flex vertical gap={8} align="center" className={className}>
      {zoomLevels.length > 1 && (
        <Flex align="center" gap={8}>
          <Text type="secondary" className="text-xs">
            Zoom to
          </Text>
          <Segmented
            size="small"
            value={current}
            onChange={(v) => setExtent(v as number)}
            options={zoomLevels.map((r) => ({
              value: r,
              label: formatArcsec(r),
            }))}
            aria-label="Map radius"
          />
        </Flex>
      )}
      <OffsetMap
        sources={sources}
        radii={radii}
        extentArcsec={current}
        hoveredKey={highlight.hoveredKey}
        selectedKey={highlight.selectedKey}
        onHover={highlight.setHoveredKey}
        onSelect={highlight.toggleSelected}
        size={size}
        className={mapClassName}
      />
      <Flex wrap gap={12} justify="center">
        {catalogs.map((c) => (
          <span key={c} className="inline-flex items-center gap-1.5 text-xs">
            <CatalogMark slug={c} />
            {getSearchCatalogLabel(c)}
            {sources.filter((s) => s.slug === c).length >=
              CONE_SEARCH_MAX_NEIGHBORS && (
              <span className="text-amber-400">
                nearest {CONE_SEARCH_MAX_NEIGHBORS} only
              </span>
            )}
            {radii[c] !== undefined && (
              <span className="text-neutral-400">
                ≤ {formatArcsec(radii[c])}
              </span>
            )}
          </span>
        ))}
        <Text type="secondary" className="text-xs">
          dashed = search radius · shaded = typical 1σ · N up, E left
        </Text>
      </Flex>
    </Flex>
  );
}
