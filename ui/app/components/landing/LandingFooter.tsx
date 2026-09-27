"use client";

import { BookOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { Divider, Space, Typography } from "antd";
import { useRouter } from "next/navigation";
import type { ReactNode } from "react";

import { SWAGGER_URL } from "@/app/lib/constants/site";

const { Link } = Typography;

interface FooterLink {
  href: string;
  label: string;
  icon?: ReactNode;
}

const LINKS: FooterLink[] = [
  { href: "/explore", label: "Explore" },
  { href: "/bulk", label: "Bulk" },
  { href: "/catalogs", label: "Catalogs" },
  { href: "/developers", label: "API" },
  { href: "/learn", label: "Learn" },
  { href: "/about", label: "About", icon: <InfoCircleOutlined /> },
];

export function LandingFooter() {
  const router = useRouter();

  return (
    <div className="fixed bottom-0 left-0 right-0 py-3 px-4 md:px-6 bg-surface border-t border-border">
      <Space
        size={[4, 2]}
        separator={<Divider orientation="vertical" />}
        className="flex justify-center flex-wrap text-sm"
        wrap
      >
        {/*
          antd's Typography.Link renders its own anchor, so wrapping it in a
          next/link would nest anchors. Keep the real href for middle-click and
          crawlers, and route client-side on a plain click.
        */}
        {LINKS.map((link) => (
          <Link
            key={link.href}
            href={link.href}
            onClick={(e) => {
              if (e.metaKey || e.ctrlKey || e.shiftKey) return;
              e.preventDefault();
              router.push(link.href);
            }}
          >
            {link.icon ? (
              <Space size={4}>
                {link.icon}
                {link.label}
              </Space>
            ) : (
              link.label
            )}
          </Link>
        ))}
        <Link href={SWAGGER_URL} target="_blank" rel="noopener noreferrer">
          <Space size={4}>
            <BookOutlined />
            Documentation
          </Space>
        </Link>
      </Space>
    </div>
  );
}
