"use client";

import type {
  CustomSeriesOption,
  CustomSeriesRenderItem,
  EChartsOption,
  ScatterSeriesOption,
} from "echarts";
import dynamic from "next/dynamic";

import { formatMjdDate, formatScientific } from "@/app/lib/utils/format";
import type { GammaPoint } from "@/app/lib/utils/gamma";

import { TIME_AXIS_LABEL, TIME_AXIS_NAME_GAP, TIME_GRID } from "./timeAxis";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

interface GammaLightCurveChartProps {
  points: GammaPoint[];
  color: string;
  /** Y axis title, with units. */
  yLabel: string;
  /** Log axis for fluxes; BAT rates can be negative, so they stay linear. */
  log?: boolean;
  symbolSize?: number;
  /** Shared time axis with the other light curves on the page. */
  mjdRange?: { min: number; max: number };
}

type Datum = {
  value: [number, number];
  point: GammaPoint;
  symbol: string;
  symbolRotate: number;
  itemStyle: { color: string; borderColor: string; borderWidth: number };
};

/** One gamma-ray or hard X-ray light curve with 1σ bars and 95% limits. */
export function GammaLightCurveChart({
  points,
  color,
  yLabel,
  log = false,
  symbolSize = 6,
  mjdRange,
}: GammaLightCurveChartProps) {
  const data: Datum[] = points.map((p) => ({
    value: [p.mjd, p.value],
    point: p,
    // Upper limits point down and are hollow, like the X-ray chart's.
    symbol: p.upperLimit ? "triangle" : "circle",
    symbolRotate: p.upperLimit ? 180 : 0,
    itemStyle: {
      color: p.upperLimit ? "transparent" : color,
      borderColor: color,
      borderWidth: 1.25,
    },
  }));

  const scatter: ScatterSeriesOption = {
    type: "scatter",
    symbolSize,
    z: 3,
    data,
  };

  const errorData = points.flatMap((p) =>
    !p.upperLimit && p.lo !== undefined && p.hi !== undefined
      ? // A log axis cannot draw ≤ 0.
        [{ value: [p.mjd, log ? Math.max(p.lo, p.value / 10) : p.lo, p.hi] }]
      : []
  );

  const renderErrorBar: CustomSeriesRenderItem = (_params, api) => {
    const x = api.value(0) as number;
    const [loX, loY] = api.coord([x, api.value(1) as number]);
    const [hiX, hiY] = api.coord([x, api.value(2) as number]);
    if (loX === undefined || loY === undefined) return null;
    if (hiX === undefined || hiY === undefined) return null;
    if (Math.abs(loY - hiY) < symbolSize) return null;
    return {
      type: "line",
      shape: { x1: loX, y1: loY, x2: hiX, y2: hiY },
      style: { stroke: color, lineWidth: 1, opacity: 0.6 },
    };
  };

  const errorBars: CustomSeriesOption = {
    type: "custom",
    renderItem: renderErrorBar,
    data: errorData,
    encode: { x: 0, y: [1, 2] },
    silent: true,
    z: 2,
  };

  // The page's shared range, widened only if a point falls outside it.
  const mjds = points.map((p) => p.mjd);
  const xRange =
    mjdRange && mjds.length > 0
      ? {
          min: Math.min(mjdRange.min, Math.min(...mjds)),
          max: Math.max(mjdRange.max, Math.max(...mjds)),
        }
      : mjdRange;

  const axisCommon = {
    nameLocation: "middle" as const,
    nameTextStyle: { color: "#bfbfbf" },
    axisLine: { lineStyle: { color: "#303030" } },
    splitLine: { lineStyle: { color: "#202020" } },
    axisLabel: { color: "#d9d9d9", hideOverlap: true },
  };

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { ...TIME_GRID, top: 16 },
    xAxis: {
      ...axisCommon,
      type: "value",
      name: "MJD",
      nameGap: TIME_AXIS_NAME_GAP,
      scale: true,
      axisLabel: TIME_AXIS_LABEL,
      ...(xRange ? { min: xRange.min, max: xRange.max } : {}),
    },
    yAxis: {
      ...axisCommon,
      type: log ? "log" : "value",
      name: yLabel,
      nameGap: 55,
      scale: !log,
      axisLabel: {
        color: "#d9d9d9",
        formatter: (v: number) => (v === 0 ? "0" : v.toExponential(0)),
      },
    },
    tooltip: {
      trigger: "item",
      formatter: (params: unknown) => {
        const p = (params as { data?: Partial<Datum> }).data?.point;
        if (!p) return "";
        const value = p.upperLimit
          ? `< ${formatScientific(p.value)} (95% upper limit)`
          : `${formatScientific(p.value)}${
              p.lo !== undefined && p.hi !== undefined
                ? ` (${formatScientific(p.lo)} – ${formatScientific(p.hi)})`
                : ""
            }`;
        return `${value}<br/>${formatMjdDate(p.mjd)} (MJD ${p.mjd.toFixed(1)})`;
      },
    },
    series: [errorBars, scatter],
  };

  return <ReactECharts option={option} className="h-56 w-full" />;
}
