"use client";

import {
  CodeOutlined,
  DownloadOutlined,
  FileExcelOutlined,
  FileTextOutlined,
} from "@ant-design/icons";
import type { MenuProps } from "antd";
import { Button, Dropdown, Flex, Space, Tag, Typography } from "antd";

import { useSearchParams } from "@/app/hooks/useSearchParamsSync";
import { useCrossmatchState } from "@/app/store/crossmatch-context";

import { EmptyState } from "./EmptyState";
import { LoadingState } from "./LoadingState";
import { ResultsTable, sampleData } from "./ResultsTable";

const { Title, Text } = Typography;

// Sample data stats - in real implementation, this would come from the API
const SAMPLE_MATCH_COUNT = 12;
const SAMPLE_CATALOG_COUNT = 4;

function ResultsHeader() {
  const { ra, dec, radius, unit } = useSearchParams();

  const handleExport = (_format: string) => {
    // TODO: Implement export functionality
  };

  const exportMenuItems: MenuProps["items"] = [
    {
      key: "csv",
      label: "CSV",
      icon: <FileTextOutlined />,
      onClick: () => handleExport("csv"),
    },
    {
      key: "votable",
      label: "VOTable",
      icon: <FileExcelOutlined />,
      onClick: () => handleExport("votable"),
    },
    {
      key: "json",
      label: "JSON",
      icon: <CodeOutlined />,
      onClick: () => handleExport("json"),
    },
  ];

  const searchContext =
    ra && dec ? (
      <Flex align="center" gap={8} wrap="wrap">
        <Text type="secondary" className="text-[13px]">
          Search:
        </Text>
        <Tag className="m-0 font-mono text-xs">RA {ra}</Tag>
        <Tag className="m-0 font-mono text-xs">Dec {dec}</Tag>
        <Tag className="m-0 font-mono text-xs">
          r = {radius} {unit}
        </Tag>
      </Flex>
    ) : null;

  return (
    <Flex
      justify="space-between"
      align="flex-start"
      className="p-5 px-6 border-b border-border bg-surface"
    >
      <Flex vertical gap={8}>
        <Flex align="baseline" gap={12}>
          <Title level={4} className="!m-0">
            Crossmatch Results
          </Title>
          <Text type="secondary" className="text-sm">
            {SAMPLE_MATCH_COUNT} matches across {SAMPLE_CATALOG_COUNT} catalogs
          </Text>
        </Flex>
        {searchContext}
      </Flex>

      <Space>
        <Dropdown menu={{ items: exportMenuItems }} placement="bottomRight">
          <Button icon={<DownloadOutlined />}>Export</Button>
        </Dropdown>
      </Space>
    </Flex>
  );
}

function ResultsContent() {
  // Sample data for now - will be replaced with actual results
  const data = sampleData;

  return (
    <div className="p-6">
      {/* <SkyPlot data={data} /> */}
      <ResultsTable data={data} />
    </div>
  );
}

export function ResultsPanel() {
  const { state } = useCrossmatchState();

  switch (state.resultsState) {
    case "loading":
      return <LoadingState />;
    case "success":
      return (
        <Flex vertical className="h-full">
          <ResultsHeader />
          <ResultsContent />
        </Flex>
      );
    case "error":
      return <EmptyState />;
    case "empty":
    default:
      return <EmptyState />;
  }
}
