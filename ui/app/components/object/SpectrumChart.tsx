"use client";

import { LineChartOutlined } from "@ant-design/icons";
import { Empty, Spin, Typography } from "antd";
import type { EChartsOption, LineSeriesOption } from "echarts";
import dynamic from "next/dynamic";

import { SKY_BANDS, SPECTRAL_LINES } from "@/app/lib/constants/spectralLines";
import { buildDesiSpectrumUrl } from "@/app/lib/utils/urls";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const { Text, Link } = Typography;

// Dark-theme adaptation of the Legacy Survey palette.
const COLOR_DATA = "#ff4d4f"; // observed flux
const COLOR_MODEL = "#e6e6e6"; // pipeline (redrock) best fit
const COLOR_NOISE = "#4096ff"; // 1/√ivar
const COLOR_EMISSION = "#d3adf7"; // violet line markers (brightened for legibility)
const COLOR_ABSORPTION = "#95de64"; // green line markers (brightened for legibility)

interface SpectrumChartProps {
  wavelength: number[] | undefined;
  flux: number[] | undefined;
  loading?: boolean;
  error?: Error | null;
  notFound?: boolean;
  targetid?: string | null;
  spectype?: string | null;
  redshift?: number | null;
  model?: number[] | null;
  ivar?: number[] | null;
}

