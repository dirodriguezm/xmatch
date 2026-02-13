"use client";

import { ThunderboltOutlined } from "@ant-design/icons";
import { Alert, Flex, Space, Switch, Typography } from "antd";
import { useState } from "react";

const { Title, Text } = Typography;

export function PrecisionMode() {
  const [isHighPrecision, setIsHighPrecision] = useState(false);

  return (
    <div className="sidebar-section">
      <Flex vertical gap="middle">
        <Space align="center">
          <ThunderboltOutlined className="text-primary" />
          <Title level={5} className="!m-0">
            Precision Mode
          </Title>
        </Space>

        <Flex vertical gap="small">
          <Space>
            <Switch
              checked={isHighPrecision}
              onChange={setIsHighPrecision}
              checkedChildren="High"
              unCheckedChildren="Standard"
            />
            <Text>{isHighPrecision ? "High Precision" : "Standard Mode"}</Text>
          </Space>

          <Alert
            type="info"
            showIcon={false}
            message={
              isHighPrecision
                ? "High precision mode uses exact spherical trigonometry for cross-matching. Recommended for small search radii."
                : "Standard mode uses optimized algorithms suitable for most use cases."
            }
            className="bg-surface-elevated border border-border"
          />
        </Flex>
      </Flex>
    </div>
  );
}
