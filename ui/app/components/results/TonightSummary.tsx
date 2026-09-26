"use client";

import { FieldTimeOutlined } from "@ant-design/icons";
import { useQuery } from "@tanstack/react-query";
import { Flex, Select, Typography } from "antd";
import { useState } from "react";

import {
  DEFAULT_OBSERVATORY_ID,
  OBSERVATORIES,
  OBSERVATORY_TIME_ZONE,
} from "@/app/lib/constants/observatories";

const { Text } = Typography;

interface TonightSummaryProps {
  ra: number;
  dec: number;
}

/**
 * One line of visibility for the searched position. Every match lies within
 * the search radius (≤ 2′) of it, so one line covers them all; per-row values
 * would be identical. astronomy-engine is imported on demand.
 */
export function TonightSummary({ ra, dec }: TonightSummaryProps) {
  const [siteId, setSiteId] = useState(DEFAULT_OBSERVATORY_ID);
  const site = OBSERVATORIES.find((o) => o.id === siteId) ?? OBSERVATORIES[0];

  const { data, isError } = useQuery({
    queryKey: ["tonight-summary", ra, dec, siteId],
    queryFn: async () => {
      const m = await import("@/app/lib/utils/observability");
      const night = m.tonightInTimeZone(OBSERVATORY_TIME_ZONE);
      return m.summarizeVisibility(
        m.computeNightVisibility(ra, dec, site, night),
        site.label
      );
    },
    staleTime: 10 * 60 * 1000,
  });

  return (
    <Flex align="center" gap={8} wrap="wrap">
      <FieldTimeOutlined className="text-border" />
      <Text type="secondary" className="text-sm">
        Tonight from
      </Text>
      <Select
        size="small"
        variant="borderless"
        value={siteId}
        onChange={setSiteId}
        popupMatchSelectWidth={false}
        options={OBSERVATORIES.map((o) => ({ label: o.label, value: o.id }))}
        aria-label="Observatory"
      />
      <Text className="text-sm">
        {isError ? "visibility unavailable" : (data ?? "…")}
      </Text>
    </Flex>
  );
}
