"use client";

import { InfoCircleOutlined } from "@ant-design/icons";
import type { TableProps } from "antd";
import { Flex, Table, Tooltip, Typography } from "antd";
import Link from "next/link";
import { useEffect, useMemo, useRef } from "react";

import { BasketButton } from "@/app/components/basket/BasketButton";
import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";
import { toCsv } from "@/app/lib/utils/csv";

import { CatalogMark } from "./CatalogMark";
import { SeparationBar } from "./SeparationBar";
import {
  compassPoint,
  type NearbySource,
  scrollRowIntoView,
  TARGET_SIGMA_ARCSEC,
} from "./shared";
import type { LinkedHighlight } from "./useLinkedHighlight";

const { Text } = Typography;

interface LinkedTableProps {
  sources: NearbySource[];
  radii: Record<string, number>;
  highlight: LinkedHighlight;
  /** Hide RA/Dec to fit beside a map. */
  dense?: boolean;
  /** Body height before the table scrolls (px or a CSS length). */
  scrollY?: number | string;
}

/** CSV of the matches, as offered by the results toolbar. */
export function nearbyCsv(sources: NearbySource[]) {
  return toCsv(
    [
      "object_id",
      "catalog",
      "ra_deg",
      "dec_deg",
      "angular_distance_arcsec",
      "position_angle_deg",
      "sigma_ratio",
      "mag",
      "band",
    ],
    sources.map((s) => [
      s.objectId,
      s.slug,
      s.ra,
      s.dec,
      s.sepArcsec,
      Number(s.pa.toFixed(1)),
      Number((s.sepArcsec / s.sigma).toFixed(2)),
      s.photometry?.mag,
      s.photometry?.band,
    ])
  );
}

