"use client";

import { type ReactNode, useId } from "react";

/**
 * Loading indicators, one per context so the motion hints at what is coming:
 *
 * - `arcs`: the default — one arc per catalog. Legible down to 14px, so it
 *   is the one for buttons, inline states and small panels.
 * - `xwave`: the logo's two waves. First load of a page (Suspense fallbacks).
 * - `crossmatch`: two detections merging on the reticle. Cone search results.
 * - `healpix`: a base pixel lighting up in nested (Z-order) order. Bulk runs.
 * - `sed`: blue-to-red bars shaped like an SED. Photometry.
 * - `lightcurve`: a trace drawn with error-barred points. Light curve panels.
 *
 * Only `arcs` reads below 32px. Never show two different variants on screen
 * at once; the lesser one uses `arcs`. Keyframes live in globals.css (`xws-`),
 * including the reduced-motion fallback.
 */
export type XWaveSpinnerVariant =
  | "arcs"
  | "xwave"
  | "crossmatch"
  | "healpix"
  | "sed"
  | "lightcurve";

interface XWaveSpinnerProps {
  variant?: XWaveSpinnerVariant;
  /** Rendered width and height in px. */
  size?: number;
  /** Visible text under the spinner; also what screen readers announce. */
  label?: ReactNode;
  className?: string;
}

const GAIA = "#4096ff";
const ALLWISE = "#9254de";
const EROSITA = "#eb2f96";
const FAINT = "#3a3d47";

const f = (n: number) => n.toFixed(2);

// xwave: two periods of a sine, slid by one period so the loop is seamless.
function sinePath(amp: number, phase: number) {
  let d = "";
  for (let x = 0; x <= 96; x += 1.5) {
    const y = 32 + amp * Math.sin((x / 32) * 2 * Math.PI + phase);
    d += `${x ? "L" : "M"}${f(x)} ${f(y)}`;
  }
  return d;
}
const WAVE_1 = sinePath(10, 0);
const WAVE_2 = sinePath(-8, 1.2);

// lightcurve: fixed points with error bars.
const LC_POINTS = Array.from({ length: 12 }, (_, i) => ({
  x: 6 + i * 4.7,
  y: 32 + 9 * Math.sin(i * 0.9) + ((i % 3) - 1) * 2.2,
  err: 2 + (i % 4),
}));
const LC_PATH = LC_POINTS.map(
  (p, i) => `${i ? "L" : "M"}${f(p.x)} ${f(p.y)}`
).join("");

// sed: heights follow a blackbody-ish hump, colors run blue to red.
const SED_BARS = [
  [0.55, "#4096ff"],
  [0.8, "#36cfc9"],
  [1, "#73d13d"],
  [0.9, "#fadb14"],
  [0.7, "#ffa940"],
  [0.48, "#ff7a45"],
  [0.3, "#f5222d"],
] as const;

// healpix: 4×4 cells of a base pixel, lit in nested (Morton) index order.
function morton(i: number, j: number) {
  let z = 0;
  for (let b = 0; b < 2; b++) {
    z |= (((i >> b) & 1) << (2 * b + 1)) | (((j >> b) & 1) << (2 * b));
  }
  return z;
}
// Sorted by nested index: the nth-child delays in globals.css follow it.
const HPX_CELLS = Array.from({ length: 16 }, (_, k) => {
  const i = Math.floor(k / 4);
  const j = k % 4;
  return { i, j, z: morton(i, j) };
}).sort((a, b) => a.z - b.z);

function Arcs() {
  const rings = [
    { r: 27, color: GAIA, cls: "xws-arc-1" },
    { r: 19.5, color: ALLWISE, cls: "xws-arc-2" },
    { r: 12, color: EROSITA, cls: "xws-arc-3" },
  ];
  return (
    <>
      {rings.map(({ r }) => (
        <circle
          key={`t${r}`}
          cx="32"
          cy="32"
          r={r}
          fill="none"
          stroke={FAINT}
          strokeOpacity={0.5}
          strokeWidth="4"
        />
      ))}
      {rings.map(({ r, color, cls }) => (
        <circle
          key={r}
          cx="32"
          cy="32"
          r={r}
          fill="none"
          stroke={color}
          strokeWidth="4.4"
          strokeLinecap="round"
          strokeDasharray="28 72"
          pathLength={100}
          className={cls}
        />
      ))}
    </>
  );
}

