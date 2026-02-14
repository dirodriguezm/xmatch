"use client";

import { DownloadOutlined, LineChartOutlined } from "@ant-design/icons";
import { Button, Card, Empty, Flex, Space, Tooltip, Typography } from "antd";

import type { components } from "@/types/xwave-api";

import { LightCurveChart } from "./LightCurveChart";

const { Text, Title } = Typography;

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

interface SurveyConfig {
  id: string;
  name: string;
  colorClass: string;
}

const surveys: SurveyConfig[] = [
  { id: "ztf", name: "ZTF", colorClass: "bg-green-500" },
  { id: "neowise", name: "NEOWISE", colorClass: "bg-orange-500" },
  { id: "swift", name: "Swift", colorClass: "bg-blue-500" },
  { id: "vlass", name: "VLASS", colorClass: "bg-purple-500" },
];

interface MultiSurveyLightCurvesProps {
  data?: Record<string, Lightcurve | null>;
  loading?: boolean;
}

function SurveyLightCurveRow({
  survey,
  data,
  loading,
  isLast,
}: {
  survey: SurveyConfig;
  data?: Lightcurve | null;
  loading?: boolean;
  isLast?: boolean;
}) {
  const hasData = data && data.detections && data.detections.length > 0;

  return (
    <div
      className={`py-2 min-h-[120px] ${!isLast ? "border-b border-border" : ""}`}
    >
      <Flex justify="space-between" align="center" className="mb-1 px-2">
        <Space>
          <div className={`w-3 h-3 rounded-full ${survey.colorClass}`} />
          <Text strong>{survey.name}</Text>
        </Space>
        <Tooltip title={hasData ? "Download data" : "No data available"}>
          <Button
            type="text"
            size="small"
            icon={<DownloadOutlined />}
            disabled={!hasData}
          />
        </Tooltip>
      </Flex>
      {hasData ? (
        <div className="h-24">
          <LightCurveChart data={data} loading={loading} />
        </div>
      ) : (
        <div className="h-24 flex items-center justify-center">
          <Empty
            image={<LineChartOutlined className="text-2xl text-border" />}
            description={
              <Text type="secondary" className="text-xs">
                No data
              </Text>
            }
            styles={{ image: { height: 24 } }}
          />
        </div>
      )}
    </div>
  );
}

export function MultiSurveyLightCurves({
  data = {},
  loading,
}: MultiSurveyLightCurvesProps) {
  const hasAnyData = Object.values(data).some(
    (d) => d && d.detections && d.detections.length > 0
  );

  return (
    <Card
      className="bg-surface h-full"
      styles={{ body: { padding: "12px 0" } }}
      title={
        <Flex justify="space-between" align="center">
          <Title level={5} className="!m-0">
            Light Curves
          </Title>
          <Tooltip
            title={hasAnyData ? "Download all surveys" : "No data to download"}
          >
            <Button
              size="small"
              icon={<DownloadOutlined />}
              disabled={!hasAnyData}
            >
              Download All
            </Button>
          </Tooltip>
        </Flex>
      }
    >
      <div className="px-2">
        {surveys.map((survey, idx) => (
          <SurveyLightCurveRow
            key={survey.id}
            survey={survey}
            data={data[survey.id]}
            loading={loading}
            isLast={idx === surveys.length - 1}
          />
        ))}
      </div>

      {/* Shared time axis label */}
      <Flex justify="center" className="mt-2 pt-2 border-t border-border">
        <Text type="secondary" className="text-xs">
          Time (MJD / Year)
        </Text>
      </Flex>
    </Card>
  );
}
