"use client";

import { AimOutlined } from "@ant-design/icons";
import { Typography } from "antd";
import Link from "next/link";

import { buildSearchHref } from "@/app/lib/utils/commands";

const { Text } = Typography;

/** M31's centre, a field every catalog has matches in. */
const EXAMPLE = { name: "M31", ra: 10.6847, dec: 41.269 };

export function EmptyState() {
  return (
    <div className="flex max-w-sm flex-col items-center gap-1 text-center">
      <span className="mb-3 flex size-11 items-center justify-center rounded-full border border-border bg-surface text-lg text-primary">
        <AimOutlined />
      </span>
      <Text strong className="text-base">
        No search yet
      </Text>
      <Text type="secondary">
        Enter a position in the search form, or{" "}
        <Link href={buildSearchHref(EXAMPLE.ra, EXAMPLE.dec)}>
          try {EXAMPLE.name}
        </Link>
        .
      </Text>
    </div>
  );
}
