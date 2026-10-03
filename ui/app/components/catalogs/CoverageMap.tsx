"use client";

import { useId } from "react";

import {
  CATALOG_FILL_CLASSES,
  type CatalogOption,
} from "@/app/lib/constants/catalogs";

/** Which part of the sky a catalog covers, in Galactic coordinates. */
export type CoverageKind = "full" | "western-galactic";

export const CATALOG_COVERAGE: Record<CatalogOption, CoverageKind> = {
  gaia: "full",
  allwise: "full",
  erosita: "western-galactic",
};

// Mollweide frame: 2:1 ellipse centred in a 200×100 view box.
const CX = 100;
const CY = 50;
const A = 96; // half-width
const B = A / 2; // half-height

/** Auxiliary angle θ with 2θ + sin 2θ = π sin φ (Newton iterations). */
function mollweideTheta(lat: number): number {
  if (Math.abs(lat) >= Math.PI / 2) return Math.sign(lat) * (Math.PI / 2);
  let t = lat;
  for (let i = 0; i < 20; i++) {
    const f = 2 * t + Math.sin(2 * t) - Math.PI * Math.sin(lat);
    const d = 2 + 2 * Math.cos(2 * t);
    if (d === 0) break;
    t -= f / d;
  }
  return t;
}

/** Galactic latitude parallels drawn as grid lines, in degrees. */
const PARALLELS = [-60, -30, 30, 60].map((deg) => {
  const theta = mollweideTheta((deg * Math.PI) / 180);
  return { deg, y: CY - B * Math.sin(theta), halfWidth: A * Math.cos(theta) };
});

/** Galactic longitude meridians every 60°, as ellipse half-widths. */
const MERIDIANS = [1, 2].map((k) => (A * k) / 3);

interface CoverageMapProps {
  catalog: CatalogOption;
  label: string;
  className?: string;
}

/**
 * All-sky sketch of a catalog's footprint in a Mollweide projection of
 * Galactic coordinates (l = 0° at the centre, increasing to the left), so a
 * Galactic-hemisphere split is a straight vertical line.
 */
export function CoverageMap({ catalog, label, className }: CoverageMapProps) {
  const clipId = useId();
  const kind = CATALOG_COVERAGE[catalog];
  const fill = CATALOG_FILL_CLASSES[catalog] ?? "fill-neutral-500";
  const description =
    kind === "full"
      ? `${label} covers the whole sky.`
      : `${label} public data covers the western Galactic hemisphere, Galactic longitude 180° to 360°.`;

  return (
    <figure className={`m-0 ${className ?? ""}`}>
      <svg
        viewBox="0 0 200 108"
        role="img"
        aria-label={`Sky coverage of ${label}. ${description}`}
        className="w-full h-auto"
      >
        <defs>
          <clipPath id={clipId}>
            <ellipse cx={CX} cy={CY} rx={A} ry={B} />
          </clipPath>
        </defs>

        <ellipse
          cx={CX}
          cy={CY}
          rx={A}
          ry={B}
          className="fill-neutral-900 stroke-neutral-700"
          strokeWidth={0.6}
        />

        <g clipPath={`url(#${clipId})`}>
          {kind === "full" ? (
            <rect
              x={CX - A}
              y={CY - B}
              width={2 * A}
              height={2 * B}
              className={fill}
              fillOpacity={0.45}
            />
          ) : (
            <rect
              x={CX}
              y={CY - B}
              width={A}
              height={2 * B}
              className={fill}
              fillOpacity={0.45}
            />
          )}
        </g>

        {/* Graticule */}
        <g
          className="stroke-neutral-500"
          strokeWidth={0.3}
          strokeOpacity={0.5}
          fill="none"
        >
          {MERIDIANS.map((rx) => (
            <ellipse key={rx} cx={CX} cy={CY} rx={rx} ry={B} />
          ))}
          <line x1={CX} y1={CY - B} x2={CX} y2={CY + B} />
          {PARALLELS.map((p) => (
            <line
              key={p.deg}
              x1={CX - p.halfWidth}
              x2={CX + p.halfWidth}
              y1={p.y}
              y2={p.y}
            />
          ))}
        </g>

        {/* Galactic plane */}
        <line
          x1={CX - A}
          x2={CX + A}
          y1={CY}
          y2={CY}
          className="stroke-neutral-300"
          strokeWidth={0.5}
          strokeDasharray="2 1.5"
        />

        {kind === "western-galactic" && (
          <line
            x1={CX}
            y1={CY - B}
            x2={CX}
            y2={CY + B}
            className="stroke-neutral-200"
            strokeWidth={0.8}
          />
        )}

        <g className="fill-neutral-400" fontSize={5} fontFamily="monospace">
          <text x={CX} y={CY + B + 6} textAnchor="middle">
            l = 0°
          </text>
          <text x={CX - A} y={CY + B + 6} textAnchor="start">
            180°
          </text>
          <text x={CX + A} y={CY + B + 6} textAnchor="end">
            180°
          </text>
        </g>
        {kind === "western-galactic" && (
          <g fontSize={5.5}>
            <text
              x={CX + A / 2}
              y={CY - 6}
              textAnchor="middle"
              className="fill-neutral-100"
            >
              public
            </text>
            <text
              x={CX + A / 2}
              y={CY + 10}
              textAnchor="middle"
              className="fill-neutral-300"
              fontFamily="monospace"
              fontSize={4.5}
            >
              180° &lt; l &lt; 360°
            </text>
            <text
              x={CX - A / 2}
              y={CY - 6}
              textAnchor="middle"
              className="fill-neutral-500"
            >
              not released
            </text>
          </g>
        )}
      </svg>
      <figcaption className="text-xs text-neutral-500 mt-1">
        {kind === "full"
          ? "All sky · Galactic Mollweide, plane dashed"
          : "Western Galactic hemisphere · Galactic Mollweide, plane dashed"}
      </figcaption>
    </figure>
  );
}
