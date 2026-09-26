"use client";

import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { Button, Flex, Select, Typography } from "antd";
import type { EChartsOption, LineSeriesOption } from "echarts";
import dynamic from "next/dynamic";
import { useMemo, useState } from "react";

import {
  DEFAULT_OBSERVATORY_ID,
  OBSERVATORIES,
  OBSERVATORY_TIME_ZONE,
} from "@/app/lib/constants/observatories";
import {
  computeNightVisibility,
  DEFAULT_MIN_ALTITUDE_DEG,
  formatChileTime,
  summarizeVisibility,
  tonightInTimeZone,
} from "@/app/lib/utils/observability";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const { Text, Link } = Typography;

const COLOR_TARGET = "#4096ff";
const COLOR_MOON = "#8c8c8c";
const COLOR_DAY = "rgba(255, 255, 255, 0.07)";
const COLOR_TWILIGHT = "rgba(255, 255, 255, 0.035)";

/** Altitudes where the airmass guide lines sit (X = 1.5, 2, 3). */
const AIRMASS_GUIDES = [
  { altitude: 41.8, label: "X 1.5" },
  { altitude: 30, label: "X 2" },
  { altitude: 19.3, label: "X 3" },
];

const dateFormat = new Intl.DateTimeFormat("en-GB", {
  weekday: "short",
  day: "numeric",
  month: "short",
  year: "numeric",
});

const hhmm = formatChileTime;

function addDays(date: Date, days: number): Date {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate() + days);
}

function sameDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  );
}

interface ObservabilityPanelProps {
  ra: number;
  dec: number;
}

