"use client";

import { Layout, Space, Tag, Typography } from "antd";
import Link from "next/link";

import { Logo } from "@/app/components/common";

const { Header } = Layout;
const { Title } = Typography;

export function AppHeader() {
  return (
    <Header className="flex items-center px-6 border-b border-border h-16 leading-[64px]">
      <Space size="middle" align="center">
        <Link href="/" className="flex items-center gap-3 no-underline">
          <Logo />
          <Title level={4} className="!m-0 text-foreground">
            XWave
          </Title>
        </Link>
        <Tag color="blue">v1.0</Tag>
      </Space>
    </Header>
  );
}
