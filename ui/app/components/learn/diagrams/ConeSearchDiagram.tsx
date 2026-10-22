"use client";

import { Slider } from "antd";
import { useMemo, useState } from "react";

import {
  type CatalogOption,
  getSearchCatalogColor,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";

import {
  DiagramFrame,
  formatNumber,
  pickCatalog,
  seededRandom,
  Stat,
} from "./shared";

const FIELD = 40; // arcsec across
const SIZE = 400;
const SCALE = SIZE / FIELD;
const C = SIZE / 2;

interface Source {
  id: string;
  x: number;
  y: number;
  catalog: CatalogOption;
}

function exampleField(): Source[] {
  const rnd = seededRandom(7);
  const sources: Source[] = [];
  for (let i = 0; i < 46; i++) {
    const x = rnd() * FIELD - FIELD / 2;
    const y = rnd() * FIELD - FIELD / 2;
    sources.push({ id: `${100 + i}`, x, y, catalog: pickCatalog(rnd()) });
  }
  // The target itself, as each catalog sees it with its own positional error.
  sources.push(
    { id: "146", x: 0.15, y: -0.1, catalog: "gaia" },
    { id: "147", x: 0.9, y: 0.6, catalog: "allwise" },
    { id: "148", x: -3.2, y: 2.4, catalog: "erosita" }
  );
  return sources;
}

/** Part 1: a cone search on an example field, with a live result list. */
export function ConeSearchDiagram() {
  const [radius, setRadius] = useState(6);
  const [nearest, setNearest] = useState(5);
  const sources = useMemo(() => exampleField(), []);

  const ranked = useMemo(
    () =>
      sources
        .map((s) => ({ ...s, d: Math.hypot(s.x, s.y) }))
        .filter((s) => s.d <= radius)
        .sort((a, b) => a.d - b.d),
    [sources, radius]
  );
  const kept = new Set(ranked.slice(0, nearest).map((s) => s.id));

  const stage = (
    <svg
      viewBox={`0 0 ${SIZE} ${SIZE}`}
      className="h-auto w-full max-w-[420px]"
      role="img"
      aria-label={`Example 40 arcsecond field with ${sources.length} sources from three catalogs and a ${radius} arcsecond search circle`}
    >
      {[-20, -10, 0, 10, 20].map((k) => (
        <g key={k} stroke="#1b1d24">
          <line x1={C + k * SCALE} x2={C + k * SCALE} y1={0} y2={SIZE} />
          <line y1={C + k * SCALE} y2={C + k * SCALE} x1={0} x2={SIZE} />
        </g>
      ))}
      <circle
        cx={C}
        cy={C}
        r={radius * SCALE}
        fill="rgba(64,150,255,.08)"
        stroke="#4096ff"
        strokeWidth={1.5}
        strokeDasharray="5 4"
      />
      {sources.map((s) => {
        const inside = Math.hypot(s.x, s.y) <= radius;
        return (
          <circle
            key={s.id}
            cx={C + s.x * SCALE}
            cy={C - s.y * SCALE}
            r={inside ? 5 : 3.5}
            fill={getSearchCatalogColor(s.catalog)}
            fillOpacity={inside ? 1 : 0.35}
            stroke={kept.has(s.id) ? "#fff" : "none"}
            strokeWidth={1.5}
          />
        );
      })}
      <path
        d={`M${C} ${C - 12}V${C - 6}M${C} ${C + 6}V${C + 12}M${C - 12} ${C}H${C - 6}M${C + 6} ${C}H${C + 12}`}
        stroke="#fff"
        strokeWidth={1.5}
      />
      <line
        x1={16}
        x2={16 + 10 * SCALE}
        y1={384}
        y2={384}
        stroke="#868a96"
        strokeWidth={2}
      />
      <text x={16} y={376} fill="#868a96" fontSize={11} fontFamily="monospace">
        10″
      </text>
    </svg>
  );

  const side = (
    <>
      <div className="flex gap-6">
        <Stat value={ranked.length} label="inside the radius" />
        <Stat value={sources.length} label="sources in the field" />
      </div>
      <div>
        <div className="mb-2 text-sm font-semibold text-foreground">
          Result, nearest first
        </div>
        <div
          className="flex flex-col gap-1 font-mono text-xs"
          aria-live="polite"
        >
          {ranked.length === 0 && (
            <span className="text-neutral-400">
              Nothing inside the radius. Try a larger one.
            </span>
          )}
          {ranked.slice(0, 10).map((s, i) => (
            <div
              key={s.id}
              className={`grid grid-cols-[12px_1fr_auto] items-center gap-2 rounded bg-surface-elevated px-2 py-1 ${
                i >= nearest ? "opacity-50" : ""
              }`}
            >
              <svg width={9} height={9} aria-hidden="true">
                <circle
                  cx={4.5}
                  cy={4.5}
                  r={4.5}
                  fill={getSearchCatalogColor(s.catalog)}
                />
              </svg>
              <span className="text-foreground">
                {getSearchCatalogLabel(s.catalog)} {s.id}
              </span>
              <span className="text-neutral-400">
                {formatNumber(s.d, 2)}″{i >= nearest ? " · beyond N" : ""}
              </span>
            </div>
          ))}
        </div>
      </div>
      <p className="m-0 text-xs text-neutral-500">
        Example sources. Each catalog has its own typical position error (Gaia
        ~0.1″, AllWISE ~0.5″, eROSITA ~5″), which is why XWave suggests a
        different radius per catalog.
      </p>
    </>
  );

  const controls = (
    <>
      <label className="flex items-center gap-3">
        Radius
        <Slider
          className="!w-36"
          min={1}
          max={18}
          step={0.5}
          value={radius}
          onChange={setRadius}
          tooltip={{ open: false }}
          aria-label="Search radius in arcseconds"
        />
        <span className="w-10 font-mono text-foreground">{radius}″</span>
      </label>
      <label className="flex items-center gap-3">
        Neighbours (N)
        <Slider
          className="!w-28"
          min={1}
          max={10}
          value={nearest}
          onChange={setNearest}
          tooltip={{ open: false }}
          aria-label="Number of neighbours to return"
        />
        <span className="w-6 font-mono text-foreground">{nearest}</span>
      </label>
    </>
  );

  return <DiagramFrame stage={stage} side={side} controls={controls} />;
}