export function ObservabilityPanel({ ra, dec }: ObservabilityPanelProps) {
  const tonight = useMemo(() => tonightInTimeZone(OBSERVATORY_TIME_ZONE), []);
  const [siteId, setSiteId] = useState(DEFAULT_OBSERVATORY_ID);
  const [night, setNight] = useState(tonight);

  const site = OBSERVATORIES.find((o) => o.id === siteId) ?? OBSERVATORIES[0];
  const visibility = useMemo(
    () => computeNightVisibility(ra, dec, site, night),
    [ra, dec, site, night]
  );

  const start = visibility.samples[0].time;
  const end = visibility.samples[visibility.samples.length - 1].time;
  const moonPct = Math.round(visibility.moon.illumination * 100);

  const targetSeries: LineSeriesOption = {
    name: "Target",
    color: COLOR_TARGET,
    type: "line",
    showSymbol: false,
    lineStyle: { width: 2, color: COLOR_TARGET },
    itemStyle: { color: COLOR_TARGET },
    data: visibility.samples.map((s) => [s.time.getTime(), s.altitude]),
    z: 3,
    markArea: {
      silent: true,
      data: [
        [
          { xAxis: start.getTime(), itemStyle: { color: COLOR_DAY } },
          { xAxis: visibility.sunset.getTime() },
        ],
        [
          {
            xAxis: visibility.sunset.getTime(),
            itemStyle: { color: COLOR_TWILIGHT },
          },
          { xAxis: visibility.duskAstronomical.getTime() },
        ],
        [
          {
            xAxis: visibility.dawnAstronomical.getTime(),
            itemStyle: { color: COLOR_TWILIGHT },
          },
          { xAxis: visibility.sunrise.getTime() },
        ],
        [
          {
            xAxis: visibility.sunrise.getTime(),
            itemStyle: { color: COLOR_DAY },
          },
          { xAxis: end.getTime() },
        ],
      ],
    },
    markLine: {
      silent: true,
      symbol: "none",
      data: AIRMASS_GUIDES.map((g) => ({
        yAxis: g.altitude,
        lineStyle: { color: "#434343", type: "dashed" as const, width: 1 },
        label: {
          formatter: g.label,
          position: "end" as const,
          color: "#8c8c8c",
          fontSize: 10,
        },
      })),
    },
  };

  const moonSeries: LineSeriesOption = {
    name: `Moon (${moonPct}%)`,
    color: COLOR_MOON,
    type: "line",
    showSymbol: false,
    lineStyle: { width: 1, color: COLOR_MOON, type: "dashed" },
    itemStyle: { color: COLOR_MOON },
    data: visibility.samples.map((s) => [s.time.getTime(), s.moonAltitude]),
    z: 2,
  };

  const axisStyle = {
    axisLine: { lineStyle: { color: "#303030" } },
    splitLine: { lineStyle: { color: "#202020" } },
  };

  const option: EChartsOption = {
    backgroundColor: "transparent",
    grid: { left: 50, right: 50, top: 36, bottom: 40 },
    legend: {
      top: 4,
      right: 8,
      // A thin bar reads as a line key; ECharts' "line" icon renders empty here.
      icon: "rect",
      itemWidth: 16,
      itemHeight: 2,
      textStyle: { color: "#999" },
    },
    xAxis: {
      ...axisStyle,
      type: "time",
      min: start.getTime(),
      max: end.getTime(),
      name: "Chile time",
      nameLocation: "middle",
      nameGap: 26,
      nameTextStyle: { color: "#bfbfbf" },
      axisLabel: {
        color: "#d9d9d9",
        formatter: (value: number) => hhmm(new Date(value)),
      },
      splitLine: { show: false },
    },
    yAxis: {
      ...axisStyle,
      type: "value",
      min: 0,
      max: 90,
      interval: 15,
      name: "Altitude (°)",
      nameLocation: "middle",
      nameGap: 32,
      nameTextStyle: { color: "#bfbfbf" },
      axisLabel: { color: "#d9d9d9" },
    },
    tooltip: {
      trigger: "axis",
      formatter: (params: unknown) => {
        const arr = params as {
          value?: [number, number];
          seriesIndex?: number;
        }[];
        const t = arr[0]?.value?.[0];
        if (t == null) return "";
        const sample = visibility.samples.find((s) => s.time.getTime() === t);
        if (!sample) return "";
        const am = sample.airmass != null ? sample.airmass.toFixed(2) : "—";
        return [
          `<b>${hhmm(sample.time)}</b>`,
          `Target: ${sample.altitude.toFixed(1)}° · airmass ${am}`,
          `Moon: ${sample.moonAltitude.toFixed(1)}°`,
          `Sun: ${sample.sunAltitude.toFixed(1)}°`,
        ].join("<br/>");
      },
    },
    series: [targetSeries, moonSeries],
  };

  const neverUp = visibility.maxPossibleAltitude < DEFAULT_MIN_ALTITUDE_DEG;

  return (
    <Flex vertical gap={12}>
      <Flex wrap="wrap" gap={8} align="center" justify="space-between">
        <Select
          size="small"
          value={siteId}
          onChange={setSiteId}
          className="min-w-[220px]"
          options={OBSERVATORIES.map((o) => ({ label: o.label, value: o.id }))}
        />
        <Flex align="center" gap={4}>
          <Button
            size="small"
            type="text"
            icon={<LeftOutlined />}
            aria-label="Previous night"
            onClick={() => setNight((d) => addDays(d, -1))}
          />
          <Text className="text-sm tabular-nums min-w-[130px] text-center">
            Night of {dateFormat.format(night)}
          </Text>
          <Button
            size="small"
            type="text"
            icon={<RightOutlined />}
            aria-label="Next night"
            onClick={() => setNight((d) => addDays(d, 1))}
          />
          <Button
            size="small"
            disabled={sameDay(night, tonight)}
            onClick={() => setNight(tonight)}
          >
            Tonight
          </Button>
        </Flex>
      </Flex>

      <div>
        <Text className="block">
          {summarizeVisibility(visibility, site.label)}
        </Text>
        <Text type="secondary" className="text-xs block">
          Astronomical night {hhmm(visibility.duskAstronomical)}–
          {hhmm(visibility.dawnAstronomical)} · Moon {moonPct}% illuminated,{" "}
          {visibility.moon.separationDeg.toFixed(0)}° away · times in Chile time
          <br />
          Computed in your browser with{" "}
          <Link
            href="https://github.com/cosinekitty/astronomy"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs"
          >
            Astronomy Engine ↗
          </Link>
          ; J2000 coordinates precessed to date, refraction included.
        </Text>
      </div>

      {!neverUp && <ReactECharts option={option} className="h-64 w-full" />}
    </Flex>
  );
}
