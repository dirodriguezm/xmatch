"use client";

import {
  DatabaseOutlined,
  DownOutlined,
  NodeIndexOutlined,
  StarOutlined,
  ThunderboltOutlined,
} from "@ant-design/icons";
import { Button, Dropdown, Skeleton, Space, Typography } from "antd";
import Link from "next/link";
import type { ReactNode } from "react";
import { useSyncExternalStore } from "react";

import {
  FeaturedCard,
  ObjectOfTheWeek,
} from "@/app/components/explore/FeaturedCard";
import { useRandomObject } from "@/app/components/explore/useRandomObject";
import { PageShell } from "@/app/components/layout";
import {
  FEATURED_OBJECTS,
  isoWeek,
  objectOfTheWeek,
} from "@/app/lib/constants/featured";

const { Title, Paragraph } = Typography;

const noopSubscribe = () => () => {};

/**
 * ISO week, read on the client only so a prerendered page never pins a stale
 * "object of the week" (and hydration never mismatches).
 */
function useIsoWeek(): number | null {
  return useSyncExternalStore(
    noopSubscribe,
    () => isoWeek(new Date()),
    () => null
  );
}

function RandomButton() {
  const { go, loading } = useRandomObject();
  return (
    <Space.Compact>
      <Button
        type="primary"
        icon={<ThunderboltOutlined />}
        loading={loading}
        onClick={() => go("sky")}
      >
        Random object
      </Button>
      <Dropdown
        placement="bottomRight"
        menu={{
          items: [
            {
              key: "featured",
              icon: <StarOutlined />,
              label: "Random featured object",
            },
          ],
          onClick: () => go("featured"),
        }}
      >
        <Button
          type="primary"
          icon={<DownOutlined />}
          aria-label="More random options"
        />
      </Dropdown>
    </Space.Compact>
  );
}

function LinkCard({
  href,
  icon,
  title,
  children,
}: {
  href: string;
  icon: ReactNode;
  title: string;
  children: ReactNode;
}) {
  return (
    <Link
      href={href}
      className="block rounded-lg border border-border bg-surface p-5 no-underline transition-colors hover:bg-surface-elevated"
    >
      <div className="flex items-center gap-2 text-foreground font-medium">
        {icon}
        {title}
      </div>
      <p className="mt-2 mb-0 text-sm text-neutral-400">{children}</p>
    </Link>
  );
}

export default function ExplorePage() {
  const week = useIsoWeek();
  const featured = week === null ? null : objectOfTheWeek(new Date());

  return (
    <PageShell
      width="wide"
      title="Explore"
      description="Curated targets that show off what the object pages can do, plus a random jump anywhere on the sky."
      actions={<RandomButton />}
    >
      {featured && week !== null ? (
        <ObjectOfTheWeek object={featured} week={week} />
      ) : (
        <div className="rounded-xl border border-border bg-surface-elevated p-8">
          <Skeleton active paragraph={{ rows: 4 }} />
        </div>
      )}

      <Title level={4} className="mt-12! mb-4! text-foreground">
        Featured objects
      </Title>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {FEATURED_OBJECTS.map((o) => (
          <FeaturedCard key={o.slug} object={o} />
        ))}
      </div>

      <Title level={4} className="mt-12! mb-2! text-foreground">
        Go further
      </Title>
      <Paragraph className="text-neutral-400 mb-4!">
        Star any object to collect it in your basket, then cross-match the whole
        list at once.
      </Paragraph>
      <div className="grid gap-4 sm:grid-cols-2 pb-12">
        <LinkCard
          href="/bulk"
          icon={<NodeIndexOutlined />}
          title="Bulk cross-match"
        >
          Upload or paste a list of positions and match them against every
          catalog in one go.
        </LinkCard>
        <LinkCard href="/catalogs" icon={<DatabaseOutlined />} title="Catalogs">
          What Gaia DR3, AllWISE and eROSITA cover, their astrometric accuracy
          and which radius to match with.
        </LinkCard>
      </div>
    </PageShell>
  );
}
