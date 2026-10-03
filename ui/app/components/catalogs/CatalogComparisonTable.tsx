"use client";

import { Table } from "antd";
import type { ColumnsType } from "antd/es/table";

import {
  CATALOG_META,
  type CatalogMeta,
} from "@/app/lib/constants/catalogMeta";
import {
  CATALOG_COLOR_CLASSES,
  CATALOG_OPTIONS,
} from "@/app/lib/constants/catalogs";

import { adsUrl, formatDefaultRadius } from "./CatalogCard";

const ROWS = CATALOG_OPTIONS.map((slug) => CATALOG_META[slug]);

const COLUMNS: ColumnsType<CatalogMeta> = [
  {
    title: "Catalog",
    key: "name",
    fixed: "left",
    render: (_, m) => (
      <a
        href={`#${m.slug}`}
        className="inline-flex items-center gap-2 !text-foreground whitespace-nowrap"
      >
        <span
          aria-hidden
          className={`w-2 h-2 rounded-full inline-block ${CATALOG_COLOR_CLASSES[m.slug] ?? "bg-gray-500"}`}
        />
        {m.name}
      </a>
    ),
  },
  { title: "Release", dataIndex: "release", key: "release" },
  { title: "Wavelength", dataIndex: "wavelength", key: "wavelength" },
  { title: "Coverage", dataIndex: "coverage", key: "coverage" },
  {
    title: "Sources",
    dataIndex: "sources",
    key: "sources",
    className: "whitespace-nowrap",
  },
  { title: "Astrometry", dataIndex: "astrometry", key: "astrometry" },
  {
    title: "Default radius",
    key: "radius",
    align: "right",
    render: (_, m) => (
      <span className="font-mono">{formatDefaultRadius(m.slug)}</span>
    ),
  },
  {
    title: "Reference",
    key: "reference",
    render: (_, m) => (
      <a
        href={adsUrl(m.reference.bibcode)}
        target="_blank"
        rel="noopener noreferrer"
      >
        {m.reference.label}
      </a>
    ),
  },
];

export function CatalogComparisonTable() {
  return (
    <Table<CatalogMeta>
      rowKey="slug"
      size="small"
      columns={COLUMNS}
      dataSource={ROWS}
      pagination={false}
      scroll={{ x: 960 }}
      className="border border-border rounded-lg overflow-hidden"
    />
  );
}
