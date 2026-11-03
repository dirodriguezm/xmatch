import { formatMjdTick } from "@/app/lib/utils/format";

/**
 * Shared by every time-series chart on the object page. Charts are stacked
 * on one MJD range; the same left and right margins keep their plot areas,
 * and so the same date, at the same x.
 */
export const TIME_GRID = { left: 70, right: 20, bottom: 64 } as const;

/** MJD tick labels with the year underneath. */
export const TIME_AXIS_LABEL = {
  color: "#d9d9d9",
  formatter: (value: number) => formatMjdTick(Math.round(value)),
  lineHeight: 14,
  // The shared range's ends aren't round numbers; their labels would
  // collide with the neighbouring ticks.
  showMinLabel: false,
  showMaxLabel: false,
  hideOverlap: true,
} as const;

/** Gap between the axis line and its "MJD" title, below the two-line ticks. */
export const TIME_AXIS_NAME_GAP = 44;
