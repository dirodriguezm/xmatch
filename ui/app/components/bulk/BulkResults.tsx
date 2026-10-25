"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button, Segmented, Table, Tag, Typography } from "antd";
import type { ColumnsType, TableProps } from "antd/es/table";
import Link from "next/link";
import { useMemo, useState } from "react";

import { CodeSnippetButton } from "@/app/components/code/CodeSnippetButton";
import {
  CATALOG_COLOR_CLASSES,
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import {
  type BulkResultRow,
  bulkResultsToCsv,
  summarizeBulkResults,
} from "@/app/lib/utils/bulkInput";
import {
  type BulkFilter,
  type BulkSortKey,
  type BulkView,
  DEFAULT_BULK_VIEW,
  flattenBulkMatches,
  orderBulkRows,
  saveBulkNav,
  saveBulkView,
} from "@/app/lib/utils/bulkSession";
import { downloadCsv } from "@/app/lib/utils/csv";
import type { ApiRequest } from "@/app/lib/utils/snippets";
import { buildObjectUrl } from "@/app/lib/utils/urls";

const { Text } = Typography;

interface BulkResultsProps {
  rows: BulkResultRow[];
  /** Arcsec, as sent. */
  radius: number;
  request: ApiRequest;
  /** Filter, sort and page to start from (restored when coming back). */
  initialView?: BulkView;
}

function Stat({
  label,
  value,
}: {
  label: React.ReactNode;
  value: React.ReactNode;
}) {
  return (
    <div className="px-4 py-2 rounded border border-border bg-surface min-w-24">
      <div className="text-xs text-neutral-400">{label}</div>
      <div className="text-lg font-semibold text-foreground tabular-nums">
        {value}
      </div>
    </div>
  );
}

export function BulkResults({
  rows,
  radius,
  request,
  initialView = DEFAULT_BULK_VIEW,
}: BulkResultsProps) {
  const [view, setView] = useState<BulkView>(initialView);
  const summary = useMemo(() => summarizeBulkResults(rows), [rows]);

  // Sorted here rather than by antd so the object page can step through the
  // matches in exactly the order the table shows.
  const visible = useMemo(() => orderBulkRows(rows, view), [rows, view]);

  // Leaving for an object page: remember where the table was and which
  // matches ‹ › should walk through.
  const rememberView = () => {
    saveBulkView(view);
    saveBulkNav(flattenBulkMatches(visible));
  };

  const sortOrderFor = (key: BulkSortKey) =>
    view.sortKey === key ? view.sortOrder : null;

  const onTableChange: TableProps<BulkResultRow>["onChange"] = (
    pagination,
    _filters,
    sorter
  ) => {
    const s = Array.isArray(sorter) ? sorter[0] : sorter;
    const sortKey = s?.order ? (s.columnKey as BulkSortKey) : null;
    const sortOrder = s?.order ?? null;
    setView((v) => {
      const sortChanged = v.sortKey !== sortKey || v.sortOrder !== sortOrder;
      return {
        ...v,
        sortKey,
        sortOrder,
        page: sortChanged ? 1 : (pagination.current ?? 1),
        pageSize: pagination.pageSize ?? v.pageSize,
      };
    });
  };

  const columns: ColumnsType<BulkResultRow> = [
    {
      title: "#",
      key: "row",
      width: 56,
      render: (_, r) => (
        <span className="text-neutral-400 tabular-nums">
          {r.target.index + 1}
        </span>
      ),
      sorter: true,
      sortOrder: sortOrderFor("row"),
    },
    {
      title: "Input",
      key: "input",
      width: 260,
      render: (_, { target }) => (
        <div className="flex flex-col">
          <span className="text-foreground">
            {target.name || (
              <span className="text-neutral-400">line {target.line}</span>
            )}
          </span>
          <span className="font-mono text-xs text-neutral-400">
            {target.ra.toFixed(6)}, {target.dec >= 0 ? "+" : ""}
            {target.dec.toFixed(6)}
          </span>
        </div>
      ),
      sorter: true,
      sortOrder: sortOrderFor("input"),
    },
    {
      title: "Matches",
      key: "matches",
      render: (_, { matches }) =>
        matches.length === 0 ? (
          <Text className="text-neutral-400 text-xs">
            No match within {radius}″
          </Text>
        ) : (
          <ul className="m-0 p-0 list-none flex flex-col gap-1">
            {matches.map((m) => (
              <li
                key={`${m.catalog}:${m.id}`}
                className="flex items-center gap-2 min-w-0"
              >
                <Tag
                  color={getSearchCatalogColor(m.catalog)}
                  className="!m-0 font-medium shrink-0"
                >
                  {getSearchCatalogLabel(m.catalog)}
                </Tag>
                <Link
                  href={buildObjectUrl(m.id, m.catalog, "bulk")}
                  onClick={rememberView}
                  className="font-mono text-xs truncate"
                >
                  {m.id}
                </Link>
                <span className="ml-auto pl-2 font-mono text-xs text-neutral-400 tabular-nums shrink-0">
                  {m.separation === null ? "—" : `${m.separation.toFixed(2)}″`}
                </span>
              </li>
            ))}
          </ul>
        ),
      sorter: true,
      sortOrder: sortOrderFor("matches"),
    },
    {
      title: "Nearest",
      key: "nearest",
      width: 100,
      align: "right",
      render: (_, { matches }) =>
        matches[0]?.separation != null ? (
          <span className="font-mono text-xs tabular-nums">
            {matches[0].separation.toFixed(2)}″
          </span>
        ) : (
          <span className="text-neutral-400">—</span>
        ),
      sorter: true,
      sortOrder: sortOrderFor("nearest"),
    },
  ];

  const catalogs = Object.entries(summary.perCatalog).sort(
    (a, b) => b[1] - a[1]
  );

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-2">
        <Stat label="Inputs" value={summary.inputs} />
        <Stat label="Matched" value={summary.matched} />
        <Stat label="Unmatched" value={summary.unmatched} />
        <Stat label="Total matches" value={summary.totalMatches} />
        {catalogs.map(([cat, count]) => (
          <Stat
            key={cat}
            label={
              <span className="flex items-center gap-1.5">
                <span
                  className={`inline-block w-2 h-2 rounded-full ${
                    CATALOG_COLOR_CLASSES[cat] ?? "bg-neutral-500"
                  }`}
                />
                {getSearchCatalogLabel(cat)}
              </span>
            }
            value={count}
          />
        ))}
      </div>

      <div className="flex flex-wrap items-center justify-between gap-2">
        <Segmented<BulkFilter>
          value={view.filter}
          onChange={(filter) => setView((v) => ({ ...v, filter, page: 1 }))}
          options={[
            { value: "all", label: `All (${summary.inputs})` },
            { value: "matched", label: `Matched (${summary.matched})` },
            { value: "unmatched", label: `Unmatched (${summary.unmatched})` },
          ]}
        />
        <div className="flex flex-wrap gap-2">
          <CodeSnippetButton
            title="Bulk cross-match request"
            requests={[{ label: "Bulk cone search", request }]}
          />
          <Button
            icon={<DownloadOutlined />}
            onClick={() =>
              downloadCsv("xwave-bulk-crossmatch.csv", bulkResultsToCsv(rows))
            }
          >
            Download CSV
          </Button>
        </div>
      </div>

      <Table<BulkResultRow>
        size="small"
        rowKey={(r) => r.target.index}
        columns={columns}
        dataSource={visible}
        onChange={onTableChange}
        pagination={{
          current: view.page,
          pageSize: view.pageSize,
          showSizeChanger: true,
          pageSizeOptions: [25, 50, 100, 250],
          hideOnSinglePage: true,
        }}
        rowClassName={(r) => (r.matches.length === 0 ? "opacity-70" : "")}
        scroll={{ x: 640 }}
        locale={{ emptyText: "No rows for this filter" }}
      />
    </div>
  );
}
