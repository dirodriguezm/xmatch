"use client";

import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Empty, Layout, Spin } from "antd";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
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
  const raStr = searchParams.get("ra");
  const decStr = searchParams.get("dec");

  // Fetch metadata if we have a catalog
  const { data: metadata, isLoading } = useMetadata(
    catalog ? { id: decodedObjectId, catalog } : null
  );

  // Construct the object from URL params and metadata
  const object = useMemo<CrossmatchResult | null>(() => {
    if (!catalog || !raStr || !decStr) return null;

    return {
      key: decodedObjectId,
      objectId: decodedObjectId,
      ra: parseFloat(raStr),
      dec: parseFloat(decStr),
      catalog: catalog,
      angularDistance: 0, // Not needed for detail view, but required by type
    };
  }, [decodedObjectId, catalog, raStr, decStr]);

  return (
    <Layout className="min-h-screen">
      <AppHeader />
      <Content className="bg-background min-h-[calc(100vh-64px)]">
        <div className="p-4 border-b border-border bg-surface">
          <Link href="/search">
            <Button icon={<ArrowLeftOutlined />} type="text">
              Back to Results
            </Button>
          </Link>
        </div>

        {isLoading ? (
          <div className="flex items-center justify-center h-[calc(100vh-128px)]">
            <Spin size="large" />
          </div>
        ) : object ? (
          <ObjectDetail object={object} />
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
