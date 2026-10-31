"use client";

import { LinkOutlined, WifiOutlined } from "@ant-design/icons";
import { Button, Segmented, Tag, Typography } from "antd";
import { useState } from "react";

import { PageShell } from "@/app/components/layout";
import { CHANGELOG, type ChangelogTag } from "@/app/lib/constants/changelog";

const { Title, Paragraph } = Typography;

const TAG_COLORS: Record<ChangelogTag, string> = {
  New: "blue",
  Improved: "green",
  Catalogs: "purple",
  API: "cyan",
  Fixed: "orange",
};

const FILTERS = ["All", ...Object.keys(TAG_COLORS)] as const;
type Filter = (typeof FILTERS)[number];

function formatDate(iso: string) {
  return new Date(`${iso}T12:00:00Z`).toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  });
}

export default function ChangelogPage() {
  const [filter, setFilter] = useState<Filter>("All");
  const entries =
    filter === "All"
      ? CHANGELOG
      : CHANGELOG.filter((e) => e.tags.includes(filter as ChangelogTag));

  return (
    <PageShell
      title="Changelog"
      description="New features, catalogs and API changes in XWave."
      actions={
        <Button icon={<WifiOutlined />} href="/changelog/rss.xml">
          RSS
        </Button>
      }
    >
      <Segmented
        options={[...FILTERS]}
        value={filter}
        onChange={(v) => setFilter(v as Filter)}
        className="mb-8"
      />

      <ol className="relative list-none pl-0 m-0 mb-12 border-l border-border ml-2">
        {entries.map((e) => (
          <li
            key={e.id}
            id={e.id}
            className="relative pl-8 pb-10 last:pb-2 scroll-mt-24 group"
          >
            <span className="absolute -left-[7px] top-1.5 w-3 h-3 rounded-full bg-surface-elevated border-2 border-primary" />
            <div className="flex flex-wrap items-center gap-2 mb-1">
              <time
                dateTime={e.date}
                className="font-mono text-xs text-neutral-400"
              >
                {formatDate(e.date)}
              </time>
              {e.tags.map((t) => (
                <Tag key={t} color={TAG_COLORS[t]} className="!mr-0">
                  {t}
                </Tag>
              ))}
            </div>
            <Title level={4} className="!mt-0 !mb-2 text-foreground">
              <a
                href={`#${e.id}`}
                className="!text-foreground hover:!text-primary inline-flex items-center gap-2"
              >
                {e.title}
                <LinkOutlined className="text-sm opacity-0 group-hover:opacity-60 transition-opacity" />
              </a>
            </Title>
            <Paragraph className="!text-neutral-300 !mb-3">
              {e.summary}
            </Paragraph>
            <ul className="text-neutral-400 space-y-1.5 pl-5 list-disc marker:text-neutral-600 m-0">
              {e.items.map((item) => (
                <li key={item}>{item}</li>
              ))}
            </ul>
          </li>
        ))}
      </ol>
    </PageShell>
  );
}
