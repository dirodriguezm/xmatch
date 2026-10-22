"use client";

import { Button, Slider } from "antd";
import { useEffect, useMemo, useState } from "react";

import {
  type CatalogOption,
  getSearchCatalogColor,
} from "@/app/lib/constants/catalogs";
import { MAX_BULK_ROWS } from "@/app/lib/utils/bulkInput";

import {
  DiagramFrame,
  formatNumber,
  pickCatalog,
  seededRandom,
  Stat,
} from "./shared";

const W = 880;
const H = 380;
const FW = 7200; // field width in arcsec (2°)
const FH = (FW * H) / W;
const S = W / FW;
const SKY_ARCSEC2 = 41_252.96 * 3600 * 3600;
const DUPLICATE = "#ffc53d";

interface Star {
  x: number;
  y: number;
  catalog: CatalogOption;
}
interface Point {
  x: number;
  y: number;
}
type Mode = "idle" | "bulk" | "tile";

function exampleSky() {
  const rnd = seededRandom(2024);
  const stars: Star[] = Array.from({ length: 1400 }, () => ({
    x: rnd() * FW,
    y: rnd() * FH,
    catalog: pickCatalog(rnd()),
  }));
  // A target list of 60 positions, most of them on a catalog source.
  const targets: Point[] = Array.from({ length: 60 }, () => {
    if (rnd() < 0.8) {
      const s = stars[Math.floor(rnd() * stars.length)];
      return { x: s.x + (rnd() - 0.5) * 20, y: s.y + (rnd() - 0.5) * 20 };
    }
    return { x: 300 + rnd() * (FW - 600), y: 300 + rnd() * (FH - 600) };
  });
  return { stars, targets };
}

/** Circles laid out on a hexagonal grid so the field has no gaps. */
function hexTiles(radius: number): Point[] {
  const out: Point[] = [];
  const dx = radius * Math.sqrt(3);
  const dy = radius * 1.5;
  for (let row = 0, y = 0; y < FH + dy; row++, y += dy)
    for (let x = ((row % 2) * dx) / 2; x < FW + dx; x += dx) out.push({ x, y });
  return out;
}

const percent = (v: number) =>
  v >= 10
    ? `${formatNumber(v)}%`
    : v >= 0.1
      ? `${formatNumber(v, 1)}%`
      : "< 0.1%";

function usePrefersReducedMotion() {
  const [reduce] = useState(
    () =>
      typeof window !== "undefined" &&
      window.matchMedia("(prefers-reduced-motion: reduce)").matches
  );
  return reduce;
}

