/**
 * Night-time visibility of a fixed target from a ground observatory.
 *
 * Positions come from astronomy-engine. The target's J2000 (ICRS) RA/Dec is
 * precessed to the equator of date before converting to altitude/azimuth;
 * skipping that step would be off by ~0.3° in 2026.
 */

import {
  AngleBetween,
  Body,
  Equator,
  EquatorFromVector,
  Horizon,
  Illumination,
  Observer,
  RotateVector,
  Rotation_EQJ_EQD,
  SearchAltitude,
  SearchRiseSet,
  Spherical,
  VectorFromSphere,
} from "astronomy-engine";

import {
  type Observatory,
  OBSERVATORY_TIME_ZONE,
} from "@/app/lib/constants/observatories";

/** Sun altitude at which astronomical twilight begins/ends. */
const ASTRONOMICAL_TWILIGHT_DEG = -18;

/** Default altitude a target must reach to count as observable (airmass 2). */
export const DEFAULT_MIN_ALTITUDE_DEG = 30;

const SAMPLE_STEP_MIN = 5;
/** Plot a little past sunset/sunrise so the twilight bands have context. */
const PADDING_MIN = 45;

export interface AltitudeSample {
  time: Date;
  altitude: number;
  /** Null when the target is at or below the horizon. */
  airmass: number | null;
  moonAltitude: number;
  sunAltitude: number;
}

export interface NightVisibility {
  sunset: Date;
  sunrise: Date;
  /** Start and end of astronomical night (Sun below −18°). */
  duskAstronomical: Date;
  dawnAstronomical: Date;
  samples: AltitudeSample[];
  /**
   * Stretch of astronomical night with the target above the minimum
   * altitude; null when it never gets there during the night.
   */
  window: { start: Date; end: Date } | null;
  /** Highest point the target reaches during astronomical night. */
  best: { time: Date; altitude: number; airmass: number | null } | null;
  /** Highest altitude the target ever reaches from this site (at transit). */
  maxPossibleAltitude: number;
  moon: {
    /** Illuminated fraction at local midnight, 0–1. */
    illumination: number;
    /** Angular distance to the target at local midnight, degrees. */
    separationDeg: number;
  };
}

/**
 * Relative airmass from Kasten & Young (1989); accurate to the horizon, where
 * sec(z) diverges. Null at or below the horizon.
 */
export function airmass(altitudeDeg: number): number | null {
  if (altitudeDeg <= 0) return null;
  const h = altitudeDeg;
  return (
    1 /
    (Math.sin((h * Math.PI) / 180) + 0.50572 * Math.pow(h + 6.07995, -1.6364))
  );
}

function toObserver(site: Observatory): Observer {
  return new Observer(site.latitude, site.longitude, site.heightM);
}

/** Target altitude in degrees for J2000 coordinates at `time`. */
function targetAltitude(
  raDeg: number,
  decDeg: number,
  time: Date,
  observer: Observer
): number {
  const j2000 = VectorFromSphere(new Spherical(decDeg, raDeg, 1), time);
  const ofDate = EquatorFromVector(RotateVector(Rotation_EQJ_EQD(time), j2000));
  return Horizon(time, observer, ofDate.ra, ofDate.dec, "normal").altitude;
}

function bodyAltitude(body: Body, time: Date, observer: Observer): number {
  const eq = Equator(body, time, observer, true, true);
  return Horizon(time, observer, eq.ra, eq.dec, "normal").altitude;
}

/**
 * Visibility during the night that starts on the evening of `nightOf`
 * (a calendar date; only its year/month/day are used).
 */
