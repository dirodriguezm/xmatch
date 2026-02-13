"use client";

import { LineChartOutlined } from "@ant-design/icons";
import { Empty, Spin, Typography } from "antd";
import type { EChartsOption } from "echarts";
import dynamic from "next/dynamic";

import type { components } from "@/types/xwave-api";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const { Text } = Typography;

type Lightcurve = components["schemas"]["lightcurve.Lightcurve"];

interface LightCurveChartProps {
  data: Lightcurve | null | undefined;
  loading?: boolean;
  error?: Error | null;
}

// Detection point structure (assumed based on common light curve formats)
interface DetectionPoint {
  mjd?: number;
  mag?: number;
  magerr?: number;
  fid?: number;
  band?: string;
  [key: string]: unknown;
}

// Band colors for different filters
const bandColors: Record<string, string> = {
  g: "#52c41a", // green
  r: "#ff4d4f", // red
  i: "#722ed1", // purple
  z: "#fa8c16", // orange
  "1": "#52c41a", // ZTF g
  "2": "#ff4d4f", // ZTF r
  "3": "#722ed1", // ZTF i
};

function parseDetections(detections: unknown[]): DetectionPoint[] {
  return detections.map((d) => {
    if (typeof d === "object" && d !== null) {
      return d as DetectionPoint;
    }
    return {};
  });
}

function getBandName(point: DetectionPoint): string {
  if (point.band) return point.band;
  if (point.fid !== undefined) {
    const fidMap: Record<number, string> = { 1: "g", 2: "r", 3: "i" };
    return fidMap[point.fid] || `fid${point.fid}`;
  }
  return "unknown";
}

export function LightCurveChart({
  data,
  loading,
  error,
}: LightCurveChartProps) {
  if (loading) {
    return (
      <div className="h-48 flex items-center justify-center">
        <Spin spinning>{null}</Spin>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-48 flex items-center justify-center">
        <Empty
          image={<LineChartOutlined className="text-4xl text-border" />}
          description={
            <Text type="secondary">Error loading light curve data</Text>
          }
          styles={{ image: { height: 40 } }}
        />
      </div>
    );
  }

  if (!data || !data.detections || data.detections.length === 0) {
    return (
      <div className="h-48 flex items-center justify-center">
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            <Text type="secondary">No light curve data available</Text>
          }
        />
      </div>
    );
  }

  const detections = parseDetections(data.detections);

  // Group by band
  const bandData: Record<string, [number, number, number][]> = {};
  detections.forEach((d) => {
    if (d.mjd !== undefined && d.mag !== undefined) {
      const band = getBandName(d);
      if (!bandData[band]) bandData[band] = [];
      bandData[band].push([d.mjd, d.mag, d.magerr || 0]);
    }
  });

  // Sort each band by MJD
  Object.values(bandData).forEach((points) =>
    points.sort((a, b) => a[0] - b[0])
  );

  // Create series for each band
  const series: EChartsOption["series"] = Object.entries(bandData).map(
    ([band, points]) => ({
      name: band,
      type: "scatter" as const,
      data: points.map((p) => [p[0], p[1]]),
      symbolSize: 6,
      itemStyle: {
        color: bandColors[band] || "#1890ff",
      },
    })
  );

  // Calculate axis bounds
  const allMjd = detections
    .filter((d) => d.mjd !== undefined)
    .map((d) => d.mjd!);
  const allMag = detections
    .filter((d) => d.mag !== undefined)
    .map((d) => d.mag!);
  const mjdMin = Math.min(...allMjd);
  const mjdMax = Math.max(...allMjd);
  const magMin = Math.min(...allMag);
  const magMax = Math.max(...allMag);

  const mjdPadding = (mjdMax - mjdMin) * 0.05 || 1;
  const magPadding = (magMax - magMin) * 0.1 || 0.1;

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { left: 50, right: 20, top: 30, bottom: 50 },
    legend: {
      show: Object.keys(bandData).length > 1,
      top: 0,
      textStyle: { color: "#999" },
    },
    xAxis: {
      type: "value",
      name: "MJD",
      nameLocation: "middle",
      nameGap: 30,
      min: mjdMin - mjdPadding,
      max: mjdMax + mjdPadding,
      axisLabel: {
        formatter: (value: number) => value.toFixed(0),
      },
      axisLine: { lineStyle: { color: "#303030" } },
      splitLine: { lineStyle: { color: "#202020" } },
    },
    yAxis: {
      type: "value",
      name: "Magnitude",
      nameLocation: "middle",
      nameGap: 40,
      inverse: true, // Brighter = lower magnitude
      min: magMin - magPadding,
      max: magMax + magPadding,
      axisLabel: {
        formatter: (value: number) => value.toFixed(1),
      },
      axisLine: { lineStyle: { color: "#303030" } },
      splitLine: { lineStyle: { color: "#202020" } },
    },
    tooltip: {
      trigger: "item",
      formatter: (params: unknown) => {
        const p = params as { seriesName: string; value: [number, number] };
        return `<b>${p.seriesName}</b><br/>MJD: ${p.value[0].toFixed(2)}<br/>Mag: ${p.value[1].toFixed(3)}`;
      },
    },
    dataZoom: [
      { type: "inside", xAxisIndex: 0 },
      { type: "inside", yAxisIndex: 0 },
    ],
    series,
  };

  return <ReactECharts option={option} className="h-64 w-full" />;
}
