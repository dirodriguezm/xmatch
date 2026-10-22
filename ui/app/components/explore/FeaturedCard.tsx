"use client";

import { ArrowRightOutlined } from "@ant-design/icons";
import { Button, Tag, Typography } from "antd";
import Link from "next/link";

import { BasketButton } from "@/app/components/basket/BasketButton";
import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import type { FeaturedObject } from "@/app/lib/constants/featured";
import { buildObjectUrl } from "@/app/lib/utils/urls";

const { Title, Paragraph, Text } = Typography;

function formatCoords(o: FeaturedObject) {
  const sign = o.dec >= 0 ? "+" : "−";
  return `${o.ra.toFixed(4)}°, ${sign}${Math.abs(o.dec).toFixed(4)}°`;
}

/** Grid card for one curated target. */
export function FeaturedCard({ object }: { object: FeaturedObject }) {
  const href = buildObjectUrl(object.objectId, object.catalog);
  return (
    <div className="flex flex-col h-full rounded-lg border border-border bg-surface p-5 transition-colors hover:bg-surface-elevated">
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0">
          <Text className="text-xs uppercase tracking-wide text-neutral-400!">
            {object.kind}
          </Text>
          <Title level={5} className="mt-1! mb-0! text-foreground">
            <Link href={href} className="text-foreground! hover:underline">
              {object.name}
            </Link>
          </Title>
        </div>
        <BasketButton
          size="small"
          item={{
            objectId: object.objectId,
            catalog: object.catalog,
            ra: object.ra,
            dec: object.dec,
          }}
        />
      </div>
      <Paragraph className="mt-3! mb-3! text-sm text-neutral-400 flex-1">
        {object.blurb}
      </Paragraph>
      <div className="flex flex-wrap items-center gap-2">
        <Tag color={getSearchCatalogColor(object.catalog)} className="m-0!">
          {getSearchCatalogLabel(object.catalog)}
        </Tag>
        <Text className="font-mono text-xs text-neutral-400!">
          {formatCoords(object)}
        </Text>
      </div>
    </div>
  );
}

/** Hero card for the deterministic "Object of the week". */
export function ObjectOfTheWeek({
  object,
  week,
}: {
  object: FeaturedObject;
  week: number;
}) {
  const href = buildObjectUrl(object.objectId, object.catalog);
  return (
    <section className="rounded-xl border border-border bg-surface-elevated p-6 md:p-8">
      <Text className="text-xs uppercase tracking-wide text-neutral-400!">
        Object of the week · week {week}
      </Text>
      <Title level={3} className="mt-2! mb-1! text-foreground">
        {object.name}
      </Title>
      <Text className="text-neutral-400!">{object.kind}</Text>
      <Paragraph className="mt-4! mb-4! text-base text-foreground max-w-3xl">
        {object.blurb}
      </Paragraph>
      <div className="mb-5">
        <Text className="text-sm text-neutral-400!">What to look at:</Text>
        <ul className="mt-1 mb-0 pl-5 text-sm text-neutral-400">
          {object.lookFor.map((item) => (
            <li key={item}>{item}</li>
          ))}
        </ul>
      </div>
      <div className="flex flex-wrap items-center gap-3">
        <Link href={href}>
          <Button
            type="primary"
            icon={<ArrowRightOutlined />}
            iconPlacement="end"
          >
            Open object page
          </Button>
        </Link>
        <BasketButton
          item={{
            objectId: object.objectId,
            catalog: object.catalog,
            ra: object.ra,
            dec: object.dec,
          }}
        />
        <Text className="font-mono text-xs text-neutral-400!">
          {object.objectId}
        </Text>
      </div>
    </section>
  );
}
