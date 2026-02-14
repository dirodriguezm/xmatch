/**
 * Photometry band definitions and colors for astronomical surveys
 */

export interface BandConfig {
  band: string;
  survey: string;
  wavelength: string;
  color: string;
}

/**
 * Band colors for different filters/surveys
 */
export const BAND_COLORS: Record<string, string> = {
  // ZTF bands
  g: "#52c41a", // green
  r: "#ff4d4f", // red
  i: "#722ed1", // purple
  z: "#fa8c16", // orange

  // ZTF filter IDs
  "1": "#52c41a", // ZTF g
  "2": "#ff4d4f", // ZTF r
  "3": "#722ed1", // ZTF i

  // Gaia bands
  G: "#1890ff", // blue
  BP: "#52c41a", // green
  RP: "#ff4d4f", // red

  // 2MASS bands
  J: "#fa8c16", // orange
  H: "#722ed1", // purple
  K: "#13c2c2", // cyan

  // WISE bands
  W1: "#1890ff", // blue
  W2: "#52c41a", // green
  W3: "#fa8c16", // orange
  W4: "#ff4d4f", // red
};

/**
 * Photometry bands with survey and wavelength information
 */
export const PHOTOMETRY_BANDS: BandConfig[] = [
  { band: "G", survey: "Gaia", wavelength: "0.64 μm", color: BAND_COLORS.G },
  { band: "BP", survey: "Gaia", wavelength: "0.51 μm", color: BAND_COLORS.BP },
  { band: "RP", survey: "Gaia", wavelength: "0.78 μm", color: BAND_COLORS.RP },
  { band: "J", survey: "2MASS", wavelength: "1.24 μm", color: BAND_COLORS.J },
  { band: "H", survey: "2MASS", wavelength: "1.66 μm", color: BAND_COLORS.H },
  { band: "K", survey: "2MASS", wavelength: "2.16 μm", color: BAND_COLORS.K },
  { band: "W1", survey: "WISE", wavelength: "3.4 μm", color: BAND_COLORS.W1 },
  { band: "W2", survey: "WISE", wavelength: "4.6 μm", color: BAND_COLORS.W2 },
  { band: "W3", survey: "WISE", wavelength: "12 μm", color: BAND_COLORS.W3 },
  { band: "W4", survey: "WISE", wavelength: "22 μm", color: BAND_COLORS.W4 },
];

/**
 * Get the color for a band
 * @param band - Band name (e.g., "g", "G", "W1")
 * @returns Hex color string
 */
export function getBandColor(band: string): string {
  return BAND_COLORS[band] || "#1890ff";
}
