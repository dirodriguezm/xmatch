"use client";

import { Typography } from "antd";
import Link from "next/link";

import { buildSearchHref } from "@/app/lib/utils/commands";

const { Text } = Typography;

/** M31's centre, a field every catalog has matches in. */
const EXAMPLE = { name: "M31", ra: 10.6847, dec: 41.269 };

export function EmptyState() {
  return (
    <div className="flex h-full max-w-sm flex-col justify-center gap-1">
      <Text strong className="text-base">
        No search yet
      </Text>
      <Text type="secondary">
        Enter a position and radius in the search form, or{" "}
        <Link href={buildSearchHref(EXAMPLE.ra, EXAMPLE.dec)}>
          try {EXAMPLE.name}
        </Link>
        .
      </Text>
    </div>
  );
}
