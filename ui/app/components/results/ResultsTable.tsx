"use client";

import type { TableProps } from "antd";
import { Table, Tag } from "antd";
import { useRouter } from "next/navigation";

// Example data type - adjust based on actual API response
export interface CrossmatchResult {
  key: string;
  objectId: string;
  ra: number;
  dec: number;
  angularDistance: number;
  catalog: string;
  gMag?: number;
}

export interface ResultsTableProps {
  data?: CrossmatchResult[];
  loading?: boolean;
}

// Catalog color mapping
const catalogColors: Record<string, string> = {
  "GAIA DR3": "#1890ff",
  SIMBAD: "#52c41a",
  "2MASS": "#fa8c16",
  WISE: "#722ed1",
  AllWISE: "#722ed1",
};

const getCatalogColor = (catalog: string): string => {
  return catalogColors[catalog] || "#8c8c8c";
};

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
    title: "G Mag",
    dataIndex: "gMag",
    key: "gMag",
    width: 90,
    align: "right" as const,
    sorter: (a, b) => (a.gMag ?? 0) - (b.gMag ?? 0),
    render: (value?: number) => (
      <span className={`font-mono ${value != null ? "" : "text-border"}`}>
        {value != null ? value.toFixed(2) : "—"}
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
    gMag: 3.44,
  },
  {
    key: "2",
    objectId: "M31",
    ra: 10.684708,
    dec: 41.269028,
    angularDistance: 0.045,
    catalog: "SIMBAD",
    gMag: 3.47,
  },
  {
    key: "3",
    objectId: "2MASS J00424433+4116069",
    ra: 10.684722,
    dec: 41.268583,
    angularDistance: 0.089,
    catalog: "2MASS",
    gMag: 4.12,
  },
  {
    key: "4",
    objectId: "WISEA J004244.35+411608.5",
    ra: 10.684792,
    dec: 41.269028,
    angularDistance: 0.102,
    catalog: "WISE",
    gMag: 4.89,
  },
  {
    key: "5",
    objectId: "Gaia DR3 4657433693235",
    ra: 10.685,
    dec: 41.2695,
    angularDistance: 0.156,
    catalog: "GAIA DR3",
    gMag: 5.21,
  },
  {
    key: "6",
    objectId: "NGC 224",
    ra: 10.684375,
    dec: 41.268889,
    angularDistance: 0.178,
    catalog: "SIMBAD",
  },
  {
    key: "7",
    objectId: "2MASS J00424429+4116054",
    ra: 10.684542,
    dec: 41.268167,
    angularDistance: 0.201,
    catalog: "2MASS",
    gMag: 6.34,
  },
  {
    key: "8",
    objectId: "WISEA J004244.28+411605.2",
    ra: 10.6845,
    dec: 41.268111,
    angularDistance: 0.234,
    catalog: "WISE",
    gMag: 6.78,
  },
  {
    key: "9",
    objectId: "Gaia DR3 4657433693236",
    ra: 10.685208,
    dec: 41.269861,
    angularDistance: 0.267,
    catalog: "GAIA DR3",
    gMag: 7.12,
  },
  {
    key: "10",
    objectId: "Andromeda Galaxy",
    ra: 10.684167,
    dec: 41.2685,
    angularDistance: 0.289,
    catalog: "SIMBAD",
    gMag: 3.44,
  },
  {
    key: "11",
    objectId: "2MASS J00424445+4116092",
    ra: 10.685208,
    dec: 41.269222,
    angularDistance: 0.312,
    catalog: "2MASS",
    gMag: 8.45,
  },
  {
    key: "12",
    objectId: "WISEA J004244.51+411609.3",
    ra: 10.685458,
    dec: 41.26925,
    angularDistance: 0.345,
    catalog: "WISE",
    gMag: 9.01,
  },
];

export function ResultsTable({
  data = sampleData,
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
          router.push(`/object/${encodeURIComponent(record.objectId)}`),
        className: "cursor-pointer",
      })}
      rowClassName="hover:bg-surface-elevated transition-colors"
    />
  );
}
