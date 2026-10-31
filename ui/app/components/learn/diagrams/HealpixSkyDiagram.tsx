"use client";

import { Segmented, Switch } from "antd";
import { useId, useMemo, useState } from "react";

import {
  FACE_CENTERS,
  facePosition,
  nestIndex,
  QUARTER,
  ringIndex,
} from "./healpix";
import { DiagramFrame, formatNumber, HIT, MUTED } from "./shared";

const W = 860;
const H = 470;
const PAD = 30;
const SX = (W - 2 * PAD) / (2 * Math.PI);
const SY = (H - 2 * PAD) / Math.PI;
const X = (x: number) => PAD + x * SX;
const Y = (y: number) => H / 2 - y * SY;
const TWO_PI = 2 * Math.PI;
const wrap = (x: number) => ((x % TWO_PI) + TWO_PI) % TWO_PI;

const diamond = (cx: number, cy: number, h: number) =>
  `M${X(cx)} ${Y(cy - h)}L${X(cx + h)} ${Y(cy)}L${X(cx)} ${Y(cy + h)}L${X(cx - h)} ${Y(cy)}Z`;

const faceFill = (f: number) =>
  f < 4 ? "#2a3550" : f < 8 ? "#2b2a48" : "#352a45";

type Scheme = "nest" | "ring";

/** Part 2: the whole sky in the HEALPix projection, by order and numbering scheme. */
export function HealpixSkyDiagram() {
  const [order, setOrder] = useState(1);
  const [scheme, setScheme] = useState<Scheme>("nest");
  const [showPath, setShowPath] = useState(true);
  const clipId = useId();

  const nside = 2 ** order;
  const h = QUARTER / nside;

  const cells = useMemo(() => {
    const out: { key: string; x: number; y: number; id: number }[] = [];
    for (let f = 0; f < 12; f++)
      for (let ix = 0; ix < nside; ix++)
        for (let iy = 0; iy < nside; iy++) {
          const [x, y] = facePosition(
            f,
            (ix + 0.5) / nside,
            (iy + 0.5) / nside
          );
          out.push({
            key: `${f}-${ix}-${iy}`,
            x: wrap(x),
            y,
            id:
              scheme === "nest"
                ? nestIndex(f, ix, iy, order)
                : ringIndex(f, ix, iy, order),
          });
        }
    return out;
  }, [nside, order, scheme]);

  const path = useMemo(() => {
    const sorted = [...cells].sort((a, b) => a.id - b.id);
    return sorted.slice(1).flatMap((b, i) => {
      const a = sorted[i];
      if (!a) return [];
      const jump = Math.hypot(a.x - b.x, a.y - b.y) > 3.3 * h;
      return [{ a, b, jump }];
    });
  }, [cells, h]);

  const fontSize = nside === 1 ? 20 : nside === 2 ? 14 : nside === 4 ? 9.5 : 0;
  // The equatorial face at φ = 0 straddles the seam; draw it again at 2π.
  const copies = (x: number, reach: number) => (x <= reach ? [0, TWO_PI] : [0]);

  const stage = (
    <svg
      viewBox={`0 0 ${W} ${H}`}
      className="h-auto w-full"
      role="img"
      aria-label={`The whole sky in the HEALPix projection at order ${order}: ${formatNumber(12 * nside * nside)} cells numbered in the ${scheme === "nest" ? "NESTED" : "RING"} scheme`}
    >
      <defs>
        <clipPath id={clipId}>
          <rect
            x={X(0)}
            y={Y(Math.PI / 2)}
            width={X(TWO_PI) - X(0)}
            height={Y(-Math.PI / 2) - Y(Math.PI / 2)}
          />
        </clipPath>
      </defs>
      <g clipPath={`url(#${clipId})`}>
        {FACE_CENTERS.map(([cx, cy], f) =>
          copies(cx, 0).map((off) => (
            <path
              key={`${f}-${off}`}
              d={diamond(cx + off, cy, QUARTER)}
              fill={faceFill(f)}
              stroke="#4a5470"
              strokeWidth={1.5}
            />
          ))
        )}
        {nside > 1 &&
          cells.map((c) =>
            copies(c.x, h * 1.01).map((off) => (
              <path
                key={`${c.key}-${off}`}
                d={diamond(c.x + off, c.y, h)}
                fill="none"
                stroke="#3a4258"
                strokeWidth={0.8}
              />
            ))
          )}
        {showPath &&
          path.map(({ a, b, jump }) => (
            <line
              key={`${a.key}>${b.key}`}
              x1={X(a.x)}
              y1={Y(a.y)}
              x2={X(b.x)}
              y2={Y(b.y)}
              stroke={HIT}
              strokeOpacity={jump ? 0.18 : 0.85}
              strokeWidth={jump ? 1 : 1.6}
              strokeDasharray={jump ? "3 3" : undefined}
            />
          ))}
        {fontSize > 0 &&
          cells.map((c) =>
            copies(c.x, 1e-6).map((off) => (
              <text
                key={`t${c.key}-${off}`}
                x={X(c.x + off)}
                y={Y(c.y) + fontSize * 0.36}
                textAnchor="middle"
                fontSize={fontSize}
                fill="#e8e9ec"
                fontFamily="monospace"
                fontWeight={nside === 1 ? 600 : 500}
              >
                {c.id}
              </text>
            ))
          )}
      </g>
      <line
        x1={X(0)}
        x2={X(TWO_PI)}
        y1={Y(0)}
        y2={Y(0)}
        stroke={MUTED}
        strokeOpacity={0.35}
        strokeDasharray="2 4"
      />
      <text
        x={X(Math.PI)}
        y={18}
        textAnchor="middle"
        fill={MUTED}
        fontSize={11}
      >
        north pole
      </text>
      <text
        x={X(Math.PI)}
        y={H - 6}
        textAnchor="middle"
        fill={MUTED}
        fontSize={11}
      >
        south pole
      </text>
      <text
        x={X(TWO_PI) - 4}
        y={Y(0) - 6}
        textAnchor="end"
        fill={MUTED}
        fontSize={11}
      >
        equator
      </text>
    </svg>
  );

  const controls = (
    <>
      <span className="flex items-center gap-3">
        Order
        <Segmented
          size="small"
          value={order}
          onChange={(v) => setOrder(v as number)}
          options={[0, 1, 2, 3]}
        />
      </span>
      <span className="flex items-center gap-3">
        Numbering
        <Segmented
          size="small"
          value={scheme}
          onChange={(v) => setScheme(v as Scheme)}
          options={[
            { label: "NESTED", value: "nest" },
            { label: "RING", value: "ring" },
          ]}
        />
      </span>
      <label className="flex items-center gap-3">
        Path through the numbers
        <Switch size="small" checked={showPath} onChange={setShowPath} />
      </label>
      <span className="text-xs text-neutral-500">
        Order {order} · N<sub>side</sub> {nside} ·{" "}
        {formatNumber(12 * nside * nside)} cells
        {fontSize ? "" : " (numbers hidden)"}
      </span>
    </>
  );

  return <DiagramFrame stage={stage} controls={controls} wide />;
}
