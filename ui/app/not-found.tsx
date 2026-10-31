"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { PageShell } from "@/app/components/layout";

const PLACES = [
  { href: "/search", label: "Search", body: "a position or object name" },
  { href: "/bulk", label: "Bulk", body: "cross-match a list of positions" },
  { href: "/catalogs", label: "Catalogs", body: "what each survey covers" },
  { href: "/", label: "Home", body: "start over" },
];

export default function NotFound() {
  const pathname = usePathname();

  return (
    <PageShell>
      <div className="max-w-xl pb-24">
        <p className="mb-3 font-mono text-sm text-neutral-500">
          404 <span className="text-neutral-600">·</span>{" "}
          <span className="break-all">{pathname}</span>
        </p>
        <h1 className="mb-3 text-2xl font-semibold text-foreground">
          Nothing at this address
        </h1>
        <p className="mb-8 text-base text-neutral-400">
          The link may be old or mistyped. Object pages need the full catalog
          ID, so a search by position usually finds what you were after.
        </p>
        <ul className="divide-y divide-border border-y border-border">
          {PLACES.map((place) => (
            <li key={place.href}>
              <Link
                href={place.href}
                className="flex items-baseline gap-3 py-3 text-foreground hover:text-primary-text"
              >
                <span className="w-20 shrink-0 font-medium">{place.label}</span>
                <span className="text-sm text-neutral-500">{place.body}</span>
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </PageShell>
  );
}
