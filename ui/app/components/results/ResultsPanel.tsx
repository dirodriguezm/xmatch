"use client";

import { Flex } from "antd";

import { ResultsActions } from "@/app/components/actions/ResultsActions";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { EmptyState } from "./EmptyState";
import { ErrorState } from "./ErrorState";
import { LoadingState } from "./LoadingState";
import type { EnrichedResult, PhotometryStatus } from "./ResultsTable";
import { ResultsTable } from "./ResultsTable";
import { TonightSummary } from "./TonightSummary";

export interface ResultsPanelProps {
  data?: EnrichedResult[];
  loading?: boolean;
  /** Message from the failed cone search, shown in the error state. */
  errorMessage?: string;
  onRetry?: () => void;
  photometryStatus?: PhotometryStatus;
  /** Searched position; every match lies within the radius of it. */
  target?: { ra: number; dec: number } | null;
}

export function ResultsPanel({
  data = [],
  loading = false,
  errorMessage,
  onRetry,
  photometryStatus = "idle",
  target = null,
}: ResultsPanelProps) {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "success":
      return (
        <Flex vertical className="h-full">
          <div className="px-8 py-6">
            {target && (
              <div className="mb-3 flex flex-wrap items-start justify-between gap-3">
                <TonightSummary ra={target.ra} dec={target.dec} />
                <ResultsActions
                  target={target}
                  matchedCatalogs={[
                    ...new Set(
                      data.map((r) =>
                        (r.catalogSlug ?? r.catalog).toLowerCase()
                      )
                    ),
                  ]}
                />
              </div>
            )}
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
