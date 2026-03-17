"use client";

import { Layout, Spin } from "antd";
import { Suspense, useEffect, useMemo } from "react";

import { AppHeader, AppSidebar } from "@/app/components/layout";
import { ResultsPanel } from "@/app/components/results";
import { SidebarSearchForm } from "@/app/components/sidebar";
import { useParallelConeSearch } from "@/app/hooks/queries";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { decodeCatalogRadii } from "@/app/lib/constants/search";
import { mapConeSearchResults } from "@/app/lib/utils/mapConeSearchResults";
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

  const allData = useMemo(
    () => queryResults.flatMap((r) => r.data ?? []),
    [queryResults]
  );

  const isLoading = queryResults.some((r) => r.isLoading);
  const isError = queryResults.some((r) => r.isError);
  const isSuccess =
    queryResults.length > 0 && queryResults.every((r) => r.isSuccess);

  const mappedResults = useMemo(
    () => (allData.length > 0 ? mapConeSearchResults(allData) : []),
    [allData]
  );

  useEffect(() => {
    if (isLoading) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
    } else if (isError) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "error" });
    } else if (isSuccess) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "success" });
    } else if (!base) {
      dispatch({ type: "SET_RESULTS_STATE", payload: "empty" });
    }
  }, [isLoading, isError, isSuccess, base, dispatch]);

  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Layout>
        <AppSidebar>
          <SidebarSearchForm />
        </AppSidebar>
        <Content className="bg-background min-h-[calc(100vh-64px)] overflow-auto">
          <ResultsPanel data={mappedResults} loading={isLoading} />
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