function Xwave({ id }: { id: string }) {
  return (
    <>
      <defs>
        <linearGradient id={`${id}w1`} x1="0" x2="1">
          <stop offset="0" stopColor="#6366f1" />
          <stop offset="1" stopColor="#a78bfa" />
        </linearGradient>
        <linearGradient id={`${id}w2`} x1="0" x2="1">
          <stop offset="0" stopColor="#06b6d4" />
          <stop offset="1" stopColor="#3b82f6" />
        </linearGradient>
        <clipPath id={`${id}wc`}>
          <circle cx="32" cy="32" r="27" />
        </clipPath>
      </defs>
      <circle cx="32" cy="32" r="28" fill="#12131c" stroke={FAINT} />
      <g clipPath={`url(#${id}wc)`}>
        <g transform="translate(-16 0)">
          <path
            d={WAVE_1}
            fill="none"
            stroke={`url(#${id}w1)`}
            strokeWidth="3.2"
            strokeLinecap="round"
            className="xws-wave-1"
          />
          <path
            d={WAVE_2}
            fill="none"
            stroke={`url(#${id}w2)`}
            strokeWidth="3.2"
            strokeLinecap="round"
            className="xws-wave-2"
          />
        </g>
      </g>
      <circle cx="32" cy="32" r="2.4" fill="#fff" />
    </>
  );
}

function Crossmatch() {
  return (
    <>
      <circle
        cx="32"
        cy="32"
        r="8"
        fill="none"
        stroke={FAINT}
        strokeDasharray="2 2"
      />
      <path
        d="M32 6V20M32 44V58M6 32H20M44 32H58"
        stroke={FAINT}
        strokeWidth="1.2"
      />
      <circle
        cx="32"
        cy="32"
        r="14"
        fill="none"
        stroke="#fff"
        strokeWidth="1.5"
        className="xws-merge-ring"
      />
      <circle cx="32" cy="32" r="4" fill={GAIA} className="xws-merge-a" />
      <circle
        cx="32"
        cy="32"
        r="4"
        fill={ALLWISE}
        fillOpacity={0.9}
        className="xws-merge-b"
      />
    </>
  );
}

function Healpix() {
  return (
    <g transform="rotate(45 32 32) translate(32 32) scale(.92) translate(-32 -32)">
      {HPX_CELLS.map(({ i, j, z }) => (
        <rect
          key={z}
          x={12.5 + j * 10}
          y={12.5 + i * 10}
          width="9"
          height="9"
          rx="1.2"
          fill={`hsl(${220 + z * 4} 90% 64%)`}
          className="xws-cell"
        />
      ))}
    </g>
  );
}

function Sed() {
  return (
    <>
      <path d="M4 56H60" stroke={FAINT} strokeWidth="1.2" />
      {SED_BARS.map(([h, color], i) => (
        <rect
          key={color}
          x={6 + i * 7.6}
          y={54 - 44 * h}
          width="5.2"
          height={44 * h}
          rx="1.5"
          fill={color}
          className="xws-bar"
        />
      ))}
    </>
  );
}

function Lightcurve({ id }: { id: string }) {
  return (
    <>
      <defs>
        {/* Reveals the points as the trace reaches them. */}
        <clipPath id={`${id}lc`}>
          <rect x="0" y="0" width="64" height="64" className="xws-lc-reveal" />
        </clipPath>
      </defs>
      <path d="M4 56H60M4 8V56" stroke={FAINT} strokeWidth="1.2" fill="none" />
      <path
        d={LC_PATH}
        pathLength={100}
        fill="none"
        stroke={GAIA}
        strokeOpacity={0.7}
        strokeWidth="1.6"
        strokeDasharray="100"
        className="xws-lc-draw"
      />
      <g clipPath={`url(#${id}lc)`} className="xws-lc-fade">
        {LC_POINTS.map(({ x, y, err }) => (
          <g key={x}>
            <line
              x1={x}
              x2={x}
              y1={y - err}
              y2={y + err}
              stroke="#8b8f9a"
              strokeWidth="1"
            />
            <circle cx={x} cy={y} r="1.9" fill={GAIA} />
          </g>
        ))}
      </g>
    </>
  );
}

export function XWaveSpinner({
  variant = "arcs",
  size = 32,
  label,
  className,
}: XWaveSpinnerProps) {
  const id = useId().replace(/:/g, "");

  return (
    <div
      role="status"
      aria-live="polite"
      className={`xws inline-flex flex-col items-center gap-3 ${className ?? ""}`}
    >
      <svg
        width={size}
        height={size}
        viewBox="0 0 64 64"
        aria-hidden="true"
        className="block shrink-0 overflow-visible"
      >
        {variant === "arcs" && <Arcs />}
        {variant === "xwave" && <Xwave id={id} />}
        {variant === "crossmatch" && <Crossmatch />}
        {variant === "healpix" && <Healpix />}
        {variant === "sed" && <Sed />}
        {variant === "lightcurve" && <Lightcurve id={id} />}
      </svg>
      {label ? (
        <span className="text-sm text-foreground/60">{label}</span>
      ) : (
        <span className="sr-only">Loading</span>
      )}
    </div>
  );
}
