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

export const CATALOG_OPTIONS = ["allwise", "gaia"] as const;

export type CatalogOption = (typeof CATALOG_OPTIONS)[number];

export const CATALOG_LABELS: Record<CatalogOption, string> = {
  allwise: "AllWISE",
  gaia: "Gaia",
};

export const CATALOG_SELECT_OPTIONS = CATALOG_OPTIONS.map((catalog) => ({
  value: catalog,
  label: CATALOG_LABELS[catalog],
}));

export const CATALOG_DEFAULT_RADII: Record<
  CatalogOption,
  { radius: number; unit: RadiusUnit }
> = {
  allwise: { radius: 5, unit: "deg" },
  gaia: { radius: 10, unit: "deg" },
};

export const CATALOG_COLOR_CLASSES: Record<string, string> = {
  allwise: "bg-purple-600",
  gaia: "bg-blue-500",
  erosita: "bg-red-500",
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
    color: "#1890ff",
    antdColor: "blue",
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
    color: "#722ed1",
    antdColor: "purple",
  },
  ALLWISE: {
    id: "allwise",
    label: "AllWISE",
    color: "#722ed1",
    antdColor: "purple",
  },
} as const;

/**
 * Get catalog options for form selects/checkboxes
 */
export function getCatalogOptions(): CatalogConfig[] {
  // Exclude ALLWISE from options since it's a variant of WISE
  return [CATALOGS.GAIA_DR3, CATALOGS.SIMBAD, CATALOGS.TWOMASS, CATALOGS.WISE];
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

/**
 * Get the Ant Design color name for a catalog by its label
 */
export function getCatalogAntdColor(catalogLabel: string): string {
  const normalizedLabel = catalogLabel.toUpperCase().replace(/\s+/g, "");

  for (const config of Object.values(CATALOGS)) {
    if (config.label.toUpperCase().replace(/\s+/g, "") === normalizedLabel) {
      return config.antdColor;
    }
  }

  return "default";
}
