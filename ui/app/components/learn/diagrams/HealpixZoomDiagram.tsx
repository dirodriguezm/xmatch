"use client";

import { Button } from "antd";
import { type ReactNode, useMemo, useState } from "react";

import { angToNest, childDigits, FACE_CENTERS, QUARTER } from "./healpix";
import { DiagramFrame, formatNumber, HIT, Stat } from "./shared";

// θ¹ Ori F in the Trapezium; its Gaia DR3 source_id carries the order-12 pixel.
const STAR = { name: "θ¹ Ori F", ra: 83.81972, dec: -5.39034 };
const DEEP = 18;
const LEVELS = [0, 1, 2, 3, 14, 18];
const STAR_COLOR = "#ffc53d";

const VIEW = 420;
const CX = 210;
const CY = 210;
const HALF = 175;
/** Local (a, b) ∈ [0, 1]² of the zoomed cell → view coordinates (a north-east, b north-west). */
const P = (a: number, b: number) =>
  [CX + (a - b) * HALF, CY + HALF - (a + b) * HALF] as const;
const diamond = (a0: number, b0: number, s: number) =>
  "M" +
  [P(a0, b0), P(a0 + s, b0), P(a0 + s, b0 + s), P(a0, b0 + s)]
    .map((p) => p.join(" "))
    .join("L") +
  "Z";

