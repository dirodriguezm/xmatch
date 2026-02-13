"use client";

import { Layout } from "antd";

const { Sider } = Layout;

interface AppSidebarProps {
  children: React.ReactNode;
}

export function AppSidebar({ children }: AppSidebarProps) {
  return (
    <Sider
      width={320}
      trigger={null}
      className="bg-surface border-r border-border overflow-auto h-[calc(100vh-64px)]"
    >
      <div className="h-full overflow-auto">{children}</div>
    </Sider>
  );
}
