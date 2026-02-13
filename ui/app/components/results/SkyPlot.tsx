"use client";

import type { EChartsOption } from "echarts";
import dynamic from "next/dynamic";

import { CrossmatchResult } from "./ResultsTable";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

interface SkyPlotProps {
  data: CrossmatchResult[];
}

export function SkyPlot({ data }: SkyPlotProps) {
  const plotData = data.map((item) => [
    item.ra,
    item.dec,
    item.objectId,
    item.catalog,
  ]);

  // Calculate data bounds with padding
  const raValues = data.map((d) => d.ra);
  const decValues = data.map((d) => d.dec);
  const raMin = Math.min(...raValues);
  const raMax = Math.max(...raValues);
  const decMin = Math.min(...decValues);
  const decMax = Math.max(...decValues);

  // Add 10% padding around data range
  const raPadding = (raMax - raMin) * 0.1 || 0.001;
  const decPadding = (decMax - decMin) * 0.1 || 0.001;

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { left: 60, right: 20, top: 20, bottom: 50 },
    xAxis: {
      type: "value",
      name: "RA (deg)",
      nameLocation: "middle",
      nameGap: 30,
      min: raMin - raPadding,
      max: raMax + raPadding,
      axisLabel: {
        formatter: (value: number) => value.toFixed(4),
      },
    },
    yAxis: {
      type: "value",
      name: "DEC (deg)",
      nameLocation: "middle",
      nameGap: 50,
      min: decMin - decPadding,
      max: decMax + decPadding,
      axisLabel: {
        formatter: (value: number) => value.toFixed(4),
      },
    },
    tooltip: {
      trigger: "item",
      formatter: (params: unknown) => {
        const p = params as { value: [number, number, string, string] };
        return `<b>${p.value[2]}</b><br/>Catalog: ${p.value[3]}<br/>RA: ${p.value[0].toFixed(6)}°<br/>DEC: ${p.value[1].toFixed(6)}°`;
      },
    },
    dataZoom: [
      {
        type: "inside",
        xAxisIndex: 0,
        filterMode: "none",
      },
      {
        type: "inside",
        yAxisIndex: 0,
        filterMode: "none",
      },
    ],
    series: [
      {
        type: "scatter",
        data: plotData,
        symbolSize: 12,
        itemStyle: { color: "#1890ff" },
      },
    ],
  };

  return <ReactECharts option={option} className="h-[350px] w-full" />;
}
