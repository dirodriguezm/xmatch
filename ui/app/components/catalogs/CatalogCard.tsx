"use client";

import { ExportOutlined, WarningOutlined } from "@ant-design/icons";
import { Card, Collapse, Tag, Typography } from "antd";
import type { ReactNode } from "react";

import { CATALOG_FIELD_DESCRIPTIONS } from "@/app/lib/constants/catalogFields";
import type { CatalogMeta } from "@/app/lib/constants/catalogMeta";
import {
  CATALOG_COLOR_CLASSES,
  CATALOG_DEFAULT_RADII,
} from "@/app/lib/constants/catalogs";

import { CoverageMap } from "./CoverageMap";

const { Title, Text, Link, Paragraph } = Typography;

export function adsUrl(bibcode: string): string {
  return `https://ui.adsabs.harvard.edu/abs/${encodeURIComponent(bibcode)}`;
}

export function formatDefaultRadius(slug: CatalogMeta["slug"]): string {
  const { radius, unit } = CATALOG_DEFAULT_RADII[slug];
  const symbol = unit === "arcsec" ? "″" : unit === "arcmin" ? "′" : "°";
  return `${radius}${symbol}`;
}

function Fact({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="min-w-0">
      <dt className="text-xs uppercase tracking-wider text-neutral-500 mb-0.5">
        {label}
      </dt>
      <dd className="m-0 text-foreground text-sm">{children}</dd>
    </div>
  );
}

function ColumnList({ slug }: { slug: string }) {
  const fields = Object.entries(CATALOG_FIELD_DESCRIPTIONS[slug] ?? {});
  if (fields.length === 0) {
    return (
      <Text className="text-neutral-400 text-sm">
        No documented columns yet.
      </Text>
    );
  }
  return (
    <dl className="m-0 divide-y divide-border">
      {fields.map(([field, description]) => (
        <div
          key={field}
          className="grid grid-cols-1 sm:grid-cols-[14rem_1fr] gap-x-4 gap-y-0.5 py-1.5"
        >
          <dt className="font-mono text-xs text-foreground break-all">
            {field}
          </dt>
          <dd className="m-0 text-xs text-neutral-400">{description}</dd>
        </div>
      ))}
    </dl>
  );
}

export function CatalogCard({ meta }: { meta: CatalogMeta }) {
  const fieldCount = Object.keys(
    CATALOG_FIELD_DESCRIPTIONS[meta.slug] ?? {}
  ).length;

  return (
    <Card
      id={meta.slug}
      className="bg-surface border-border scroll-mt-24"
      styles={{ body: { padding: 0 } }}
    >
      <div className="p-5 md:p-6">
        <div className="flex flex-wrap items-center gap-3 mb-1">
          <span
            aria-hidden
            className={`w-3 h-3 rounded-full inline-block ${CATALOG_COLOR_CLASSES[meta.slug] ?? "bg-gray-500"}`}
          />
          <Title level={3} className="!m-0 text-foreground">
            <a href={`#${meta.slug}`} className="!text-foreground">
              {meta.name}
            </a>
          </Title>
          <Tag className="!m-0 !bg-surface-elevated !border-border font-mono">
            {meta.release}
          </Tag>
        </div>
        <Paragraph className="text-neutral-400 !mb-5">
          {meta.wavelength}
        </Paragraph>

        <div className="grid grid-cols-1 md:grid-cols-[1fr_16rem] gap-6">
          <div className="min-w-0">
            <dl className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4 m-0">
              <Fact label="Coverage">{meta.coverage}</Fact>
              <Fact label="Published sources">{meta.sources}</Fact>
              <Fact label="Astrometric precision">{meta.astrometry}</Fact>
              <Fact label="Recommended match radius">
                <span className="font-mono">
                  {formatDefaultRadius(meta.slug)}
                </span>{" "}
                <span className="text-neutral-500 text-xs">
                  (search default)
                </span>
              </Fact>
              <Fact label="Reference">
                <Link
                  href={adsUrl(meta.reference.bibcode)}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {meta.reference.label}
                </Link>
                <div className="font-mono text-xs text-neutral-500 break-all">
                  {meta.reference.bibcode}
                  {meta.reference.doi && ` · doi:${meta.reference.doi}`}
                </div>
              </Fact>
              <Fact label="Homepage">
                <Link
                  href={meta.homepage}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="break-all"
                >
                  {meta.homepage.replace(/^https?:\/\//, "")}{" "}
                  <ExportOutlined className="text-xs" />
                </Link>
              </Fact>
            </dl>

            {meta.knownIssues.length > 0 && (
              <div className="mt-5">
                <div className="text-xs uppercase tracking-wider text-neutral-500 mb-2 flex items-center gap-1.5">
                  <WarningOutlined /> Known issues
                </div>
                <ul className="text-sm text-neutral-400 space-y-1.5 pl-5 list-disc marker:text-border m-0">
                  {meta.knownIssues.map((issue) => (
                    <li key={issue}>{issue}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>

          <CoverageMap catalog={meta.slug} label={meta.name} />
        </div>
      </div>

      <Collapse
        ghost
        className="border-t border-border !rounded-none"
        items={[
          {
            key: "columns",
            label: (
              <span className="text-neutral-300 text-sm">
                Columns returned by /metadata{" "}
                <span className="text-neutral-500">({fieldCount})</span>
              </span>
            ),
            children: <ColumnList slug={meta.slug} />,
          },
        ]}
      />
    </Card>
  );
}
