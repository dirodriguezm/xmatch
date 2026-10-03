"use client";

import { Typography } from "antd";
import type { ReactNode } from "react";

const { Title } = Typography;

/** Titled page section, matching the About page's visual language. */
export function Section({
  id,
  title,
  icon,
  children,
}: {
  id?: string;
  title: string;
  icon?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section id={id} className="mb-10 scroll-mt-24">
      <Title
        level={3}
        className="!mb-4 text-foreground flex items-center gap-2"
      >
        {icon}
        {title}
      </Title>
      {children}
    </section>
  );
}

/** Numbered step with a circled index, as used on the About page. */
export function Step({
  index,
  title,
  children,
}: {
  index: number;
  title: string;
  children: ReactNode;
}) {
  return (
    <div className="flex gap-4">
      <div className="shrink-0 w-8 h-8 rounded-full bg-surface-elevated border border-border flex items-center justify-center font-mono text-sm text-foreground">
        {index}
      </div>
      <div className="min-w-0 flex-1 pb-6">
        <Title level={5} className="!mt-0 !mb-2 text-foreground">
          {title}
        </Title>
        <div className="text-neutral-400">{children}</div>
      </div>
    </div>
  );
}
