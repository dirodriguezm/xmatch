/**
 * Observatories offered in the observability panel. Coordinates are geodetic
 * (WGS84), longitude east-positive, height in metres. Times are shown in each
 * site's own time zone.
 */

export type ObservatoryRegion = "Chile" | "Rest of the world";

export interface Observatory {
  id: string;
  label: string;
  latitude: number;
  longitude: number;
  heightM: number;
  /** IANA time zone the site's local times are shown in. */
  timeZone: string;
  region: ObservatoryRegion;
}

const CHILE = "America/Santiago";

export const OBSERVATORIES: Observatory[] = [
  {
    id: "paranal",
    label: "Paranal (VLT)",
    latitude: -24.6272,
    longitude: -70.4042,
    heightM: 2635,
    timeZone: CHILE,
    region: "Chile",
  },
  {
    id: "pachon",
    label: "Cerro Pachón (Rubin, Gemini-S)",
    latitude: -30.2446,
    longitude: -70.7494,
    heightM: 2715,
    timeZone: CHILE,
    region: "Chile",
  },
  {
    id: "lascampanas",
    label: "Las Campanas (Magellan)",
    latitude: -29.0146,
    longitude: -70.6926,
    heightM: 2380,
    timeZone: CHILE,
    region: "Chile",
  },
  {
    id: "lasilla",
    label: "La Silla (ESO)",
    latitude: -29.2567,
    longitude: -70.7377,
    heightM: 2400,
    timeZone: CHILE,
    region: "Chile",
  },
  {
    id: "tololo",
    label: "Cerro Tololo (CTIO, SOAR)",
    latitude: -30.169,
    longitude: -70.8063,
    heightM: 2207,
    timeZone: CHILE,
    region: "Chile",
  },
  {
    id: "maunakea",
    label: "Maunakea (Keck, Subaru, Gemini-N)",
    latitude: 19.8206,
    longitude: -155.4681,
    heightM: 4207,
    timeZone: "Pacific/Honolulu",
    region: "Rest of the world",
  },
  {
    id: "lapalma",
    label: "Roque de los Muchachos (GTC, La Palma)",
    latitude: 28.7606,
    longitude: -17.8816,
    heightM: 2396,
    timeZone: "Atlantic/Canary",
    region: "Rest of the world",
  },
  {
    id: "kittpeak",
    label: "Kitt Peak (KPNO, DESI)",
    latitude: 31.9583,
    longitude: -111.5967,
    heightM: 2096,
    timeZone: "America/Phoenix",
    region: "Rest of the world",
  },
  {
    id: "mtgraham",
    label: "Mt. Graham (LBT)",
    latitude: 32.7016,
    longitude: -109.8891,
    heightM: 3221,
    timeZone: "America/Phoenix",
    region: "Rest of the world",
  },
  {
    id: "apachepoint",
    label: "Apache Point (SDSS)",
    latitude: 32.7803,
    longitude: -105.8203,
    heightM: 2788,
    timeZone: "America/Denver",
    region: "Rest of the world",
  },
  {
    id: "palomar",
    label: "Palomar (ZTF)",
    latitude: 33.3563,
    longitude: -116.865,
    heightM: 1712,
    timeZone: "America/Los_Angeles",
    region: "Rest of the world",
  },
  {
    id: "sanpedromartir",
    label: "San Pedro Mártir (OAN)",
    latitude: 31.0447,
    longitude: -115.4637,
    heightM: 2830,
    timeZone: "America/Tijuana",
    region: "Rest of the world",
  },
  {
    id: "calaralto",
    label: "Calar Alto (CAHA)",
    latitude: 37.2236,
    longitude: -2.5463,
    heightM: 2168,
    timeZone: "Europe/Madrid",
    region: "Rest of the world",
  },
  {
    id: "sutherland",
    label: "Sutherland (SALT, SAAO)",
    latitude: -32.3794,
    longitude: 20.8107,
    heightM: 1798,
    timeZone: "Africa/Johannesburg",
    region: "Rest of the world",
  },
  {
    id: "sidingspring",
    label: "Siding Spring (AAT)",
    latitude: -31.2733,
    longitude: 149.0644,
    heightM: 1164,
    timeZone: "Australia/Sydney",
    region: "Rest of the world",
  },
];

export const DEFAULT_OBSERVATORY_ID = "paranal";

/** Select options grouped by region, Chile first. */
export const OBSERVATORY_OPTIONS = (
  ["Chile", "Rest of the world"] as ObservatoryRegion[]
).map((region) => ({
  label: region,
  title: region,
  options: OBSERVATORIES.filter((o) => o.region === region).map((o) => ({
    label: o.label,
    value: o.id,
  })),
}));

export function getObservatory(id: string): Observatory {
  return OBSERVATORIES.find((o) => o.id === id) ?? OBSERVATORIES[0];
}
