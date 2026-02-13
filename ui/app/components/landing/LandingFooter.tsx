"use client";

import { BookOutlined, CheckCircleOutlined } from "@ant-design/icons";
import { Divider, Space, Typography } from "antd";

const { Text, Link } = Typography;

export function LandingFooter() {
  return (
    <div className="fixed bottom-0 left-0 right-0 py-4 px-6 bg-surface border-t border-border">
      <Space
        separator={<Divider orientation="vertical" />}
        className="flex justify-center flex-wrap"
      >
        <Text type="secondary">XWave v1.0</Text>

        <Space size="small">
          <CheckCircleOutlined className="text-success" />
          <Text type="secondary">All systems operational</Text>
        </Space>

        <Link
          href="https://ifa.uv.cl/xwave/swagger/index.html#"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Space size={4}>
            <BookOutlined />
            Documentation
          </Space>
        </Link>
      </Space>
    </div>
  );
}
