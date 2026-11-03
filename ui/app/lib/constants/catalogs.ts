/**
 * Centralized catalog configuration for consistent naming and colors across the app
 */

import type { CatalogRadiusConfig, RadiusUnit } from "./search";

export interface CatalogConfig {
  id: string;
  label: string;
  color: string;
  antdColor: string;
}

export const CATALOG_OPTIONS = ["allwise", "gaia", "erosita"] as const;

export type CatalogOption = (typeof CATALOG_OPTIONS)[number];

export const CATALOG_LABELS: Record<CatalogOption, string> = {
  allwise: "AllWISE",
  gaia: "Gaia",
  erosita: "eROSITA",
};

export const CATALOG_SELECT_OPTIONS = CATALOG_OPTIONS.map((catalog) => ({
  value: catalog,
  label: CATALOG_LABELS[catalog],
}));

export const CATALOG_DEFAULT_RADII: Record<
  CatalogOption,
  { radius: number; unit: RadiusUnit }
> = {
  allwise: { radius: 3, unit: "arcsec" },
  gaia: { radius: 3, unit: "arcsec" },
  erosita: { radius: 20, unit: "arcsec" },
};

export const CATALOG_COLOR_CLASSES: Record<string, string> = {
  allwise: "bg-cat-allwise",
  gaia: "bg-cat-gaia",
  erosita: "bg-cat-erosita",
};

/** Same colours as CATALOG_COLOR_CLASSES, for SVG marks. */
export const CATALOG_FILL_CLASSES: Record<string, string> = {
  allwise: "fill-cat-allwise",
  gaia: "fill-cat-gaia",
  erosita: "fill-cat-erosita",
};

export function buildDefaultCatalogConfigs(): CatalogRadiusConfig[] {
  return CATALOG_OPTIONS.map((catalog) => ({
    catalog,
    radius: CATALOG_DEFAULT_RADII[catalog].radius,
    unit: CATALOG_DEFAULT_RADII[catalog].unit,
    enabled: true,
  }));
}

export const CATALOGS: Record<string, CatalogConfig> = {
  GAIA_DR3: {
    id: "gaia_dr3",
    label: "GAIA DR3",
    color: "#e0b84c",
    antdColor: "gold",
  },
  SIMBAD: {
    id: "simbad",
    label: "SIMBAD",
    color: "#52c41a",
    antdColor: "green",
  },
  TWOMASS: {
    id: "2mass",
    label: "2MASS",
    color: "#fa8c16",
    antdColor: "orange",
  },
  WISE: {
    id: "wise",
    label: "WISE",
    color: "#e8684a",
    antdColor: "volcano",
  },
  ALLWISE: {
    id: "allwise",
    label: "AllWISE",
    color: "#e8684a",
    antdColor: "volcano",
  },
  EROSITA: {
    id: "erosita",
    label: "eROSITA",
    color: "#5cc3e8",
    antdColor: "cyan",
  },
} as const;

/**
 * Resolve a search-catalog slug (`allwise`, `gaia`, `erosita` — what the API
 * actually returns) to its display label.
 *
 * Distinct from {@link import("../utils/lightcurve").getCatalogLabel}, which
 * resolves light-curve *survey* names (neowise, ztf, …) from its own table.
 */
export function getSearchCatalogLabel(slug: string): string {
  const key = slug.toLowerCase() as CatalogOption;
  return CATALOG_LABELS[key] ?? slug.toUpperCase();
}

/**
 * Hex color per search-catalog slug, matching the dot colors in
 * {@link CATALOG_COLOR_CLASSES} so the sidebar and the results table agree.
 */
const SEARCH_CATALOG_COLORS: Record<CatalogOption, string> = {
  allwise: "#e8684a",
  gaia: "#e0b84c",
  erosita: "#5cc3e8",
};

/**
 * Hex color for a search-catalog slug.
 *
 * {@link getCatalogColor} matches on the legacy `CATALOGS` *label*, so `"gaia"`
 * resolved to neither `"GAIA DR3"` nor `"Gaia"` and Gaia tags rendered gray.
 * Keying on the slug removes the label round-trip entirely.
 */
export function getSearchCatalogColor(slug: string): string {
  return (
    SEARCH_CATALOG_COLORS[slug.toLowerCase() as CatalogOption] ?? "#8c8c8c"
  );
}

/**
 * Get the hex color for a catalog by its label
 */
export function getCatalogColor(catalogLabel: string): string {
  const normalizedLabel = catalogLabel.toUpperCase().replace(/\s+/g, "");

  for (const config of Object.values(CATALOGS)) {
    if (config.label.toUpperCase().replace(/\s+/g, "") === normalizedLabel) {
      return config.color;
    }
  }

  return "#8c8c8c"; // Default gray
}
