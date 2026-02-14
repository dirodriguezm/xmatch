/**
 * Coordinate conversion and formatting utilities for astronomical coordinates
 */

/**
 * Convert Right Ascension (in degrees) to Hours:Minutes:Seconds format
 * @param ra - Right Ascension in degrees (0-360)
 * @returns Formatted string like "12h 30m 45.67s"
 */
export function toHMS(ra: number): string {
  const hours = ra / 15;
  const h = Math.floor(hours);
  const m = Math.floor((hours - h) * 60);
  const s = ((hours - h) * 60 - m) * 60;
  return `${h}h ${m}m ${s.toFixed(2)}s`;
}

/**
 * Convert Declination (in degrees) to Degrees:Minutes:Seconds format
 * @param dec - Declination in degrees (-90 to +90)
 * @returns Formatted string like "+45° 30′ 15.00″"
 */
export function toDMS(dec: number): string {
  const sign = dec >= 0 ? "+" : "-";
  const absDec = Math.abs(dec);
  const d = Math.floor(absDec);
  const m = Math.floor((absDec - d) * 60);
  const s = ((absDec - d) * 60 - m) * 60;
  return `${sign}${d}° ${m}′ ${s.toFixed(2)}″`;
}

/**
 * Format a coordinate value with specified decimal places
 * @param value - The coordinate value
 * @param decimals - Number of decimal places (default: 6)
 * @returns Formatted string
 */
export function formatCoordinate(value: number, decimals = 6): string {
  return value.toFixed(decimals);
}

/**
 * Parse HMS (hours:minutes:seconds) string to degrees
 * @param hms - String in format "HH:MM:SS" or "HH:MM:SS.ss"
 * @returns Degrees, or null if parsing fails
 */
export function parseHMS(hms: string): number | null {
  const parts = hms.split(":");
  if (parts.length !== 3) {
    return null;
  }

  const hours = parseFloat(parts[0]);
  const minutes = parseFloat(parts[1]);
  const seconds = parseFloat(parts[2]);

  if (isNaN(hours) || isNaN(minutes) || isNaN(seconds)) {
    return null;
  }

  // Convert to degrees (1 hour = 15 degrees)
  return (hours + minutes / 60 + seconds / 3600) * 15;
}

/**
 * Parse DMS (degrees:minutes:seconds) string to degrees
 * @param dms - String in format "[+-]DD:MM:SS" or "[+-]DD:MM:SS.ss"
 * @returns Degrees, or null if parsing fails
 */
export function parseDMS(dms: string): number | null {
  // Handle sign
  let sign = 1;
  let cleanDms = dms.trim();

  if (cleanDms.startsWith("-")) {
    sign = -1;
    cleanDms = cleanDms.substring(1);
  } else if (cleanDms.startsWith("+")) {
    cleanDms = cleanDms.substring(1);
  }

  const parts = cleanDms.split(":");
  if (parts.length !== 3) {
    return null;
  }

  const degrees = parseFloat(parts[0]);
  const minutes = parseFloat(parts[1]);
  const seconds = parseFloat(parts[2]);

  if (isNaN(degrees) || isNaN(minutes) || isNaN(seconds)) {
    return null;
  }

  return sign * (degrees + minutes / 60 + seconds / 3600);
}

/**
 * Try to parse a string as coordinates (RA Dec format)
 * Supports: "12:30:00 -45:00:00" (HMS/DMS) or "187.5 -45.0" (degrees)
 * @returns Parsed coordinates in degrees, or null if not valid coordinates
 */
export function parseCoordinates(
  input: string
): { ra: number; dec: number } | null {
  const trimmed = input.trim();

  // Try to split by common separators
  const parts = trimmed.split(/[\s,]+/);

  if (parts.length < 2) {
    return null;
  }

  const raPart = parts[0];
  const decPart = parts.slice(1).join(" ");

  // Try parsing as HMS/DMS (contains colons)
  if (raPart.includes(":") || decPart.includes(":")) {
    const ra = parseHMS(raPart);
    const dec = parseDMS(decPart);

    if (ra !== null && dec !== null) {
      return { ra, dec };
    }
  }

  // Try parsing as decimal degrees
  const ra = parseFloat(raPart);
  const dec = parseFloat(decPart);

  if (
    !isNaN(ra) &&
    !isNaN(dec) &&
    ra >= 0 &&
    ra <= 360 &&
    dec >= -90 &&
    dec <= 90
  ) {
    return { ra, dec };
  }

  return null;
}
