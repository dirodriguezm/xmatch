"use client";

import { ArrowLeftOutlined, MenuOutlined } from "@ant-design/icons";
import { Button, Drawer, Layout, Space, Tag, Typography } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState } from "react";

import { HeaderBasket } from "@/app/components/basket/HeaderBasket";
import { CommandPaletteTrigger } from "@/app/components/command/CommandPalette";
import { Logo } from "@/app/components/common";
import {
  FOOTER_GROUPS,
  isExternalHref,
  NAV_ITEMS,
} from "@/app/lib/constants/site";

const { Header } = Layout;
const { Title } = Typography;

interface AppHeaderProps {
  onBack?: () => void;
  backLabel?: string;
}

export function AppHeader({
  onBack,
  backLabel = "Back to Results",
}: AppHeaderProps) {
  const pathname = usePathname();
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <Header className="flex items-center px-4 md:px-6 border-b border-border h-16 leading-[64px] gap-4">
      {onBack && (
        <Button
          icon={<ArrowLeftOutlined />}
          type="text"
          onClick={onBack}
          className="text-muted hover:text-foreground shrink-0"
          aria-label={backLabel}
        >
          <span className="hidden lg:inline">{backLabel}</span>
        </Button>
      )}
      <Space size="middle" align="center" className="shrink-0">
        <Link href="/" className="flex items-center gap-3 no-underline">
          <Logo />
          <Title level={4} className="!m-0 text-foreground">
            XWave
          </Title>
        </Link>
        <Tag color="blue" className="hidden sm:inline-block">
          v1.0
        </Tag>
      </Space>
      <nav
        aria-label="Main"
        className="hidden md:flex items-center gap-1 min-w-0 overflow-x-auto"
      >
        {NAV_ITEMS.map((item) => {
          const active =
            pathname === item.href || pathname?.startsWith(`${item.href}/`);
          return (
            <Link
              key={item.href}
              href={item.href}
              aria-current={active ? "page" : undefined}
              className={`px-3 leading-8 rounded-md text-sm no-underline whitespace-nowrap transition-colors ${
                active
                  ? "bg-surface-elevated !text-foreground"
                  : "!text-neutral-400 hover:!text-foreground hover:bg-surface-elevated"
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>
      <div className="ml-auto flex items-center gap-2 shrink-0">
        <CommandPaletteTrigger />
        <HeaderBasket />
        <Button
          type="text"
          icon={<MenuOutlined />}
          aria-label="Open menu"
          className="md:!hidden"
          onClick={() => setMenuOpen(true)}
        />
      </div>
      {/* The nav links are hidden below md; this drawer replaces them. */}
      <Drawer
        open={menuOpen}
        onClose={() => setMenuOpen(false)}
        placement="right"
        size={280}
        title="XWave"
      >
        <nav aria-label="Mobile" className="flex flex-col gap-6">
          {FOOTER_GROUPS.map((group) => (
            <div key={group.title}>
              <div className="text-xs uppercase tracking-wider text-neutral-500 mb-2">
                {group.title}
              </div>
              <ul className="list-none m-0 p-0 flex flex-col gap-1">
                {group.items.map((item) => {
                  const external = isExternalHref(item.href);
                  const active = pathname === item.href;
                  return (
                    <li key={item.href}>
                      <Link
                        href={item.href}
                        onClick={() => setMenuOpen(false)}
                        {...(external
                          ? { target: "_blank", rel: "noopener noreferrer" }
                          : {})}
                        aria-current={active ? "page" : undefined}
                        className={`block rounded-md px-2 py-1.5 text-sm no-underline ${
                          active
                            ? "bg-surface-elevated !text-foreground"
                            : "!text-neutral-300 hover:!text-foreground"
                        }`}
                      >
                        {item.label}
                        {external && " ↗"}
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </div>
          ))}
        </nav>
      </Drawer>
    </Header>
  );
}
