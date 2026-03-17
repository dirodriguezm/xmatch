"use client";

import { Flex } from "antd";

import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { EmptyState } from "./EmptyState";
import { LoadingState } from "./LoadingState";
import type { CrossmatchResult } from "./ResultsTable";
import { ResultsTable } from "./ResultsTable";

export interface ResultsPanelProps {
  data?: CrossmatchResult[];
  loading?: boolean;
}

export function ResultsPanel({
  data = [],
  loading = false,
}: ResultsPanelProps) {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "success":
      return (
        <Flex vertical className="h-full">
          <div className="px-8 py-6">
            <ResultsTable data={data} loading={loading} />
          </div>
        </Flex>
      );
    case "error":
      return <EmptyState />;
    case "empty":
    default:
      return <EmptyState />;
  }
}
