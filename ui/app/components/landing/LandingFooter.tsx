"use client";

import Link from "next/link";

import { SWAGGER_URL } from "@/app/lib/constants/site";

const LINKS = [
  { href: "/explore", label: "Explore" },
  { href: "/bulk", label: "Bulk" },
  { href: "/catalogs", label: "Catalogs" },
  { href: "/developers", label: "API" },
  { href: "/learn", label: "Learn" },
  { href: "/about", label: "About" },
];

const LINK_CLASS = "text-sm !text-neutral-400 hover:!text-foreground";

/** Compact footer pinned under the landing hero; same link style as SiteFooter. */
export function LandingFooter() {
  return (
    <footer className="fixed bottom-0 left-0 right-0 bg-surface border-t border-border">
      <nav
        aria-label="Footer"
        className="max-w-6xl mx-auto px-4 md:px-6 py-3 flex flex-wrap justify-center gap-x-6 gap-y-1"
      >
        {LINKS.map((link) => (
          <Link key={link.href} href={link.href} className={LINK_CLASS}>
            {link.label}
          </Link>
        ))}
        <a
          href={SWAGGER_URL}
          target="_blank"
          rel="noopener noreferrer"
          className={LINK_CLASS}
        >
          API reference ↗
        </a>
      </nav>
    </footer>
  );
}
