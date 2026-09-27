"use client";

import {
  CommentOutlined,
  FileTextOutlined,
  LineChartOutlined,
  TableOutlined,
} from "@ant-design/icons";
import { Button, Card, Tag, Typography } from "antd";
import type { ReactNode } from "react";

import {
  CatalogCard,
  CatalogComparisonTable,
  EXTERNAL_SOURCES,
} from "@/app/components/catalogs";
import { PageShell } from "@/app/components/layout";
import { CATALOG_META } from "@/app/lib/constants/catalogMeta";
import {
  CATALOG_COLOR_CLASSES,
  CATALOG_OPTIONS,
} from "@/app/lib/constants/catalogs";
import { DISCUSSIONS_URL, REPO_URL } from "@/app/lib/constants/site";

const { Title, Paragraph, Link, Text } = Typography;

const ADD_CATALOG_GUIDE_URL = `${REPO_URL}/blob/main/ADD_NEW_CATALOG.md`;

function Section({
  id,
  title,
  icon,
  lede,
  children,
}: {
  id: string;
  title: string;
  icon?: ReactNode;
  lede?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section id={id} className="mt-14 scroll-mt-24">
      <Title
        level={3}
        className="!mb-2 text-foreground flex items-center gap-2"
      >
        {icon}
        {title}
      </Title>
      {lede && <Paragraph className="text-neutral-400 !mb-5">{lede}</Paragraph>}
      {children}
    </section>
  );
}

export default function CatalogsPage() {
  return (
    <PageShell
      width="wide"
      title="Catalogs"
      description="What XWave indexes, how precise each catalog is, and which radius to match with. Every catalog is indexed on its own HEALPix grid and searched with its own radius."
      actions={
        <Button href="#request" icon={<CommentOutlined />}>
          Request a catalog
        </Button>
      }
    >
      <nav
        aria-label="Catalogs on this page"
        className="flex flex-wrap gap-2 mb-8"
      >
        {CATALOG_OPTIONS.map((slug) => (
          <a
            key={slug}
            href={`#${slug}`}
            className="inline-flex items-center gap-2 px-3 py-1.5 rounded border border-border bg-surface !text-foreground hover:bg-surface-elevated"
          >
            <span
              aria-hidden
              className={`w-2 h-2 rounded-full inline-block ${CATALOG_COLOR_CLASSES[slug] ?? "bg-gray-500"}`}
            />
            {CATALOG_META[slug].name}
            <span className="text-neutral-500 text-xs">
              {CATALOG_META[slug].release}
            </span>
          </a>
        ))}
        <a
          href="#compare"
          className="inline-flex items-center px-3 py-1.5 rounded border border-border !text-neutral-400 hover:!text-foreground"
        >
          Compare
        </a>
        <a
          href="#external"
          className="inline-flex items-center px-3 py-1.5 rounded border border-border !text-neutral-400 hover:!text-foreground"
        >
          External sources
        </a>
      </nav>

      <div className="space-y-6">
        {CATALOG_OPTIONS.map((slug) => (
          <CatalogCard key={slug} meta={CATALOG_META[slug]} />
        ))}
      </div>

      <Section
        id="compare"
        title="Side by side"
        icon={<TableOutlined />}
        lede="The default radius is what the search form starts with. Widen it for faint sources or older epochs, narrow it in crowded fields."
      >
        <CatalogComparisonTable />
      </Section>

      <Section
        id="external"
        title="Light-curve & external sources"
        icon={<LineChartOutlined />}
        lede="These are not indexed by XWave. The object page queries them live for the position you are looking at, so their coverage and availability are theirs."
      >
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          {EXTERNAL_SOURCES.map((source) => (
            <div
              key={source.name}
              className="rounded-lg border border-border bg-surface p-4"
            >
              <div className="flex items-start justify-between gap-3 mb-1">
                <Link
                  href={source.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="font-medium"
                >
                  {source.name}
                </Link>
                <Tag className="!m-0 shrink-0 !bg-surface-elevated !border-border text-xs">
                  {source.kind}
                </Tag>
              </div>
              <Text className="text-neutral-400 text-sm">{source.usedFor}</Text>
            </div>
          ))}
        </div>
      </Section>

      <Section
        id="request"
        title="Request a catalog"
        icon={<CommentOutlined />}
      >
        <Card className="bg-surface-elevated border-border">
          <div className="flex flex-col md:flex-row md:items-center gap-4 justify-between">
            <div className="max-w-2xl">
              <Paragraph className="text-foreground !mb-2">
                Missing a survey you cross-match against all the time? Tell us
                which one and what you would use it for — requests with a
                concrete science case are the easiest to prioritise.
              </Paragraph>
              <Paragraph className="text-neutral-400 !mb-0">
                Adding a catalog is a self-contained change: a migration, a few
                queries and a catalog adapter. If you would rather contribute it
                yourself, the step-by-step guide is{" "}
                <Link
                  href={ADD_CATALOG_GUIDE_URL}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  <Text code>ADD_NEW_CATALOG.md</Text>
                </Link>{" "}
                in the repository.
              </Paragraph>
            </div>
            <div className="flex flex-wrap gap-2 shrink-0">
              <Button
                type="primary"
                href={DISCUSSIONS_URL}
                target="_blank"
                rel="noopener noreferrer"
                icon={<CommentOutlined />}
              >
                Start a discussion
              </Button>
              <Button
                href={ADD_CATALOG_GUIDE_URL}
                target="_blank"
                rel="noopener noreferrer"
                icon={<FileTextOutlined />}
              >
                Read the guide
              </Button>
            </div>
          </div>
        </Card>
      </Section>
    </PageShell>
  );
}
