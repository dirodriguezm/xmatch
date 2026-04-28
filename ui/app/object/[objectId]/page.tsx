"use client";

import { Empty, Layout, Spin } from "antd";
import { useRouter, useSearchParams } from "next/navigation";
import { use, useEffect, useMemo, useRef } from "react";

import { AppHeader } from "@/app/components/layout";
import { ObjectDetail } from "@/app/components/object";
import type { CrossmatchResult } from "@/app/components/results/ResultsTable";
import { useMetadata } from "@/app/hooks/queries";

const { Content } = Layout;

interface ObjectPageProps {
  params: Promise<{ objectId: string }>;
}

export default function ObjectPage({ params }: ObjectPageProps) {
  const { objectId } = use(params);
  const decodedObjectId = decodeURIComponent(objectId);
  const router = useRouter();
  const titleSet = useRef(false);

  useEffect(() => {
    const title = `XWave | ${decodedObjectId}`;
    document.title = title;
    titleSet.current = true;

    // Keep checking and restoring the title if something changes it
    const interval = setInterval(() => {
      if (document.title !== title) {
        document.title = title;
      }
    }, 100);

    // Clean up after a short time once the page is stable
    const cleanup = setTimeout(() => clearInterval(interval), 2000);

    return () => {
      clearInterval(interval);
      clearTimeout(cleanup);
    };
  }, [decodedObjectId]);

  const searchParams = useSearchParams();
  const catalog = searchParams.get("catalog");

  // Fetch metadata if we have a catalog
  const { data: metadata, isLoading } = useMetadata(
    catalog ? { id: decodedObjectId, catalog } : null
  );

  // Construct the object from metadata response (ra/dec come from the API
  // payload but aren't part of the per-catalog typed schema, hence the cast).
  const object = useMemo<CrossmatchResult | null>(() => {
    const meta = metadata as
      | ({ ra?: number; dec?: number } & typeof metadata)
      | undefined;
    if (!catalog || meta?.ra == null || meta?.dec == null) return null;

    return {
      key: decodedObjectId,
      objectId: decodedObjectId,
      ra: meta.ra,
      dec: meta.dec,
      catalog: catalog,
      angularDistance: 0, // Not needed for detail view, but required by type
    };
  }, [decodedObjectId, catalog, metadata]);

  return (
    <Layout className="min-h-screen">
      <AppHeader onBack={() => router.back()} />
      <Content className="bg-background min-h-[calc(100vh-64px)]">
        {isLoading ? (
          <div className="flex items-center justify-center h-[calc(100vh-128px)]">
            <Spin size="large" />
          </div>
        ) : object ? (
          <ObjectDetail object={object} metadata={metadata} />
        ) : (
          <div className="flex items-center justify-center h-[calc(100vh-128px)]">
            <Empty
              description={`Object "${decodedObjectId}" not found or missing parameters`}
            />
          </div>
        )}
      </Content>
    </Layout>
  );
}