/** Match list whose hover and selection are shared with a map. */
export function LinkedTable({
  sources,
  radii,
  highlight,
  dense = false,
  scrollY = 420,
}: LinkedTableProps) {
  const ref = useRef<HTMLDivElement>(null);
  const { hoveredKey, selectedKey, setHoveredKey, toggleSelected } = highlight;

  // Follow selections made on the map.
  useEffect(() => {
    if (selectedKey) scrollRowIntoView(ref.current, selectedKey);
  }, [selectedKey]);
  useEffect(() => {
    if (hoveredKey) scrollRowIntoView(ref.current, hoveredKey);
  }, [hoveredKey]);

  const catalogCounts = useMemo(() => {
    const counts = new Map<string, number>();
    for (const s of sources) counts.set(s.slug, (counts.get(s.slug) ?? 0) + 1);
    return [...counts.entries()];
  }, [sources]);

  const columns = useMemo<TableProps<NearbySource>["columns"]>(() => {
    const cols: TableProps<NearbySource>["columns"] = [
      {
        title: "",
        key: "basket",
        width: 40,
        align: "center",
        render: (_, s) => (
          <span onClick={(e) => e.stopPropagation()}>
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
        ),
      },
      {
        title: "Catalog",
        dataIndex: "slug",
        width: 100,
        filters: [...new Set(sources.map((s) => s.slug))].map((c) => ({
          text: getSearchCatalogLabel(c),
          value: c,
        })),
        onFilter: (value, s) => s.slug === value,
        render: (slug: string) => (
          <span className="inline-flex items-center gap-2">
            <CatalogMark slug={slug} />
            {getSearchCatalogLabel(slug)}
          </span>
        ),
      },
      {
        title: "Object ID",
        dataIndex: "objectId",
        ellipsis: true,
        width: 210,
        render: (id: string, s) => (
          <Link
            href={s.href}
            className="font-mono text-xs"
            onClick={(e) => e.stopPropagation()}
          >
            {id}
          </Link>
        ),
      },
      {
        title: (
          <Tooltip
            title={`Separation from the searched position. The bar spans the catalog's search radius; green is 2σ and amber 3σ of the combined error (typical catalog error plus ${TARGET_SIGMA_ARCSEC}″ for the searched position).`}
          >
            <span className="cursor-help underline decoration-dotted underline-offset-2">
              Separation
            </span>
          </Tooltip>
        ),
        key: "sep",
        width: dense ? 170 : 210,
        sorter: (a, b) => a.sepArcsec - b.sepArcsec,
        defaultSortOrder: "ascend",
        render: (_, s) => (
          <SeparationBar
            source={s}
            scaleArcsec={radii[s.slug] ?? s.sepArcsec * 1.2}
            width={dense ? 64 : 96}
          />
        ),
      },
    ];
    cols.push({
      title: "PA",
      dataIndex: "pa",
      width: 80,
      align: "right",
      render: (pa: number) => (
        <span className="font-mono text-xs whitespace-nowrap">
          {pa.toFixed(0)}°{" "}
          <span className="inline-block min-w-[1.5rem] text-left text-neutral-400">
            {compassPoint(pa)}
          </span>
        </span>
      ),
    });
    cols.push({
      title: "Mag",
      key: "mag",
      width: 100,
      align: "right",
      sorter: (a, b) =>
        (a.photometry?.mag ?? Infinity) - (b.photometry?.mag ?? Infinity),
      render: (_, s) =>
        s.photometry ? (
          <span className="font-mono text-xs whitespace-nowrap">
            {s.photometry.mag.toFixed(2)}
            <span className="ml-1 text-neutral-400">{s.photometry.band}</span>
          </span>
        ) : (
          <span className="text-border">—</span>
        ),
    });
    if (!dense) {
      cols.push(
        {
          title: "RA (°)",
          dataIndex: "ra",
          width: 110,
          align: "right",
          render: (v: number) => (
            <span className="font-mono text-xs">{v.toFixed(6)}</span>
          ),
        },
        {
          title: "Dec (°)",
          dataIndex: "dec",
          width: 110,
          align: "right",
          render: (v: number) => (
            <span className="font-mono text-xs">{v.toFixed(6)}</span>
          ),
        }
      );
    }
    return cols;
  }, [sources, radii, dense]);

  return (
    <Flex vertical gap={8} ref={ref} className="min-w-0">
      {/* The match count is the page's main fact, so it gets the weight;
          the per-catalog breakdown and the hint stay quiet. */}
      <Flex align="baseline" gap={12} wrap="wrap">
        <Text strong className="text-base">
          {sources.length} {sources.length === 1 ? "match" : "matches"}
        </Text>
        <Flex gap={12} wrap="wrap">
          {catalogCounts.map(([slug, n]) => (
            <span
              key={slug}
              className="inline-flex items-center gap-1.5 text-xs text-neutral-400"
            >
              <CatalogMark slug={slug} size={8} />
              {getSearchCatalogLabel(slug)} {n}
            </span>
          ))}
        </Flex>
        <Tooltip title="Hover or click a row to find it on the map, and vice versa">
          <InfoCircleOutlined className="text-neutral-500 text-xs" />
        </Tooltip>
      </Flex>
      <Table<NearbySource>
        size="small"
        rowKey="key"
        dataSource={sources}
        columns={columns}
        pagination={false}
        scroll={{ y: scrollY, x: dense ? 700 : 950 }}
        onRow={(s) => ({
          tabIndex: 0,
          "aria-selected": s.key === selectedKey,
          onMouseEnter: () => setHoveredKey(s.key),
          onMouseLeave: () => setHoveredKey(null),
          onFocus: () => setHoveredKey(s.key),
          onBlur: () => setHoveredKey(null),
          onClick: () => toggleSelected(s.key),
          onKeyDown: (e) => {
            if (e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              toggleSelected(s.key);
            }
          },
        })}
        rowClassName={(s) =>
          [
            "cursor-pointer transition-colors outline-none focus-visible:ring-1 focus-visible:ring-primary",
            s.key === selectedKey
              ? "bg-primary/15"
              : s.key === hoveredKey
                ? "bg-foreground/5"
                : "",
          ].join(" ")
        }
      />
    </Flex>
  );
}
