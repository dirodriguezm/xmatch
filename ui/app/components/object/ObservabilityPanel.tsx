"use client";

import { LeftOutlined, RightOutlined } from "@ant-design/icons";
import { Button, DatePicker, Flex, Select, Typography } from "antd";
import dayjs, { type Dayjs } from "dayjs";
import type { EChartsOption, LineSeriesOption } from "echarts";
import dynamic from "next/dynamic";
import { type ReactNode, useMemo, useState } from "react";

import {
  DEFAULT_OBSERVATORY_ID,
  OBSERVATORIES,
  OBSERVATORY_TIME_ZONE,
} from "@/app/lib/constants/observatories";
import {
  computeNightVisibility,
  DEFAULT_MIN_ALTITUDE_DEG,
  formatChileTime,
  observableHours,
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

/** Calendar dot per night: how long the target is above the minimum altitude. */
function hoursDotClass(hours: number): string | null {
  if (hours >= 4) return "bg-green-500";
  if (hours >= 1) return "bg-amber-400";
  if (hours > 0) return "bg-amber-700";
  return null;
}

function CalendarLegend() {
  const items = [
    { cls: "bg-green-500", label: "≥ 4 h" },
    { cls: "bg-amber-400", label: "1–4 h" },
    { cls: "bg-amber-700", label: "< 1 h" },
  ];
  return (
    <Flex gap={12} wrap className="py-1 text-xs text-neutral-400">
      <span>Above {DEFAULT_MIN_ALTITUDE_DEG}° at night:</span>
      {items.map((i) => (
        <span key={i.label} className="inline-flex items-center gap-1">
          <span className={`inline-block w-1.5 h-1.5 rounded-full ${i.cls}`} />
          {i.label}
        </span>
      ))}
      <span>no dot: not observable</span>
    </Flex>
  );
}

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

  // Hours observable per night for the month the calendar shows, plus the
  // neighbouring days it greys in (~1.5 ms a night, so ~80 ms a month).
  // Stable identity: a fresh dayjs each render would make the picker snap its
  // panel back to this month whenever the user browses to another.
  const nightValue = useMemo(() => dayjs(night), [night]);
  const [panelMonth, setPanelMonth] = useState(() =>
    dayjs(night).startOf("month")
  );
  const monthHours = useMemo(() => {
    const hours = new Map<string, number>();
    const first = panelMonth.subtract(7, "day");
    const last = panelMonth.endOf("month").add(14, "day");
    for (let d = first; !d.isAfter(last, "day"); d = d.add(1, "day")) {
      hours.set(
        d.format("YYYY-MM-DD"),
        observableHours(computeNightVisibility(ra, dec, site, d.toDate()))
      );
    }
    return hours;
  }, [ra, dec, site, panelMonth]);
  const hoursOn = (d: Dayjs): number =>
    monthHours.get(d.format("YYYY-MM-DD")) ??
    observableHours(computeNightVisibility(ra, dec, site, d.toDate()));

  const renderCell = (current: Dayjs, originNode: ReactNode): ReactNode => {
    const hours = hoursOn(current);
    const dot = hoursDotClass(hours);
    return (
      <div
        className="relative"
        title={
          hours > 0
            ? `${hours.toFixed(1)} h above ${DEFAULT_MIN_ALTITUDE_DEG}°`
            : `Not above ${DEFAULT_MIN_ALTITUDE_DEG}° at night`
        }
      >
        {originNode}
        {dot && (
          <span
            className={`absolute left-1/2 -translate-x-1/2 -bottom-0.5 w-1 h-1 rounded-full ${dot}`}
          />
        )}
      </div>
    );
  };

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
          <DatePicker
            size="small"
            value={nightValue}
            onChange={(d) =>
              d && setNight(new Date(d.year(), d.month(), d.date()))
            }
            allowClear={false}
            format={(d) => `Night of ${d.format("ddd, D MMM YYYY")}`}
            cellRender={(current, info) =>
              info.type === "date" && dayjs.isDayjs(current)
                ? renderCell(current, info.originNode)
                : info.originNode
            }
            renderExtraFooter={() => <CalendarLegend />}
            onOpenChange={(open) =>
              open && setPanelMonth(dayjs(night).startOf("month"))
            }
            onPanelChange={(d) => setPanelMonth(d.startOf("month"))}
            aria-label="Choose a night"
            className="!w-[210px]"
          />
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
