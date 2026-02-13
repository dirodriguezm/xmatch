"use client";

import { Space, Typography } from "antd";

import { StatusIndicator } from "@/app/components/common";

const { Text } = Typography;

interface StatusBarProps {
  connected?: boolean;
  message?: string;
}

export function StatusBar({
  connected = true,
  message = "API Connected",
}: StatusBarProps) {
  return (
    <footer className="status-bar flex items-center justify-between h-10">
      <Space size="small">
        <StatusIndicator status={connected ? "success" : "error"} />
        <Text type="secondary" className="text-sm">
          {message}
        </Text>
      </Space>
    </footer>
  );
}