export function computeNightVisibility(
  raDeg: number,
  decDeg: number,
  site: Observatory,
  nightOf: Date,
  minAltitudeDeg = DEFAULT_MIN_ALTITUDE_DEG
): NightVisibility {
  const observer = toObserver(site);

  // Start every search at local solar noon so "the next sunset" is always
  // this evening's, whatever the time zone.
  const solarNoon = new Date(
    Date.UTC(nightOf.getFullYear(), nightOf.getMonth(), nightOf.getDate(), 12) -
      (site.longitude / 15) * 3_600_000
  );

  const sunset = SearchRiseSet(Body.Sun, observer, -1, solarNoon, 1);
  const dusk = SearchAltitude(
    Body.Sun,
    observer,
    -1,
    solarNoon,
    1,
    ASTRONOMICAL_TWILIGHT_DEG
  );
  if (!sunset || !dusk) {
    throw new Error(`No sunset from ${site.label} on this date`);
  }
  const dawn = SearchAltitude(
    Body.Sun,
    observer,
    +1,
    dusk,
    1,
    ASTRONOMICAL_TWILIGHT_DEG
  );
  const sunrise = SearchRiseSet(Body.Sun, observer, +1, sunset, 1);
  if (!dawn || !sunrise) {
    throw new Error(`No sunrise from ${site.label} on this date`);
  }

  const start = sunset.date.getTime() - PADDING_MIN * 60_000;
  const end = sunrise.date.getTime() + PADDING_MIN * 60_000;
  const samples: AltitudeSample[] = [];
  for (let t = start; t <= end; t += SAMPLE_STEP_MIN * 60_000) {
    const time = new Date(t);
    const altitude = targetAltitude(raDeg, decDeg, time, observer);
    samples.push({
      time,
      altitude,
      airmass: airmass(altitude),
      moonAltitude: bodyAltitude(Body.Moon, time, observer),
      sunAltitude: bodyAltitude(Body.Sun, time, observer),
    });
  }

  const night = samples.filter(
    (s) => s.time >= dusk.date && s.time <= dawn.date
  );
  // Altitude is single-peaked over one night, so the samples above the
  // limit form one contiguous run.
  const up = night.filter((s) => s.altitude >= minAltitudeDeg);
  const window =
    up.length > 0 ? { start: up[0].time, end: up[up.length - 1].time } : null;

  const peak = night.reduce<AltitudeSample | null>(
    (best, s) => (best === null || s.altitude > best.altitude ? s : best),
    null
  );

  const midnight = new Date(
    (sunset.date.getTime() + sunrise.date.getTime()) / 2
  );
  const target = VectorFromSphere(new Spherical(decDeg, raDeg, 1), midnight);
  const moonEq = Equator(Body.Moon, midnight, observer, false, true);

  return {
    sunset: sunset.date,
    sunrise: sunrise.date,
    duskAstronomical: dusk.date,
    dawnAstronomical: dawn.date,
    samples,
    window,
    best: peak
      ? { time: peak.time, altitude: peak.altitude, airmass: peak.airmass }
      : null,
    maxPossibleAltitude: 90 - Math.abs(site.latitude - decDeg),
    moon: {
      illumination: Illumination(Body.Moon, midnight).phase_fraction,
      separationDeg: AngleBetween(target, moonEq.vec),
    },
  };
}

/**
 * Calendar date of "tonight" in `timeZone`. Before 08:00 local the night in
 * progress (which started yesterday evening) is still the relevant one.
 */
export function tonightInTimeZone(timeZone: string, now = new Date()): Date {
  const shifted = new Date(now.getTime() - 8 * 3_600_000);
  const parts = new Intl.DateTimeFormat("en-CA", {
    timeZone,
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).formatToParts(shifted);
  const get = (type: string) =>
    Number(parts.find((p) => p.type === type)?.value);
  return new Date(get("year"), get("month") - 1, get("day"));
}

const chileTime = new Intl.DateTimeFormat("en-GB", {
  timeZone: OBSERVATORY_TIME_ZONE,
  hour: "2-digit",
  minute: "2-digit",
});

/** HH:mm in Chile time, where every offered observatory is. */
export function formatChileTime(date: Date): string {
  return chileTime.format(date);
}

/** One-line description of a night's visibility, times in Chile time. */
export function summarizeVisibility(
  v: NightVisibility,
  siteLabel: string,
  minAltitudeDeg = DEFAULT_MIN_ALTITUDE_DEG
): string {
  const hhmm = formatChileTime;
  if (v.window && v.best) {
    const am =
      v.best.airmass != null ? ` (airmass ${v.best.airmass.toFixed(2)})` : "";
    return `Observable ${hhmm(v.window.start)}–${hhmm(v.window.end)} above ${minAltitudeDeg}° · peaks at ${v.best.altitude.toFixed(0)}°${am} at ${hhmm(v.best.time)}`;
  }
  if (v.maxPossibleAltitude < minAltitudeDeg) {
    return v.maxPossibleAltitude <= 0
      ? `Never rises from ${siteLabel}.`
      : `Never rises above ${minAltitudeDeg}° from ${siteLabel} (max ${v.maxPossibleAltitude.toFixed(0)}° at transit).`;
  }
  const best = v.best
    ? ` Best: ${v.best.altitude.toFixed(0)}° at ${hhmm(v.best.time)}.`
    : "";
  return `Not above ${minAltitudeDeg}° during astronomical night on this date.${best}`;
}
