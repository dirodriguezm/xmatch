"use client";

import { Typography } from "antd";

const { Text } = Typography;

export function LightCurveSkeleton({
  message = "Loading light curves…",
}: {
  message?: string;
}) {
  return (
    <div className="w-full">
      <div className="relative h-48 w-full overflow-hidden rounded border border-border bg-surface-elevated/40">
        {/* Faint left axis */}
        <div className="absolute bottom-8 left-12 top-4 w-px bg-border/60" />
        {/* Faint bottom axis */}
        <div className="absolute bottom-8 left-12 right-4 h-px bg-border/60" />
        {/* Quiet horizontal gridlines */}
        <div className="absolute left-12 right-4 top-[28%] h-px bg-border/25" />
        <div className="absolute left-12 right-4 top-[58%] h-px bg-border/25" />
        {/* Shimmer sweep */}
        <div className="lc-skeleton-shimmer" />
      </div>
      <div className="mt-2 text-center">
        <Text type="secondary" className="text-xs">
          {message}
        </Text>
      </div>
    </div>
  );
}
