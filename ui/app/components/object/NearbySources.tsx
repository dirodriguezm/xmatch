"use client";

import { Empty, Table, Typography } from "antd";
import Link from "next/link";

import { type Neighbor } from "@/app/hooks/queries";
import {
  CATALOG_COLOR_CLASSES,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import { buildObjectUrl } from "@/app/lib/utils/urls";

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
  return (
    <>
      <Table<Neighbor>
        size="small"
        loading={loading && neighbors.length === 0}
        dataSource={neighbors}
        rowKey={(n) => `${n.catalog}:${n.id}`}
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
            filters: [...new Set(neighbors.map((n) => n.catalog))].map((c) => ({
              text: getSearchCatalogLabel(c),
              value: c,
            })),
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
            title: "Separation",
            dataIndex: "separationArcsec",
            align: "right",
            sorter: (a, b) => a.separationArcsec - b.separationArcsec,
            render: (sep: number) => (
              <span className="font-mono">{sep.toFixed(2)}″</span>
            ),
          },
          {
            title: "PA",
            dataIndex: "positionAngle",
            align: "right",
            render: (pa: number) => (
              <span className="font-mono">
                {pa.toFixed(0)}°{" "}
                <Text type="secondary" className="text-xs">
                  {compassPoint(pa)}
                </Text>
              </span>
            ),
          },
        ]}
      />
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
