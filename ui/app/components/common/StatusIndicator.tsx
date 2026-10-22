"use client";

import { Badge, Space, Typography } from "antd";

const { Text } = Typography;

interface StatusIndicatorProps {
  status: "success" | "processing" | "error" | "warning" | "default";
  text?: string;
  showDot?: boolean;
}

export function StatusIndicator({
  status,
  text,
  showDot = true,
}: StatusIndicatorProps) {
  return (
    <Space size="small">
      {showDot && <Badge status={status} />}
      {text && (
        <Text type="secondary" className="text-sm">
          {text}
        </Text>
      )}
    </Space>
  );
}
