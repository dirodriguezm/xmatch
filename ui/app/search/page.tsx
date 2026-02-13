"use client";

import { Layout, Spin } from "antd";
import { Suspense, useEffect } from "react";

import { AppHeader, AppSidebar } from "@/app/components/layout";
import { ResultsPanel } from "@/app/components/results";
import { SidebarSearchForm } from "@/app/components/sidebar";
import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import {
  CrossmatchProvider,
  useCrossmatchState,
} from "@/app/store/crossmatch-context";

const { Content } = Layout;

function SearchContent() {
  const { state, dispatch } = useCrossmatchState();
  const { ra, dec, radius, unit, catalogs, isValid } = useSearchParams();

  // Trigger search when page loads with valid params
  useEffect(() => {
    if (isValid && state.resultsState === "empty") {
      dispatch({ type: "SET_RESULTS_STATE", payload: "loading" });
      // Simulate API call - in real implementation, this would use the actual crossmatch API
      setTimeout(() => {
        dispatch({ type: "SET_RESULTS_STATE", payload: "success" });
      }, 1500);
    }
  }, [isValid, ra, dec, radius, unit, catalogs, dispatch, state.resultsState]);

  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Layout>
        <AppSidebar>
          <SidebarSearchForm />
        </AppSidebar>
        <Content className="bg-background min-h-[calc(100vh-64px)] overflow-auto">
          <ResultsPanel />
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
