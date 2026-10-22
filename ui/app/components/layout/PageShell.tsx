"use client";

import { Layout, Typography } from "antd";
import type { ReactNode } from "react";

import { AppHeader } from "./AppHeader";
import { SiteFooter } from "./SiteFooter";

const { Content } = Layout;
const { Title, Paragraph } = Typography;

interface PageShellProps {
  title?: ReactNode;
  /** One-line lede under the title. */
  description?: ReactNode;
  /** Right-aligned actions next to the title. */
  actions?: ReactNode;
  /** Max content width; defaults to a reading width for docs pages. */
  width?: "narrow" | "wide";
  children: ReactNode;
}

/** Standard frame for content pages: header, titled container, site footer. */
export function PageShell({
  title,
  description,
  actions,
  width = "narrow",
  children,
}: PageShellProps) {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background">
        <div
          className={`mx-auto px-4 md:px-6 pt-10 ${
            width === "wide" ? "max-w-6xl" : "max-w-4xl"
          }`}
        >
          {(title || actions) && (
            <div className="flex flex-wrap items-end justify-between gap-4 mb-8">
              <div className="min-w-0">
                {title && (
                  <Title level={2} className="!mb-2 text-foreground">
                    {title}
                  </Title>
                )}
                {description && (
                  <Paragraph className="!mb-0 text-neutral-400 text-base max-w-2xl">
                    {description}
                  </Paragraph>
                )}
              </div>
              {actions && <div className="flex gap-2">{actions}</div>}
            </div>
          )}
          {children}
        </div>
        <SiteFooter />
      </Content>
    </Layout>
  );
}
