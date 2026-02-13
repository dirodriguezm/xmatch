"use client";

import { AimOutlined } from "@ant-design/icons";
import { Flex, Input, Space, Typography } from "antd";

import { useSearchParams } from "@/app/hooks/useSearchParamsSync";

const { Title, Link, Text } = Typography;

interface TargetCoordinatesProps {
  onResolveClick?: () => void;
}

export function TargetCoordinates({ onResolveClick }: TargetCoordinatesProps) {
  const { ra, dec, setRa, setDec } = useSearchParams();

  return (
    <div className="sidebar-section">
      <Flex vertical gap="middle">
        <Space align="center">
          <AimOutlined className="text-primary" />
          <Title level={5} className="!m-0">
            Target Coordinates
          </Title>
        </Space>

        <Flex vertical gap="small">
          <Flex vertical gap={4}>
            <Text type="secondary">RA</Text>
            <Input
              placeholder="e.g., 180.5 or 12:00:00"
              value={ra}
              onChange={(e) => setRa(e.target.value)}
            />
          </Flex>
          <Flex vertical gap={4}>
            <Text type="secondary">Dec</Text>
            <Input
              placeholder="e.g., -45.0 or -45:00:00"
              value={dec}
              onChange={(e) => setDec(e.target.value)}
            />
          </Flex>
          <Link onClick={onResolveClick}>Resolve Target Name</Link>
        </Flex>
      </Flex>
    </div>
  );
}
