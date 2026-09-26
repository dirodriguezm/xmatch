/**
 * Observatories offered in the observability panel. Coordinates are geodetic
 * (WGS84), longitude east-positive, height in metres.
 */

export interface Observatory {
  id: string;
  label: string;
  latitude: number;
  longitude: number;
  heightM: number;
}

export const OBSERVATORIES: Observatory[] = [
  {
    id: "paranal",
    label: "Paranal (VLT)",
    latitude: -24.6272,
    longitude: -70.4042,
    heightM: 2635,
  },
  {
    id: "pachon",
    label: "Cerro Pachón (Rubin, Gemini-S)",
    latitude: -30.2446,
    longitude: -70.7494,
    heightM: 2715,
  },
  {
    id: "lascampanas",
    label: "Las Campanas (Magellan)",
    latitude: -29.0146,
    longitude: -70.6926,
    heightM: 2380,
  },
  {
    id: "lasilla",
    label: "La Silla (ESO)",
    latitude: -29.2567,
    longitude: -70.7377,
    heightM: 2400,
  },
];

export const DEFAULT_OBSERVATORY_ID = "paranal";

/** Every observatory above is in Chile, so times are shown in Chile time. */
export const OBSERVATORY_TIME_ZONE = "America/Santiago";
