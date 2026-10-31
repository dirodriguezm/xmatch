/**
 * Data processing utilities for charts and visualizations
 */

export interface Bounds {
  min: number;
  max: number;
}

/**
 * Calculate min/max bounds from an array of numbers with optional padding
 * @param values - Array of numeric values
 * @param paddingPercent - Percentage of range to add as padding (default: 0.1 = 10%)
 * @param minPadding - Minimum padding value when range is zero or very small
 * @returns Object with min and max values including padding
 */
export function calculateAxisBounds(
  values: number[],
  paddingPercent = 0.1,
  minPadding = 0.001
): Bounds {
  if (values.length === 0) {
    return { min: 0, max: 1 };
  }

  const min = Math.min(...values);
  const max = Math.max(...values);
  const range = max - min;
  const padding = range * paddingPercent || minPadding;

  return {
    min: min - padding,
    max: max + padding,
  };
}
