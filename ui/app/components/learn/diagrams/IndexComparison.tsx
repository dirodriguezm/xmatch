"use client";

import { Slider } from "antd";
import { useState } from "react";

import { DiagramFrame, formatNumber, HIT, Stat } from "./shared";

const GAIA_DR3_SOURCES = 1_811_709_771;
const SKY_ARCSEC2 = 41_252.96 * 3600 * 3600;
const DENSITY = GAIA_DR3_SOURCES / SKY_ARCSEC2;
const CELL_AREA = SKY_ARCSEC2 / (12 * 4 ** 18); // order 18
const CELL = Math.sqrt(CELL_AREA);

const W = 420;
const BAR = 380;
const LOG_MAX = Math.log10(GAIA_DR3_SOURCES) + 3;
const barWidth = (v: number) =>
  Math.max(3, (BAR * Math.max(0, Math.log10(Math.max(v, 1e-3)) + 3)) / LOG_MAX);

const compact = (n: number) =>
  n >= 1e6
    ? `${formatNumber(n / 1e6)} million`
    : formatNumber(n, n < 10 ? 1 : 0);

/** Part 3: how much work a cone search does with and without the HEALPix index. */
export function IndexComparison() {
  const [radius, setRadius] = useState(3);
  // Cells touched by the disc, padded by half a cell diagonal; candidates at
  // the all-sky average Gaia density.
  const cells = Math.ceil((Math.PI * (radius + CELL * 0.71) ** 2) / CELL_AREA);
  const candidates = cells * CELL_AREA * DENSITY;

  const bars: [string, number, string][] = [
    ["Without an index: measure all of Gaia", GAIA_DR3_SOURCES, "#ff7875"],
    ["With HEALPix: cells looked up", cells, HIT],
    ["With HEALPix: candidates to measure", candidates, "#4096ff"],
  ];

  const stage = (
    <svg
      viewBox={`0 0 ${W} 300`}
      className="h-auto w-full max-w-[440px]"
      role="img"
      aria-label={`For a ${radius} arcsecond search, Gaia DR3 without an index needs ${compact(GAIA_DR3_SOURCES)} distance checks; with HEALPix it looks up ${cells} cells and measures about ${compact(candidates)} candidates`}
    >
      {bars.map(([label, value, color], i) => {
        const y = 30 + i * 84;
        return (
          <g key={label}>
            <text
              x={20}
              y={y}
              fill="#b9bcc6"
              fontSize={12}
              fontFamily="monospace"
            >
              {label}
            </text>
            <rect
              x={20}
              y={y + 10}
              width={BAR}
              height={22}
              rx={4}
              fill="#15171d"
            />
            <rect
              x={20}
              y={y + 10}
              width={barWidth(value)}
              height={22}
              rx={4}
              fill={color}
              fillOpacity={0.85}
            />
            <text
              x={20}
              y={y + 52}
              fill="#e8e9ec"
              fontSize={13}
              fontFamily="monospace"
            >
              {compact(value)}
            </text>
          </g>
        );
      })}
      <text x={400} y={292} textAnchor="end" fill="#868a96" fontSize={10}>
        log scale
      </text>
    </svg>
  );

  const side = (
    <>
      <div className="text-sm font-semibold text-foreground">
        Without an index vs with HEALPix
      </div>
      <p className="m-0 text-sm text-neutral-400">
        The same cone search in Gaia DR3. Without an index, every source in the
        catalog needs its distance measured. With HEALPix only the sources in
        the cells the circle touches are read.
      </p>
      <div className="flex flex-wrap gap-6">
        <Stat value={cells} label="cells looked up" />
        <Stat
          value={candidates < 1 ? "< 1" : compact(candidates)}
          label="candidates on average"
        />
      </div>
      <p className="m-0 text-xs text-neutral-500">
        Estimate at Gaia DR3&apos;s all-sky average density (~0.0034 sources per
        square arcsecond). The Galactic plane is 10 to 100 times denser, and
        even then it is tens of candidates.
      </p>
    </>
  );

  const controls = (
    <label className="flex items-center gap-3">
      Radius
      <Slider
        className="!w-40"
        min={1}
        max={30}
        value={radius}
        onChange={setRadius}
        tooltip={{ open: false }}
        aria-label="Search radius in arcseconds"
      />
      <span className="w-10 font-mono text-foreground">{radius}″</span>
    </label>
  );

  return <DiagramFrame stage={stage} side={side} controls={controls} />;
}
