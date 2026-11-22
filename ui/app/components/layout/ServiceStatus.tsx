"use client";

import { Tooltip } from "antd";
import Link from "next/link";

import { useFailingServices } from "@/app/lib/serviceHealth";

/**
 * Header hint that something the page relies on is failing, linking to the
 * Status page. Shows nothing while every call in this session succeeds.
 */
export function ServiceStatusIndicator() {
  const failing = useFailingServices();
  if (failing.length === 0) return null;

  return (
    <Tooltip
      title={`Not responding: ${failing.join(", ")}. Open the Status page for live checks.`}
    >
      <Link
        href="/status"
        aria-label={`Service issues: ${failing.join(", ")}`}
        className="inline-flex items-center gap-2 rounded-md px-2 leading-8 text-xs text-amber-400 no-underline hover:bg-foreground/5"
      >
        <span className="h-2 w-2 rounded-full bg-amber-400" aria-hidden />
        <span className="hidden lg:inline">Service issues</span>
      </Link>
    </Tooltip>
  );
}

/** Inline link for error messages caused by an external service. */
export function ServiceStatusLink({ className = "" }: { className?: string }) {
  return (
    <Link href="/status" className={`text-xs ${className}`}>
      Service status
    </Link>
  );
}
