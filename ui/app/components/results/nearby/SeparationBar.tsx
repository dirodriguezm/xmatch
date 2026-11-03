import {
  agreement,
  AGREEMENT_LABEL,
  AGREEMENT_TEXT_CLASSES,
  formatArcsec,
  type NearbySource,
  sigmaRatio,
} from "./shared";

interface SeparationBarProps {
  source: NearbySource;
  /** Arcseconds the full bar stands for (usually the catalog's radius). */
  scaleArcsec: number;
  width?: number;
}

const H = 8;

/**
 * Separation drawn against the positional error: the green band is 2σ, the
 * amber one 3σ, the tick is where the source lies. The σ multiple is spelled
 * out next to it so colour is not the only cue.
 */
export function SeparationBar({
  source,
  scaleArcsec,
  width = 96,
}: SeparationBarProps) {
  const px = (arcsec: number) =>
    Math.min(width, (arcsec / Math.max(scaleArcsec, 1e-9)) * width);
  const tick = px(source.sepArcsec);
  const ag = agreement(source);
  const n = sigmaRatio(source);

  return (
    <span className="inline-flex items-center gap-2 whitespace-nowrap">
      <svg
        width={width}
        height={H}
        viewBox={`0 0 ${width} ${H}`}
        aria-hidden
        className="shrink-0"
      >
        <rect width={width} height={H} rx={2} className="fill-neutral-800" />
        <rect
          width={px(3 * source.sigma)}
          height={H}
          rx={2}
          className="fill-amber-500/30"
        />
        <rect
          width={px(2 * source.sigma)}
          height={H}
          rx={2}
          className="fill-green-500/40"
        />
        <rect
          x={Math.max(0, tick - 1.5)}
          width={3}
          height={H}
          className="fill-neutral-100"
        />
      </svg>
      <span className="font-mono text-xs">
        {formatArcsec(source.sepArcsec)}
      </span>
      <span
        className={`font-mono text-xs ${AGREEMENT_TEXT_CLASSES[ag]}`}
        title={`${n.toFixed(1)} × the combined 1σ error of ${formatArcsec(source.sigma)} — ${AGREEMENT_LABEL[ag]}`}
      >
        {n < 10 ? n.toFixed(1) : n.toFixed(0)}σ
      </span>
    </span>
  );
}
