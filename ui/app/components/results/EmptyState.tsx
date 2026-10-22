"use client";

import { CompassOutlined } from "@ant-design/icons";
import { Empty, Typography } from "antd";

const { Title, Text } = Typography;

export function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center h-full">
      <Empty
        image={<CompassOutlined className="text-[64px] text-border" />}
        description={
          <div className="text-center">
            <Title level={4} className="text-foreground mb-2">
              Ready for Input
            </Title>
            <Text type="secondary">
              Enter target coordinates and search parameters in the sidebar to
              begin cross-matching.
            </Text>
          </div>
        }
      />
    </div>
  );
}
