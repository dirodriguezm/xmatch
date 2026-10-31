"use client";

import { Button, Slider } from "antd";
import {
  type MouseEvent,
  type ReactNode,
  useMemo,
  useRef,
  useState,
} from "react";

import {
  type CatalogOption,
  getSearchCatalogColor,
} from "@/app/lib/constants/catalogs";

import { interleave, toRanges } from "./healpix";
import { DiagramFrame, HIT, pickCatalog, seededRandom, Stat } from "./shared";

const GRID = 16; // 16 × 16 order-18 cells = one order-14 cell
const CELL = 0.8; // arcsec, order 18
const PX = 400 / GRID;
const OFF = 10;
const toPx = (arcsec: number) => OFF + (arcsec / CELL) * PX;
/** Local NESTED number of a cell inside the patch (column on even bits). */
const local = (row: number, col: number) => interleave(col, row, 4);

interface Source {
  x: number;
  y: number;
  catalog: CatalogOption;
  cell: number;
}

const STEPS: { title: string; body: ReactNode }[] = [
  {
    title: "Indexing: every source gets its cell",
    body: (
      <>
        When a catalog is loaded, each source&apos;s position becomes the number
        of the order-18 HEALPix cell that contains it. That number is stored
        with the source as <code>ipix</code>, in an indexed column. It is
        computed once, not on every search.
      </>
    ),
  },
  {
    title: "Your search: a position and a radius",
    body: (
      <>
        You ask for a position and a radius. The service does not look at any
        source yet: first it works out which part of the sky to read.
      </>
    ),
  },
  {
    title: "Cells the circle touches",
    body: (
      <>
        HEALPix lists every cell the circle touches, even at a corner (an{" "}
        <em>inclusive</em> query). It may include a few too many, but it never
        leaves out a source that is inside. Thanks to NESTED numbering those
        cells form a few <strong>ranges</strong> of consecutive numbers: see the
        line below.
      </>
    ),
  },
  {
    title: "Fetch candidates by range",
    body: (
      <>
        The database reads the sources whose <code>ipix</code> falls in each
        range. The column is indexed, so each range is a direct read. Only the
        sources in those cells come back: the candidates.
      </>
    ),
  },
  {
    title: "Measure for real and rank",
    body: (
      <>
        Each candidate gets its true distance on the sphere (haversine). Those
        beyond the radius are dropped: they came from cells the circle only
        grazed. The rest are sorted nearest first, cut to N, and returned with
        their separation and the total per catalog.
      </>
    ),
  },
];

