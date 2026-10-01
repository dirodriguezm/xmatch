"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button, Flex } from "antd";
import { useSearchParams } from "next/navigation";
import { useMemo } from "react";

import { ResultsActions } from "@/app/components/actions/ResultsActions";
import { downloadCsv } from "@/app/lib/utils/csv";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { EmptyState } from "./EmptyState";
import { ErrorState } from "./ErrorState";
import { LoadingState } from "./LoadingState";
import {
  nearbyCsv,
  NearbySplit,
  radiiFromParam,
  toNearbySources,
} from "./nearby";
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

  // Matches are shown as a sky map of offsets linked to the table.
  const params = useSearchParams();
  const catalogRadii = params.get("catalogRadii") ?? "";
  const radii = useMemo(() => radiiFromParam(catalogRadii), [catalogRadii]);
  const sources = useMemo(
    () => (target ? toNearbySources(data, target) : []),
    [data, target]
  );
  const showNearby = target !== null && data.length > 0;

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "success":
      return (
        <Flex vertical className="h-full">
          <div className="px-4 py-4 md:px-8 md:py-6">
            {target && (
              <div className="mb-3 flex flex-wrap items-start justify-between gap-3">
                <TonightSummary ra={target.ra} dec={target.dec} />
                <ResultsActions
                  target={target}
                  extra={
                    showNearby && (
                      <Button
                        size="small"
                        icon={<DownloadOutlined />}
                        onClick={() =>
                          downloadCsv(
                            "xwave_crossmatch.csv",
                            nearbyCsv(sources)
                          )
                        }
                      >
                        Export CSV
                      </Button>
                    )
                  }
                />
              </div>
            )}
            {showNearby ? (
              <NearbySplit sources={sources} radii={radii} />
            ) : (
              <ResultsTable
                data={data}
                loading={loading}
                photometryStatus={photometryStatus}
              />
            )}
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
