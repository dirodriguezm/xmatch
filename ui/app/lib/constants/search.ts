export const RADIUS_UNITS = ["arcsec", "arcmin", "deg"] as const;

export type RadiusUnit = (typeof RADIUS_UNITS)[number];

export const DEFAULT_RADIUS_UNIT: RadiusUnit = "deg";

export const RADIUS_UNIT_OPTIONS = RADIUS_UNITS.map((unit) => ({
  value: unit,
  label: unit,
}));

export interface CatalogRadiusConfig {
  catalog: string;
  radius: number;
  unit: RadiusUnit;
  enabled: boolean;
}

/**
 * Encode catalog radius configs to a URL-safe string.
 * Format: "allwise:5:arcsec:1,gaia:2:arcsec:1,erosita:10:arcmin:0"
 */
export function encodeCatalogRadii(configs: CatalogRadiusConfig[]): string {
  return configs
    .map((c) => `${c.catalog}:${c.radius}:${c.unit}:${c.enabled ? "1" : "0"}`)
    .join(",");
}

/**
 * Decode a URL-safe string back to catalog radius configs.
 */
export function decodeCatalogRadii(str: string): CatalogRadiusConfig[] {
  if (!str) return [];
  return str
    .split(",")
    .map((segment) => {
      const [catalog, radiusStr, unit, enabledStr] = segment.split(":");
      const radius = parseFloat(radiusStr);
      if (!catalog || isNaN(radius) || !unit) return null;
      return {
        catalog,
        radius,
        unit: (RADIUS_UNITS.includes(unit as RadiusUnit)
          ? unit
          : "arcsec") as RadiusUnit,
        enabled: enabledStr === "1",
      };
    })
    .filter((c): c is CatalogRadiusConfig => c !== null);
}

/**
 * Convert a radius value to degrees for API calls.
 */
export function convertRadiusToDegrees(
  radius: number,
  unit: RadiusUnit
): number {
  switch (unit) {
    case "arcsec":
      return radius / 3600;
    case "arcmin":
      return radius / 60;
    case "deg":
    default:
      return radius;
  }
}
