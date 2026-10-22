"use client";

import { SearchOutlined } from "@ant-design/icons";
import { Empty, Input, Tag, Typography } from "antd";
import { useMemo, useState } from "react";

import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";

import { buildGlossary, filterGlossary, groupByCatalog } from "./glossary";

const { Text, Title } = Typography;

const ALL_ENTRIES = buildGlossary();

/** Searchable glossary of catalog columns, grouped by catalog. */
export function GlossarySection() {
  const [query, setQuery] = useState("");
  const groups = useMemo(
    () => groupByCatalog(filterGlossary(ALL_ENTRIES, query)),
    [query]
  );
  const count = groups.reduce((n, [, list]) => n + list.length, 0);

  return (
    <div>
      <Input
        allowClear
        size="large"
        prefix={<SearchOutlined className="text-neutral-500" />}
        placeholder="Filter columns, e.g. parallax, ruwe, W1, flux…"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        aria-label="Filter glossary"
        className="mb-2"
      />
      <Text className="text-neutral-500 text-xs block mb-6">
        {count} of {ALL_ENTRIES.length} columns
      </Text>

      {groups.length === 0 && (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={`No columns match “${query}”`}
        />
      )}

      {groups.map(([catalog, entries]) => (
        <div
          key={catalog}
          id={`glossary-${catalog}`}
          className="mb-8 scroll-mt-24"
        >
          <Title
            level={5}
            className="!mb-3 text-foreground flex items-center gap-2"
          >
            {getSearchCatalogLabel(catalog)}
            <Tag className="!bg-surface !border-border !text-neutral-400">
              {entries.length}
            </Tag>
          </Title>
          <dl className="rounded border border-border divide-y divide-border bg-surface m-0">
            {entries.map((e) => (
              <div
                key={e.field}
                className="grid gap-1 sm:grid-cols-[14rem_1fr] sm:gap-4 px-4 py-2.5"
              >
                <dt className="font-mono text-sm text-foreground break-all">
                  {e.field}
                </dt>
                <dd className="text-sm text-neutral-400 m-0">
                  {e.description}
                </dd>
              </div>
            ))}
          </dl>
        </div>
      ))}
    </div>
  );
}
