"use client";

import Link from "next/link";

import { dataVersionLine } from "@/app/lib/constants/catalogMeta";
import { FOOTER_GROUPS, isExternalHref } from "@/app/lib/constants/site";

export function SiteFooter() {
  return (
    <footer className="border-t border-border bg-surface mt-16">
      <div className="max-w-6xl mx-auto px-4 md:px-6 py-10 grid grid-cols-2 md:grid-cols-4 gap-8">
        {FOOTER_GROUPS.map((group) => (
          <div key={group.title}>
            <div className="text-xs uppercase tracking-wider text-neutral-500 mb-3">
              {group.title}
            </div>
            <ul className="space-y-2 list-none p-0 m-0">
              {group.items.map((item) =>
                isExternalHref(item.href) ? (
                  <li key={item.href}>
                    <a
                      href={item.href}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-sm !text-neutral-400 hover:!text-foreground"
                    >
                      {item.label} ↗
                    </a>
                  </li>
                ) : (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      className="text-sm !text-neutral-400 hover:!text-foreground"
                    >
                      {item.label}
                    </Link>
                  </li>
                )
              )}
            </ul>
          </div>
        ))}
      </div>
      <div className="border-t border-border">
        <div className="max-w-6xl mx-auto px-4 md:px-6 py-4 flex flex-wrap gap-x-6 gap-y-1 justify-between text-xs text-neutral-500">
          <span>Data: {dataVersionLine()}</span>
          <span>Universidad Diego Portales · Apache-2.0</span>
        </div>
      </div>
    </footer>
  );
}
