"use client";

import { RadiusSettingOutlined } from "@ant-design/icons";
import { Flex, InputNumber, Select, Space, Typography } from "antd";

import { useSearchParams } from "@/app/hooks/useSearchParamsSync";

const { Title, Text } = Typography;

const unitOptions = [
  { value: "arcsec", label: "arcsec" },
  { value: "arcmin", label: "arcmin" },
  { value: "deg", label: "deg" },
];

export function SearchParameters() {
  const { radius, unit, setRadius, setUnit } = useSearchParams();

  return (
    <div className="sidebar-section">
      <Flex vertical gap="middle">
        <Space align="center">
          <RadiusSettingOutlined className="text-primary" />
          <Title level={5} className="!m-0">
            Search Parameters
          </Title>
        </Space>

        <Flex vertical gap="small">
          <Text type="secondary">Search Radius</Text>
          <Space.Compact className="w-full">
            <InputNumber
              value={radius}
              min={0}
              step={0.1}
              onChange={(value) => setRadius(value ?? 1)}
              className="flex-1"
            />
            <Select
              value={unit}
              options={unitOptions}
              onChange={(value) => setUnit(value)}
              className="w-[100px]"
            />
          </Space.Compact>
        </Flex>
      </Flex>
    </div>
  );
}
