import { Tooltip } from "antd";

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
  /** The scale is the catalog's search radius, so the tooltip can say so. */
  scaleIsRadius?: boolean;
  width?: number;
}

const H = 8;

/** Plain-language key to the bar, with this row's own numbers. */
function barTooltip(
  source: NearbySource,
  scaleArcsec: number,
  scaleIsRadius: boolean
) {
  const n = sigmaRatio(source);
  const sigma = formatArcsec(source.sigma);
  return (
    <div className="flex flex-col gap-1 text-xs">
      <span>
        {formatArcsec(source.sepArcsec)} from the searched position,{" "}
        {n.toFixed(1)}σ ({AGREEMENT_LABEL[agreement(source)]})
      </span>
      <span>
        Green band: within 2σ ({formatArcsec(2 * source.sigma)}). Amber: 2 to 3σ
        ({formatArcsec(3 * source.sigma)}). The tick marks this source.
      </span>
      <span className="text-neutral-400">
        Bar spans 0 to {formatArcsec(scaleArcsec)}
        {scaleIsRadius ? ", the search radius" : ""}. 1σ = {sigma}, the typical
        survey error combined with the searched position&apos;s.
      </span>
    </div>
  );
}

/**
 * Separation drawn against the positional error: the green band is 2σ, the
 * amber one 3σ, the tick is where the source lies. The σ multiple is spelled
 * out next to it so colour is not the only cue; the tooltip decodes the bar,
 * since nothing else on the page does.
 */
export function SeparationBar({
  source,
  scaleArcsec,
  scaleIsRadius = false,
  width = 96,
}: SeparationBarProps) {
  const px = (arcsec: number) =>
    Math.min(width, (arcsec / Math.max(scaleArcsec, 1e-9)) * width);
  const tick = px(source.sepArcsec);
  const ag = agreement(source);
  const n = sigmaRatio(source);

  return (
    <Tooltip
      title={barTooltip(source, scaleArcsec, scaleIsRadius)}
      placement="top"
      mouseEnterDelay={0.3}
    >
      <span className="inline-flex cursor-help items-center gap-2 whitespace-nowrap">
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
        <span className={`font-mono text-xs ${AGREEMENT_TEXT_CLASSES[ag]}`}>
          {n < 10 ? n.toFixed(1) : n.toFixed(0)}σ
        </span>
      </span>
    </Tooltip>
  );
}
