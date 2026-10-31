"use client";

import { Segmented, Slider, Tag } from "antd";
import { useState } from "react";

import { MAX_RADIUS_ARCSEC } from "@/app/lib/constants/search";

import { DiagramFrame, formatNumber, MUTED, Stat } from "./shared";

const CATALOG_SIZE = 1.8e9; // a Gaia-sized catalog
const ARCSEC2_PER_SR = ((180 / Math.PI) * 3600) ** 2;
const CELL_AREA = (4 * Math.PI * ARCSEC2_PER_SR) / (12 * 4 ** 18);
const R_MIN = 1;
const R_MAX = 180 * 3600; // the whole sky
const SLIDER_MAX = 1000;
const BYTES_PER_SOURCE = 100;

/** Slider position ↔ radius on a log scale from 1″ to the whole sky. */
const radiusAt = (t: number) => R_MIN * (R_MAX / R_MIN) ** (t / SLIDER_MAX);
const sliderAt = (r: number) =>
  (SLIDER_MAX * Math.log(r / R_MIN)) / Math.log(R_MAX / R_MIN);

/** Area of a spherical cap of the given radius, in square arcseconds. */
const capArea = (r: number) =>
  2 *
  Math.PI *
  (1 - Math.cos((Math.min(r / 3600, 180) * Math.PI) / 180)) *
  ARCSEC2_PER_SR;

const TYPICAL_DENSITY = 0.0034;
const ZONES = [
  { label: "Typical sky", density: TYPICAL_DENSITY },
  { label: "High latitude", density: 0.0006 },
  { label: "Galactic plane", density: 0.04 },
];

const VERDICTS = [
  { max: 1e4, label: "Light", color: "success", fill: "#73d13d" },
  { max: 1e7, label: "Heavy", color: "warning", fill: "#ffc53d" },
  { max: Infinity, label: "Impractical", color: "error", fill: "#ff7875" },
] as const;

const radiusLabel = (r: number) =>
  r < 60
    ? `${formatNumber(r, r < 10 ? 1 : 0)}″`
    : r < 3600
      ? `${formatNumber(r / 60, r < 600 ? 1 : 0)}′`
      : r >= 179.9 * 3600
        ? "whole sky"
        : `${formatNumber(r / 3600, r < 36000 ? 1 : 0)}°`;

const compact = (n: number) =>
  n >= 1e9
    ? `${formatNumber(n / 1e9, 1)} billion`
    : n >= 1e6
      ? `${formatNumber(n / 1e6, n < 1e7 ? 1 : 0)} million`
      : formatNumber(n, n < 10 ? 1 : 0);

const bytes = (n: number) => {
  const b = n * BYTES_PER_SOURCE;
  if (b >= 1e9) return `${formatNumber(b / 1e9, 1)} GB`;
  if (b >= 1e6) return `${formatNumber(b / 1e6, 1)} MB`;
  if (b >= 1e3) return `${formatNumber(b / 1e3)} kB`;
  return `${Math.round(b)} B`;
};

const L = 54;
const R = 446;
const T = 16;
const B = 280;
const Y_MIN = -2;
const Y_MAX = Math.log10(CATALOG_SIZE);
const X = (r: number) => L + ((R - L) * sliderAt(r)) / SLIDER_MAX;
const Y = (v: number) =>
  B -
  ((B - T) * (Math.log10(Math.max(v, 10 ** Y_MIN)) - Y_MIN)) / (Y_MAX - Y_MIN);

