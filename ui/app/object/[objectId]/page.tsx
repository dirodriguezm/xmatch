"use client";

import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Empty, Layout } from "antd";
import Link from "next/link";
import { use } from "react";

import { AppHeader } from "@/app/components/layout";
import { ObjectDetail } from "@/app/components/object";
import { sampleData } from "@/app/components/results/ResultsTable";

const { Content } = Layout;

interface ObjectPageProps {
  params: Promise<{ objectId: string }>;
}

export default function ObjectPage({ params }: ObjectPageProps) {
  const { objectId } = use(params);
  const decodedObjectId = decodeURIComponent(objectId);

  // Find the object in sample data (in real implementation, this would fetch from API)
  const object = sampleData.find((item) => item.objectId === decodedObjectId);

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

        {object ? (
          <ObjectDetail object={object} />
        ) : (
          <div className="flex items-center justify-center h-[calc(100vh-128px)]">
            <Empty description={`Object "${decodedObjectId}" not found`} />
          </div>
        )}
      </Content>
    </Layout>
  );
}