export function SpectrumChart({
  wavelength,
  flux,
  loading,
  error,
  notFound,
  targetid,
  spectype,
  redshift,
  model,
  ivar,
}: SpectrumChartProps) {
  if (loading) {
    return (
      <div className="h-64 flex items-center justify-center">
        <Spin spinning>{null}</Spin>
      </div>
    );
  }

  // Upstream (SPARCL) unreachable or errored — distinct from "no spectrum".
  // Offer the external viewer as a fallback so the data is still reachable.
  if (error) {
    return (
      <div className="h-64 flex flex-col items-center justify-center gap-2">
        <Empty
          image={<LineChartOutlined className="text-4xl text-border" />}
          description={
            <Text type="secondary">DESI spectrum temporarily unavailable</Text>
          }
          styles={{ image: { height: 40 } }}
        />
        {targetid && (
          <Link
            href={buildDesiSpectrumUrl(targetid)}
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs"
          >
            Open in DESI spectrum viewer ↗
          </Link>
        )}
      </div>
    );
  }

  if (notFound || !wavelength || !flux || wavelength.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center">
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description={
            <Text type="secondary">No DESI spectrum for this object</Text>
          }
        />
      </div>
    );
  }

  // Zip wavelength/flux into [x, y] pairs for the observed-data series.
  const fluxPoints: [number, number][] = wavelength.map((w, i) => [w, flux[i]]);

  // Pipeline best-fit (redrock model) — drawn on top of the noisy data.
  const hasModel = Array.isArray(model) && model.length === wavelength.length;
  const modelPoints: [number, number][] | null = hasModel
    ? wavelength.map((w, i) => [w, model![i]])
    : null;

  // Noise = 1/√ivar; gap (null) where ivar is unusable so the trace breaks cleanly.
  const hasIvar = Array.isArray(ivar) && ivar.length === wavelength.length;
  const noisePoints: [number, number | null][] | null = hasIvar
    ? wavelength.map((w, i) => {
        const v = ivar![i];
        return [w, v && v > 0 ? 1 / Math.sqrt(v) : null];
      })
    : null;

  // Spectral-line markers at observed λ = rest × (1 + z), within the plotted range.
  const wMin = wavelength[0];
  const wMax = wavelength[wavelength.length - 1];
  const lineMarkers =
    redshift != null
      ? SPECTRAL_LINES.map((line) => ({
          ...line,
          obs: line.restWavelength * (1 + redshift),
        }))
          .filter((line) => line.obs >= wMin && line.obs <= wMax)
          .map((line) => ({
            xAxis: line.obs,
            lineStyle: {
              color:
                line.kind === "emission" ? COLOR_EMISSION : COLOR_ABSORPTION,
              type: "dashed" as const,
              width: 1,
              opacity: 0.65,
            },
            label: {
              show: true,
              formatter: line.label,
              position: "end" as const,
              rotate: 90,
              fontSize: 10,
              color:
                line.kind === "emission" ? COLOR_EMISSION : COLOR_ABSORPTION,
            },
          }))
      : [];

  const dataSeries: LineSeriesOption = {
    name: "data",
    type: "line",
    data: fluxPoints,
    showSymbol: false,
    lineStyle: { width: 1, color: COLOR_DATA },
    sampling: "lttb",
    z: 2,
    // Markers + sky bands ride along on the base data series.
    markLine: {
      silent: true,
      symbol: "none",
      data: lineMarkers,
    },
    markArea: {
      silent: true,
      itemStyle: { color: "rgba(64, 144, 255, 0.06)" },
      data: SKY_BANDS.map(([start, end]) => [{ xAxis: start }, { xAxis: end }]),
    },
  };

  const series: LineSeriesOption[] = [dataSeries];
  if (noisePoints) {
    series.push({
      name: "noise",
      type: "line",
      data: noisePoints,
      showSymbol: false,
      lineStyle: { width: 1, color: COLOR_NOISE, opacity: 0.7 },
      sampling: "lttb",
      z: 1,
    });
  }
  if (modelPoints) {
    series.push({
      name: "pipeline fit",
      type: "line",
      data: modelPoints,
      showSymbol: false,
      lineStyle: { width: 1.2, color: COLOR_MODEL },
      sampling: "lttb",
      z: 3,
    });
  }

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { left: 60, right: 20, top: 40, bottom: 50 },
    legend: {
      top: 4,
      right: 8,
      icon: "line",
      itemWidth: 16,
      itemHeight: 8,
      textStyle: { color: "#999" },
      data: series.map((s) => s.name as string),
    },
    xAxis: {
      type: "value",
      name: "Wavelength (Å)",
      nameLocation: "middle",
      nameGap: 30,
      nameTextStyle: { color: "#bfbfbf" },
      min: "dataMin",
      max: "dataMax",
      axisLabel: {
        color: "#d9d9d9",
        formatter: (value: number) => value.toFixed(0),
      },
      axisLine: { lineStyle: { color: "#303030" } },
      splitLine: { lineStyle: { color: "#202020" } },
    },
    yAxis: {
      type: "value",
      name: "Flux (10⁻¹⁷ erg s⁻¹ cm⁻² Å⁻¹)",
      nameLocation: "middle",
      nameGap: 45,
      nameTextStyle: { color: "#bfbfbf" },
      scale: true,
      axisLabel: {
        color: "#d9d9d9",
        formatter: (value: number) => value.toFixed(0),
      },
      axisLine: { lineStyle: { color: "#303030" } },
      splitLine: { lineStyle: { color: "#202020" } },
    },
    tooltip: {
      trigger: "axis",
      formatter: (params: unknown) => {
        const arr = params as {
          seriesName?: string;
          value?: [number, number | null];
          color?: string;
        }[];
        if (!arr?.length) return "";
        const wl = arr[0]?.value?.[0];
        if (wl == null) return "";
        const rows = arr
          .filter((p) => p.value?.[1] != null)
          .map(
            (p) =>
              `<span style="color:${p.color}">●</span> ${p.seriesName}: ${(
                p.value![1] as number
              ).toFixed(3)}`
          )
          .join("<br/>");
        return `λ: ${wl.toFixed(1)} Å<br/>${rows}`;
      },
    },
    dataZoom: [
      { type: "inside", xAxisIndex: 0 },
      { type: "inside", yAxisIndex: 0 },
    ],
    series,
  };

  return (
    <div>
      {(spectype || redshift != null) && (
        <Text type="secondary" className="text-xs block mb-2">
          {spectype ? `Type: ${spectype}` : null}
          {spectype && redshift != null ? "  ·  " : null}
          {redshift != null ? `z = ${redshift.toFixed(4)}` : null}
        </Text>
      )}
      <ReactECharts option={option} className="h-64 w-full" />
    </div>
  );
}
