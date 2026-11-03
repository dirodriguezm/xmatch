"use client";

import { Typography } from "antd";
import type { ReactNode } from "react";

import { EmptyState, ErrorState, LoadingState } from "@/app/components/results";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

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
}

/** The page's loading / error / empty / no-match states around the results. */
export function MatchesBody({
  target,
  matchCount,
  errorMessage,
  onRetry,
  searchForm,
  children,
}: MatchesBodyProps) {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "error":
      return <ErrorState message={errorMessage} onRetry={onRetry} />;
    case "success":
      if (target && matchCount > 0) return <>{children}</>;
      return (
        <div className="flex flex-1 flex-col items-center justify-center gap-2 p-8 text-center">
          <Text strong className="text-base">
            No sources within the searched radii
          </Text>
          <Text type="secondary">
            Widen a catalog&apos;s radius with Edit, or try another position.
          </Text>
        </div>
      );
    case "empty":
    default:
      if (searchForm === null) return <EmptyState />;
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
