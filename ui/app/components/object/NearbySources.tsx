"use client";

import { Empty, Flex, Table, Tooltip, Typography } from "antd";
import Link from "next/link";
import { useState } from "react";

import { type Neighbor } from "@/app/hooks/queries";
import {
  CATALOG_COLOR_CLASSES,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import { buildObjectUrl } from "@/app/lib/utils/urls";

import { neighborKey, NeighborMap } from "./NeighborMap";

const { Text } = Typography;

interface NearbySourcesProps {
  neighbors: Neighbor[];
  radiusArcsec: number;
  loading: boolean;
  failedCatalogs: string[];
  truncatedCatalogs: string[];
}

const COMPASS = ["N", "NE", "E", "SE", "S", "SW", "W", "NW"];

function compassPoint(pa: number): string {
  return COMPASS[Math.round(pa / 45) % 8];
}

const catalogList = (cats: string[]) =>
  cats.map(getSearchCatalogLabel).join(", ");

/** Every indexed source around the object, nearest first. */
export function NearbySources({
  neighbors,
  radiusArcsec,
  loading,
  failedCatalogs,
  truncatedCatalogs,
}: NearbySourcesProps) {
  const [hoveredKey, setHoveredKey] = useState<string | null>(null);

  return (
    <>
      <Flex gap={32} wrap align="center" justify="center">
        {neighbors.length > 0 && (
          <NeighborMap
            neighbors={neighbors}
            radiusArcsec={radiusArcsec}
            hoveredKey={hoveredKey}
            onHover={setHoveredKey}
          />
        )}
        <Table<Neighbor>
          className="flex-1 min-w-[320px] max-w-[720px]"
          size="small"
          loading={loading && neighbors.length === 0}
          dataSource={neighbors}
          rowKey={neighborKey}
          onRow={(n) => ({
            onMouseEnter: () => setHoveredKey(neighborKey(n)),
            onMouseLeave: () => setHoveredKey(null),
          })}
          rowClassName={(n) =>
            neighborKey(n) === hoveredKey ? "bg-foreground/5" : ""
          }
          pagination={neighbors.length > 10 ? { pageSize: 10 } : false}
          locale={{
            emptyText: (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description={`No other sources within ${radiusArcsec}″`}
              />
            ),
          }}
          columns={[
            {
              title: "Catalog",
              dataIndex: "catalog",
              filters: [...new Set(neighbors.map((n) => n.catalog))].map(
                (c) => ({
                  text: getSearchCatalogLabel(c),
                  value: c,
                })
              ),
              onFilter: (value, n) => n.catalog === value,
              render: (catalog: string) => (
                <span className="inline-flex items-center gap-2">
                  <span
                    className={`inline-block w-2 h-2 rounded-full ${CATALOG_COLOR_CLASSES[catalog]}`}
                  />
                  {getSearchCatalogLabel(catalog)}
                </span>
              ),
            },
            {
              title: "ID",
              dataIndex: "id",
              render: (id: string, n) => (
                <Link
                  href={buildObjectUrl(id, n.catalog)}
                  className="font-mono text-xs"
                >
                  {id}
                </Link>
              ),
            },
            {
              title: (
                <Tooltip title="Angular distance from this object, in arcseconds">
                  <span className="cursor-help underline decoration-dotted underline-offset-2">
                    Separation
                  </span>
                </Tooltip>
              ),
              dataIndex: "separationArcsec",
              align: "right",
              sorter: (a, b) => a.separationArcsec - b.separationArcsec,
              render: (sep: number) => (
                <span className="font-mono">{sep.toFixed(2)}″</span>
              ),
            },
            {
              title: (
                <Tooltip title="Position angle: direction from this object, measured from north (0°) through east (90°), south (180°) and west (270°)">
                  <span className="cursor-help underline decoration-dotted underline-offset-2">
                    PA
                  </span>
                </Tooltip>
              ),
              dataIndex: "positionAngle",
              align: "right",
              render: (pa: number) => (
                <span className="font-mono whitespace-nowrap">
                  {pa.toFixed(0)}°{" "}
                  {/* Fixed width so the degrees line up whether the
                      compass point has one letter or two. */}
                  <Text
                    type="secondary"
                    className="inline-block min-w-[1.75rem] text-left font-mono text-xs"
                  >
                    {compassPoint(pa)}
                  </Text>
                </span>
              ),
            },
          ]}
        />
      </Flex>
      {(failedCatalogs.length > 0 || truncatedCatalogs.length > 0) && (
        <Text type="secondary" className="text-xs block mt-2">
          {failedCatalogs.length > 0 &&
            `Search failed for ${catalogList(failedCatalogs)}. `}
          {truncatedCatalogs.length > 0 &&
            `Only the nearest sources are shown for ${catalogList(truncatedCatalogs)}.`}
        </Text>
      )}
    </>
  );
}
