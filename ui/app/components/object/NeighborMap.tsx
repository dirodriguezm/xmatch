"use client";

import { useRouter } from "next/navigation";

import { type Neighbor } from "@/app/hooks/queries";
import {
  CATALOG_FILL_CLASSES,
  getSearchCatalogLabel,
} from "@/app/lib/constants/catalogs";
import { buildObjectUrl } from "@/app/lib/utils/urls";

interface NeighborMapProps {
  neighbors: Neighbor[];
  radiusArcsec: number;
  /** Row key (`catalog:id`) currently hovered here or in the table. */
  hoveredKey: string | null;
  onHover: (key: string | null) => void;
}

export const neighborKey = (n: Neighbor) => `${n.catalog}:${n.id}`;

const SIZE = 260;
/** Room around the outer ring for the N/E labels, in SVG units. */
const PAD = 16;

/**
 * Offsets of the nearby sources from the object, on the sky: north up, east
 * left, as in Sky View. Rings every 10″ out to the search radius.
 */
export function NeighborMap({
  neighbors,
  radiusArcsec,
  hoveredKey,
  onHover,
}: NeighborMapProps) {
  const router = useRouter();
  const half = SIZE / 2;
  const scale = (half - PAD) / radiusArcsec;
  const rings = Array.from(
    { length: Math.floor(radiusArcsec / 10) },
    (_, i) => (i + 1) * 10
  );

  // Draw the hovered source last so it sits on top of any overlapping mark.
  const ordered = [...neighbors].sort(
    (a, b) =>
      Number(neighborKey(a) === hoveredKey) -
      Number(neighborKey(b) === hoveredKey)
  );

  return (
    <svg
      viewBox={`${-half} ${-half} ${SIZE} ${SIZE}`}
      width={SIZE}
      height={SIZE}
      className="shrink-0 text-foreground"
      role="img"
      aria-label={`Positions of ${neighbors.length} nearby sources relative to the object`}
    >
      {rings.map((r) => (
        <g key={r}>
          <circle
            r={r * scale}
            fill="none"
            stroke="currentColor"
            strokeOpacity={0.15}
            strokeDasharray={r === radiusArcsec ? undefined : "2 3"}
          />
          <text
            x={r * scale * Math.SQRT1_2 + 2}
            y={-r * scale * Math.SQRT1_2 - 2}
            fontSize={9}
            fill="currentColor"
            fillOpacity={0.4}
          >
            {r}″
          </text>
        </g>
      ))}
      <text
        y={-half + 11}
        textAnchor="middle"
        fontSize={10}
        fill="currentColor"
        fillOpacity={0.6}
      >
        N
      </text>
      <text
        x={-half + 4}
        y={4}
        fontSize={10}
        fill="currentColor"
        fillOpacity={0.6}
      >
        E
      </text>

      {/* The object itself */}
      <g stroke="currentColor" strokeOpacity={0.7}>
        <line x1={-7} x2={-3} />
        <line x1={3} x2={7} />
        <line y1={-7} y2={-3} />
        <line y1={3} y2={7} />
      </g>

      {ordered.map((n) => {
        const key = neighborKey(n);
        const pa = (n.positionAngle * Math.PI) / 180;
        // East is left, so the east offset maps to −x; north is up (−y).
        const x = -n.separationArcsec * Math.sin(pa) * scale;
        const y = -n.separationArcsec * Math.cos(pa) * scale;
        const hovered = key === hoveredKey;
        return (
          <circle
            key={key}
            cx={x}
            cy={y}
            r={hovered ? 6 : 4}
            className={`${CATALOG_FILL_CLASSES[n.catalog]} cursor-pointer`}
            fillOpacity={hoveredKey && !hovered ? 0.35 : 0.85}
            stroke={hovered ? "currentColor" : "none"}
            strokeWidth={1.5}
            onMouseEnter={() => onHover(key)}
            onMouseLeave={() => onHover(null)}
            onClick={() => router.push(buildObjectUrl(n.id, n.catalog))}
          >
            <title>
              {`${getSearchCatalogLabel(n.catalog)} ${n.id}\n${n.separationArcsec.toFixed(2)}″ at PA ${n.positionAngle.toFixed(0)}°`}
            </title>
          </circle>
        );
      })}
    </svg>
  );
}
