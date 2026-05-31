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

/**
 * Build an SDSS DR19 SkyServer Navigate URL
 */
export function buildSdssNavigateUrl(ra: number, dec: number): string {
  return `https://skyserver.sdss.org/dr19/VisualTools/navi?ra=${ra}&dec=${dec}`;
}

/**
 * Build a DESI Legacy Survey Sky Viewer URL (imaging coverage incl. DESI/DECaLS/BASS/MzLS)
 */
export function buildLegacySurveyViewerUrl(
  ra: number,
  dec: number,
  layer = "ls-dr10",
  zoom = 16
): string {
  return `https://www.legacysurvey.org/viewer?ra=${ra}&dec=${dec}&layer=${layer}&zoom=${zoom}&mark=${ra},${dec}`;
}

/**
 * Build a NASA/IPAC Extragalactic Database (NED) cone search URL
 * @param radius - Search radius in arcminutes (NED's "Radius" expects arcmin; default 0.033 ≈ 2″)
 */
export function buildNedUrl(ra: number, dec: number, radius = 0.033): string {
  return `https://ned.ipac.caltech.edu/cgi-bin/objsearch?search_type=Near+Position+Search&in_csys=Equatorial&in_equinox=J2000.0&lon=${ra}d&lat=${dec}d&radius=${radius}&obj_sort=Distance+to+search+center`;
}

/**
 * Build a Pan-STARRS1 image cutout URL
 */
export function buildPanstarrsUrl(ra: number, dec: number): string {
  return `https://ps1images.stsci.edu/cgi-bin/ps1cutouts?pos=${ra}+${dec}&filter=color&filetypes=stack&size=240&output_size=0&verbose=0&autoscale=99.500000&catlist=`;
}

/**
 * Build a Legacy Survey DESI DR1 spectrum viewer URL.
 *
 * Requires a TARGETID — not derivable from (RA, Dec) by URL manipulation.
 * Use the `/api/desi-target` route (server-side SPARCL cone search) to get
 * the TARGETID for a given position.
 */
export function buildDesiSpectrumUrl(targetid: string | number): string {
  return `https://www.legacysurvey.org/viewer/desi-spectrum/dr1/targetid${targetid}`;
}
