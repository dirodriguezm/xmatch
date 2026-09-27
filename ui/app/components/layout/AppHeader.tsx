"use client";

import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Layout, Space, Tag, Typography } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { HeaderBasket } from "@/app/components/basket/HeaderBasket";
import { CommandPaletteTrigger } from "@/app/components/command/CommandPalette";
import { Logo } from "@/app/components/common";
import { NAV_ITEMS } from "@/app/lib/constants/site";

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
      </div>
    </Header>
  );
}