/** Part 4: the service's cone search on a real-scale patch of order-18 cells. */
export function IpixPipelineDiagram() {
  const [step, setStep] = useState(0);
  const [radius, setRadius] = useState(2);
  const [center, setCenter] = useState({ x: 6.6, y: 6.0 });
  const svgRef = useRef<SVGSVGElement>(null);

  const sources = useMemo<Source[]>(() => {
    const rnd = seededRandom(42);
    return Array.from({ length: 44 }, () => {
      const catalog = pickCatalog(rnd());
      const x = rnd() * GRID * CELL;
      const y = rnd() * GRID * CELL;
      const col = Math.min(GRID - 1, Math.floor(x / CELL));
      const row = Math.min(GRID - 1, Math.floor(y / CELL));
      return { x, y, catalog, cell: local(row, col) };
    });
  }, []);

  const { selected, ranges } = useMemo(() => {
    const sel = new Set<number>();
    for (let row = 0; row < GRID; row++)
      for (let col = 0; col < GRID; col++) {
        const x0 = col * CELL;
        const y0 = row * CELL;
        const nx = Math.max(x0, Math.min(center.x, x0 + CELL));
        const ny = Math.max(y0, Math.min(center.y, y0 + CELL));
        if (Math.hypot(nx - center.x, ny - center.y) <= radius)
          sel.add(local(row, col));
      }
    return { selected: sel, ranges: toRanges([...sel].sort((a, b) => a - b)) };
  }, [center, radius]);

  const candidates = sources.filter((s) => selected.has(s.cell));
  const hits = candidates.filter(
    (s) => Math.hypot(s.x - center.x, s.y - center.y) <= radius
  );

  const moveCenter = (e: MouseEvent<SVGSVGElement>) => {
    const svg = svgRef.current;
    const ctm = svg?.getScreenCTM();
    if (!svg || !ctm) return;
    const pt = svg.createSVGPoint();
    pt.x = e.clientX;
    pt.y = e.clientY;
    const p = pt.matrixTransform(ctm.inverse());
    const clamp = (n: number) => Math.max(0, Math.min(GRID * CELL, n));
    setCenter({
      x: clamp(((p.x - OFF) / PX) * CELL),
      y: clamp(((p.y - OFF) / PX) * CELL),
    });
    if (step < 1) setStep(1);
  };

  const cx = toPx(center.x);
  const cy = toPx(center.y);
  const lineW = 400 / 256;

  const stage = (
    <>
      <svg
        ref={svgRef}
        viewBox="0 0 420 420"
        className="h-auto w-full max-w-[420px] cursor-crosshair"
        role="img"
        aria-label="A 16 by 16 patch of order-18 HEALPix cells with example sources, the search circle and the selected cells"
        onClick={moveCenter}
      >
        {Array.from({ length: GRID * GRID }, (_, i) => {
          const row = Math.floor(i / GRID);
          const col = i % GRID;
          const z = local(row, col);
          const on = step >= 2 && selected.has(z);
          return (
            <g key={i}>
              <rect
                x={OFF + col * PX + 0.5}
                y={OFF + row * PX + 0.5}
                width={PX - 1}
                height={PX - 1}
                rx={2}
                fill={on ? "rgba(54,207,201,.22)" : "#15171d"}
                stroke={on ? HIT : "#22252d"}
                strokeWidth={on ? 1 : 0.8}
              />
              <text
                x={OFF + col * PX + PX / 2}
                y={OFF + row * PX + PX / 2 + 3}
                textAnchor="middle"
                fontSize={7.5}
                fontFamily="monospace"
                fill={on ? "#bff3ef" : "#3d4150"}
              >
                {z}
              </text>
            </g>
          );
        })}
        {step >= 1 && (
          <>
            <circle
              cx={cx}
              cy={cy}
              r={(radius / CELL) * PX}
              fill="rgba(64,150,255,.07)"
              stroke="#4096ff"
              strokeWidth={1.6}
              strokeDasharray="5 4"
            />
            <path
              d={`M${cx} ${cy - 9}V${cy - 4}M${cx} ${cy + 4}V${cy + 9}M${cx - 9} ${cy}H${cx - 4}M${cx + 4} ${cy}H${cx + 9}`}
              stroke="#fff"
              strokeWidth={1.6}
            />
          </>
        )}
        {sources.map((s, i) => {
          const isCandidate = selected.has(s.cell);
          const isHit =
            isCandidate && Math.hypot(s.x - center.x, s.y - center.y) <= radius;
          const x = toPx(s.x);
          const y = toPx(s.y);
          return (
            <g key={i}>
              <circle
                cx={x}
                cy={y}
                r={step >= 3 && isCandidate ? 5 : 4}
                fill={getSearchCatalogColor(s.catalog)}
                fillOpacity={step >= 3 && !isCandidate ? 0.18 : 1}
                stroke={step === 4 && isHit ? "#fff" : "none"}
                strokeWidth={1.5}
              />
              {step === 4 && isCandidate && !isHit && (
                <path
                  d={`M${x - 6} ${y - 6}L${x + 6} ${y + 6}M${x + 6} ${y - 6}L${x - 6} ${y + 6}`}
                  stroke="#ff7875"
                  strokeWidth={2}
                />
              )}
              {step === 0 && (
                <text
                  x={x + 6}
                  y={y - 5}
                  fontSize={9}
                  fontFamily="monospace"
                  fill="#e8e9ec"
                >
                  {s.cell}
                </text>
              )}
            </g>
          );
        })}
      </svg>
      <svg
        viewBox="0 0 420 64"
        className="h-auto w-full max-w-[420px]"
        role="img"
        aria-label={
          step >= 2
            ? `${selected.size} cells form ${ranges.length} ranges of ipix`
            : "ipix of the 256 cells in the patch"
        }
      >
        <rect
          x={10}
          y={14}
          width={400}
          height={16}
          rx={3}
          fill="#15171d"
          stroke="#22252d"
        />
        {step >= 2 && (
          <>
            {[...selected].map((z) => (
              <rect
                key={z}
                x={10 + z * lineW}
                y={14}
                width={Math.max(lineW, 1.2)}
                height={16}
                fill={HIT}
                fillOpacity={0.9}
              />
            ))}
            {ranges.map(([a, b]) => (
              <rect
                key={a}
                x={10 + a * lineW - 0.5}
                y={11}
                width={(b - a + 1) * lineW + 1}
                height={22}
                rx={2}
                fill="none"
                stroke="#e8e9ec"
              />
            ))}
          </>
        )}
        {[0, 64, 128, 192, 255].map((t) => (
          <text
            key={t}
            x={10 + t * lineW + (t === 255 ? lineW : 0)}
            y={48}
            textAnchor={t === 0 ? "start" : t === 255 ? "end" : "middle"}
            fontSize={10}
            fill="#868a96"
          >
            {t}
          </text>
        ))}
        <text x={210} y={62} textAnchor="middle" fontSize={10} fill="#868a96">
          {step >= 2
            ? `${selected.size} cells → ${ranges.length} range${ranges.length === 1 ? "" : "s"} of ipix`
            : "ipix of the 256 cells in the patch"}
        </text>
      </svg>
    </>
  );

  const side = (
    <>
      <ol className="m-0 flex list-none flex-col gap-1 p-0">
        {STEPS.map((s, i) => (
          <li key={s.title}>
            <button
              type="button"
              onClick={() => setStep(i)}
              aria-current={i === step ? "step" : undefined}
              className={`grid w-full cursor-pointer grid-cols-[20px_1fr] gap-2 rounded-md border-0 px-2 py-1.5 text-left text-sm ${
                i === step
                  ? "bg-surface-elevated text-foreground ring-1 ring-border"
                  : "bg-transparent text-neutral-400 hover:text-neutral-200"
              }`}
            >
              <span className="font-mono text-xs leading-6">{i + 1}</span>
              <span>{s.title}</span>
            </button>
          </li>
        ))}
      </ol>
      <p
        className="m-0 min-h-[7.5em] text-sm text-neutral-300"
        aria-live="polite"
      >
        {STEPS[step]?.body}
      </p>
      <div className="flex flex-wrap gap-6">
        {step >= 2 && <Stat value={selected.size} label="cells" />}
        {step >= 2 && <Stat value={ranges.length} label="ranges" />}
        {step >= 3 && (
          <Stat
            value={`${candidates.length}/${sources.length}`}
            label="candidates"
          />
        )}
        {step >= 4 && <Stat value={hits.length} label="inside the radius" />}
      </div>
      <div className="flex gap-2">
        <Button
          onClick={() => setStep((s) => Math.max(0, s - 1))}
          disabled={step === 0}
        >
          Back
        </Button>
        <Button
          type="primary"
          onClick={() => setStep((s) => Math.min(STEPS.length - 1, s + 1))}
          disabled={step === STEPS.length - 1}
        >
          Next
        </Button>
      </div>
    </>
  );

  const controls = (
    <>
      <label className="flex items-center gap-3">
        Radius
        <Slider
          className="!w-36"
          min={0.6}
          max={4}
          step={0.1}
          value={radius}
          onChange={setRadius}
          tooltip={{ open: false }}
          aria-label="Search radius in arcseconds"
        />
        <span className="w-10 font-mono text-foreground">
          {radius.toFixed(1)}″
        </span>
      </label>
      <span className="text-xs text-neutral-500">
        Click the patch to move the search position
      </span>
    </>
  );

  return <DiagramFrame stage={stage} side={side} controls={controls} />;
}