/** Part 6: a bulk run covers small circles around a list, not the sky between them. */
export function BulkCoverageDiagram() {
  const { stars, targets } = useMemo(() => exampleSky(), []);
  const [radius, setRadius] = useState(120);
  const [mode, setMode] = useState<Mode>("idle");
  const [shown, setShown] = useState(0);
  const reduce = usePrefersReducedMotion();

  const circles = useMemo(
    () => (mode === "tile" ? hexTiles(radius) : targets),
    [mode, radius, targets]
  );

  // Advance the animation one batch per frame until every circle is shown.
  useEffect(() => {
    if (mode === "idle" || shown >= circles.length) return;
    const batch =
      mode === "tile" ? Math.max(1, Math.ceil(circles.length / 140)) : 1;
    const timer = window.setTimeout(
      () => setShown((n) => Math.min(circles.length, n + batch)),
      mode === "tile" ? 12 : 70
    );
    return () => window.clearTimeout(timer);
  }, [mode, shown, circles.length]);

  const start = (next: Mode) => {
    setMode(next);
    const total = next === "tile" ? hexTiles(radius).length : targets.length;
    setShown(reduce ? total : 0);
  };

  // How many visible circles each sample source falls in.
  const hitCounts = useMemo(() => {
    const counts = new Array<number>(stars.length).fill(0);
    for (const c of circles.slice(0, mode === "idle" ? 0 : shown))
      stars.forEach((s, k) => {
        if (Math.hypot(s.x - c.x, s.y - c.y) <= radius) counts[k]++;
      });
    return counts;
  }, [circles, shown, stars, radius, mode]);

  const found = hitCounts.filter((n) => n > 0).length;
  const repeated = hitCounts.filter((n) => n > 1).length;
  const visible = mode === "idle" ? [] : circles.slice(0, shown);
  const covered = Math.min(
    100,
    (100 * visible.length * Math.PI * radius * radius) / (FW * FH)
  );
  const skyPositions =
    SKY_ARCSEC2 / (((3 * Math.sqrt(3)) / 2) * radius * radius);

  const stage = (
    <svg
      viewBox={`0 0 ${W} ${H}`}
      className="h-auto w-full"
      role="img"
      aria-label="A 2 by 1.4 degree field with example catalog sources, a list of target positions and their search circles"
    >
      <rect width={W} height={H} fill="#0b0c10" />
      {stars.map((s, k) => {
        const n = hitCounts[k];
        return (
          <circle
            key={k}
            cx={s.x * S}
            cy={s.y * S}
            r={n ? 2.4 : 1.4}
            fill={
              n > 1
                ? DUPLICATE
                : n
                  ? getSearchCatalogColor(s.catalog)
                  : "#4b4f5c"
            }
            fillOpacity={n ? 1 : 0.55}
          />
        );
      })}
      {visible.map((c, i) => (
        <circle
          key={i}
          cx={c.x * S}
          cy={c.y * S}
          r={radius * S}
          fill={
            mode === "tile" ? "rgba(255,120,117,.06)" : "rgba(64,150,255,.10)"
          }
          stroke={mode === "tile" ? "#ff7875" : "#4096ff"}
          strokeOpacity={mode === "tile" ? 0.45 : 0.9}
        />
      ))}
      {mode !== "tile" &&
        targets.map((t, i) => {
          const on = mode === "bulk" && i < shown;
          return (
            <path
              key={i}
              d={`M${t.x * S - 4} ${t.y * S}H${t.x * S + 4}M${t.x * S} ${t.y * S - 4}V${t.y * S + 4}`}
              stroke={on ? "#fff" : "#9aa0b0"}
              strokeWidth={1.4}
              strokeOpacity={on ? 1 : 0.7}
            />
          );
        })}
      <line
        x1={16}
        x2={16 + 1800 * S}
        y1={H - 14}
        y2={H - 14}
        stroke="#868a96"
        strokeWidth={2}
      />
      <text x={16} y={H - 20} fontSize={11} fill="#868a96">
        30′
      </text>
    </svg>
  );

  const side = (
    <div className="grid gap-4 md:grid-cols-2">
      <div className="flex flex-wrap gap-6">
        {mode === "tile" ? (
          <>
            <Stat
              value={`${shown}/${circles.length}`}
              label="circles to cover this field"
            />
            <Stat value={found} label="sources found" />
            <Stat value={repeated} label="repeated (in 2+ circles)" />
          </>
        ) : (
          <>
            <Stat
              value={`${mode === "bulk" ? shown : 0}/${targets.length}`}
              label="positions searched"
            />
            <Stat value={found} label="sources found" />
            <Stat value={percent(covered)} label="of this field read" />
          </>
        )}
      </div>
      <p className="m-0 text-sm text-neutral-300" aria-live="polite">
        {mode === "idle" && (
          <>
            The crosses are your positions; the grey dots are what the catalogs
            hold. Press <strong>Run bulk</strong>.
          </>
        )}
        {mode === "bulk" && (
          <>
            Each position read only its own {radius}″ circle. The coloured
            sources fell inside: those are your matches. The rest of the field,
            almost all of it, <strong>was never read</strong>.
          </>
        )}
        {mode === "tile" && (
          <>
            To leave no gaps the circles must overlap. This 2° × 1.4° field
            alone needs <strong>{circles.length} circles</strong>, and sources
            in the overlaps{" "}
            <span className="text-[#ffc53d]">come back repeated</span>. It is
            reading the whole field, in pieces.
          </>
        )}
      </p>
    </div>
  );

  const controls = (
    <>
      <Button type="primary" onClick={() => start("bulk")}>
        Run bulk
      </Button>
      <Button onClick={() => start("tile")}>Try to cover this field</Button>
      <Button
        onClick={() => {
          setMode("idle");
          setShown(0);
        }}
      >
        Reset
      </Button>
      <label className="flex items-center gap-3">
        Radius
        <Slider
          className="!w-36"
          min={30}
          max={120}
          step={10}
          value={radius}
          onChange={(v) => {
            setRadius(v);
            if (mode === "tile") setShown(hexTiles(v).length);
          }}
          tooltip={{ open: false }}
          aria-label="Search radius in arcseconds"
        />
        <span className="w-10 font-mono text-foreground">{radius}″</span>
      </label>
    </>
  );

  return (
    <>
      <DiagramFrame stage={stage} side={side} controls={controls} wide />
      <div className="rounded-lg border border-border bg-surface p-4">
        <div className="text-sm font-semibold text-foreground">
          What if I run bulks until the whole sky is covered?
        </div>
        <p className="m-0 mt-2 text-sm text-neutral-400">
          It works, but it saves nothing. With {radius}″ circles, covering the
          whole sky without gaps takes about{" "}
          <strong className="text-foreground">
            {formatNumber(skyPositions / 1e6, skyPositions < 1e7 ? 1 : 0)}{" "}
            million positions
          </strong>
          , or about{" "}
          <strong className="text-foreground">
            {formatNumber(skyPositions / MAX_BULK_ROWS)} bulk runs
          </strong>{" "}
          of {formatNumber(MAX_BULK_ROWS)}. In the end the whole catalog is read
          anyway, plus the repeats from the overlaps. That job is not a
          cross-match: download the catalog or run a catalog query (for example,
          TAP on the Gaia archive).
        </p>
      </div>
    </>
  );
}
