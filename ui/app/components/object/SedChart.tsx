"use client";

import { DownloadOutlined } from "@ant-design/icons";
import { Button, Flex, Switch, Tooltip, Typography } from "antd";
import type {
  CustomSeriesOption,
  CustomSeriesRenderItem,
  EChartsOption,
  ScatterSeriesOption,
} from "echarts";
import dynamic from "next/dynamic";
import { useState } from "react";

import {
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import { equatorialToGalactic } from "@/app/lib/utils/coordinates";
import { downloadCsv } from "@/app/lib/utils/csv";
import { dereddenFactor, dereddeningBlocker } from "@/app/lib/utils/extinction";
import type { GalacticReddening } from "@/app/lib/utils/irsaDust";
import type { SedPoint } from "@/app/lib/utils/sed";
import { sedToCsv } from "@/app/lib/utils/sedCsv";
import { buildVizierSedViewerUrl } from "@/app/lib/utils/urls";
import type { VizierSedPoint } from "@/app/lib/utils/vizierSed";

const ReactECharts = dynamic(() => import("echarts-for-react"), { ssr: false });

const { Text, Link } = Typography;

/** An SED needs at least two points to show a shape. */
export const MIN_SED_POINTS = 2;

const SYMBOL_SIZE = 9;
const VIZIER_SYMBOL_SIZE = 6;
const COLOR_VIZIER = "#8c8c8c";

/** Where each of our own catalogs' SED bands comes from. */
const OWN_SOURCE_LABELS: Record<string, string> = {
  gaia: "Gaia DR3",
  allwise: "AllWISE (with its 2MASS associations)",
};

export interface VizierSedState {
  points: VizierSedPoint[];
  rowCount: number;
  loading: boolean;
  error: Error | null;
}

export interface ReddeningState {
  value?: GalacticReddening;
  loading: boolean;
  error: Error | null;
}

interface SedChartProps {
  ra: number;
  dec: number;
  points: SedPoint[];
  vizier: VizierSedState;
  reddening: ReddeningState;
  /** Download filename without extension. */
  filenameStem: string;
}

/** Why dereddening is unavailable, or null when it can be applied. */
function reddeningBlocker(
  reddening: ReddeningState,
  galacticLatitude: number
): string | null {
  if (reddening.loading) return "loading E(B−V)…";
  if (reddening.error || !reddening.value) return "E(B−V) unavailable";
  return dereddeningBlocker(reddening.value.ebvSF11, galacticLatitude);
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

type VizierDatum = { value: [number, number]; vizier: VizierSedPoint };

function formatVizierPoint(p: VizierSedPoint): string {
  const shown = p.tables.slice(0, 4).join(", ");
  const more = p.tables.length > 4 ? ` +${p.tables.length - 4} more` : "";
  const rejected =
    p.nRejected > 0
      ? ` (${p.nRejected} outlier${p.nRejected > 1 ? "s" : ""} dropped)`
      : "";
  return [
    `<b>${p.filter}</b> · ${p.wavelengthUm.toFixed(2)} µm`,
    `νFν: ${p.nuFnu.toExponential(2)} erg s⁻¹ cm⁻²`,
    ...(p.inconsistent
      ? [`<span style="color:#faad14">>30× off neighbouring filters</span>`]
      : []),
    `median of ${p.nMeasurements} measurement${p.nMeasurements > 1 ? "s" : ""}${rejected}`,
    `<span style="color:#999">VizieR: ${shown}${more}</span>`,
  ].join("<br/>");
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

export function SedChart({
  ra,
  dec,
  points: observedPoints,
  vizier,
  reddening,
  filenameStem,
}: SedChartProps) {
  const [showVizier, setShowVizier] = useState(true);
  const [showInconsistent, setShowInconsistent] = useState(false);
  const [deredden, setDeredden] = useState(false);

  const galacticLatitude = equatorialToGalactic(ra, dec).b;
  const blocker = reddeningBlocker(reddening, galacticLatitude);
  const ebv = blocker ? null : reddening.value!.ebvSF11;
  const applied = deredden && ebv != null;

  // Dereddening scales flux and error alike; the chart then plots intrinsic νFν.
  const correct = <
    T extends { wavelengthUm: number; nuFnu: number; nuFnuErr?: number },
  >(
    p: T
  ): T => {
    if (!applied) return p;
    const f = dereddenFactor(p.wavelengthUm, ebv!);
    return {
      ...p,
      nuFnu: p.nuFnu * f,
      nuFnuErr: p.nuFnuErr != null ? p.nuFnuErr * f : undefined,
    };
  };
  const points = observedPoints.map(correct);

  const inconsistentCount = vizier.points.filter((p) => p.inconsistent).length;
  const vizierShown = showVizier
    ? vizier.points
        .filter((p) => showInconsistent || !p.inconsistent)
        .map(correct)
    : [];
  const sources = (
    <SedSources
      ra={ra}
      dec={dec}
      points={points}
      vizier={vizier}
      reddening={reddening}
      filenameStem={filenameStem}
      inconsistentCount={showVizier ? inconsistentCount : 0}
      showInconsistent={showInconsistent}
      onToggleInconsistent={() => setShowInconsistent((v) => !v)}
    />
  );

  if (points.length + vizierShown.length < MIN_SED_POINTS) {
    return (
      <Flex vertical gap={4}>
        <Text type="secondary" className="text-xs block">
          {vizier.loading
            ? "Not enough photometry for an SED yet — loading VizieR photometry…"
            : "Not enough photometry for an SED."}
        </Text>
        {sources}
      </Flex>
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

  // Drawn small, grey and behind our own photometry so XWave's data keeps
  // the visual weight.
  const vizierSeries: ScatterSeriesOption | null =
    vizierShown.length > 0
      ? {
          name: "VizieR",
          type: "scatter",
          symbolSize: VIZIER_SYMBOL_SIZE,
          z: 1,
          itemStyle: { color: COLOR_VIZIER, opacity: 0.7 },
          data: vizierShown.map(
            (p): VizierDatum => ({
              value: [p.wavelengthUm, p.nuFnu],
              vizier: p,
            })
          ),
        }
      : null;

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
      data: [
        ...scatter.map((s) => s.name as string),
        ...(vizierSeries ? ["VizieR"] : []),
      ],
    },
    xAxis: {
      ...axisCommon,
      type: "log",
      name: "Wavelength (µm)",
      nameGap: 30,
      axisLabel: { color: "#d9d9d9", hideOverlap: true },
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
        const data = (params as { data?: Partial<Datum & VizierDatum> }).data;
        if (data?.point) return formatPoint(data.point);
        if (data?.vizier) return formatVizierPoint(data.vizier);
        return "";
      },
    },
    series: [errorBars, ...scatter, ...(vizierSeries ? [vizierSeries] : [])],
  };

  let extinctionNote: string;
  if (applied)
    extinctionNote = `Corrected for Galactic extinction: E(B−V) = ${ebv!.toFixed(3)} (Schlafly & Finkbeiner 2011), CCM89 law with R_V = 3.1.`;
  else if (ebv != null)
    extinctionNote = `Not corrected for extinction (E(B−V) = ${ebv.toFixed(3)} here).`;
  else if (reddening.loading || reddening.error)
    extinctionNote = "Not corrected for extinction.";
  else extinctionNote = `Not corrected for extinction: ${blocker}.`;

  return (
    <Flex vertical gap={4}>
      <Flex justify="space-between" align="center" wrap="wrap" gap={8}>
        <Button
          size="small"
          type="text"
          icon={<DownloadOutlined />}
          onClick={() =>
            downloadCsv(
              `${filenameStem}_sed.csv`,
              sedToCsv(observedPoints, vizier.points, ebv)
            )
          }
        >
          SED as CSV
        </Button>
        <Flex align="center" gap={16} wrap="wrap">
          <Tooltip title={blocker ?? undefined}>
            <Flex align="center" gap={8}>
              <Text type="secondary" className="text-xs">
                Deredden
              </Text>
              <Switch
                size="small"
                checked={applied}
                disabled={ebv == null}
                onChange={setDeredden}
                aria-label="Correct for Galactic extinction"
              />
            </Flex>
          </Tooltip>
          {vizier.points.length > 0 && (
            <Flex align="center" gap={8}>
              <Text type="secondary" className="text-xs">
                VizieR photometry ({vizier.points.length} filters)
              </Text>
              <Switch
                size="small"
                checked={showVizier}
                onChange={setShowVizier}
                aria-label="Show VizieR photometry"
              />
            </Flex>
          )}
        </Flex>
      </Flex>
      <ReactECharts option={option} className="h-64 w-full" />
      <Text type="secondary" className="text-xs block">
        Filled: this object · hollow: nearest counterpart · ▼ upper limit
        {vizierShown.length > 0 ? " · grey: VizieR" : ""}. {extinctionNote}
      </Text>
      {sources}
    </Flex>
  );
}

/** Attribution for every dataset drawn in the SED. */
function SedSources({
  ra,
  dec,
  points,
  vizier,
  reddening,
  inconsistentCount,
  showInconsistent,
  onToggleInconsistent,
}: SedChartProps & {
  inconsistentCount: number;
  showInconsistent: boolean;
  onToggleInconsistent: () => void;
}) {
  const own = [...new Set(points.map((p) => p.catalog))]
    .map((c) => OWN_SOURCE_LABELS[c] ?? getSearchCatalogLabel(c))
    .join(", ");
  const tableCount = new Set(vizier.points.flatMap((p) => p.tables)).size;

  let vizierText: string;
  if (vizier.loading) vizierText = "loading…";
  else if (vizier.error) vizierText = "unavailable right now";
  else if (vizier.points.length === 0)
    vizierText = "no published photometry within 2″";
  else
    vizierText = `${vizier.rowCount} measurements from ${tableCount} catalogs, merged per filter`;

  return (
    <Text type="secondary" className="text-xs block">
      Sources: {own ? `${own} via XWave · ` : ""}
      <Link
        href={buildVizierSedViewerUrl(ra, dec)}
        target="_blank"
        rel="noopener noreferrer"
        className="text-xs"
      >
        VizieR SED, CDS Strasbourg ↗
      </Link>{" "}
      ({vizierText})
      {reddening.value && (
        <>
          {" · E(B−V): "}
          <Link
            href="https://irsa.ipac.caltech.edu/applications/DUST/"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs"
          >
            IRSA Galactic Dust ↗
          </Link>
        </>
      )}
      {inconsistentCount > 0 && (
        <>
          {" · "}
          {inconsistentCount} filter{inconsistentCount > 1 ? "s" : ""}{" "}
          {showInconsistent ? "shown" : "hidden"} as &gt;30× off neighbouring
          filters (
          <Link onClick={onToggleInconsistent} className="text-xs">
            {showInconsistent ? "hide" : "show"}
          </Link>
          )
        </>
      )}
    </Text>
  );
}
