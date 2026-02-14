/**
 * URL builders for external astronomical services
 */

/**
 * Build a SIMBAD cone search URL
 * @param ra - Right Ascension in degrees
 * @param dec - Declination in degrees
 * @param radius - Search radius in arcseconds (default: 2)
 */
export function buildSimbadUrl(ra: number, dec: number, radius = 2): string {
  return `https://simbad.cds.unistra.fr/simbad/sim-coo?Coord=${ra}+${dec}&Radius=${radius}&Radius.unit=arcsec`;
}

/**
 * Build a VizieR cone search URL
 * @param ra - Right Ascension in degrees
 * @param dec - Declination in degrees
 * @param radius - Search radius in arcseconds (default: 2)
 */
export function buildVizierUrl(ra: number, dec: number, radius = 2): string {
  return `https://vizier.cds.unistra.fr/viz-bin/VizieR?-c=${ra}+${dec}&-c.rs=${radius}`;
}

/**
 * Build an Aladin Lite URL
 * @param ra - Right Ascension in degrees
 * @param dec - Declination in degrees
 * @param fov - Field of view in degrees (default: 0.1)
 */
export function buildAladinUrl(ra: number, dec: number, fov = 0.1): string {
  return `https://aladin.cds.unistra.fr/AladinLite/?target=${ra}+${dec}&fov=${fov}`;
}
