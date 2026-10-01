"use client";

import { SettingOutlined } from "@ant-design/icons";
import { Collapse, Grid, Layout } from "antd";

const { Sider } = Layout;

interface AppSidebarProps {
  children: React.ReactNode;
  /** Label of the collapsible panel the sidebar becomes on small screens. */
  mobileLabel?: string;
}

/**
 * Fixed 320px sidebar on desktop. Below `lg` it would leave the results a
 * sliver of the screen, so it becomes a collapsible panel above them instead.
 */
export function AppSidebar({
  children,
  mobileLabel = "Search parameters",
}: AppSidebarProps) {
  const screens = Grid.useBreakpoint();
  // Unknown before hydration: render the desktop layout, as the server does.
  const isSmall = screens.lg === false;

  if (isSmall) {
    return (
      <div className="w-full bg-surface border-b border-border">
        <Collapse
          ghost
          items={[
            {
              key: "params",
              label: (
                <span className="inline-flex items-center gap-2 text-foreground">
                  <SettingOutlined />
                  {mobileLabel}
                </span>
              ),
              children,
            },
          ]}
        />
      </div>
    );
  }

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
