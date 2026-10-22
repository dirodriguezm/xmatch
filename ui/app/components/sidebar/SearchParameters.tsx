"use client";

import { RadiusSettingOutlined } from "@ant-design/icons";
import { Flex, Space, Typography } from "antd";

const { Title, Text } = Typography;

export function SearchParameters() {
  return (
    <div className="sidebar-section">
      <Flex vertical gap="middle">
        <Space align="center">
          <RadiusSettingOutlined className="text-primary" />
          <Title level={5} className="m-0!">
            Search Parameters
          </Title>
        </Space>
        <Text type="secondary">
          Configure search radius per catalog in the Catalogs &amp; Radii
          section.
        </Text>
      </Flex>
    </div>
  );
}
