import { fillClass, type MarkShape, shapeOf } from "./shared";

interface MarkPathProps {
  shape: MarkShape;
  /** Half-size of the mark, in SVG units. */
  r: number;
  x?: number;
  y?: number;
  className?: string;
  fillOpacity?: number;
  stroke?: string;
  strokeWidth?: number;
}

/** One catalog glyph inside an existing <svg>. */
export function MarkPath({
  shape,
  r,
  x = 0,
  y = 0,
  className,
  fillOpacity,
  stroke = "none",
  strokeWidth,
}: MarkPathProps) {
  const common = { className, fillOpacity, stroke, strokeWidth };
  switch (shape) {
    case "square":
      return (
        <rect
          x={x - r * 0.85}
          y={y - r * 0.85}
          width={r * 1.7}
          height={r * 1.7}
          {...common}
        />
      );
    case "diamond":
      return (
        <polygon
          points={`${x},${y - r * 1.15} ${x + r * 1.15},${y} ${x},${y + r * 1.15} ${x - r * 1.15},${y}`}
          {...common}
        />
      );
    case "triangle":
      return (
        <polygon
          points={`${x},${y - r * 1.1} ${x + r},${y + r * 0.8} ${x - r},${y + r * 0.8}`}
          {...common}
        />
      );
    default:
      return <circle cx={x} cy={y} r={r} {...common} />;
  }
}

/** Inline catalog glyph (shape + colour) for legends and table cells. */
export function CatalogMark({
  slug,
  size = 10,
}: {
  slug: string;
  size?: number;
}) {
  return (
    <svg
      viewBox="-6 -6 12 12"
      width={size}
      height={size}
      aria-hidden
      className="inline-block shrink-0"
    >
      <MarkPath shape={shapeOf(slug)} r={4.5} className={fillClass(slug)} />
    </svg>
  );
}
