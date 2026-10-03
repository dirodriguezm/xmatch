"use client";

import { Space, Typography } from "antd";
import Link from "next/link";

import { StatusIndicator } from "@/app/components/common";
import { useApiHealth } from "@/app/components/feedback/useApiHealth";

const { Text } = Typography;

/** Live XWave API status from /api/health; links to the Status page. */
export function StatusBar() {
  const { data, isPending, isError } = useApiHealth();

  const up = data?.status === "up";
  const status = isPending ? "processing" : up ? "success" : "error";
  const message = isPending
    ? "Checking API…"
    : isError || !data
      ? "API status unknown"
      : up
        ? `API online · ${data.latencyMs} ms`
        : "API unreachable";

  return (
    <footer className="status-bar flex items-center justify-between h-10">
      <Link href="/status" className="hover:opacity-80">
        <Space size="small">
          <StatusIndicator status={status} />
          <Text type="secondary" className="text-sm">
            {message}
          </Text>
        </Space>
      </Link>
    </footer>
  );
}
