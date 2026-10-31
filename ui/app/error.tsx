"use client";

import { Button } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { buildIssueUrl } from "@/app/components/feedback/githubIssue";
import { PageShell } from "@/app/components/layout";

interface ErrorPageProps {
  error: Error & { digest?: string };
  reset: () => void;
}

export default function ErrorPage({ error, reset }: ErrorPageProps) {
  const pathname = usePathname();
  // Server errors arrive with a generic message and a digest that matches
  // the server log; client errors keep their own message.
  const detail = error.digest ? `digest ${error.digest}` : error.message;
  const reportHref = buildIssueUrl({
    title: `Error on ${pathname}`,
    body: `**Page:** ${pathname}\n**Error:** ${detail}\n\nWhat I was doing:\n`,
    labels: ["bug"],
  });

  return (
    <PageShell>
      <div className="max-w-xl pb-24">
        <h1 className="mb-3 text-2xl font-semibold text-foreground">
          This page stopped working
        </h1>
        <p className="mb-4 text-base text-neutral-400">
          Retrying usually fixes it when the API was briefly out of reach. If it
          keeps failing, a report with the line below helps track it down.
        </p>
        {detail && (
          <pre className="mb-8 overflow-x-auto rounded border border-border bg-surface px-3 py-2 font-mono text-sm text-neutral-400">
            {detail}
          </pre>
        )}
        <div className="flex flex-wrap items-center gap-x-5 gap-y-3">
          <Button type="primary" onClick={reset}>
            Retry
          </Button>
          <a
            href={reportHref}
            target="_blank"
            rel="noopener noreferrer"
            className="text-primary-text hover:underline"
          >
            Report it on GitHub ↗
          </a>
          <Link href="/" className="text-neutral-400 hover:text-foreground">
            Home
          </Link>
        </div>
      </div>
    </PageShell>
  );
}
