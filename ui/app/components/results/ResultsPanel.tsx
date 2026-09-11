"use client";

import { Flex } from "antd";

import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { EmptyState } from "./EmptyState";
import { ErrorState } from "./ErrorState";
import { LoadingState } from "./LoadingState";
import type { EnrichedResult, PhotometryStatus } from "./ResultsTable";
import { ResultsTable } from "./ResultsTable";

export interface ResultsPanelProps {
  data?: EnrichedResult[];
  loading?: boolean;
  /** Message from the failed cone search, shown in the error state. */
  errorMessage?: string;
  onRetry?: () => void;
  photometryStatus?: PhotometryStatus;
}

export function ResultsPanel({
  data = [],
  loading = false,
  errorMessage,
  onRetry,
  photometryStatus = "idle",
}: ResultsPanelProps) {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "success":
      return (
        <Flex vertical className="h-full">
          <div className="px-8 py-6">
            <ResultsTable
              data={data}
              loading={loading}
              photometryStatus={photometryStatus}
            />
          </div>
        </Flex>
      );
    case "error":
      return <ErrorState message={errorMessage} onRetry={onRetry} />;
    case "empty":
    default:
      return <EmptyState />;
  }
}
