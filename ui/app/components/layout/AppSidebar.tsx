"use client";

import {
  DoubleLeftOutlined,
  SearchOutlined,
  SettingOutlined,
} from "@ant-design/icons";
import { Button, Collapse, Grid, Layout, Tooltip } from "antd";
import { useEffect, useState } from "react";

const { Sider } = Layout;

const STORAGE_KEY = "xwave:sidebar-collapsed";
const WIDTH = 320;
const RAIL_WIDTH = 48;

interface AppSidebarProps {
  children: React.ReactNode;
  /** Label of the collapsible panel the sidebar becomes on small screens. */
  mobileLabel?: string;
}

/**
 * The search form in a 320px sidebar on desktop, which folds to a 48px rail to
 * give the results more room (remembered per browser). Below `lg` it becomes a
 * collapsible panel above the results instead.
 */
export function AppSidebar({
  children,
  mobileLabel = "Search parameters",
}: AppSidebarProps) {
  const screens = Grid.useBreakpoint();
  // Unknown before hydration: render the desktop layout, as the server does.
  const isSmall = screens.lg === false;
  const [collapsed, setCollapsed] = useState(false);

  useEffect(() => {
    try {
      // Restoring a saved UI preference after hydration.
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setCollapsed(localStorage.getItem(STORAGE_KEY) === "1");
    } catch {
      // Storage unavailable: keep the default.
    }
  }, []);

  const toggle = (next: boolean) => {
    setCollapsed(next);
    try {
      localStorage.setItem(STORAGE_KEY, next ? "1" : "0");
    } catch {
      // Not persisted; fine.
    }
  };

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
      width={WIDTH}
      collapsed={collapsed}
      collapsedWidth={RAIL_WIDTH}
      trigger={null}
      className="bg-surface border-r border-border overflow-auto h-[calc(100vh-64px)]"
    >
      {collapsed ? (
        <div className="flex h-full flex-col items-center gap-2 pt-4">
          <Tooltip title="Show search" placement="right">
            <Button
              type="text"
              icon={<SearchOutlined />}
              aria-label="Show search"
              aria-expanded={false}
              onClick={() => toggle(false)}
            />
          </Tooltip>
        </div>
      ) : (
        <div className="relative h-full overflow-auto">
          <Tooltip title="Hide search" placement="right">
            <Button
              type="text"
              size="small"
              icon={<DoubleLeftOutlined />}
              aria-label="Hide search"
              aria-expanded
              onClick={() => toggle(true)}
              className="!absolute right-2 top-4 z-10"
            />
          </Tooltip>
          {children}
        </div>
      )}
    </Sider>
  );
}