/** Part 5: how the work of a cone search grows with its radius (illustrative). */
export function RadiusScaleDiagram() {
  const [t, setT] = useState(Math.round(sliderAt(15)));
  const [density, setDensity] = useState(TYPICAL_DENSITY);

  const work = (r: number) => Math.min(CATALOG_SIZE, density * capArea(r));
  const radius = radiusAt(t);
  const candidates = work(radius);
  const cells = capArea(radius) / CELL_AREA;
  const zone = VERDICTS.findIndex((v) => candidates < v.max);
  const verdict = VERDICTS[zone] ?? VERDICTS[2];

  let curve = "";
  for (let s = 0; s <= SLIDER_MAX; s += 10) {
    const r = radiusAt(s);
    curve += `${s ? "L" : "M"}${X(r).toFixed(1)} ${Y(work(r)).toFixed(1)}`;
  }

  const stage = (
    <svg
      viewBox="0 0 460 330"
      className="h-auto w-full max-w-[480px]"
      role="img"
      aria-label={`Candidates to measure against search radius on log scales; at ${radiusLabel(radius)} about ${compact(candidates)}`}
    >
      {VERDICTS.map((v, i) => {
        const lo = VERDICTS[i - 1]?.max ?? 10 ** Y_MIN;
        const hi = Math.min(v.max, CATALOG_SIZE);
        return (
          <g key={v.label}>
            <rect
              x={L}
              y={Y(hi)}
              width={R - L}
              height={Y(lo) - Y(hi)}
              fill={v.fill}
              fillOpacity={0.07}
            />
            <text
              x={R - 6}
              y={Y(hi) + 14}
              textAnchor="end"
              fontSize={10.5}
              fill={v.fill}
            >
              {v.label.toLowerCase()}
            </text>
          </g>
        );
      })}
      {[0, 3, 6, 9].map((e) => (
        <g key={e}>
          <line
            x1={L}
            x2={R}
            y1={Y(10 ** e)}
            y2={Y(10 ** e)}
            stroke="#23252d"
          />
          <text
            x={L - 6}
            y={Y(10 ** e) + 4}
            textAnchor="end"
            fontSize={10}
            fill={MUTED}
          >
            {["1", "1k", "1M", "1B"][e / 3]}
          </text>
        </g>
      ))}
      {(
        [
          [1, "1″"],
          [60, "60″"],
          [3600, "1°"],
          [36000, "10°"],
          [R_MAX, "sky"],
        ] as const
      ).map(([v, label]) => (
        <g key={label}>
          <line x1={X(v)} x2={X(v)} y1={B} y2={B + 5} stroke={MUTED} />
          <text
            x={X(v)}
            y={B + 18}
            textAnchor="middle"
            fontSize={10.5}
            fill={MUTED}
          >
            {label}
          </text>
        </g>
      ))}
      <text x={L} y={B + 40} fontSize={10.5} fill={MUTED}>
        radius (log scale)
      </text>
      <text x={4} y={T - 4} fontSize={10.5} fill={MUTED}>
        candidates to measure
      </text>
      <line
        x1={X(MAX_RADIUS_ARCSEC)}
        x2={X(MAX_RADIUS_ARCSEC)}
        y1={T}
        y2={B}
        stroke="#4096ff"
        strokeDasharray="4 4"
        strokeOpacity={0.7}
      />
      <text
        x={X(MAX_RADIUS_ARCSEC) + 4}
        y={T + 12}
        fontSize={10.5}
        fill="#4096ff"
      >
        XWave limit: {MAX_RADIUS_ARCSEC}″
      </text>
      <path d={curve} fill="none" stroke="#e8e9ec" strokeWidth={2} />
      <circle
        cx={X(radius)}
        cy={Y(candidates)}
        r={6}
        fill="#ffc53d"
        stroke="#0a0a0a"
        strokeWidth={2}
      />
    </svg>
  );

  const side = (
    <>
      <Tag color={verdict.color} className="!m-0 self-start">
        {verdict.label}
      </Tag>
      <div className="flex flex-wrap gap-6">
        <Stat value={compact(cells)} label="HEALPix cells" />
        <Stat
          value={candidates < 1 ? "< 1" : compact(candidates)}
          label="candidates to measure"
        />
        <Stat value={bytes(candidates)} label="data to read and sort" />
      </div>
      <p className="m-0 text-sm text-neutral-300" aria-live="polite">
        {zone === 0 &&
          "The circle is tiny: a few cells and hardly any candidates. This is what an everyday search looks like."}
        {zone === 1 &&
          "Now it is thousands to millions of sources to read, measure and sort, and a huge response to send to the browser. Areas this size call for a catalog query, not a cross-match."}
        {zone === 2 &&
          `At this size it is no longer a cross-match: it means reading ${
            radius >= 179.9 * 3600
              ? "the whole catalog"
              : "a large share of the catalog"
          }. An index skips what is far away, but it cannot skip what is inside the circle.`}
      </p>
      <p className="m-0 text-xs text-neutral-500">
        Illustrative example: a catalog of 1.8 billion sources at typical
        densities for each part of the sky. Not a measurement of the service.
      </p>
    </>
  );

  const controls = (
    <>
      <label className="flex items-center gap-3">
        Radius
        <Slider
          className="!w-48"
          min={0}
          max={SLIDER_MAX}
          value={t}
          onChange={setT}
          tooltip={{ open: false }}
          aria-label="Search radius, log scale from 1 arcsecond to the whole sky"
        />
        <span className="w-20 font-mono text-foreground">
          {radiusLabel(radius)}
        </span>
      </label>
      <span className="flex items-center gap-3">
        Part of the sky
        <Segmented
          size="small"
          value={density}
          onChange={(v) => setDensity(v as number)}
          options={ZONES.map((z) => ({ label: z.label, value: z.density }))}
        />
      </span>
    </>
  );

  return <DiagramFrame stage={stage} side={side} controls={controls} />;
}
