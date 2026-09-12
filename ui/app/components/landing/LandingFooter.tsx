"use client";

import { BookOutlined, InfoCircleOutlined } from "@ant-design/icons";
import { Divider, Space, Typography } from "antd";
import { useRouter } from "next/navigation";

const { Link } = Typography;

export function LandingFooter() {
  const router = useRouter();

  return (
    <div className="fixed bottom-0 left-0 right-0 py-4 px-6 bg-surface border-t border-border">
      <Space
        separator={<Divider orientation="vertical" />}
        className="flex justify-center flex-wrap"
      >
        {/*
          antd's Typography.Link renders its own anchor, so wrapping it in a
          next/link would nest anchors. Keep the real href for middle-click and
          crawlers, and route client-side on a plain click.
        */}
        <Link
          href="/about"
          onClick={(e) => {
            if (e.metaKey || e.ctrlKey || e.shiftKey) return;
            e.preventDefault();
            router.push("/about");
          }}
        >
          <Space size={4}>
            <InfoCircleOutlined />
            About
          </Space>
        </Link>
        <Link
          href="https://xwave-astro.udp.cl/swagger/index.html"
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
