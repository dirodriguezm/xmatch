/**
 * Sesame name resolver - converts astronomical object names to coordinates
 * Uses CDS Sesame service: http://cds.u-strasbg.fr/cgi-bin/nph-sesame
 */

export interface ResolvedCoordinates {
  ra: number;
  dec: number;
  objectName: string;
}

/**
 * Resolve an astronomical object name to RA/Dec coordinates using Sesame
 * Uses local API route to avoid CORS issues
 * @param name - Object name (e.g., "M31", "NGC 1234", "Crab Nebula")
 * @returns Resolved coordinates or null if not found
 */
export async function resolveObjectName(
  name: string
): Promise<ResolvedCoordinates | null> {
  const encodedName = encodeURIComponent(name.trim());
  const url = `/api/sesame?name=${encodedName}`;

  try {
    const response = await fetch(url);
    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data as ResolvedCoordinates;
  } catch (error) {
    console.error("Sesame resolution failed:", error);
    return null;
  }
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

/**
 * Parse HMS (hours:minutes:seconds) to degrees
 */
function parseHMS(hms: string): number | null {
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
 * Parse DMS (degrees:minutes:seconds) to degrees
 */
function parseDMS(dms: string): number | null {
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
