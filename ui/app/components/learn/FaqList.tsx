"use client";

import { LinkOutlined } from "@ant-design/icons";
import { Collapse, Table, Typography } from "antd";
import Link from "next/link";

import { CATALOG_META } from "@/app/lib/constants/catalogMeta";
import {
  CATALOG_DEFAULT_RADII,
  CATALOG_OPTIONS,
} from "@/app/lib/constants/catalogs";
import { FAQ_ITEMS, type FaqItem } from "@/app/lib/constants/faq";
import { isExternalHref } from "@/app/lib/constants/site";

const { Paragraph } = Typography;

const RADIUS_ROWS = CATALOG_OPTIONS.map((slug) => ({
  key: slug,
  catalog:
    `${CATALOG_META[slug].name} ${CATALOG_META[slug].release.startsWith(CATALOG_META[slug].name) ? "" : CATALOG_META[slug].release}`.trim(),
  astrometry: CATALOG_META[slug].astrometry,
  radius: `${CATALOG_DEFAULT_RADII[slug].radius}${CATALOG_DEFAULT_RADII[slug].unit === "arcsec" ? "″" : ` ${CATALOG_DEFAULT_RADII[slug].unit}`}`,
}));

function RadiusTable() {
  return (
    <Table
      size="small"
      pagination={false}
      className="mb-3"
      dataSource={RADIUS_ROWS}
      columns={[
        { title: "Catalog", dataIndex: "catalog" },
        { title: "Typical positional accuracy", dataIndex: "astrometry" },
        { title: "Default radius", dataIndex: "radius", width: 120 },
      ]}
    />
  );
}

function Answer({ item }: { item: FaqItem }) {
  return (
    <div className="text-neutral-400">
      {item.answer.map((p, i) => (
        <Paragraph key={i} className="!text-neutral-300">
          {p}
        </Paragraph>
      ))}
      {item.extra === "radius-table" && <RadiusTable />}
      <div className="flex flex-wrap gap-x-4 gap-y-1">
        {item.links?.map((l) =>
          isExternalHref(l.href) ? (
            <a key={l.href} href={l.href} target="_blank" rel="noreferrer">
              {l.label} ↗
            </a>
          ) : (
            <Link key={l.href} href={l.href}>
              {l.label} →
            </Link>
          )
        )}
        <a
          href={`#faq-${item.id}`}
          className="!text-neutral-500 hover:!text-neutral-300 text-xs inline-flex items-center gap-1 ml-auto"
          aria-label="Link to this question"
        >
          <LinkOutlined /> link
        </a>
      </div>
    </div>
  );
}

/** FAQ accordion. `openIds` are expanded (e.g. from a #faq-<id> hash). */
export function FaqList({
  openIds,
  onChange,
}: {
  openIds: string[];
  onChange: (ids: string[]) => void;
}) {
  return (
    <Collapse
      activeKey={openIds}
      onChange={(keys) => onChange(Array.isArray(keys) ? keys : [keys])}
      className="bg-surface"
      items={FAQ_ITEMS.map((item) => ({
        key: item.id,
        label: (
          <span
            id={`faq-${item.id}`}
            className="scroll-mt-24 text-foreground font-medium"
          >
            {item.question}
          </span>
        ),
        children: <Answer item={item} />,
      }))}
    />
  );
}
