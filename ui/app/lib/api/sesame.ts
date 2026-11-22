/**
 * Sesame name resolver - converts astronomical object names to coordinates
 * Uses CDS Sesame service: http://cds.u-strasbg.fr/cgi-bin/nph-sesame
 */

import { reportServiceResult } from "@/app/lib/serviceHealth";

// Re-export coordinate parsing functions from the shared utility
export { parseCoordinates } from "@/app/lib/utils/coordinates";

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
    // 404 means "no such object": Sesame itself answered.
    reportServiceResult("Sesame", response.ok || response.status === 404);
    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data as ResolvedCoordinates;
  } catch (error) {
    console.error("Sesame resolution failed:", error);
    reportServiceResult("Sesame", false);
    return null;
  }
}
