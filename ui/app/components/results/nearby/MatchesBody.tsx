"use client";

import { Typography } from "antd";
import type { ReactNode } from "react";

import { EmptyState, ErrorState, LoadingState } from "@/app/components/results";
import { getSearchCatalogLabel } from "@/app/lib/constants/catalogs";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { CatalogMark } from "./CatalogMark";
import { formatArcsec } from "./shared";

const { Text } = Typography;

interface MatchesBodyProps {
  target: { ra: number; dec: number } | null;
  matchCount: number;
  errorMessage?: string;
  onRetry?: () => void;
  /** Shown when there is no search yet, beside the empty state; null when
   * the page already shows the form elsewhere (e.g. a sidebar). */
  searchForm: ReactNode | null;
  /** Rendered once the search succeeded with at least one match. */
  children: ReactNode;
  /** Catalogs searched, listed in the no-match state so it is clear which
   * came back empty and which failed. */
  searchedCatalogs?: {
    slug: string;
    radius?: number;
    error?: string;
    pending?: boolean;
  }[];
}

/** Centres a state message in the results area, where LoadingState's
 * spinner sits, so every state appears in the same place. */
function Centered({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-[calc(100dvh-64px)] items-center justify-center p-8">
      {children}
    </div>
  );
}

/** The page's loading / error / empty / no-match states around the results. */
export function MatchesBody({
  target,
  matchCount,
  errorMessage,
  onRetry,
  searchForm,
  children,
  searchedCatalogs = [],
}: MatchesBodyProps) {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "error":
      return (
        <Centered>
          {/* No retry means the request itself was rejected (a bad position,
              a radius over the limit): something to fix, not an outage. */}
          <ErrorState
            title={onRetry ? undefined : "Check your search"}
            message={errorMessage}
            onRetry={onRetry}
          />
        </Centered>
      );
    case "success": {
      if (target && matchCount > 0) return <>{children}</>;
      const anyFailed = searchedCatalogs.some((c) => c.error !== undefined);
      const anyPending = searchedCatalogs.some((c) => c.pending);
      return (
        <Centered>
          <div className="flex max-w-sm flex-col items-center gap-1 text-center">
            <Text strong className="text-base">
              {anyPending
                ? "No sources found yet"
                : anyFailed
                  ? "No sources found"
                  : "No sources within the searched radii"}
            </Text>
            <Text type="secondary">
              Widen a catalog&apos;s radius under Catalogs &amp; Radii, or try
              another position.
            </Text>
            {searchedCatalogs.length > 0 && (
              // Name each catalog so "nothing found" is never mistaken for
              // "not searched" or for a search that quietly failed.
              <ul className="m-0 mt-3 flex list-none flex-col gap-1 p-0 text-xs text-neutral-400">
                {searchedCatalogs.map((c) => (
                  <li
                    key={c.slug}
                    className="inline-flex items-center justify-center gap-2"
                  >
                    <span className="inline-flex opacity-60">
                      <CatalogMark slug={c.slug} size={9} />
                    </span>
                    <span className="text-neutral-300">
                      {getSearchCatalogLabel(c.slug)}
                    </span>
                    {c.pending ? (
                      <span className="animate-pulse motion-reduce:animate-none">
                        searching…
                      </span>
                    ) : c.error !== undefined ? (
                      <span className="text-amber-400">search failed</span>
                    ) : (
                      <span>
                        0 within{" "}
                        {c.radius !== undefined ? formatArcsec(c.radius) : ""}
                      </span>
                    )}
                  </li>
                ))}
              </ul>
            )}
            {/* Re-submitting the form would not retry: the URL is unchanged. */}
            {anyFailed && onRetry && (
              <button
                type="button"
                onClick={onRetry}
                className="mt-2 cursor-pointer rounded px-1 py-0.5 text-xs text-primary hover:underline focus-visible:ring-1 focus-visible:ring-primary focus-visible:outline-none"
              >
                Retry failed search
              </button>
            )}
          </div>
        </Centered>
      );
    }
    case "empty":
    default:
      if (searchForm === null)
        return (
          <Centered>
            <EmptyState />
          </Centered>
        );
      return (
        <div className="flex flex-1 flex-col items-center justify-center gap-6 p-6 lg:flex-row lg:items-start lg:pt-16">
          <div className="w-full max-w-sm overflow-hidden rounded-lg border border-border bg-surface">
            {searchForm}
          </div>
          <div className="hidden lg:block lg:pt-24">
            <EmptyState />
          </div>
        </div>
      );
  }
}
