"use client";

import { Layout, Spin } from "antd";
import { Suspense, useCallback, useEffect, useMemo } from "react";

import { AppHeader, AppSidebar } from "@/app/components/layout";
import { ResultsPanel } from "@/app/components/results";
import type { PhotometryStatus } from "@/app/components/results/ResultsTable";
import { SidebarSearchForm } from "@/app/components/sidebar";
import { useParallelConeSearch } from "@/app/hooks/queries";
import { useBulkMetadata } from "@/app/hooks/queries/useBulkMetadata";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { decodeCatalogRadii } from "@/app/lib/constants/search";
import { mapConeSearchResults } from "@/app/lib/utils/mapConeSearchResults";
import {
  buildBulkMetadataRequests,
  buildMagIndex,
  enrichWithPhotometry,
} from "@/app/lib/utils/mergePhotometry";
import {
  CrossmatchProvider,
  useCrossmatchState,
} from "@/app/store/crossmatch-context";

const { Content } = Layout;

function SearchContent() {
  const { dispatch } = useCrossmatchState();
  const { ra, dec, catalogRadii: catalogRadiiStr, isValid } = useSearchParams();

  const raNum = ra ? parseFloat(ra) : null;
  const decNum = dec ? parseFloat(dec) : null;

  const base = useMemo(() => {
    if (
      !isValid ||
      raNum === null ||
      decNum === null ||
      isNaN(raNum) ||
      isNaN(decNum)
    ) {
      return null;
    }
    return { ra: raNum, dec: decNum };
  }, [isValid, raNum, decNum]);

  const catalogConfigs = useMemo(
    () => decodeCatalogRadii(catalogRadiiStr),
    [catalogRadiiStr]
  );

  const queryResults = useParallelConeSearch(base, catalogConfigs);

  // queryResults[i] pairs with enabledCatalogs[i]: useParallelConeSearch builds
  // its queries from this same `enabled` filter, in this same order.
  const enabledCatalogs = useMemo(
    () => catalogConfigs.filter((c) => c.enabled).map((c) => c.catalog),
    [catalogConfigs]
  );

  const isLoading = queryResults.some((r) => r.isLoading);
  const isError = queryResults.some((r) => r.isError);
  const isSuccess =
    queryResults.length > 0 && queryResults.every((r) => r.isSuccess);

  const errorMessage = queryResults.find((r) => r.isError)?.error?.message;

  const mappedResults = useMemo(
    () =>
      queryResults.flatMap((r, i) =>
        r.data ? mapConeSearchResults(r.data, enabledCatalogs[i]) : []
      ),
    [queryResults, enabledCatalogs]
  );

  // Photometry enrichment. Deliberately kept out of the effect below: the table
  // renders from mappedResults immediately and the Mag column fills in later,
  // so a slow or failed metadata request must never blank the results.
  const bulkRequests = useMemo(
    () => buildBulkMetadataRequests(mappedResults),
    [mappedResults]
  );
  const {
    groups: metadataGroups,
    isFetching: photometryFetching,
    isError: photometryError,
  } = useBulkMetadata(bulkRequests);

  const enrichedResults = useMemo(
    () => enrichWithPhotometry(mappedResults, buildMagIndex(metadataGroups)),
    [mappedResults, metadataGroups]
  );

  const photometryStatus: PhotometryStatus =
    bulkRequests.length === 0
      ? "idle"
      : photometryFetching
        ? "pending"
        : photometryError
          ? "error"
          : "ready";

  const handleRetry = useCallback(() => {
    queryResults.forEach((r) => r.refetch());
  }, [queryResults]);

  // NOTE: photometry state must never enter this effect or its dependencies —
  // ResultsPanel switches on resultsState to choose between the table and the
  // empty/error states.
  useEffect(() => {
    if (isLoading) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    } else if (isError) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "error" });
    } else if (isSuccess) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "success" });
    } else if (!base || queryResults.length === 0) {
      // Also covers a valid target with every catalog unchecked, which
      // previously dispatched nothing and left the panel on its last state.
      dispatch({ type: "SET_RESULTS_STATE", payload: "empty" });
    }
  }, [isLoading, isError, isSuccess, base, queryResults.length, dispatch]);

  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Layout>
        <AppSidebar>
          <SidebarSearchForm />
        </AppSidebar>
        <Content className="bg-background min-h-[calc(100vh-64px)] overflow-auto">
          <ResultsPanel
            data={enrichedResults}
            loading={isLoading}
            errorMessage={errorMessage}
            onRetry={handleRetry}
            photometryStatus={photometryStatus}
          />
        </Content>
      </Layout>
    </Layout>
  );
}

function LoadingFallback() {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background min-h-[calc(100vh-64px)] flex items-center justify-center">
        <Spin size="large" />
      </Content>
    </Layout>
  );
}

export default function SearchPage() {
  return (
    <CrossmatchProvider>
      <Suspense fallback={<LoadingFallback />}>
        <SearchContent />
      </Suspense>
    </CrossmatchProvider>
  );
}
