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
 * Convert a radius value to arcseconds for API calls.
 * The conesearch API expects the radius in arcseconds.
 */
export function convertRadiusToArcsec(
  radius: number,
  unit: RadiusUnit
): number {
  switch (unit) {
    case "arcmin":
      return radius * 60;
    case "deg":
      return radius * 3600;
    case "arcsec":
    default:
      return radius;
  }
}

/** Convert an arcsecond value back into `unit`, for display in the form. */
export function convertArcsecToUnit(arcsec: number, unit: RadiusUnit): number {
  switch (unit) {
    case "arcmin":
      return arcsec / 60;
    case "deg":
      return arcsec / 3600;
    case "arcsec":
    default:
      return arcsec;
  }
}

/**
 * Largest radius the conesearch backend answers in reasonable time.
 * Measured against https://xwave-astro.udp.cl/v1: 120" replies in ~46ms, while
 * 150", 180" and 240" never reply (still open after 60s). The unit selector
 * makes it trivially easy to cross that line — 3 arcmin is already 180" — so we
 * refuse client-side instead of hanging on a request that will not come back.
 */
export const MAX_RADIUS_ARCSEC = 120;

/**
 * Upper bound on objects returned per catalog.
 *
 * The backend defaults `nneighbor` to 1, and the UI never sent the parameter,
 * so a search returned a single object per catalog no matter how large the
 * radius — making the radius control almost inert. This is a *cap*, not a
 * count: the backend returns whatever falls inside the radius, up to this many.
 * Latency is flat with respect to it (~35ms measured at both n=1 and n=500),
 * so a generous cap costs nothing and hands the radius back its meaning.
 */
export const CONE_SEARCH_MAX_NEIGHBORS = 100;
