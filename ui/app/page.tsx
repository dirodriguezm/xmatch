"use client";

import { Layout } from "antd";
import { Suspense } from "react";

import { XWaveSpinner } from "@/app/components/common/XWaveSpinner";
import { LandingFooter, LandingSearch } from "@/app/components/landing";
import { AppHeader } from "@/app/components/layout";
import { CrossmatchProvider } from "@/app/store/crossmatch-context";

const { Content } = Layout;

function PageContent() {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background min-h-[calc(100vh-64px)]">
        <LandingSearch />
        <LandingFooter />
      </Content>
    </Layout>
  );
}

function LoadingFallback() {
  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background min-h-[calc(100vh-64px)] flex items-center justify-center">
        <XWaveSpinner variant="xwave" size={64} />
      </Content>
    </Layout>
  );
}

export default function Home() {
  return (
    <CrossmatchProvider>
      <Suspense fallback={<LoadingFallback />}>
        <PageContent />
      </Suspense>
    </CrossmatchProvider>
  );
}
