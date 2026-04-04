"use client";

import type { TableProps } from "antd";
import { Table, Tag } from "antd";
import { useRouter } from "next/navigation";

import { getCatalogColor } from "@/app/lib/constants/catalogs";

// Example data type - adjust based on actual API response
export interface CrossmatchResult {
  key: string;
  objectId: string;
  ra: number;
  dec: number;
  angularDistance: number;
  catalog: string;
  ipix?: number;
}

export interface ResultsTableProps {
  data?: CrossmatchResult[];
  loading?: boolean;
}

const columns: TableProps<CrossmatchResult>["columns"] = [
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
    filters: [
      { text: "GAIA DR3", value: "GAIA DR3" },
      { text: "SIMBAD", value: "SIMBAD" },
      { text: "2MASS", value: "2MASS" },
      { text: "WISE", value: "WISE" },
    ],
    onFilter: (value, record) => record.catalog === value,
    render: (value: string) => (
      <Tag color={getCatalogColor(value)} className="font-medium">
        {value}
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

// Sample data for demonstration
export const sampleData: CrossmatchResult[] = [
  {
    key: "1",
    objectId: "Gaia DR3 4657433693234",
    ra: 10.684583,
    dec: 41.269167,
    angularDistance: 0.012,
    catalog: "GAIA DR3",
    ipix: 405766747735,
  },
  {
    key: "2",
    objectId: "2MASS J00424433+4116069",
    ra: 10.684722,
    dec: 41.268583,
    angularDistance: 0.089,
    catalog: "2MASS",
    ipix: 405766747736,
  },
];

export function ResultsTable({
  data = [],
  loading = false,
}: ResultsTableProps) {
  const router = useRouter();

  return (
    <Table
      columns={columns}
      dataSource={data}
      loading={loading}
      pagination={{
        pageSize: 10,
        showSizeChanger: true,
        pageSizeOptions: ["10", "20", "50"],
        showTotal: (total, range) =>
          `${range[0]}-${range[1]} of ${total} results`,
      }}
      scroll={{ x: 800 }}
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
  );
}
