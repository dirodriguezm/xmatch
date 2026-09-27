/**
 * Curated targets for the Explore page and the landing "Surprise me" link.
 *
 * `ra`/`dec` are J2000 (ICRS) positions from SIMBAD. `objectId` + `catalog`
 * were resolved against the live API (conesearch, then GET /v1/metadata).
 * Gaia positions are at epoch 2016.0, so for the fast movers (Barnard's Star,
 * Proxima Cen) the catalog position sits arcminutes away from J2000.
 */

export interface FeaturedObject {
  slug: string;
  name: string;
  /** Short kind label, e.g. "High proper-motion star". */
  kind: string;
  ra: number;
  dec: number;
  objectId: string;
  /** Search-catalog slug: gaia, allwise, erosita. */
  catalog: string;
  blurb: string;
  /** What to look at on the object page. */
  lookFor: string[];
}

export const FEATURED_OBJECTS: FeaturedObject[] = [
  {
    slug: "barnards-star",
    name: "Barnard's Star",
    kind: "High proper-motion red dwarf",
    ra: 269.452075,
    dec: 4.693391,
    objectId: "Gaia DR3 4472832130942575872",
    catalog: "gaia",
    blurb:
      "The fastest-moving star on our sky, crossing about 10.4″ every year — a Moon diameter in under two centuries. At 1.8 pc it is the second-closest stellar system after α Centauri, which is why its Gaia parallax is over half an arcsecond.",
    lookFor: [
      "Proper motion (~10 400 mas/yr)",
      "Parallax",
      "Gaia epoch 2016 vs J2000 offset",
    ],
  },
  {
    slug: "proxima-cen",
    name: "Proxima Centauri",
    kind: "Nearest star, flare star",
    ra: 217.428953,
    dec: -62.679484,
    objectId: "Gaia DR3 5853498713190525696",
    catalog: "gaia",
    blurb:
      "Our nearest stellar neighbour at 1.30 pc, an M5.5 dwarf that hosts at least one planet in its habitable zone. It is a frequent flarer, and its 3.9″/yr motion means the Gaia source sits well away from the J2000 position.",
    lookFor: ["Parallax (768 mas)", "Proper motion", "Light-curve flares"],
  },
  {
    slug: "rr-lyrae",
    name: "RR Lyrae",
    kind: "Pulsating variable (prototype)",
    ra: 291.365882,
    dec: 42.784361,
    objectId: "Gaia DR3 2125982599343482624",
    catalog: "gaia",
    blurb:
      "The namesake of a whole class of old, horizontal-branch pulsators that serve as standard candles. It brightens by almost a magnitude every 0.567 days, so its epoch photometry is a textbook saw-tooth once folded.",
    lookFor: ["Gaia epoch light curve", "Pan-STARRS / ZTF variability", "SED"],
  },
  {
    slug: "3c-273",
    name: "3C 273",
    kind: "Quasar, X-ray source",
    ra: 187.277915,
    dec: 2.052388,
    objectId: "1eRASS J122906.7+020309",
    catalog: "erosita",
    blurb:
      "The first quasar ever identified (1963) and still the optically brightest. It sits in the western Galactic hemisphere covered by eROSITA-DE DR1, where it is one of the brightest extragalactic X-ray sources, with a one-sided jet visible in radio and X-ray imaging.",
    lookFor: [
      "eROSITA X-ray fluxes",
      "Counterparts in Gaia and AllWISE",
      "Sky view (Chandra/XMM layers)",
    ],
  },
  {
    slug: "lmc-x-3",
    name: "LMC X-3",
    kind: "Black-hole X-ray binary",
    ra: 84.733417,
    dec: -64.083889,
    objectId: "1eRASS J053856.1-640457",
    catalog: "erosita",
    blurb:
      "A persistent black-hole X-ray binary in the Large Magellanic Cloud: a ~7 M☉ black hole stripping gas from a B-type companion every 1.7 days. Its accretion disc makes it one of the brightest X-ray points in the LMC.",
    lookFor: [
      "eROSITA detection and fluxes",
      "Nearby sources in the LMC field",
    ],
  },
  {
    slug: "crab-pulsar",
    name: "Crab Pulsar",
    kind: "Young neutron star",
    ra: 83.633083,
    dec: 22.0145,
    objectId: "Gaia DR3 3403818172572314624",
    catalog: "gaia",
    blurb:
      "The spinning neutron star left by the supernova of 1054 AD, rotating 30 times a second at the heart of the Crab Nebula. Gaia even measures its parallax (~0.5 mas, about 2 kpc), a rare optical handle on a pulsar.",
    lookFor: ["Parallax", "Sky view of the nebula", "SED"],
  },
  {
    slug: "betelgeuse",
    name: "Betelgeuse",
    kind: "Red supergiant",
    ra: 88.792939,
    dec: 7.407064,
    objectId: "0884p075_ac51-016718",
    catalog: "allwise",
    blurb:
      "Orion's shoulder, a red supergiant so large it would swallow Jupiter's orbit. It is far too bright for Gaia, but AllWISE records it — saturated, with negative W1 and W4 magnitudes — a reminder to check quality flags on bright stars. It made headlines in 2019–20 with its Great Dimming.",
    lookFor: [
      "AllWISE / 2MASS magnitudes (saturated)",
      "SED shape of a cool giant",
    ],
  },
  {
    slug: "algol",
    name: "Algol",
    kind: "Eclipsing binary",
    ra: 47.042215,
    dec: 40.955648,
    objectId: "0474p408_ac51-031388",
    catalog: "allwise",
    blurb:
      "The 'Demon Star' in Perseus, the prototype eclipsing binary: every 2.87 days the dim companion passes in front of the bright B star and the system fades by more than a magnitude for about ten hours.",
    lookFor: ["Infrared photometry", "SED", "Observability tonight"],
  },
  {
    slug: "polaris",
    name: "Polaris",
    kind: "Cepheid, pole star",
    ra: 37.954561,
    dec: 89.264109,
    objectId: "0600p893_ac51-031964",
    catalog: "allwise",
    blurb:
      "The North Star is also the nearest classical Cepheid, pulsating with a 4-day period and an amplitude that shrank dramatically in the 20th century. At 89° declination it is a good stress test for RA wrap-around near the pole.",
    lookFor: ["Sky view near the pole", "AllWISE / 2MASS photometry"],
  },
  {
    slug: "m31-nucleus",
    name: "M31 nucleus",
    kind: "Galaxy nucleus",
    ra: 10.684708,
    dec: 41.26875,
    objectId: "0098p408_ac51-043708",
    catalog: "allwise",
    blurb:
      "The core of the Andromeda Galaxy, hosting a ~100 million M☉ black hole and a double nucleus. AllWISE sees the unresolved bulge light; Gaia instead breaks the region into many crowded point sources nearby.",
    lookFor: ["Nearby sources", "Mid-IR colours", "Sky view"],
  },
];

/** ISO-8601 week number (1–53) of a date, in UTC. */
export function isoWeek(date: Date): number {
  const d = new Date(
    Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate())
  );
  // Thursday of this week decides the year.
  const day = d.getUTCDay() || 7;
  d.setUTCDate(d.getUTCDate() + 4 - day);
  const yearStart = Date.UTC(d.getUTCFullYear(), 0, 1);
  return Math.ceil(((d.getTime() - yearStart) / 86_400_000 + 1) / 7);
}

/** Deterministic "Object of the week", rotating through the list. */
export function objectOfTheWeek(
  date: Date = new Date(),
  list: FeaturedObject[] = FEATURED_OBJECTS
): FeaturedObject {
  return list[isoWeek(date) % list.length];
}
