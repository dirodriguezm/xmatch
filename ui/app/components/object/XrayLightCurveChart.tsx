"use client";

import type {
  CustomSeriesOption,
  CustomSeriesRenderItem,
  EChartsOption,
  ScatterSeriesOption,
} from "echarts";
import dynamic from "next/dynamic";

import { formatFlux, formatMjdDate } from "@/app/lib/utils/format";
import {
  HILIGT_MISSIONS,
  type HiligtMission,
  type HiligtPoint,
} from "@/app/lib/utils/hiligt";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const SYMBOL_SIZE = 8;

export const MISSION_COLORS: Record<HiligtMission, string> = {
  XMMpnt: "#1677ff",
  XMMslew: "#13c2c2",
  RosatSurvey: "#fa8c16",
};

/** Same colours as Tailwind classes, for legend dots outside the chart. */
export const MISSION_DOT_CLASSES: Record<HiligtMission, string> = {
  XMMpnt: "bg-[#1677ff]",
  XMMslew: "bg-[#13c2c2]",
  RosatSurvey: "bg-[#fa8c16]",
};

type Datum = {
  value: [number, number];
  point: HiligtPoint;
  symbol: string;
  symbolRotate: number;
  itemStyle: { color: string; borderColor: string; borderWidth: number };
};

function toDatum(p: HiligtPoint): Datum {
  const color = MISSION_COLORS[p.mission];
  const isLimit = p.flux === undefined;
  return {
    value: [p.mjd, (isLimit ? p.upperLimit : p.flux)!],
    point: p,
    // Upper limits point down and are hollow, like the SED's.
    symbol: isLimit ? "triangle" : "circle",
    symbolRotate: isLimit ? 180 : 0,
    itemStyle: {
      color: isLimit ? "transparent" : color,
      borderColor: color,
      borderWidth: 1.5,
    },
  };
}

function formatPoint(p: HiligtPoint): string {
  const value =
    p.flux !== undefined
      ? `${formatFlux(p.flux)}${p.fluxErr ? ` ± ${p.fluxErr.toExponential(1)}` : ""}`
      : `< ${formatFlux(p.upperLimit!)} (${p.sigma ?? 2}σ upper limit)`;
  return [
    `<b>${HILIGT_MISSIONS[p.mission]}</b>${p.instrument ? ` · ${p.instrument}` : ""}`,
    `${p.band} keV: ${value}`,
    `${formatMjdDate(p.mjd)} (MJD ${p.mjd.toFixed(1)})`,
    `<span style="color:#999">obs ${p.obsid ?? "—"}${p.exposureS ? ` · ${Math.round(p.exposureS)} s` : ""}</span>`,
  ].join("<br/>");
}

interface XrayLightCurveChartProps {
  points: HiligtPoint[];
  /** Shared time axis with the optical light curves. */
  mjdRange?: { min: number; max: number };
}

/** Long-term X-ray fluxes and upper limits on a log flux axis. */
export function XrayLightCurveChart({
  points,
  mjdRange,
}: XrayLightCurveChartProps) {
  const missions = (Object.keys(HILIGT_MISSIONS) as HiligtMission[]).filter(
    (m) => points.some((p) => p.mission === m)
  );

  const scatter: ScatterSeriesOption[] = missions.map((mission) => ({
    name: HILIGT_MISSIONS[mission],
    type: "scatter",
    symbolSize: SYMBOL_SIZE,
    z: 3,
    itemStyle: { color: MISSION_COLORS[mission] },
    data: points.filter((p) => p.mission === mission).map(toDatum),
  }));

  const errorData = points
    .filter((p) => p.flux !== undefined && p.fluxErr)
    .map((p) => ({
      value: [
        p.mjd,
        // A log axis cannot draw ≤ 0.
        Math.max(p.flux! - p.fluxErr!, p.flux! / 10),
        p.flux! + p.fluxErr!,
      ],
      itemStyle: { color: MISSION_COLORS[p.mission] },
    }));

  const renderErrorBar: CustomSeriesRenderItem = (params, api) => {
    const x = api.value(0) as number;
    const lo = api.coord([x, api.value(1) as number]);
    const hi = api.coord([x, api.value(2) as number]);
    if (Math.abs(lo[1] - hi[1]) < SYMBOL_SIZE) return null;
    const style = { stroke: api.visual("color") as string, lineWidth: 1 };
    return {
      type: "line",
      shape: { x1: lo[0], y1: lo[1], x2: hi[0], y2: hi[1] },
      style,
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

  // Start from the optical light curves' axis so panels line up, but widen it
  // to fit every X-ray point: ROSAT (1990) predates all the optical surveys.
  const mjds = points.map((p) => p.mjd);
  const pad = 200;
  const xRange =
    mjds.length > 0
      ? {
          min: Math.min(mjdRange?.min ?? Infinity, Math.min(...mjds) - pad),
          max: Math.max(mjdRange?.max ?? -Infinity, Math.max(...mjds) + pad),
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
      type: "value",
      name: "MJD",
      nameGap: 30,
      scale: true,
      ...(xRange
        ? { min: Math.floor(xRange.min), max: Math.ceil(xRange.max) }
        : {}),
    },
    yAxis: {
      ...axisCommon,
      type: "log",
      name: "Flux (erg s⁻¹ cm⁻²)",
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

  return <ReactECharts option={option} className="h-64 w-full" />;
}
