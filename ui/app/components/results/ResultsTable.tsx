"use client";

import { DownloadOutlined } from "@ant-design/icons";
import type { TableProps } from "antd";
import { Button, Flex, Skeleton, Table, Tag, Tooltip, Typography } from "antd";
import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";

import {
  CATALOG_LABELS,
  CATALOG_OPTIONS,
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import type { MagValue } from "@/app/lib/constants/photometry";
import { downloadCsv, toCsv } from "@/app/lib/utils/csv";

const { Text } = Typography;

export interface CrossmatchResult {
  key: string;
  objectId: string;
  ra: number;
  dec: number;
  angularDistance: number;
  catalog: string;
  ipix?: number;
  /**
   * Canonical request slug, used only as a join key for photometry. Optional
   * because the object detail page builds a CrossmatchResult by hand.
   */
  catalogSlug?: string;
}

/** A search result once photometry has been merged in. */
export interface EnrichedResult extends CrossmatchResult {
  photometry?: MagValue;
}

/** Lifecycle of the photometry enrichment, kept separate from the table's own. */
export type PhotometryStatus = "idle" | "pending" | "ready" | "error";

export interface ResultsTableProps {
  data?: EnrichedResult[];
  loading?: boolean;
  photometryStatus?: PhotometryStatus;
}

const DASH = <span className="font-mono text-xs text-border">—</span>;

function renderMagCell(
  photometry: MagValue | undefined,
  status: PhotometryStatus
) {
  if (photometry) {
    const tooltip = photometry.magErr
      ? `${photometry.description} (±${photometry.magErr.toFixed(3)})`
      : photometry.description;
    return (
      <Tooltip title={tooltip}>
        <span className="font-mono">{photometry.mag.toFixed(3)}</span>
        <span className="font-mono text-xs text-border ml-1.5">
          {photometry.band}
        </span>
      </Tooltip>
    );
  }

  if (status === "pending") {
    return <Skeleton.Input active size="small" className="!w-14 !min-w-0" />;
  }

  if (status === "error") {
    return (
      <Tooltip title="Photometry unavailable — the metadata request failed">
        {DASH}
      </Tooltip>
    );
  }

  return DASH;
}

/**
 * Unknown magnitudes sink to the bottom in BOTH directions. antd negates the
 * comparator's result for descending, so the sink direction is pre-negated.
 */
function magSorter(
  a: EnrichedResult,
  b: EnrichedResult,
  sortOrder?: "ascend" | "descend" | null
): number {
  const av = a.photometry?.mag;
  const bv = b.photometry?.mag;
  if (av == null && bv == null) return 0;
  const sink = sortOrder === "descend" ? -1 : 1;
  if (av == null) return sink;
  if (bv == null) return -sink;
  return av - bv;
}

/** Slug a row is filtered by — the request slug, falling back to what the API returned. */
function rowSlug(row: CrossmatchResult): string {
  return (row.catalogSlug ?? row.catalog).toLowerCase();
}

function buildColumns(
  photometryStatus: PhotometryStatus,
  catalogFilter: string[] | null
): TableProps<EnrichedResult>["columns"] {
  return [
    {
      title: "Object ID",
      dataIndex: "objectId",
      key: "objectId",
      ellipsis: true,
      width: 200,
    },
    {
      title: "Catalog",
      dataIndex: "catalog",
      key: "catalog",
      width: 110,
      // Derived from the catalogs we actually query. These used to be hardcoded
      // to "GAIA DR3" / "SIMBAD" / "2MASS" / "WISE" while the rows hold slugs,
      // so the filter matched nothing at all.
      filters: CATALOG_OPTIONS.map((c) => ({
        text: CATALOG_LABELS[c],
        value: c,
      })),
      filteredValue: catalogFilter,
      onFilter: (value, record) => rowSlug(record) === value,
      render: (value: string) => (
        <Tag color={getSearchCatalogColor(value)} className="font-medium">
          {getSearchCatalogLabel(value)}
        </Tag>
      ),
    },
    {
      title: 'Ang. Dist (")',
      dataIndex: "angularDistance",
      key: "angularDistance",
      width: 110,
      align: "right" as const,
      sorter: (a, b) => a.angularDistance - b.angularDistance,
      render: (value: number) => (
        <span className="font-mono">{value.toFixed(3)}</span>
      ),
    },
    {
      title: "Mag",
      dataIndex: "photometry",
      key: "photometry",
      width: 130,
      align: "right" as const,
      sorter: magSorter,
      showSorterTooltip: {
        title:
          "Sorts by magnitude. Bands differ between catalogs — filter to one catalog first to compare like with like.",
      },
      render: (photometry: MagValue | undefined) =>
        renderMagCell(photometry, photometryStatus),
    },
    {
      title: "RA (°)",
      dataIndex: "ra",
      key: "ra",
      width: 120,
      align: "right" as const,
      render: (value: number) => (
        <span className="font-mono">{value.toFixed(6)}</span>
      ),
    },
    {
      title: "Dec (°)",
      dataIndex: "dec",
      key: "dec",
      width: 120,
      align: "right" as const,
      render: (value: number) => (
        <span className="font-mono">{value.toFixed(6)}</span>
      ),
    },
    {
      title: "IPix",
      dataIndex: "ipix",
      key: "ipix",
      width: 130,
      align: "right" as const,
      render: (value?: number) => (
        <span
          className={`font-mono text-xs ${value != null ? "" : "text-border"}`}
        >
          {value != null ? value : "—"}
        </span>
      ),
    },
  ];
}

function resultsToCsv(data: EnrichedResult[]): string {
  return toCsv(
    [
      "object_id",
      "catalog",
      "ra_deg",
      "dec_deg",
      "angular_distance_arcsec",
      "ipix",
      "mag",
      "band",
    ],
    data.map((r) => [
      r.objectId,
      r.catalogSlug ?? r.catalog,
      r.ra,
      r.dec,
      r.angularDistance,
      r.ipix,
      r.photometry?.mag,
      r.photometry?.band,
    ])
  );
}

export function ResultsTable({
  data = [],
  loading = false,
  photometryStatus = "idle",
}: ResultsTableProps) {
  const router = useRouter();

  // The catalog filter is controlled so the header count and the CSV export can
  // agree with what is on screen — exporting every row while the table shows a
  // filtered subset would be a lie. Derived, not synced, so it stays correct
  // when `data` changes underneath it.
  const [catalogFilter, setCatalogFilter] = useState<string[] | null>(null);

  const columns = useMemo(
    () => buildColumns(photometryStatus, catalogFilter),
    [photometryStatus, catalogFilter]
  );

  const visibleRows = useMemo(
    () =>
      catalogFilter?.length
        ? data.filter((r) => catalogFilter.includes(rowSlug(r)))
        : data,
    [data, catalogFilter]
  );

  return (
    <Flex vertical gap={12}>
      <Flex justify="space-between" align="center">
        <Text type="secondary">
          {visibleRows.length} {visibleRows.length === 1 ? "match" : "matches"}
          {visibleRows.length !== data.length && ` (of ${data.length})`}
        </Text>
        <Button
          icon={<DownloadOutlined />}
          size="small"
          disabled={visibleRows.length === 0}
          // The whole filtered set, not just the page on screen.
          onClick={() =>
            downloadCsv("xwave_crossmatch.csv", resultsToCsv(visibleRows))
          }
        >
          Export CSV
        </Button>
      </Flex>
      <Table
        columns={columns}
        dataSource={data}
        loading={loading}
        onChange={(_pagination, filters) =>
          setCatalogFilter((filters.catalog as string[] | null) ?? null)
        }
        pagination={{
          pageSize: 10,
          showSizeChanger: true,
          pageSizeOptions: ["10", "20", "50"],
          showTotal: (total, range) =>
            `${range[0]}-${range[1]} of ${total} results`,
        }}
        // Fixed widths now sum to 920; leaving this at 800 would make antd
        // compress every column instead of scrolling.
        scroll={{ x: 930 }}
        size="small"
        onRow={(record) => ({
          onClick: () =>
            router.push(
              `/object/${encodeURIComponent(record.objectId)}?catalog=${encodeURIComponent(record.catalog)}`
            ),
          className: "cursor-pointer",
        })}
        rowClassName="hover:bg-surface-elevated transition-colors"
      />
    </Flex>
  );
}
