"use client";

import { Typography } from "antd";
import type {
  CustomSeriesOption,
  CustomSeriesRenderItem,
  EChartsOption,
  ScatterSeriesOption,
} from "echarts";
import dynamic from "next/dynamic";

import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import type { SedPoint } from "@/app/lib/utils/sed";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const { Text } = Typography;

/** An SED needs at least two points to show a shape. */
export const MIN_SED_POINTS = 2;

const SYMBOL_SIZE = 9;

interface SedChartProps {
  points: SedPoint[];
}

type Datum = {
  value: [number, number];
  point: SedPoint;
  symbol: string;
  symbolRotate: number;
  itemStyle: { color: string; borderColor: string; borderWidth: number };
};

function toDatum(p: SedPoint): Datum {
  const color = getSearchCatalogColor(p.catalog);
  return {
    value: [p.wavelengthUm, p.nuFnu],
    point: p,
    // Upper limits point down; counterparts are hollow so the object's own
    // photometry reads as the primary data.
    symbol: p.upperLimit ? "triangle" : "circle",
    symbolRotate: p.upperLimit ? 180 : 0,
    itemStyle: {
      color: p.isSelf ? color : "transparent",
      borderColor: color,
      borderWidth: 1.5,
    },
  };
}

function formatPoint(p: SedPoint): string {
  const mag = p.upperLimit
    ? `> ${p.mag.toFixed(2)} (upper limit)`
    : p.magErr != null
      ? `${p.mag.toFixed(2)} ± ${p.magErr.toFixed(2)}`
      : p.mag.toFixed(2);
  const origin = p.isSelf
    ? "this object"
    : `${getSearchCatalogLabel(p.catalog)} counterpart at ${p.separationArcsec.toFixed(1)}″`;
  return [
    `<b>${p.survey} ${p.band}</b> · ${p.wavelengthUm.toFixed(2)} µm`,
    `mag (Vega): ${mag}`,
    `νFν: ${p.nuFnu.toExponential(2)} erg s⁻¹ cm⁻²`,
    `<span style="color:#999">${origin}</span>`,
  ].join("<br/>");
}

export function SedChart({ points }: SedChartProps) {
  if (points.length < MIN_SED_POINTS) {
    return (
      <Text type="secondary" className="text-xs block">
        Not enough photometry for an SED.
      </Text>
    );
  }

  const catalogs = [...new Set(points.map((p) => p.catalog))];

  const scatter: ScatterSeriesOption[] = catalogs.map((catalog) => ({
    name: getSearchCatalogLabel(catalog),
    type: "scatter",
    symbolSize: SYMBOL_SIZE,
    z: 3,
    itemStyle: { color: getSearchCatalogColor(catalog) },
    data: points.filter((p) => p.catalog === catalog).map(toDatum),
  }));

  // Vertical error bars, drawn in pixel space so they stay correct on log axes.
  const errorData = points
    .filter((p) => !p.upperLimit && p.nuFnuErr != null)
    .map((p) => ({
      value: [
        p.wavelengthUm,
        // Clamp the lower end: a log axis cannot draw ≤ 0.
        Math.max(p.nuFnu - p.nuFnuErr!, p.nuFnu / 10),
        p.nuFnu + p.nuFnuErr!,
      ],
      itemStyle: { color: getSearchCatalogColor(p.catalog) },
    }));

  const renderErrorBar: CustomSeriesRenderItem = (params, api) => {
    const x = api.value(0) as number;
    const lo = api.coord([x, api.value(1) as number]);
    const hi = api.coord([x, api.value(2) as number]);
    // A bar shorter than the marker would only draw inside it.
    if (Math.abs(lo[1] - hi[1]) < SYMBOL_SIZE) return null;
    const color = api.visual("color") as string;
    const style = { stroke: color, lineWidth: 1 };
    return {
      type: "group",
      children: [
        {
          type: "line",
          shape: { x1: lo[0], y1: lo[1], x2: hi[0], y2: hi[1] },
          style,
        },
        {
          type: "line",
          shape: { x1: lo[0] - 3, y1: lo[1], x2: lo[0] + 3, y2: lo[1] },
          style,
        },
        {
          type: "line",
          shape: { x1: hi[0] - 3, y1: hi[1], x2: hi[0] + 3, y2: hi[1] },
          style,
        },
      ],
    };
  };

  const errorBars: CustomSeriesOption = {
    type: "custom",
    name: "errors",
    renderItem: renderErrorBar,
    data: errorData,
    encode: { x: 0, y: [1, 2] },
    silent: true,
    z: 2,
  };

  const axisCommon = {
    nameLocation: "middle" as const,
    nameTextStyle: { color: "#bfbfbf" },
    axisLine: { lineStyle: { color: "#303030" } },
    splitLine: { lineStyle: { color: "#202020" } },
  };

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { left: 70, right: 20, top: 36, bottom: 50 },
    legend: {
      top: 4,
      right: 8,
      itemWidth: 10,
      itemHeight: 10,
      textStyle: { color: "#999" },
      data: scatter.map((s) => s.name as string),
    },
    xAxis: {
      ...axisCommon,
      type: "log",
      name: "Wavelength (µm)",
      nameGap: 30,
      axisLabel: { color: "#d9d9d9" },
    },
    yAxis: {
      ...axisCommon,
      type: "log",
      name: "νFν (erg s⁻¹ cm⁻²)",
      nameGap: 55,
      axisLabel: {
        color: "#d9d9d9",
        formatter: (v: number) => v.toExponential(0),
      },
    },
    tooltip: {
      trigger: "item",
      formatter: (params: unknown) => {
        const data = (params as { data?: Partial<Datum> }).data;
        return data?.point ? formatPoint(data.point) : "";
      },
    },
    series: [errorBars, ...scatter],
  };

  return (
    <div>
      <ReactECharts option={option} className="h-64 w-full" />
      <Text type="secondary" className="text-xs block">
        Filled: this object · hollow: nearest counterpart · ▼ upper limit. Not
        corrected for extinction.
      </Text>
    </div>
  );
}