/** Part 3: zoom down the HEALPix tree to the order-18 cell of one real star. */
export function HealpixZoomDiagram() {
  const [step, setStep] = useState(0);
  const deep = useMemo(() => angToNest(STAR.ra, STAR.dec, DEEP), []);
  const digits = useMemo(() => childDigits(deep, DEEP), [deep]);
  const ipixAt = (k: number) => Math.floor(deep.ipix / 4 ** (DEEP - k));
  const nside = 2 ** DEEP;
  const u = (deep.ix + 0.5) / nside;
  const v = (deep.iy + 0.5) / nside;
  const cellOrigin = (k: number) => {
    const n = 2 ** k;
    return [Math.floor(u * n) / n, Math.floor(v * n) / n, 1 / n] as const;
  };

  const k = LEVELS[step];
  let drawing: ReactNode;
  let caption: ReactNode;

  if (k === 0) {
    const sx = 380 / (2 * Math.PI);
    const sy = 190 / Math.PI;
    const X = (x: number) => 20 + x * sx;
    const Y = (y: number) => 210 - y * sy;
    const [fx, fy] = FACE_CENTERS[deep.face];
    const px = fx + (u - v) * QUARTER;
    const py = fy - QUARTER + (u + v) * QUARTER;
    drawing = (
      <>
        {FACE_CENTERS.flatMap(([cx, cy], f) =>
          (cx === 0 ? [0, 2 * Math.PI] : [0]).map((off) => {
            const on = f === deep.face;
            return (
              <g key={`${f}-${off}`}>
                <path
                  d={`M${X(cx + off)} ${Y(cy - QUARTER)}L${X(cx + off + QUARTER)} ${Y(cy)}L${X(cx + off)} ${Y(cy + QUARTER)}L${X(cx + off - QUARTER)} ${Y(cy)}Z`}
                  fill={on ? "rgba(54,207,201,.25)" : "#1f2433"}
                  stroke={on ? HIT : "#3a4258"}
                  strokeWidth={on ? 2 : 1}
                />
                <text
                  x={X(cx + off)}
                  y={Y(cy) + 5}
                  textAnchor="middle"
                  fontSize={14}
                  fontWeight={600}
                  fill={on ? "#e8e9ec" : "#6b7080"}
                >
                  {f}
                </text>
              </g>
            );
          })
        )}
        <circle cx={X(px)} cy={Y(py)} r={4} fill={STAR_COLOR} />
      </>
    );
    caption = (
      <>
        The whole sky in 12 base cells. {STAR.name} (the yellow dot, RA{" "}
        {STAR.ra.toFixed(2)}°, Dec {STAR.dec.toFixed(2)}°) falls in{" "}
        <strong>cell {deep.face}</strong>, one of the four equatorial ones.
        First digit of its address: <strong>{deep.face}</strong>.
      </>
    );
  } else if (k <= 3) {
    const [u0, v0, s] = cellOrigin(k - 1);
    const d = digits[k - 1];
    const [sx, sy] = P((u - u0) / s, (v - v0) / s);
    drawing = (
      <>
        {[0, 1, 2, 3].map((dd) => {
          const dx = dd % 2;
          const dy = Math.floor(dd / 2);
          const on = dd === d;
          const [lx, ly] = P(dx / 2 + 0.25, dy / 2 + 0.25);
          return (
            <g key={dd}>
              <path
                d={diamond(dx / 2, dy / 2, 0.5)}
                fill={on ? "rgba(54,207,201,.25)" : "#1a1d26"}
                stroke={on ? HIT : "#3a4258"}
                strokeWidth={on ? 2 : 1}
              />
              <text
                x={lx}
                y={ly + 10}
                textAnchor="middle"
                fontSize={30}
                fontWeight={600}
                fill={on ? "#e8e9ec" : "#4c5263"}
              >
                {dd}
              </text>
            </g>
          );
        })}
        <circle cx={sx} cy={sy} r={5} fill={STAR_COLOR} />
      </>
    );
    caption = (
      <>
        Zoom into cell {formatNumber(ipixAt(k - 1))} (order {k - 1}). It splits
        into 4 children and the star falls in <strong>child {d}</strong>. Its
        number at order {k} is{" "}
        <strong>
          {formatNumber(ipixAt(k - 1))} × 4 + {d} = {formatNumber(ipixAt(k))}
        </strong>
        .
      </>
    );
  } else {
    const [u0, v0, s] = cellOrigin(14);
    const [sx, sy] = P((u - u0) / s, (v - v0) / s);
    const local = deep.ipix % 256;
    const grid = Array.from({ length: 15 }, (_, i) => (i + 1) / 16);
    drawing = (
      <>
        <path
          d={diamond(0, 0, 1)}
          fill="#1a1d26"
          stroke={k === 14 ? HIT : "#3a4258"}
          strokeWidth={k === 14 ? 2 : 1.5}
        />
        {grid.map((t) => (
          <g key={t} stroke="#2c3140" strokeWidth={0.8}>
            <line
              x1={P(t, 0)[0]}
              y1={P(t, 0)[1]}
              x2={P(t, 1)[0]}
              y2={P(t, 1)[1]}
            />
            <line
              x1={P(0, t)[0]}
              y1={P(0, t)[1]}
              x2={P(1, t)[0]}
              y2={P(1, t)[1]}
            />
          </g>
        ))}
        {k === 18 && (
          <path
            d={diamond((deep.ix % 16) / 16, (deep.iy % 16) / 16, 1 / 16)}
            fill="rgba(54,207,201,.45)"
            stroke={HIT}
            strokeWidth={2}
          />
        )}
        <circle cx={sx} cy={sy} r={4} fill={STAR_COLOR} />
      </>
    );
    caption =
      k === 14 ? (
        <>
          Eleven more levels lead to order 14: a 12.8″ cell, like the patch
          below. On the way, at order 12, comes{" "}
          <strong>{formatNumber(ipixAt(12))}</strong>, the same number Gaia
          stores inside this star&apos;s <code>source_id</code>. This cell
          splits into 16 × 16 cells of order 18.
        </>
      ) : (
        <>
          Four more levels reach the 0.8″ cell of order 18. Its number,{" "}
          <strong>{formatNumber(deep.ipix)}</strong>, is the <code>ipix</code>{" "}
          XWave stores with the star. Inside the patch it is cell{" "}
          <strong>{local}</strong>: the last 4 digits of the address.
        </>
      );
  }

  const shown = digits.slice(0, k);
  const address = (
    <span className="break-all font-mono text-base text-foreground">
      <span className="text-[#36cfc9]">{deep.face}</span>
      {k > 0 && " · "}
      {shown.slice(0, 3).join(" ")}
      {k > 3 && (
        <span className="text-neutral-500">
          {" "}
          {shown.slice(3, 14).join(" ")}
        </span>
      )}
      {k === 18 && (
        <span className="text-[#36cfc9]"> {shown.slice(14).join(" ")}</span>
      )}
    </span>
  );

  const stage = (
    <svg
      viewBox={`0 0 ${VIEW} ${VIEW}`}
      className="h-auto w-full max-w-[420px]"
      role="img"
      aria-label={`HEALPix zoom on ${STAR.name} at order ${k}`}
    >
      {drawing}
    </svg>
  );

  const side = (
    <>
      <span className="font-mono text-xs text-neutral-500">
        Order {k} · {formatNumber(12 * 4 ** k)} cells
      </span>
      <p
        className="m-0 min-h-[8.5em] text-sm text-neutral-300"
        aria-live="polite"
      >
        {caption}
      </p>
      <div className="flex flex-col gap-1">
        <span className="text-xs text-neutral-500">Address (base 4)</span>
        {address}
      </div>
      <Stat
        value={formatNumber(ipixAt(k))}
        label={`cell number (ipix) at order ${k}`}
      />
      <div className="flex gap-2">
        <Button
          onClick={() => setStep((s) => Math.max(0, s - 1))}
          disabled={step === 0}
        >
          Back
        </Button>
        <Button
          type="primary"
          onClick={() => setStep((s) => Math.min(LEVELS.length - 1, s + 1))}
          disabled={step === LEVELS.length - 1}
        >
          Next
        </Button>
      </div>
    </>
  );

  return <DiagramFrame stage={stage} side={side} />;
}
